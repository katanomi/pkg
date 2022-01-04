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

	kerrors "github.com/katanomi/pkg/errors"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
)

type gitCommitStatusLister struct {
	impl client.GitCommitStatusLister
	tags []string
}

// NewGitCommitStatusLister list git Commit status route with plugin client
func NewGitCommitStatusLister(impl client.GitCommitStatusLister) Route {
	return &gitCommitStatusLister{
		tags: []string{"git", "repositories", "commit status"},
		impl: impl,
	}
}

// Register route
func (a *gitCommitStatusLister) Register(ws *restful.WebService) {
	repositoryParam := ws.PathParameter("repository", "commit belong to repository")
	shaParam := ws.PathParameter("sha", "commit sha")
	projectParam := ws.PathParameter("project", "repository belong to project")
	ws.Route(
		ws.GET("/projects/{project:*}/coderepositories/{repository}/commit/{sha}/status").To(a.ListGitCommitStatus).
			Doc("ListGitCommitStatus").Param(projectParam).Param(repositoryParam).Param(shaParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.GitCommitStatusList{}),
	)
}

// ListGitCommitStatus List commit status
func (a *gitCommitStatusLister) ListGitCommitStatus(request *restful.Request, response *restful.Response) {
	option := GetListOptionsFromRequest(request)
	sha := request.PathParameter("sha")
	repo := handlePathParamHasSlash(request.PathParameter("repository"))
	project := request.PathParameter("project")
	commitOption := metav1alpha1.GitCommitOption{
		GitRepo:            metav1alpha1.GitRepo{Repository: repo, Project: project},
		GitCommitBasicInfo: metav1alpha1.GitCommitBasicInfo{SHA: &sha},
	}
	statusList, err := a.impl.ListGitCommitStatus(request.Request.Context(), commitOption, option)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, statusList)
}

type gitCommitStatusCreator struct {
	impl client.GitCommitStatusCreator
	tags []string
}

// NewGitCommitStatusCreator create git Commit status route with plugin client
func NewGitCommitStatusCreator(impl client.GitCommitStatusCreator) Route {
	return &gitCommitStatusCreator{
		tags: []string{"git", "repositories", "commit status"},
		impl: impl,
	}
}

// Register route
func (a *gitCommitStatusCreator) Register(ws *restful.WebService) {
	repositoryParam := ws.PathParameter("repository", "commit belong to repository")
	shaParam := ws.PathParameter("sha", "commit sha")
	projectParam := ws.PathParameter("project", "repository belong to project")
	ws.Route(
		ws.POST("/projects/{project:*}/coderepositories/{repository}/commit/{sha}/status").To(a.CreateGitCommitStatus).
			Doc("CreateGitCommitStatus").Param(projectParam).Param(repositoryParam).Param(shaParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.GitCommitStatus{}),
	)
}

// CreateGitCommitStatus create commit status
func (a *gitCommitStatusCreator) CreateGitCommitStatus(request *restful.Request, response *restful.Response) {
	sha := request.PathParameter("sha")
	repo := handlePathParamHasSlash(request.PathParameter("repository"))
	project := request.PathParameter("project")
	var params metav1alpha1.CreateCommitStatusParam
	if err := request.ReadEntity(&params); err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	commitComment, err := a.impl.CreateGitCommitStatus(request.Request.Context(), metav1alpha1.CreateCommitStatusPayload{
		GitRepo:                 metav1alpha1.GitRepo{Project: project, Repository: repo},
		GitCommitBasicInfo:      metav1alpha1.GitCommitBasicInfo{SHA: &sha},
		CreateCommitStatusParam: params,
	})
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, commitComment)
}
