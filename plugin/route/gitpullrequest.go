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
	"net/http"
	"strconv"

	"github.com/katanomi/pkg/plugin/path"

	kerrors "github.com/katanomi/pkg/errors"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
)

type gitPullRequestHandler struct {
	impl client.GitPullRequestHandler
	tags []string
}

// NewGitPullRequestLister get a git pr route with plugin client
func NewGitPullRequestLister(impl client.GitPullRequestHandler) Route {
	return &gitPullRequestHandler{
		tags: []string{"git", "repositories", "pull request"},
		impl: impl,
	}
}

// Register route
func (a *gitPullRequestHandler) Register(ws *restful.WebService) {
	repositoryParam := ws.PathParameter("repository", "pulls belong to repository")
	projectParam := ws.PathParameter("project", "repository belong to project")
	indexParam := ws.PathParameter("index", "pr index")
	ws.Route(
		ws.GET("/projects/{project:*}/coderepositories/{repository}/pulls").To(a.ListGitPullRequest).
			Doc("GetGitPullRequest").Param(projectParam).Param(repositoryParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.GitPullRequestList{}),
	)
	ws.Route(
		ws.POST("/projects/{project:*}/coderepositories/{repository}/pulls").To(a.CreateGitPullRequest).
			Doc("GetGitPullRequest").Param(projectParam).Param(repositoryParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.GitPullRequest{}),
	)
	ws.Route(
		ws.GET("/projects/{project:*}/coderepositories/{repository}/pulls/{index}").To(a.GetGitPullRequest).
			Doc("GetGitPullRequest").Param(projectParam).Param(repositoryParam).Param(indexParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.GitPullRequest{}),
	)
}

// ListGitPullRequest list pr info
func (a *gitPullRequestHandler) ListGitPullRequest(request *restful.Request, response *restful.Response) {
	repo := path.Parameter(request, "repository")
	project := path.Parameter(request, "project")
	_state := request.QueryParameter("state")
	state := metav1alpha1.String2PullRequestState(_state)

	option := metav1alpha1.GitPullRequestListOption{
		GitRepo: metav1alpha1.GitRepo{Repository: repo, Project: project},
		State:   state,
	}
	listOption := GetListOptionsFromRequest(request)
	prList, err := a.impl.ListGitPullRequest(request.Request.Context(), option, listOption)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, prList)
}

// GetGitPullRequest get pr info
func (a *gitPullRequestHandler) GetGitPullRequest(request *restful.Request, response *restful.Response) {
	indexStr := path.Parameter(request, "index")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	repo := path.Parameter(request, "repository")
	project := path.Parameter(request, "project")
	option := metav1alpha1.GitPullRequestOption{
		GitRepo: metav1alpha1.GitRepo{Repository: repo, Project: project},
		Index:   index,
	}
	prInfo, err := a.impl.GetGitPullRequest(request.Request.Context(), option)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, prInfo)
}

// CreateGitPullRequest create a pr
func (a *gitPullRequestHandler) CreateGitPullRequest(request *restful.Request, response *restful.Response) {
	repo := path.Parameter(request, "repository")
	project := path.Parameter(request, "project")
	var params metav1alpha1.CreatePullRequestPayload
	if err := request.ReadEntity(&params); err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	params.Source.Repository = fmt.Sprintf("%s/%s", project, repo)
	prObject, err := a.impl.CreatePullRequest(request.Request.Context(), params)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, prObject)
}

type gitPullRequestNoteCreator struct {
	impl client.GitPullRequestCommentCreator
	tags []string
}

// NewGitPullRequestNoteCreator create a git pr note route with plugin client
func NewGitPullRequestNoteCreator(impl client.GitPullRequestCommentCreator) Route {
	return &gitPullRequestNoteCreator{
		tags: []string{"git", "repositories", "pull request", "note"},
		impl: impl,
	}
}

// Register route
func (a *gitPullRequestNoteCreator) Register(ws *restful.WebService) {
	repositoryParam := ws.PathParameter("repository", "pulls belong to repository")
	projectParam := ws.PathParameter("project", "repository belong to project")
	indexParam := ws.PathParameter("index", "note belong to index")
	ws.Route(
		ws.POST("/projects/{project:*}/coderepositories/{repository}/pulls/{index}/note").To(a.CreateGitPullRequestNote).
			Doc("CreateGitPullRequestNote").Param(projectParam).Param(repositoryParam).Param(indexParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.GitPullRequestNote{}),
	)
}

