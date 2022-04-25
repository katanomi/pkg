package tracing

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// defaultSpanNameFormatter Format span name according to http.Request
func defaultSpanNameFormatter(_ string, r *http.Request) string {
	return r.Method + " " + r.URL.EscapedPath()
}

// WrapTransportForTracing Wrap Transport for Tracing
// When rt is nil, default transport will be used
func WrapTransportForTracing(rt http.RoundTripper) http.RoundTripper {
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
