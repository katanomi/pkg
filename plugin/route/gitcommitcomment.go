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
	"github.com/katanomi/pkg/plugin/path"
	"net/http"

	kerrors "github.com/katanomi/pkg/errors"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
)

type gitCommitCommentLister struct {
	impl client.GitCommitCommentLister
	tags []string
}

// NewGitCommitCommentLister get git Commit comment route with plugin client
func NewGitCommitCommentLister(impl client.GitCommitCommentLister) Route {
	return &gitCommitCommentLister{
		tags: []string{"git", "repositories", "commit comment"},
		impl: impl,
	}
}

// Register route
func (a *gitCommitCommentLister) Register(ws *restful.WebService) {
	repositoryParam := ws.PathParameter("repository", "commit belong to repository")
	shaParam := ws.PathParameter("sha", "commit sha")
	projectParam := ws.PathParameter("project", "repository belong to project")
	ws.Route(
		ws.GET("/projects/{project:*}/coderepositories/{repository}/commit/{sha}/comments").To(a.ListGitCommitComment).
			Doc("ListGitCommitComment").Param(projectParam).Param(repositoryParam).Param(shaParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.GitCommitCommentList{}),
	)
}

// ListGitCommitComment List commit info
func (a *gitCommitCommentLister) ListGitCommitComment(request *restful.Request, response *restful.Response) {
	option := GetListOptionsFromRequest(request)
	sha := path.Parameter(request, "sha")
	repo := path.Parameter(request, "repository")
	project := path.Parameter(request, "project")
	commitOption := metav1alpha1.GitCommitOption{
		GitRepo:            metav1alpha1.GitRepo{Repository: repo, Project: project},
		GitCommitBasicInfo: metav1alpha1.GitCommitBasicInfo{SHA: &sha},
	}
	commitList, err := a.impl.ListGitCommitComment(request.Request.Context(), commitOption, option)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, commitList)
}

type gitCommitCommentCreator struct {
	impl client.GitCommitCommentCreator
	tags []string
}

// NewGitCommitCommentCreator create git Commit comment route with plugin client
func NewGitCommitCommentCreator(impl client.GitCommitCommentCreator) Route {
	return &gitCommitCommentCreator{
		tags: []string{"git", "repositories", "commit comment"},
		impl: impl,
	}
}

// Register route
func (a *gitCommitCommentCreator) Register(ws *restful.WebService) {
	repositoryParam := ws.PathParameter("repository", "commit belong to repository")
	shaParam := ws.PathParameter("sha", "commit sha")
	projectParam := ws.PathParameter("project", "repository belong to project")
	ws.Route(
		ws.POST("/projects/{project:*}/coderepositories/{repository}/commit/{sha}/comments").To(a.CreateGitCommitComment).
			Doc("CreateGitCommitComment").Param(projectParam).Param(repositoryParam).Param(shaParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.GitCommitComment{}),
	)
}

// CreateGitCommitComment create commit info
func (a *gitCommitCommentCreator) CreateGitCommitComment(request *restful.Request, response *restful.Response) {
	sha := path.Parameter(request, "sha")
	repo := path.Parameter(request, "repository")
	project := path.Parameter(request, "project")
	var params metav1alpha1.CreateCommitCommentParam
	if err := request.ReadEntity(&params); err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	commitComment, err := a.impl.CreateGitCommitComment(request.Request.Context(), metav1alpha1.CreateCommitCommentPayload{
		GitRepo:                  metav1alpha1.GitRepo{Project: project, Repository: repo},
		GitCommitBasicInfo:       metav1alpha1.GitCommitBasicInfo{SHA: &sha},
		CreateCommitCommentParam: params,
	})
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, commitComment)
}
