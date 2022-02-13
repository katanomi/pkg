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
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	kerrors "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/plugin/client"
)

type issueList struct {
	impl client.IssueLister
	tags []string
}

// NewIssueList create a list issue route with plugin client
func NewIssueList(impl client.IssueLister) Route {
	return &issueList{
		tags: []string{"projects", "issues"},
		impl: impl,
	}
}

func (i *issueList) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "issue belong to integrate project")
	ws.Route(
		ListOptionsDocs(
			ws.GET("/project/{project:*}/issues").To(i.ListIssues).
				Doc("ListIssues").Param(projectParam).
				Metadata(restfulspec.KeyOpenAPITags, i.tags).
				Returns(http.StatusOK, "OK", metav1alpha1.IssueList{}),
		),
	)
}

func (i *issueList) ListIssues(request *restful.Request, response *restful.Response) {
	option := GetListOptionsFromRequest(request)
	projectId, err := strconv.ParseInt(request.PathParameter("project"), 10, 64)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	pathParams := metav1alpha1.IssueOptions{
		ProjectId: int(projectId),
	}
	issues, err := i.impl.ListIssues(request.Request.Context(), pathParams, option)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, issues)
}

type issueGetter struct {
	impl client.IssueGetter
	tags []string
}

// NewIssueGet create a get issue route with plugin client
func NewIssueGet(impl client.IssueGetter) Route {
	return &issueGetter{
		tags: []string{"projects", "issues"},
		impl: impl,
	}
}

func (i *issueGetter) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "issue belong to integrate project")
	issueParam := ws.PathParameter("issue", "issue id")
	ws.Route(
		ws.GET("/project/{project:*}/issues/{issue:*}").To(i.GetIssue).
			Doc("GetIssues").Param(projectParam).Param(issueParam).
			Metadata(restfulspec.KeyOpenAPITags, i.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.Issue{}),
	)
}

func (i *issueGetter) GetIssue(request *restful.Request, response *restful.Response) {
	projectId, err := strconv.ParseInt(request.PathParameter("project"), 10, 64)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	issueId, err := strconv.ParseInt(request.PathParameter("issue"), 10, 64)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	pathParams := metav1alpha1.IssueOptions{
		ProjectId: int(projectId),
		IssueId:   int(issueId),
	}
	issue, err := i.impl.GetIssue(request.Request.Context(), pathParams)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, issue)

}

type issueAttributeGetter struct {
	impl client.IssueAttributeGetter
	tags []string
}

// NewIssueAttributeGet create a get issue attribute route with plugin client
func NewIssueAttributeGet(impl client.IssueAttributeGetter) Route {
	return &issueAttributeGetter{
		tags: []string{"projects", "issues", "attributes"},
		impl: impl,
	}
}

func (i *issueAttributeGetter) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "integrate project")
	ws.Route(
		ws.GET("/projectmanagement/{project:*}/attributes").To(i.GetAttributes).
			Doc("GetAttributes").Param(projectParam).
			Metadata(restfulspec.KeyOpenAPITags, i.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.Attribute{}),
	)
}

func (i *issueAttributeGetter) GetAttributes(request *restful.Request, response *restful.Response) {
	projectId, err := strconv.ParseInt(request.PathParameter("project"), 10, 64)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	pathParams := metav1alpha1.IssueOptions{
		ProjectId: int(projectId),
	}
	attribute, err := i.impl.GetIssueAttribute(request.Request.Context(), pathParams)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, attribute)
}
