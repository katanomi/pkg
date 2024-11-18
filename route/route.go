/*
Copyright 2021 The AlaudaDevops Authors.

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
	"context"

	"github.com/emicklei/go-restful/v3"
)

// Route a service should implement register func to register go restful webservice
// Deprecated: replaced by ContextRoute
type Route interface {
	Register(ws *restful.WebService)
}

// ContextRoute is a service should implement register func to register go restful webservice
type ContextRoute interface {
	Register(ctx context.Context, ws *restful.WebService) error
}

// NewDefaultService default service included with metrics,pprof
func NewDefaultService(ctx context.Context) *restful.WebService {
	routes := []Route{
		NewSystem(ctx),
		NewHealthz(ctx),
	}

	ws := &restful.WebService{}
	for _, each := range routes {
		each.Register(ws)
	}

	return ws
}
