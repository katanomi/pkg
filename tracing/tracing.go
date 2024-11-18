/*
Copyright 2021 The AlaudaDevops Authors.

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
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	traceApi "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"knative.dev/pkg/logging"
)

// NewTracing construct `Tracing` instance
// `TraceOption` can be used to customize settings.
func NewTracing(logger *zap.SugaredLogger, ops ...TraceOption) *Tracing {
	tracing := &Tracing{
		logger:        logger,
		ServiceName:   defaultServiceName,
		ConfigMapName: ConfigMapName(),
	}

	for _, op := range ops {
		if op != nil {
			op(tracing)
		}
	}
	return tracing
}

// Tracing describe an entity that watching configuration file changes and
// maintain the global tracing.
type Tracing struct {
	applyOnce sync.Once
	lock      sync.Mutex
	logger    *zap.SugaredLogger

	ServiceName   string
	ConfigMapName string

	exporterConstructor      ExporterConstructor
	resourceConstructor      ResourceConstructor
	traceProviderConstructor TraceProviderConstructor

	traceProviderOptions []trace.TracerProviderOption
	Propagators          []propagation.TextMapPropagator
}

// ApplyConfig Apply configuration and reinitialize global tracing.
// This method will be triggered when the configuration changes.
// If an error occurs, the initialization process will be skipped.
func (t *Tracing) ApplyConfig(cfg *Config) {
	t.applyOnce.Do(func() {
		otel.SetTracerProvider(traceApi.NewNoopTracerProvider())

		propagators := propagation.TextMapPropagator(propagation.TraceContext{})
		if len(t.Propagators) > 0 {
			propagators = propagation.NewCompositeTextMapPropagator(t.Propagators...)
		}
		otel.SetTextMapPropagator(propagators)
	})
	t.lock.Lock()
	defer t.lock.Unlock()

	tp, err := t.traceProvider(cfg)
	if err != nil {
		return
	}

	otel.SetTracerProvider(tp)
}

// traceProvider construct traceProvider according to the specified configuration.
func (t *Tracing) traceProvider(cfg *Config) (tp traceApi.TracerProvider, err error) {
	if t.traceProviderConstructor != nil {
		tp, err = t.traceProviderConstructor(cfg)
		if err != nil {
			t.logger.Errorw("Tracing construct trace provider err",
				"err", err,
				"config", cfg,
			)
		}
		return
	}

	if cfg == nil {
		return traceApi.NewNoopTracerProvider(), nil
	}

	exp, err := t.exporter(cfg)
	if err != nil {
		return
	}

	res, err := t.resource(cfg)
	if err != nil {
		return
	}

	ops := []trace.TracerProviderOption{
		trace.WithBatcher(exp),
		trace.WithResource(res),
		trace.WithSampler(trace.TraceIDRatioBased(cfg.SamplingRatio)),
	}
	ops = append(ops, t.traceProviderOptions...)
	tp = trace.NewTracerProvider(ops...)
	return tp, nil
}

// exporter construct backend exporter according to the specified configuration.
func (t *Tracing) exporter(config *Config) (exporter trace.SpanExporter, err error) {
	if t.exporterConstructor != nil {
		exporter, err = t.exporterConstructor(config)
		if err != nil {
			t.logger.Errorw("Tracing construct exporter err",
				"err", err,
				"config", config,
			)
			return nil, err
		}
	}

	if exporter == nil {
		switch config.Backend {
		case ExporterBackendJaeger:
			exporter, err = t.constructJaegerExporter(config.Jaeger)
		case ExporterBackendZipkin:
			exporter, err = t.constructZipkinExporter(config.Zipkin)
		case ExporterBackendCustom:
			t.logger.Errorw("Use WithExporter function to customize exporter",
				"err", err,
				"config", config,
			)
		default:
			logging.FromContext(context.TODO()).Warnw("unknown tracing backend", "backend", config.Backend)
		}
	}

	return exporter, nil
}

// constructZipkinExporter construct zipkin exporter according to the specified configuration.
func (t *Tracing) constructZipkinExporter(cfg ZipkinConfig) (exporter trace.SpanExporter, err error) {
	exporter, err = zipkin.New(cfg.Url)
	if err != nil {
		t.logger.Errorw("Tracing construct zipkin exporter error",
			"err", err,
			"config", cfg,
		)
	}
	return exporter, err
}

// constructJaegerExporter construct jaeger exporter according to the specified configuration.
func (t *Tracing) constructJaegerExporter(cfg JaegerConfig) (exporter trace.SpanExporter, err error) {
	ops := make([]jaeger.AgentEndpointOption, 0)
	if cfg.Host != "" {
		ops = append(ops, jaeger.WithAgentHost(cfg.Host))
	}
	if cfg.Port != "" {
		ops = append(ops, jaeger.WithAgentPort(cfg.Port))
	}
	if cfg.MaxPacketSize > 0 {
		ops = append(ops, jaeger.WithMaxPacketSize(cfg.MaxPacketSize))
	}
	if cfg.DisableAttemptReconnecting {
		ops = append(ops, jaeger.WithDisableAttemptReconnecting())
	}
	if cfg.AttemptReconnectInterval > 0 {
		ops = append(ops, jaeger.WithAttemptReconnectingInterval(cfg.AttemptReconnectInterval))
	}
	exporter, err = jaeger.New(
		jaeger.WithAgentEndpoint(ops...),
	)
	if err != nil {
		t.logger.Errorw("Tracing construct jaeger exporter error",
			"err", err,
			"config", cfg,
		)
	}
	return exporter, err
}

// resource construct Resource according to the specified configuration.
func (t *Tracing) resource(config *Config) (r *resource.Resource, err error) {
	if t.resourceConstructor != nil {
		r, err = t.resourceConstructor(config)
		if err != nil {
			t.logger.Errorw("Tracing construct resource err",
				"err", err,
				"config", config,
			)
			return nil, err
		}
	}

	if r == nil {
		r, _ = resource.Merge(
			resource.Default(),
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(t.ServiceName),
			),
		)
	}
	return r, nil
}
