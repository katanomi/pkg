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

package sharedmain

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cminformer "knative.dev/pkg/configmap/informer"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/logging/logkey"
	"knative.dev/pkg/system"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"yunion.io/x/pkg/util/wait"

	kclient "github.com/katanomi/pkg/client"
	klogging "github.com/katanomi/pkg/logging"

	// kmanager "github.com/katanomi/pkg/manager"
	kscheme "github.com/katanomi/pkg/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
)

func GetClientManager(ctx context.Context) (context.Context, *kclient.Manager) {
	clientManager := kclient.ManagerCtx(ctx)
	if clientManager == nil {
		clientManager = kclient.NewManager(ctx, nil, nil)
		ctx = kclient.WithManager(ctx, clientManager)
	}
	return ctx, clientManager
}

func GetConfigOrDie(ctx context.Context) (context.Context, *rest.Config) {
	cfg := injection.GetConfig(ctx)
	if cfg == nil {
		cfg = ctrl.GetConfigOrDie()
		ctx = injection.WithConfig(ctx, cfg)
	}
	cfg.WrapTransport = kclient.WrapTransportForTracing
	return ctx, cfg
}

// SetupLoggerOrDie sets up the logger using the config from the given context
// and returns a logger and atomic level, or dies by calling log.Fatalf.
func SetupLoggerOrDie(ctx context.Context, component string) (*zap.SugaredLogger, zap.AtomicLevel) {
	loggingConfig, err := GetLoggingConfig(ctx)
	if err != nil {
		log.Fatal("Error reading/parsing logging configuration: ", err)
	}
	l, level := logging.NewLoggerFromConfig(loggingConfig, component)

	// If PodName is injected into the env vars, set it on the logger.
	// This is needed for HA components to distinguish logs from different
	// pods.
	if pn := os.Getenv("POD_NAME"); pn != "" {
		l = l.With(zap.String(logkey.Pod, pn))
	}

	return l, level
}

func GetLoggingConfig(ctx context.Context) (*logging.Config, error) {
	loggingConfigMap := &corev1.ConfigMap{}

	key := types.NamespacedName{Name: logging.ConfigMapName(), Namespace: system.Namespace()}

	directClt, err := client.New(injection.GetConfig(ctx), client.Options{Scheme: kscheme.Scheme(ctx)})
	if err != nil {
		return nil, err
	}

	// These timeout and retry interval are set by heuristics.
	// e.g. istio sidecar needs a few seconds to configure the pod network.
	var lastErr error
	if err := wait.PollImmediate(1*time.Second, 5*time.Second, func() (bool, error) {
		lastErr = directClt.Get(ctx, key, loggingConfigMap)
		return lastErr == nil || apierrors.IsNotFound(lastErr), nil
	}); err != nil {
		return nil, fmt.Errorf("timed out waiting for the condition: %w", lastErr)

	}

	if loggingConfigMap == nil {
		return logging.NewConfigFromMap(nil)
	}
	return logging.NewConfigFromConfigMap(loggingConfigMap)
}

// WatchLoggingConfigOrDie establishes a watch of the logging config or dies by
// calling log.Fatalw. Note, if the config does not exist, it will be defaulted
// and this method will not die.
func WatchLoggingConfigOrDie(ctx context.Context, cmw *cminformer.InformedWatcher, logger *zap.SugaredLogger, lvlMGR *klogging.LevelManager) {
	if _, err := kubeclient.Get(ctx).CoreV1().ConfigMaps(system.Namespace()).Get(ctx, logging.ConfigMapName(),
		metav1.GetOptions{}); err == nil {
		cmw.Watch(logging.ConfigMapName(), lvlMGR.Update())
	} else if !apierrors.IsNotFound(err) {
		logger.Fatalw("Error reading ConfigMap "+logging.ConfigMapName(), zap.Error(err))
	}
}
