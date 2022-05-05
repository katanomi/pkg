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
	"net/http"

	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// defaultSpanNameFormatter Format span name according to http.Request
func defaultSpanNameFormatter(_ string, r *http.Request) string {
	return r.Method + " " + r.URL.EscapedPath()
}

// WrapTransportForRestyClient Specifically wrapped for the go-resty client for tracking
// Warning: Because some methods in go-resty (such as SetTLSClientConfig、SetProxy)
// rely on *http.Transport, this method must be called after initialization.
//
// correct example:
//		restyClient := resty.New()
//		restyClient.SetTLSClientConfig(&tls.Config{
//			InsecureSkipVerify: true,
//		})
//		tracing.WrapTransportForRestyClient(restyClient)
//
// wrong example:
//		restyClient := resty.New()
//      tracing.WrapTransportForRestyClient(restyClient)
//		restyClient.SetTLSClientConfig(&tls.Config{
//			InsecureSkipVerify: true,
//		})
func WrapTransportForRestyClient(client *resty.Client) {
	if client == nil {
		return
	}
	if httpClient := client.GetClient(); httpClient != nil {
		client.SetTransport(WrapTransport(httpClient.Transport))
	} else {
		client.SetTransport(WrapTransport(nil))
	}
}

// WrapTransport Wrap Transport for Tracing
// When rt is nil, default transport will be used
func WrapTransport(rt http.RoundTripper) http.RoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return &Transport{
		originalRT: rt,
		Transport: otelhttp.NewTransport(rt,
			otelhttp.WithSpanNameFormatter(defaultSpanNameFormatter),
		),
	}
}

// Transport for tracing
type Transport struct {
	// originalRT The original RoundTripper
	// Because the original RoundTripper is wrapped, other
	// interfaces implemented by original RoundTripper will fail.
	//
	// It is possible to customize the methods and make certain
	// implementations available through assertions.
	originalRT http.RoundTripper
	*otelhttp.Transport
}

// RoundTrip creates a Span and propagates its context via the provided request's headers
// before handing the request to the configured base RoundTripper. The created span will
// end when the response body is closed or when a read from the body returns io.EOF.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.Transport.RoundTrip(req)
}

// CancelRequest cancels an in-flight request by closing its connection.
// It works when the original RoundTripper implementation canceler interface.
func (t *Transport) CancelRequest(req *http.Request) {
	type canceler interface {
		CancelRequest(*http.Request)
	}

	if rt, ok := t.originalRT.(canceler); ok {
		rt.CancelRequest(req)
	}
}
