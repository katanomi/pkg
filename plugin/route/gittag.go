/*
Copyright 2022 The Katanomi Authors.

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
	"k8s.io/apimachinery/pkg/api/errors"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	kerrors "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/plugin/client"
)

type gitTagLister struct {
	impl client.GitTagLister
	tags []string
}

// NewGitTagLister get a git tag route with plugin client
func NewGitTagLister(impl client.GitTagLister) Route {
	return &gitTagLister{
		tags: []string{"git", "repositories", "tag"},
		impl: impl,
	}
}

// Register route
func (a *gitTagLister) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "repository belong to project").DataType("string")
	repositoryParam := ws.PathParameter("repository", "tag belong to repository")
	ws.Route(
		ws.GET("/projects/{project:*}/coderepositories/{repository}/tags").To(a.ListGitTag).
			Doc("GetGitRepoList").Param(projectParam).Param(repositoryParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.GitTagList{}),
	)
}

// ListGitTag list repo tags
func (a *gitTagLister) ListGitTag(request *restful.Request, response *restful.Response) {
	repository := handlePathParamHasSlash(request.PathParameter("repository"))
	project := request.PathParameter("project")
	listOption := GetListOptionsFromRequest(request)
	tagList, err := a.impl.ListGitTag(request.Request.Context(),
		metav1alpha1.GitTagListOption{
			GitRepo: metav1alpha1.GitRepo{Repository: repository, Project: project},
		},
		listOption)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, tagList)
}

type gitTagGetter struct {
	impl client.GitTagGetter
	tags []string
}

// NewGitTagGetter get a git repo route with plugin client
func NewGitTagGetter(impl client.GitTagGetter) Route {
	return &gitTagGetter{
		tags: []string{"git", "repositories", "tag"},
		impl: impl,
	}
}

// Register route
func (a *gitTagGetter) Register(ws *restful.WebService) {
	repositoryParam := ws.PathParameter("repository", "repository name")
	projectParam := ws.PathParameter("project", "repository belong to project").DataType("string")
	tagParam := ws.PathParameter("tag", "the name of the tag")
	ws.Route(
		ws.GET("/projects/{project:*}/coderepositories/{repository}/tags/{tag}").To(a.GetGitTag).
			Doc("GetGitRepo").Param(projectParam).Param(repositoryParam).Param(tagParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.GitTag{}),
	)
}

// GetGitTag get repo tag
func (a *gitTagGetter) GetGitTag(request *restful.Request, response *restful.Response) {
	project := request.PathParameter("project")
	repo := handlePathParamHasSlash(request.PathParameter("repository"))
	tag := request.PathParameter("tag")
	repoInfo, err := a.impl.GetGitTag(request.Request.Context(),
		metav1alpha1.GitTagOption{
			GitRepo: metav1alpha1.GitRepo{Repository: repo, Project: project},
			Tag:     tag,
		},
	)
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
