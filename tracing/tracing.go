package tracing

import (
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
)

type TraceOption func(tracing *Tracing)
type ExporterConstructor func(config *Config) (trace.SpanExporter, error)
type ResourceConstructor func(config *Config) (*resource.Resource, error)
type TraceProviderConstructor func(config *Config) (*trace.TracerProvider, error)

func WithServiceName(name string) TraceOption {
	return func(tracing *Tracing) {
		tracing.ServiceName = name
	}
}

func WithExporter(f ExporterConstructor) TraceOption {
	return func(tracing *Tracing) {
		if f != nil {
			tracing.exporterConstructor = f
		}
	}
}

func WithResource(f ResourceConstructor) TraceOption {
	return func(tracing *Tracing) {
		if f != nil {
			tracing.resourceConstructor = f
		}
	}
}

func WithTraceProvider(f TraceProviderConstructor) TraceOption {
	return func(tracing *Tracing) {
		if f != nil {
			tracing.traceProviderConstructor = f
		}
	}
}

func WithTracerProviderOption(ops ...trace.TracerProviderOption) TraceOption {
	return func(tracing *Tracing) {
		tracing.traceProviderOptions = append(tracing.traceProviderOptions, ops...)
	}
}

func WithTextMapPropagator(ops ...propagation.TextMapPropagator) TraceOption {
	return func(tracing *Tracing) {
		for _, item := range ops {
			if item != nil {
				tracing.Propagators = append(tracing.Propagators, item)
			}
		}
	}
}

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

type Tracing struct {
	applyOnce sync.Once
	logger    *zap.SugaredLogger

	ServiceName   string
	ConfigMapName string

	exporterConstructor      ExporterConstructor
	resourceConstructor      ResourceConstructor
	traceProviderConstructor TraceProviderConstructor

	traceProviderOptions []trace.TracerProviderOption
	Propagators          []propagation.TextMapPropagator
}

func (t *Tracing) ApplyConfig(cfg *Config) {
	t.applyOnce.Do(func() {
		otel.SetTracerProvider(traceApi.NewNoopTracerProvider())

		propagators := propagation.TextMapPropagator(propagation.TraceContext{})
		if len(t.Propagators) > 0 {
			propagators = propagation.NewCompositeTextMapPropagator(t.Propagators...)
		}
		otel.SetTextMapPropagator(propagators)
	})

	tp, err := t.traceProvider(cfg)
	if err != nil {
		return
	}

	otel.SetTracerProvider(tp)
}

func (t *Tracing) traceProvider(cfg *Config) (tp *trace.TracerProvider, err error) {
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
	}
	ops = append(ops, t.traceProviderOptions...)
	tp = trace.NewTracerProvider(ops...)
	return tp, nil
}

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
		}
	}

	return exporter, nil
}

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
