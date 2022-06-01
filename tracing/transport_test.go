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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/net/context"
)

var _ = Describe("testing for WrapTransport", func() {
	Context("when nil transport is provided", func() {
		It("should use the default transport", func() {
			rt := WrapTransport(nil)
			tp, ok := rt.(*Transport)
			Expect(ok).Should(BeTrue())
			Expect(tp.originalRT).Should(Equal(http.DefaultTransport))
		})
	})
	Context("when special transport is provided", func() {
		var (
			originalTp *http.Transport
			rt         http.RoundTripper
		)
		BeforeEach(func() {
			originalTp = &http.Transport{}
			rt = WrapTransport(originalTp)
		})
		When("wrapping once", func() {
			It("should use the special transport", func() {
				tp, ok := rt.(*Transport)
				Expect(ok).Should(BeTrue())
				Expect(tp.originalRT).Should(Equal(originalTp))
			})
		})
		When("wrapping multiple times", func() {
			It("should wrap only once", func() {
				rt2 := WrapTransport(rt)
				tp2, ok := rt2.(*Transport)
				Expect(ok).Should(BeTrue())
				Expect(tp2.originalRT).Should(Equal(originalTp))
				Expect(tp2.originalRT).ShouldNot(Equal(rt))
			})
		})
	})
	Context("when handler request", func() {
		var shutdown []func()
		var client http.Client

		testHttpServer := func(handler func(http.ResponseWriter, *http.Request)) *http.Request {
			testServer := httptest.NewServer(http.HandlerFunc(handler))
			shutdown = append(shutdown, func() {
				testServer.Close()
			})
			testRequest, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
			Expect(err).ShouldNot(HaveOccurred())
			return testRequest
		}

		BeforeEach(func() {
			client = http.Client{Transport: WrapTransport(nil)}
			otel.SetTracerProvider(trace.NewTracerProvider(
				trace.WithSampler(trace.AlwaysSample()),
			))
			otel.SetTextMapPropagator(propagation.TextMapPropagator(propagation.TraceContext{}))
		})

		It("should have a Tracing header", func() {
			body := []byte("ok")
			req := testHttpServer(func(response http.ResponseWriter, request *http.Request) {
				Expect(request.Header.Get("Traceparent")).ShouldNot(BeEmpty())
				_, _ = response.Write(body)
			})

			resp, err := client.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			defer resp.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(ioutil.ReadAll(resp.Body)).Should(Equal(body))
		})

		It("need to cancel the request when it times out", func() {
			body := []byte("ok")
			req := testHttpServer(func(response http.ResponseWriter, request *http.Request) {
				Expect(request.Header.Get("Traceparent")).ShouldNot(BeEmpty())
				time.Sleep(2 * time.Second)
				response.Write(body)
			})

			ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*1)
			defer cancelFunc()
			req = req.WithContext(ctx)

			_, err := client.Do(req)
			Expect(strings.Contains(err.Error(), "deadline exceeded")).Should(BeTrue())
		})
	})
})

func Test_defaultSpanNameFormatter(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{
			url:  "https://a.com/a/b/c?a=b",
			want: "GET /a/b/c",
		},
		{
			url:  "https://a.com/a/b/c/?a=b",
			want: "GET /a/b/c/",
		},
	}
	g := NewGomegaWithT(t)
	for _, tt := range tests {
		url, err := url.Parse(tt.url)
		g.Expect(err).ShouldNot(HaveOccurred())
		req := &http.Request{
			URL:    url,
			Method: http.MethodGet,
		}
		spanName := defaultSpanNameFormatter("", req)
		g.Expect(spanName).Should(Equal(tt.want))
	}
}
