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
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/api/errors"

	kerrors "github.com/katanomi/pkg/errors"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
)

type gitCommitGetter struct {
	impl client.GitCommitGetter
	tags []string
}

// NewGitCommitGetter get a git Commit route with plugin client
func NewGitCommitGetter(impl client.GitCommitGetter) Route {
	return &gitCommitGetter{
		tags: []string{"git", "repositories", "commit"},
		impl: impl,
	}
}

// Register route
func (a *gitCommitGetter) Register(ws *restful.WebService) {
	repositoryParam := ws.PathParameter("repository", "commit belong to repository")
	shaParam := ws.PathParameter("sha", "commit sha")
	projectParam := ws.PathParameter("project", "repository belong to project").DataType("string")
	ws.Route(
		ws.GET("/projects/{project:*}/coderepositories/{repository}/commit/{sha}").To(a.GetCommit).
			Doc("GetGitRepoFile").Param(projectParam).Param(repositoryParam).Param(shaParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.GitCommit{}),
	)
}

// GitCommit get commit info
func (a *gitCommitGetter) GetCommit(request *restful.Request, response *restful.Response) {
	sha := request.PathParameter("sha")
	repo := handlePathParamHasSlash(request.PathParameter("repository"))
	project := request.PathParameter("project")
	commitOption := metav1alpha1.GitCommitOption{
		GitRepo:            metav1alpha1.GitRepo{Repository: repo, Project: project},
		GitCommitBasicInfo: metav1alpha1.GitCommitBasicInfo{SHA: &sha},
	}
	commitObject, err := a.impl.GetGitCommit(request.Request.Context(), commitOption)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, commitObject)
}

// HandleTimeParamInQuery Processing parameters related to time filtering in query
func HandleTimeParamInQuery(param string) (res *v1.Time, err error) {
	var timeObj time.Time
	if param == "" {
		return
	}
	timeObj, err = time.Parse(time.RFC3339, param)
	if err != nil {
		return
	}
	res = &v1.Time{Time: timeObj}
	return
}

type gitCommitLister struct {
	impl client.GitCommitLister
	tags []string
}

// NewGitCommitLister get list git Commit route with plugin client
func NewGitCommitLister(impl client.GitCommitLister) Route {
	return &gitCommitLister{
		tags: []string{"git", "repositories", "commit"},
		impl: impl,
	}
}

// Register route
func (a *gitCommitLister) Register(ws *restful.WebService) {
	repositoryParam := ws.PathParameter("repository", "commit belong to repository")
	shaParam := ws.PathParameter("sha", "commit sha")
	projectParam := ws.PathParameter("project", "repository belong to project").DataType("string")
	ws.Route(
		ws.GET("/projects/{project:*}/coderepositories/{repository}/commits").To(a.ListCommit).
			Doc("ListCodeRepositoryCommit").Param(projectParam).Param(repositoryParam).Param(shaParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.GitCommit{}),
	)
}

// ListCommit get commit info
func (a *gitCommitLister) ListCommit(request *restful.Request, response *restful.Response) {
	repo := handlePathParamHasSlash(request.PathParameter("repository"))
	project := request.PathParameter("project")
	ref := request.QueryParameter("ref")
	option := GetListOptionsFromRequest(request)
	commitOption := metav1alpha1.GitCommitListOption{
		GitRepo: metav1alpha1.GitRepo{Repository: repo, Project: project},
		Ref:     ref,
	}
	var err error
	commitOption.Since, err = HandleTimeParamInQuery(request.QueryParameter(SinceQueryKey))
	if err != nil {
		kerrors.HandleError(request, response, errors.NewBadRequest(err.Error()))
		return
	}
	commitOption.Until, err = HandleTimeParamInQuery(request.QueryParameter(UntilQueryKey))
	if err != nil {
		kerrors.HandleError(request, response, errors.NewBadRequest(err.Error()))
		return
	}
	commitObject, err := a.impl.ListGitCommit(request.Request.Context(), commitOption, option)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, commitObject)
}
