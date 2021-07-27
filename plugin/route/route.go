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
	"fmt"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/katanomi/pkg/plugin/client"
	"github.com/katanomi/pkg/plugin/component/metrics"
	"github.com/katanomi/pkg/plugin/component/tracing"
)

var filters = []restful.FilterFunction{
	tracing.Filter,
	metrics.Filter,
	client.AuthFilter,
	client.MetaFilter,
}

// Route a service should implement register func to register go restful webservice
type Route interface {
	Register(ws *restful.WebService)
}

// match math route with plugin client
func match(c client.PluginClient) []Route {
	routes := make([]Route, 0)
	if v, ok := c.(client.ProjectLister); ok {
		routes = append(routes, NewProjectList(v))
	}

	if v, ok := c.(client.ProjectCreator); ok {
		routes = append(routes, NewProjectCreate(v))
	}

	if v, ok := c.(client.ResourceLister); ok {
		routes = append(routes, NewResourceList(v))
	}

	return routes
}

// NewService new service from plugin client
func NewService(c client.PluginClient) (*restful.WebService, error) {
	routes := match(c)
	if len(routes) == 0 {
		return nil, fmt.Errorf("no route for provider %s", c.Path())
	}

	group := &restful.WebService{}
	group.Path(c.Path()).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	for _, filter := range filters {
		group.Filter(filter)
	}

	for _, r := range routes {
		r.Register(group)
	}

	return group, nil
}

// NewDefaultService default service included with metrics,pprof
func NewDefaultService() *restful.WebService {
	routes := []Route{
		NewSystem(),
	}

	ws := &restful.WebService{}
	for _, each := range routes {
		each.Register(ws)
	}

	return ws
}

//NewDocService go restful api doc
func NewDocService() *restful.WebService {
	config := restfulspec.Config{
		WebServices: restful.RegisteredWebServices(),
		APIPath:     "/openapi.json",
	}
	return restfulspec.NewOpenAPIService(config)
}
