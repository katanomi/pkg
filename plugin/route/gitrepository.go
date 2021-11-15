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
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"

	kerrors "github.com/katanomi/pkg/errors"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
)

func handlePathParamHasSlash(param string) (result string) {
	result = strings.Replace(param, "%2F", "/", -1)
	return result
}

type gitRepositoryLister struct {
	impl client.GitRepositoryLister
	tags []string
}

// NewGitRepositoryLister get a git repo route with plugin client
func NewGitRepositoryLister(impl client.GitRepositoryLister) Route {
	return &gitRepositoryLister{
		tags: []string{"git", "repositories"},
		impl: impl,
	}
}

// Register route
func (a *gitRepositoryLister) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "repository belong to project")
	keywordParam := ws.QueryParameter("keyword", "keyword for search repository")
	subtypeParam := ws.QueryParameter("subtype", "subtype for search repository")
	ws.Route(
		ws.GET("/projects/{project}/coderepositories").To(a.ListGitRepository).
			Doc("GetGitRepoList").Param(projectParam).Param(keywordParam).Param(subtypeParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.GitRepositoryList{}),
	)
}

// ListGitRepository list repo info
func (a *gitRepositoryLister) ListGitRepository(request *restful.Request, response *restful.Response) {
	project := request.PathParameter("project")
	keyword := request.QueryParameter("keyword")
	subtype := request.QueryParameter("subtype")
	kind := metav1alpha1.ProjectSubType(subtype)
	if kind != metav1alpha1.GitUserProjectSubType && kind != metav1alpha1.GitGroupProjectSubType && kind != "CodeRepository" {
		kerrors.HandleError(request, response, errors.NewBadRequest("no git source subtype"))
		return
	}
	listOption := GetListOptionsFromRequest(request)
	repoList, err := a.impl.ListGitRepository(request.Request.Context(), project, keyword, kind, listOption)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, repoList)
}

type gitRepositoryGetter struct {
	impl client.GitRepositoryGetter
	tags []string
}

// NewGitRepositoryGetter get a git repo route with plugin client
func NewGitRepositoryGetter(impl client.GitRepositoryGetter) Route {
	return &gitRepositoryGetter{
		tags: []string{"git", "repositories"},
		impl: impl,
	}
}

// Register route
func (a *gitRepositoryGetter) Register(ws *restful.WebService) {
	repositoryParam := ws.PathParameter("repository", "repository name")
	projectParam := ws.PathParameter("project", "repository belong to project")
	ws.Route(
		ws.GET("/projects/{project}/coderepositories/{repository}").To(a.GetGitRepository).
			Doc("GetGitRepo").Param(projectParam).Param(repositoryParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.GitRepository{}),
	)
}

// GetGitRepository get repo info
func (a *gitRepositoryGetter) GetGitRepository(request *restful.Request, response *restful.Response) {
	project := request.PathParameter("project")
	repo := handlePathParamHasSlash(request.PathParameter("repository"))
	repoInfo, err := a.impl.GetGitRepository(request.Request.Context(), metav1alpha1.GitRepo{Repository: repo, Project: project})
	if err != nil {
		if errors.IsNotFound(err) {
			response.WriteError(http.StatusNotFound, err)
			return
		}
		kerrors.HandleError(request, response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, repoInfo)
}