// CreateGitPullRequestNote create pr note
func (a *gitPullRequestNoteCreator) CreateGitPullRequestNote(request *restful.Request, response *restful.Response) {
	repo := path.Parameter(request, "repository")
	project := path.Parameter(request, "project")
	indexStr := path.Parameter(request, "index")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	var params metav1alpha1.CreatePullRequestCommentParam
	if err = request.ReadEntity(&params); err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	note, err := a.impl.CreatePullRequestComment(request.Request.Context(), metav1alpha1.CreatePullRequestCommentPayload{
		GitRepo:                       metav1alpha1.GitRepo{Repository: repo, Project: project},
		Index:                         index,
		CreatePullRequestCommentParam: params,
	})
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, note)
}

type gitPullRequestNoteUpdater struct {
	impl client.GitPullRequestCommentUpdater
	tags []string
}

// NewGitPullRequestNoteUpdater updates a git pr note route with plugin client
func NewGitPullRequestNoteUpdater(impl client.GitPullRequestCommentUpdater) Route {
	return &gitPullRequestNoteUpdater{
		tags: []string{"git", "repositories", "pull request", "note"},
		impl: impl,
	}
}

// Register route
func (a *gitPullRequestNoteUpdater) Register(ws *restful.WebService) {
	repositoryParam := ws.PathParameter("repository", "pulls belong to repository")
	projectParam := ws.PathParameter("project", "repository belong to project")
	indexParam := ws.PathParameter("index", "note belong to index")
	noteIDParam := ws.PathParameter("noteid", "PR comment/note id")
	ws.Route(
		ws.PUT("/projects/{project:*}/coderepositories/{repository}/pulls/{index}/note/{noteid}").To(a.UpdateGitPullRequestNote).
			Doc("UpdateGitPullRequestNote").Param(projectParam).Param(repositoryParam).Param(indexParam).Param(noteIDParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.GitPullRequestNote{}),
	)
}

// UpdateGitPullRequestNote updates pr note
func (a *gitPullRequestNoteUpdater) UpdateGitPullRequestNote(request *restful.Request, response *restful.Response) {
	repo := path.Parameter(request, "repository")
	project := path.Parameter(request, "project")
	indexStr := path.Parameter(request, "index")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	noteIdStr := path.Parameter(request, "noteid")
	noteId, err := strconv.Atoi(noteIdStr)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	var params metav1alpha1.CreatePullRequestCommentParam
	if err = request.ReadEntity(&params); err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	note, err := a.impl.UpdatePullRequestComment(request.Request.Context(), metav1alpha1.UpdatePullRequestCommentPayload{
		GitRepo:                       metav1alpha1.GitRepo{Repository: repo, Project: project},
		Index:                         index,
		CreatePullRequestCommentParam: params,
		CommentID:                     noteId,
	})
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, note)
}

type gitPullRequestCommentLister struct {
	impl client.GitPullRequestCommentLister
	tags []string
}

// NewGitPullRequestCommentLister list a git pr note route with plugin client
func NewGitPullRequestCommentLister(impl client.GitPullRequestCommentLister) Route {
	return &gitPullRequestCommentLister{
		tags: []string{"git", "repositories", "pull request", "note"},
		impl: impl,
	}
}

// Register route
func (a *gitPullRequestCommentLister) Register(ws *restful.WebService) {
	repositoryParam := ws.PathParameter("repository", "pulls belong to repository")
	projectParam := ws.PathParameter("project", "repository belong to project")
	indexParam := ws.PathParameter("index", "note belong to index")
	ws.Route(
		ws.GET("/projects/{project:*}/coderepositories/{repository}/pulls/{index}/note").To(a.ListPullRequestComment).
			Doc("ListPullRequestComment").Param(projectParam).Param(repositoryParam).Param(indexParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.GitPullRequestNoteList{}),
	)
}

// ListGitPullRequestNote List pr note
func (a *gitPullRequestCommentLister) ListPullRequestComment(request *restful.Request, response *restful.Response) {
	repo := path.Parameter(request, "repository")
	project := path.Parameter(request, "project")
	indexStr := path.Parameter(request, "index")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	option := GetListOptionsFromRequest(request)
	noteList, err := a.impl.ListPullRequestComment(request.Request.Context(), metav1alpha1.GitPullRequestOption{
		GitRepo: metav1alpha1.GitRepo{Repository: repo, Project: project},
		Index:   index,
	}, option)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, noteList)
}
