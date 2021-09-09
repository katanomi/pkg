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

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	kerrors "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/plugin/client"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
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
	ws.Route(
		ListOptionsDocs(
			ws.GET("/projects").To(p.ListProjects).
				// docs
				Doc("ListProjects").
				Metadata(restfulspec.KeyOpenAPITags, p.tags).
				Returns(http.StatusOK, "OK", metav1alpha1.ProjectList{}),
		),
	)
}

// ListProjects http handler for list project
func (p *projectList) ListProjects(request *restful.Request, response *restful.Response) {
	option := GetListOptionsFromRequest(request)
	projects, err := p.impl.ListProjects(request.Request.Context(), option)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, projects)
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
	ws.Route(ws.POST("/projects").To(p.CreateProject).
		// docs
		Doc("CreateProject").
		Metadata(restfulspec.KeyOpenAPITags, p.tags).
		Reads(metav1alpha1.Project{}, "Project").
		Returns(http.StatusCreated, "Project Created", metav1alpha1.Project{}))
}

// CreateProject http handler for create project
func (p *projectCreate) CreateProject(request *restful.Request, response *restful.Response) {
	project := &metav1alpha1.Project{}
	if err := request.ReadEntity(project); err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	resp, err := p.impl.CreateProject(request.Request.Context(), project)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusCreated, resp)
}

type projectGet struct {
	impl client.ProjectGetter
	tags []string
}

// NewProjectGet create a get project route with plugin client
func NewProjectGet(impl client.ProjectGetter) Route {
	return &projectGet{
		tags: []string{"projects"},
		impl: impl,
	}
}

func (p *projectGet) Register(ws *restful.WebService) {
	ws.Route(ws.GET("/projects/{project-id}").To(p.GetProject).
		Param(ws.PathParameter("project-id", "identifier of the project").DataType("string")).
		// docs
		Doc("GetProject").
		Metadata(restfulspec.KeyOpenAPITags, p.tags).
		Reads(metav1alpha1.Project{}, "Project").
		Returns(http.StatusOK, "Get Project Succeeded", metav1alpha1.Project{}))
}

// GetProject http handler for get project
func (p *projectGet) GetProject(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("project-id")

	resp, err := p.impl.GetProject(request.Request.Context(), id)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, resp)
}
