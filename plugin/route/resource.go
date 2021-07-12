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
	"strconv"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/katanomi/pkg/plugin/client"
)

type resourceList struct {
	impl client.ResourceLister
	tags []string
}

// NewResourceList create a list resource route with plugin client
func NewResourceList(impl client.ResourceLister) Route {
	return &resourceList{
		tags: []string{"projects"},
		impl: impl,
	}
}

func (r *resourceList) Register(ws *restful.WebService) {
	ws.Route(ws.GET("/resources").To(r.ResourceList).
		// docs
		Doc("list projects").
		Metadata(restfulspec.KeyOpenAPITags, r.tags))
}

// ResourceList http handler for list resource
func (r *resourceList) ResourceList(request *restful.Request, response *restful.Response) {
	option := r.parseQuery(request)
	resources, err := r.impl.ListResources(request.Request.Context(), option)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	response.WriteEntity(resources)
}

func (r *resourceList) parseQuery(request *restful.Request) client.ListOption {
	option := client.ListOption{
		ItemsPerPage: 10,
		Page:         1,
	}

	itemsPerPage := request.QueryParameter("itemsPerPage")
	if v, err := strconv.Atoi(itemsPerPage); err == nil {
		option.ItemsPerPage = v
	}
	page := request.QueryParameter("page")
	if v, err := strconv.Atoi(page); err == nil {
		option.Page = v
	}

	return option
}
