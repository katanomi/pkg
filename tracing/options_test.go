package tracing

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

var _ = Describe("testing for customizing service name", func() {
	Context("specify a service name", func() {
		tracing := &Tracing{ServiceName: defaultServiceName}
		It("service name should be modified", func() {
			WithServiceName("hello")(tracing)
			Expect(tracing.ServiceName).Should(Equal("hello"))
		})
	})

	Context("specify an empty service name", func() {
		tracing := &Tracing{ServiceName: defaultServiceName}
		It("service name should not be modified", func() {
			WithServiceName("")(tracing)
			Expect(tracing.ServiceName).Should(Equal(defaultServiceName))
		})
	})
})

var _ = Describe("testing for customizing exporter", func() {
	var testConfig *Config

	BeforeEach(func() {
		testConfig = &Config{Enable: true}
	})

	Context("when a valid exporter is provided", func() {
		tracing := &Tracing{ServiceName: defaultServiceName}
		exporter := func(config *Config) (trace.SpanExporter, error) {
			Expect(config).Should(Equal(testConfig))
			return nil, nil
		}
		It("exporterConstructor property should be set", func() {
			WithExporter(exporter)(tracing)
			Expect(tracing.exporterConstructor).ShouldNot(BeNil())

			tracing.ApplyConfig(testConfig)

			r1, r2 := tracing.exporterConstructor(testConfig)
			Expect(r1).Should(BeNil())
			Expect(r2).Should(BeNil())
		})
	})

	Context("when a nil exporter is provided", func() {
		tracing := &Tracing{ServiceName: defaultServiceName}
		It("exporterConstructor property should not be set", func() {
			WithExporter(nil)(tracing)
			Expect(tracing.exporterConstructor).Should(BeNil())
		})
	})
})

var _ = Describe("testing for customizing resource", func() {
	var testConfig *Config

	BeforeEach(func() {
		testConfig = &Config{Enable: true}
	})

	Context("when a valid resource is provided", func() {
		tracing := &Tracing{ServiceName: defaultServiceName}
		resource := func(config *Config) (*resource.Resource, error) {
			Expect(config).Should(Equal(testConfig))
			return nil, nil
		}
		It("resourceConstructor property should be set", func() {
			WithResource(resource)(tracing)
			Expect(tracing.resourceConstructor).ShouldNot(BeNil())

			tracing.ApplyConfig(testConfig)

			r1, r2 := tracing.resourceConstructor(testConfig)
			Expect(r1).Should(BeNil())
			Expect(r2).Should(BeNil())
		})
	})

	Context("when a nil resource is provided", func() {
		tracing := &Tracing{ServiceName: defaultServiceName}
		It("resourceConstructor property should not be set", func() {
			WithResource(nil)(tracing)
			Expect(tracing.resourceConstructor).Should(BeNil())
		})
	})
})

var _ = Describe("testing for customizing traceProvider", func() {
	var testConfig *Config

	BeforeEach(func() {
		testConfig = &Config{Enable: true}
	})

	Context("when a valid traceProvider is provided", func() {
		tracing := &Tracing{ServiceName: defaultServiceName}
		traceProvider := func(config *Config) (*trace.TracerProvider, error) {
			Expect(config).Should(Equal(testConfig))
			return nil, nil
		}
		WithTraceProvider(traceProvider)(tracing)

		It("traceProviderConstructor property should be set", func() {
			Expect(tracing.traceProviderConstructor).ShouldNot(BeNil())

			tracing.ApplyConfig(testConfig)

			r1, r2 := tracing.traceProviderConstructor(testConfig)
			Expect(r1).Should(BeNil())
			Expect(r2).Should(BeNil())
		})

		Context("with a valid exporter is provided", func() {
			exporter := func(config *Config) (trace.SpanExporter, error) {
				panic("should not be executed")
			}
			resource := func(config *Config) (*resource.Resource, error) {
				Expect(config).Should(Equal(testConfig))
				return nil, nil
			}
			It("exporterConstructor property should be set", func() {
				WithExporter(exporter)(tracing)
				Expect(tracing.exporterConstructor).ShouldNot(BeNil())

				tracing.ApplyConfig(testConfig)
			})
			It("resourceConstructor property should be set", func() {
				WithResource(resource)(tracing)
				Expect(tracing.resourceConstructor).ShouldNot(BeNil())

				tracing.ApplyConfig(testConfig)
			})
		})
	})

	Context("with a nil traceProvider", func() {
		tracing := &Tracing{ServiceName: defaultServiceName}
		It("traceProviderConstructor property should not be set", func() {
			WithTraceProvider(nil)(tracing)
			Expect(tracing.traceProviderConstructor).Should(BeNil())
		})
	})
})

var _ = Describe("testing for customizing TraceProvider option", func() {
	Context("when always sample", func() {
		tracing := NewTracing(zap.NewNop().Sugar(),
			WithTracerProviderOption(trace.WithSampler(trace.AlwaysSample())),
		)
		It("The count of TracProvider options should be 1", func() {
			Expect(len(tracing.traceProviderOptions)).Should(Equal(1))
		})
		It("Should be Sampled", func() {
			tracing.ApplyConfig(getValidTracingConfig())
			_, span := otel.Tracer("xx").Start(context.Background(), "xx")
			Expect(span.SpanContext().IsSampled()).Should(BeTrue())
		})
	})

	Context("when never sample", func() {
		tracing := NewTracing(zap.NewNop().Sugar(),
			WithTracerProviderOption(trace.WithSampler(trace.NeverSample())),
		)
		It("The count of TracProvider options should be 1", func() {
			Expect(len(tracing.traceProviderOptions)).Should(Equal(1))
		})
		It("Should not be Sampled", func() {
			tracing.ApplyConfig(getValidTracingConfig())
			_, span := otel.Tracer("xx").Start(context.Background(), "xx")
			Expect(span.SpanContext().IsSampled()).ShouldNot(BeTrue())
		})
	})
})

var _ = Describe("testing for customizing propagator option", func() {
	testPropagator := propagation.TraceContext{}
	Context("when always sample", func() {
		tracing := NewTracing(zap.NewNop().Sugar(),
			WithTextMapPropagator(testPropagator),
		)
		It("Should be Sampled", func() {
			tracing.ApplyConfig(getValidTracingConfig())
			Expect(otel.GetTextMapPropagator()).
				Should(Equal(propagation.NewCompositeTextMapPropagator(testPropagator)))
		})
	})
})
