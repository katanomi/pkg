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
			log.Debugw("==> received request")
		}
		req.Request = req.Request.WithContext(logging.WithLogger(req.Request.Context(), log))
		chain.ProcessFilter(req, resp)
		if notHealthz(req.Request.URL.Path) {
			log.Debugw("<== returned a response")
		}
	}
}

func notHealthz(path string) bool {
	return path != "/healthz" && path != "/readyz"
}
