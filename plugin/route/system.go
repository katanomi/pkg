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
	"net/http/pprof"

	"github.com/emicklei/go-restful/v3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type system struct {
}

func NewSystem() Route {
	return &system{}
}

func (s *system) Register(ws *restful.WebService) {
	// prometheus metrics
	ws.Route(ws.GET("/metrics").To(wrapperH(promhttp.Handler())))
	// golang profile
	ppprofPath := "/debug/pprof"
	ws.Route(ws.GET(ppprofPath).To(wrapperF(pprof.Index)))
	ws.Route(ws.GET(ppprofPath + "/cmdline").To(wrapperF(pprof.Index)))
	ws.Route(ws.GET(ppprofPath + "/profile").To(wrapperF(pprof.Profile)))
	ws.Route(ws.GET(ppprofPath + "/symbol").To(wrapperF(pprof.Symbol)))
	ws.Route(ws.GET(ppprofPath + "/trace").To(wrapperF(pprof.Trace)))
	ws.Route(ws.GET(ppprofPath + "/allocs").To(wrapperH(pprof.Handler("allocs"))))
	ws.Route(ws.GET(ppprofPath + "/block").To(wrapperH(pprof.Handler("block"))))
	ws.Route(ws.GET(ppprofPath + "/goroutine").To(wrapperH(pprof.Handler("goroutine"))))
	ws.Route(ws.GET(ppprofPath + "/heap").To(wrapperH(pprof.Handler("heap"))))
	ws.Route(ws.GET(ppprofPath + "/mutex").To(wrapperH(pprof.Handler("mutex"))))
	ws.Route(ws.GET(ppprofPath + "/threadcreate").To(wrapperH(pprof.Handler("threadcreate"))))
}
