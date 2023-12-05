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

// Package route contains useful functionality for the package route
package route

import (
	"context"
	"net/http/pprof"

	"github.com/emicklei/go-restful/v3"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/katanomi/pkg/config"
)

type system struct {
	Context       context.Context
	ConfigManager *config.Manager
}

func NewSystem(ctx context.Context) Route {
	return &system{
		Context:       ctx,
		ConfigManager: config.KatanomiConfigManager(ctx),
	}
}

func (s *system) Register(ws *restful.WebService) {
	// prometheus metrics
	ws.Route(ws.GET("/metrics").To(wrapperH(promhttp.Handler())))

	// set web service to accept and return JSON.
	ws.Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	// golang profile
	ppprofPath := "/debug/pprof"

	configFilter := NoOpFilter
	if s.ConfigManager != nil && s.Context != nil {
		configFilter = config.ConfigFilter(s.Context, s.ConfigManager, config.PprofEnabledKey, config.ConfigFilterNotFoundWhenNotTrue)
	}

	ws.Route(ws.GET(ppprofPath).Filter(configFilter).To(wrapperF(pprof.Index)))
	ws.Route(ws.GET(ppprofPath + "/cmdline").Filter(configFilter).To(wrapperF(pprof.Index)))
	ws.Route(ws.GET(ppprofPath + "/profile").Filter(configFilter).To(wrapperF(pprof.Profile)))
	ws.Route(ws.GET(ppprofPath + "/symbol").Filter(configFilter).To(wrapperF(pprof.Symbol)))
	ws.Route(ws.GET(ppprofPath + "/trace").Filter(configFilter).To(wrapperF(pprof.Trace)))
	ws.Route(ws.GET(ppprofPath + "/allocs").Filter(configFilter).To(wrapperH(pprof.Handler("allocs"))))
	ws.Route(ws.GET(ppprofPath + "/block").Filter(configFilter).To(wrapperH(pprof.Handler("block"))))
	ws.Route(ws.GET(ppprofPath + "/goroutine").Filter(configFilter).To(wrapperH(pprof.Handler("goroutine"))))
	ws.Route(ws.GET(ppprofPath + "/heap").Filter(configFilter).To(wrapperH(pprof.Handler("heap"))))
	ws.Route(ws.GET(ppprofPath + "/mutex").Filter(configFilter).To(wrapperH(pprof.Handler("mutex"))))
	ws.Route(ws.GET(ppprofPath + "/threadcreate").Filter(configFilter).To(wrapperH(pprof.Handler("threadcreate"))))
}
