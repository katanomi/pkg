/*
Copyright 2021 The Katanomi Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tracing

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	kscheme "github.com/katanomi/pkg/scheme"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"yunion.io/x/pkg/util/wait"

	"github.com/emicklei/go-restful/v3"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"

	"go.uber.org/zap"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	cminformer "knative.dev/pkg/configmap/informer"
	"knative.dev/pkg/system"

	"github.com/opentracing/opentracing-go/ext"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

const (
	configMapNameEnv = "CONFIG_TRACING_NAME"
)

type Manager struct {
	configManager *configManager
	tracers       sync.Map
	logger        *zap.SugaredLogger
}

type tracerBinding struct {
	tracer opentracing.Tracer
	closer io.Closer
}

func SetupTracingOrDie(ctx context.Context, logger *zap.SugaredLogger) (*Manager, error) {
	tracingConfigMap := &corev1.ConfigMap{}

	key := types.NamespacedName{Name: logging.ConfigMapName(), Namespace: system.Namespace()}
	directClt, err := client.New(injection.GetConfig(ctx), client.Options{Scheme: kscheme.Scheme(ctx)})
	if err != nil {
		return nil, err
	}

	// These timeout and retry interval are set by heuristics.
	// e.g. istio sidecar needs a few seconds to configure the pod network.
	var lastErr error
	if err := wait.PollImmediate(1*time.Second, 5*time.Second, func() (bool, error) {
		lastErr := directClt.Get(ctx, key, tracingConfigMap)
		fmt.Println("err?", lastErr, "key", key)
		return lastErr == nil || apierrors.IsNotFound(lastErr), nil
	}); err != nil {
		return nil, fmt.Errorf("timed out waiting for the condition: %w", lastErr)
	}

	configManager, err := ParseConfig(tracingConfigMap.Data)
	if err != nil {
		return nil, err
	}

	manager := &Manager{
		configManager: configManager,
		logger:        logger,
	}

	return manager, nil
}

func (t *Manager) Tracer(name string) (opentracing.Tracer, error) {
	c := t.configManager.Get(name)

	if c == nil {
		t.logger.Warnf("tracing config not found for %s", name)
		return nil, nil
	}

	v, loaded := t.tracers.Load(name)
	if loaded {
		binding := v.(tracerBinding)

		return binding.tracer, nil
	}

	cfg := jaegercfg.Configuration{
		ServiceName: name,
		Sampler: &jaegercfg.SamplerConfig{
			Type:              c.SampleType,
			Param:             c.SampleParam,
			SamplingServerURL: c.SampleServerURL,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LocalAgentHostPort: c.JaegerUrl,
			LogSpans:           c.LogSpan,
		},
	}

	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jaeger.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)

	if err != nil {
		return nil, err
	}

	binding := tracerBinding{
		tracer: tracer,
		closer: closer,
	}

	t.tracers.Store(name, binding)

	return binding.tracer, nil
}

func (t *Manager) updateTracer(configMap *corev1.ConfigMap) {
	configManager, err := ParseConfig(configMap.Data)
	if err != nil {
		t.logger.Errorf("update tracer with new config error: %s", err.Error())
		return
	}

	t.configManager = configManager
	t.Sync()
}

func (t *Manager) Sync() {
	t.tracers.Range(func(key, value interface{}) bool {
		binding := value.(tracerBinding)
		binding.closer.Close()
		return true
	})

	// clean all tracer, new tracer will be created at using
	t.tracers = sync.Map{}
}

func (t *Manager) configMapName() string {
	if cm := os.Getenv(configMapNameEnv); cm != "" {
		return cm
	}
	return "config-tracing"
}

// Filter tracing filter for go restful, follow opentracing
func Filter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	tracer := opentracing.GlobalTracer()
	spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Request.Header))

	span, ctx := opentracing.StartSpanFromContext(req.Request.Context(), "handle request", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, req.Request.URL.String())
	ext.HTTPMethod.Set(span, req.Request.Method)

	req.Request = req.Request.WithContext(ctx)

	chain.ProcessFilter(req, resp)
}

// WatchLoggingConfigOrDie establishes a watch of the logging config or dies by
// calling log.Fatalw. Note, if the config does not exist, it will be defaulted
// and this method will not die.
func (t *Manager) WatchTracingConfigOrDie(ctx context.Context, cmw *cminformer.InformedWatcher, logger *zap.SugaredLogger) {
	configMapName := t.configMapName()

	if _, err := kubeclient.Get(ctx).CoreV1().ConfigMaps(system.Namespace()).Get(ctx, configMapName,
		metav1.GetOptions{}); err == nil {
		cmw.Watch(configMapName, t.updateTracer)
	} else if !apierrors.IsNotFound(err) {
		logger.Fatalw("Error reading ConfigMap "+configMapName, zap.Error(err))
	}
}
