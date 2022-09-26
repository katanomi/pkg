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
	"github.com/katanomi/pkg/plugin/path"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
)

type gitRepoFileGetter struct {
	impl client.GitRepoFileGetter
	tags []string
}

// NewGitRepoFileGetter create a git GitRepoFile route with plugin client
func NewGitRepoFileGetter(impl client.GitRepoFileGetter) Route {
	return &gitRepoFileGetter{
		tags: []string{"git", "repositories", "file"},
		impl: impl,
	}
}

// Register route
func (a *gitRepoFileGetter) Register(ws *restful.WebService) {
	repositoryParam := ws.PathParameter("repository", "file belong to repository")
	projectParam := ws.PathParameter("project", "repository belong to project")
	pathParam := ws.PathParameter("path", "file path")
	refParam := ws.QueryParameter("ref", "file belong to commit/branch/tag name")
	ws.Route(
		ws.GET("/projects/{project:*}/coderepositories/{repository}/content/{path}").To(a.GetGitRepoFile).
			Doc("GetGitRepoFile").Param(projectParam).Param(repositoryParam).Param(pathParam).Param(refParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.GitRepoFile{}),
	)
}

// GetGitRepoFile get repo file
func (a *gitRepoFileGetter) GetGitRepoFile(request *restful.Request, response *restful.Response) {
	repo := path.Parameter(request, "repository")
	project := path.Parameter(request, "project")
	filePath := path.Parameter(request, "path")
	gitRepoFileParams := metav1alpha1.GitRepoFileOption{
		GitRepo: metav1alpha1.GitRepo{Repository: repo, Project: project},
		Ref:     request.QueryParameter("ref"),
		Path:    filePath,
	}
	fileInfo, err := a.impl.GetGitRepoFile(request.Request.Context(), gitRepoFileParams)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, fileInfo)
}

type gitRepoFileCreator struct {
	impl client.GitRepoFileCreator
	tags []string
}

// NewGitRepoFileCreator create a git GitRepoFile route with plugin client
func NewGitRepoFileCreator(impl client.GitRepoFileCreator) Route {
	return &gitRepoFileCreator{
		tags: []string{"git", "repositories", "file"},
		impl: impl,
	}
}

// Register route
func (a *gitRepoFileCreator) Register(ws *restful.WebService) {
	repositoryParam := ws.PathParameter("repository", "file belong to repository")
	projectParam := ws.PathParameter("project", "repository belong to project")
	pathParam := ws.PathParameter("filepath", "file path")
	ws.Route(
		ws.POST("/projects/{project:*}/coderepositories/{repository}/content/{filepath}").To(a.CreateGitRepoFile).
			Doc("CreateBranch").Param(projectParam).Param(repositoryParam).Param(pathParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.GitCommit{}),
	)
}

// CreateGitRepoFile create file in repo's one branch
func (a *gitRepoFileCreator) CreateGitRepoFile(request *restful.Request, response *restful.Response) {
	repo := path.Parameter(request, "repository")
	project := path.Parameter(request, "project")
	filePath := path.Parameter(request, "filepath")
	var params metav1alpha1.CreateRepoFileParams
	if err := request.ReadEntity(&params); err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	payload := metav1alpha1.CreateRepoFilePayload{
		GitRepo:              metav1alpha1.GitRepo{Repository: repo, Project: project},
		CreateRepoFileParams: params,
		FilePath:             filePath,
	}
	commitObject, err := a.impl.CreateGitRepoFile(request.Request.Context(), payload)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, commitObject)
}
