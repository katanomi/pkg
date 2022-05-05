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
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

// TraceOption optional setting of Tracing instance
type TraceOption func(tracing *Tracing)

// ExporterConstructor Construct exporter by the specified config
type ExporterConstructor func(config *Config) (trace.SpanExporter, error)

// ResourceConstructor Construct resource by the specified config
type ResourceConstructor func(config *Config) (*resource.Resource, error)

// TraceProviderConstructor Construct `TraceProvider` by the specified config
type TraceProviderConstructor func(config *Config) (*trace.TracerProvider, error)

// WithServiceName Configures the service name for Tracing instance
func WithServiceName(name string) TraceOption {
	return func(tracing *Tracing) {
		if name != "" {
			tracing.ServiceName = name
		}
	}
}

// WithExporter Configures the exporter backend for `Tracing` instance.
func WithExporter(f ExporterConstructor) TraceOption {
	return func(tracing *Tracing) {
		if f != nil {
			tracing.exporterConstructor = f
		}
	}
}

// WithResource Configures the resource for `Tracing` instance.
func WithResource(f ResourceConstructor) TraceOption {
	return func(tracing *Tracing) {
		if f != nil {
			tracing.resourceConstructor = f
		}
	}
}

// WithTraceProvider Configures the TraceProvider for `Tracing` instance
// If `TraceProvider` is specified, `WithExporter` and `WithResource` function will not work.
func WithTraceProvider(f TraceProviderConstructor) TraceOption {
	return func(tracing *Tracing) {
		if f != nil {
			tracing.traceProviderConstructor = f
		}
	}
}

// WithTracerProviderOption Configures the options for built-in traceProvider
func WithTracerProviderOption(ops ...trace.TracerProviderOption) TraceOption {
	return func(tracing *Tracing) {
		tracing.traceProviderOptions = append(tracing.traceProviderOptions, ops...)
	}
}

// WithTextMapPropagator Configures the propagator options for traceProvider
func WithTextMapPropagator(ops ...propagation.TextMapPropagator) TraceOption {
	return func(tracing *Tracing) {
		for _, item := range ops {
			if item != nil {
				tracing.Propagators = append(tracing.Propagators, item)
			}
		}
	}
}
