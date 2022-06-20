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
	"net/http"
	"net/http/httptest"

	"github.com/emicklei/go-restful/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var _ = Describe("test tracing filter", func() {
	var sc trace.SpanContext

	injectTraceContext := func(req *http.Request) {
		ctx := trace.ContextWithRemoteSpanContext(context.Background(), sc)
		ctx, _ = trace.NewNoopTracerProvider().Tracer(defaultServiceName).Start(ctx, "test")
		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	}

	BeforeEach(func() {
		sc = trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: trace.TraceID{0x01},
			SpanID:  trace.SpanID{0x01},
		})
		otel.SetTextMapPropagator(propagation.TraceContext{})
	})

	Context("testing for normal path", func() {

		It("spans should be same", func() {
			server, request := testWebService("/test-path", RestfulFilter(defaultServiceName), func(req *restful.Request, resp *restful.Response) {
				span := trace.SpanFromContext(req.Request.Context())
				Expect(sc.TraceID()).Should(Equal(span.SpanContext().TraceID()))
				resp.WriteHeader(http.StatusOK)
			})
			injectTraceContext(request)
			server.ServeHTTP(httptest.NewRecorder(), request)
		})
	})

	Context("testing for ignored path", func() {
		It("span should be empty", func() {
			server, request := testWebService("/test-path", RestfulFilter(defaultServiceName, "test-path"), func(req *restful.Request, resp *restful.Response) {
				span := trace.SpanFromContext(req.Request.Context())
				Expect(span.SpanContext()).Should(Equal(trace.SpanContext{}))
				Expect(span.IsRecording()).Should(BeFalse())
				resp.WriteHeader(http.StatusOK)
			})
			injectTraceContext(request)
			server.ServeHTTP(httptest.NewRecorder(), request)
		})
	})
})
