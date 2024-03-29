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

	"github.com/katanomi/pkg/plugin/path"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	kerrors "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/plugin/client"
)

type branchList struct {
	impl client.IssueBranchLister
	tags []string
}

// NewIssueBranchList create a list issue branch route with plugin client
func NewIssueBranchList(impl client.IssueBranchLister) Route {
	return &branchList{
		tags: []string{"projects", "issues", "branches"},
		impl: impl,
	}
}

func (b *branchList) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "issue belong to integrate project")
	issueParam := ws.PathParameter("issue", "issue id")
	ws.Route(
		ListOptionsDocs(
			ws.GET("/projects/{project:*}/issues/{issue:*}/branches").To(b.ListIssueBranches).
				Doc("ListIssurBranches").Param(projectParam).Param(issueParam).
				Metadata(restfulspec.KeyOpenAPITags, b.tags).
				Returns(http.StatusOK, "OK", metav1alpha1.BranchList{}),
		),
	)
}

func (b *branchList) ListIssueBranches(request *restful.Request, response *restful.Response) {
	option := GetListOptionsFromRequest(request)
	pathParams := metav1alpha1.IssueOptions{
		Identity: path.Parameter(request, "project"),
		IssueId:  path.Parameter(request, "issue"),
	}
	branches, err := b.impl.ListIssueBranches(request.Request.Context(), pathParams, option)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, branches)
}

type branchCreator struct {
	impl client.IssueBranchCreator
	tags []string
}

// NewIssueBranchCreate create a create issue relate branch route with plugin client
func NewIssueBranchCreate(impl client.IssueBranchCreator) Route {
	return &branchCreator{
		tags: []string{"projects", "issues", "branches"},
		impl: impl,
	}
}

func (b *branchCreator) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "issue belong to integrate project")
	issueParam := ws.PathParameter("issue", "issue id")
	ws.Route(
		ws.POST("/projects/{project:*}/issues/{issue:*}/branches").To(b.CreateIssueBranch).
			Doc("CreateIssueBranch").Param(projectParam).Param(issueParam).
			Metadata(restfulspec.KeyOpenAPITags, b.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.Branch{}),
	)
}

func (b *branchCreator) CreateIssueBranch(request *restful.Request, response *restful.Response) {
	pathParams := metav1alpha1.IssueOptions{
		Identity: path.Parameter(request, "project"),
		IssueId:  path.Parameter(request, "issue"),
	}

	payload := &metav1alpha1.Branch{}
	if err := request.ReadEntity(payload); err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	branch, err := b.impl.CreateIssueBranch(request.Request.Context(), pathParams, *payload)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, branch)
}

type branchDeleter struct {
	impl client.IssueBranchDeleter
	tags []string
}

// NewIssueBranchDelete create a create issue relate branch route with plugin client
func NewIssueBranchDelete(impl client.IssueBranchDeleter) Route {
	return &branchDeleter{
		tags: []string{"projects", "issues", "branches"},
		impl: impl,
	}
}

func (b *branchDeleter) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "issue belong to integrate project")
	issueParam := ws.PathParameter("issue", "issue id")
	ws.Route(
		ws.DELETE("/projects/{project:*}/issues/{issue:*}/branches/{branch}").To(b.DeleteIssueBranch).
			Doc("DeleteIssueBranch").Param(projectParam).Param(issueParam).
			Metadata(restfulspec.KeyOpenAPITags, b.tags).
			Returns(http.StatusOK, "OK", nil),
	)
}

func (b *branchDeleter) DeleteIssueBranch(request *restful.Request, response *restful.Response) {
	option := GetListOptionsFromRequest(request)
	pathParams := metav1alpha1.IssueOptions{
		Identity: path.Parameter(request, "project"),
		IssueId:  path.Parameter(request, "issue"),
		Branch:   path.Parameter(request, "branch"),
	}

	if err := b.impl.DeleteIssueBranch(request.Request.Context(), pathParams, option); err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, "OK")
}
