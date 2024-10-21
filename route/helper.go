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

package route

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
)

// wrapperF go restful wrapper func for http.HandlerFunc
func wrapperF(handler http.HandlerFunc) restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		handler(response.ResponseWriter, request.Request)
	}
}

// wrapperH go restful wrapper func for http.Handler
func wrapperH(handler http.Handler) restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		handler.ServeHTTP(response.ResponseWriter, request.Request)
	}
}

// NoOpFilter creates a default no operation filter
func NoOpFilter(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
	chain.ProcessFilter(req, res)
}
