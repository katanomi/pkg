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

// Package tracing contains useful functionality for tracing.
package tracing

import (
	"io"

	"github.com/emicklei/go-restful/v3"
	"github.com/katanomi/pkg/plugin/config"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

// Config global opentracing config
func Config(c *config.TraceConfig) (io.Closer, error) {
	if !c.Enable {
		return nil, nil
	}

	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:              c.SampleType,
			Param:             c.SampleParam,
			SamplingServerURL: c.SampleServerURL,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LocalAgentHostPort: c.JaegerUrl,
			LogSpans:           true,
		},
	}

	// Initialize tracer with a logger and a metrics factory
	closer, err := cfg.InitGlobalTracer(
		"plugin",
		jaegercfg.Logger(jaeger.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
	if err != nil {
		return nil, err
	}

	return closer, nil
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
