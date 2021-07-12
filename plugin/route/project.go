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

type projectList struct {
	impl client.ProjectLister
	tags []string
}

//NewProjectList create a list project route with plugin client
func NewProjectList(impl client.ProjectLister) Route {
	return &projectList{
		tags: []string{"projects"},
		impl: impl,
	}
}

func (p *projectList) Register(ws *restful.WebService) {
	ws.Route(ws.GET("/projects").To(p.ListProjects).
		// docs
		Doc("list projects").
		Metadata(restfulspec.KeyOpenAPITags, p.tags))
}

// ListProjects http handler for list project
func (p *projectList) ListProjects(request *restful.Request, response *restful.Response) {
	option := p.parseQuery(request)
	projects, err := p.impl.ListProjects(request.Request.Context(), option)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	response.WriteEntity(projects)
}

func (p *projectList) parseQuery(request *restful.Request) client.ListOption {
	option := client.ListOption{
		Keyword:      request.QueryParameter("keyword"),
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

type projectCreate struct {
	impl client.ProjectCreator
	tags []string
}

// NewProjectCreate create a create project route with plugin client
func NewProjectCreate(impl client.ProjectCreator) Route {
	return &projectCreate{
		tags: []string{"projects"},
		impl: impl,
	}
}

func (p *projectCreate) Register(ws *restful.WebService) {
	ws.Route(ws.GET("/projects").To(p.CreateProject).
		// docs
		Doc("list projects").
		Metadata(restfulspec.KeyOpenAPITags, p.tags))
}

// CreateProject http handler for create project
func (p *projectCreate) CreateProject(request *restful.Request, response *restful.Response) {
	project := &client.Project{}
	if err := request.ReadEntity(project); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	resp, err := p.impl.CreateProject(request.Request.Context(), project)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	response.WriteEntity(resp)
}
