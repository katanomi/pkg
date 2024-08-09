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

package logging

import (
	"context"
	"fmt"
	"strings"

	"github.com/katanomi/pkg/client"
	"github.com/katanomi/pkg/plugin/client/base"

	"go.uber.org/zap/zapcore"

	"net/http"
	"net/http/httputil"

	"github.com/emicklei/go-restful/v3"
	"go.uber.org/zap"
	"knative.dev/pkg/logging"
)

// Filter prometheus metric filter for go restful
func Filter(logger *zap.SugaredLogger) func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		// injects logger into context

		log := logger.With("verb", req.Request.Method, "path", req.Request.URL.Path, "from", req.Request.RemoteAddr)

		if notHealthz(req.Request.URL.Path) {
			log.Debugw(fmt.Sprintf("=> %s %s", req.Request.Method, req.Request.URL.String()),
				"path-params", req.PathParameters(), "query-params", req.Request.URL.Query())
		}
		req.Request = req.Request.WithContext(logging.WithLogger(req.Request.Context(), log))
		chain.ProcessFilter(req, resp)
		if notHealthz(req.Request.URL.Path) {
			if resp != nil {
				log = log.With("statusCode", resp.StatusCode())
			}
			log.Debugf("<= %s %d %s", req.Request.Method, resp.StatusCode(), req.Request.URL.String())

			if resp.StatusCode() > 300 {
				if log.Level().Enabled(zapcore.DebugLevel) {
					dumpRequestAndResponseWriter(req.Request.Context(), req.Request, resp)
				}
			}
		}
	}
}

func dumpRequestAndResponseWriter(ctx context.Context, req *http.Request, resp *restful.Response) {
	logger := logging.FromContext(ctx)

	dumpBody := true
	if !isTextRequest(req) {
		dumpBody = false
	}

	// dump request
	requestBts, err := httputil.DumpRequest(req, dumpBody)
	if err != nil {
		logger.Errorw("dump request error", err)
		return
	}
	reqStr := string(requestBts)

	// redaction
	items := []string{
		base.PluginSecretHeader,
		client.AuthorizationHeader,
		"Cookie",
	}
	for _, item := range items {
		value := req.Header.Get(item)
		if value != "" {
			reqStr = strings.Replace(reqStr, value, "******", -1)
		}
	}

	fmt.Printf(`=====================DUMP REQUEST and Response============================
~~~ REQUEST ~~~
%s
--------------------------------------------------------
~~~ RESPONSE ~~~
STATUS: %d
HEADERS: %#v
==========================================================================
`, reqStr, resp.StatusCode(), resp.Header())
}

func isTextRequest(req *http.Request) bool {
	text := false

	contentType := strings.ToLower(req.Header.Get("content-type"))
	txtTypes := []string{
		"text/",
		"application/json",
		"application/xml",
		"application/x-www-form-urlencoded",
	}

	for _, t := range txtTypes {
		if strings.HasPrefix(contentType, t) {
			text = true
		}
	}

	return text
}

func notHealthz(path string) bool {
	return path != "/healthz" && path != "/readyz"
}
