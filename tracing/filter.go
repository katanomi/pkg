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
	"strings"

	"github.com/emicklei/go-restful/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

// RestfulFilter Set the tracing middleware for go-restful web service framework.
// If ignorePaths param is specified, these paths will not be sampled.
func RestfulFilter(serviceName string, ignorePaths ...string) restful.FilterFunction {
	for index, item := range ignorePaths {
		ignorePaths[index] = strings.TrimPrefix(item, "/")
	}
	isIgnoredPath := func(req *restful.Request) bool {
		routePath := strings.TrimPrefix(req.SelectedRoutePath(), "/")
		for _, item := range ignorePaths {
			if item == routePath {
				return true
			}
		}
		return false
	}
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		if isIgnoredPath(req) {
			chain.ProcessFilter(req, resp)
			return
		}
		tracer := otel.GetTracerProvider().Tracer(serviceName)
		propagator := otel.GetTextMapPropagator()

		r := req.Request
		ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		route := req.SelectedRoutePath()
		spanName := r.Method + " " + route

		ctx, span := tracer.Start(ctx, spanName,
			trace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", r)...),
			trace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(serviceName, route, r)...),
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		// pass the span through the request context
		req.Request = req.Request.WithContext(ctx)

		chain.ProcessFilter(req, resp)

		attrs := semconv.HTTPAttributesFromHTTPStatusCode(resp.StatusCode())
		spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(resp.StatusCode())
		span.SetAttributes(attrs...)
		span.SetStatus(spanStatus, spanMessage)
	}
}
