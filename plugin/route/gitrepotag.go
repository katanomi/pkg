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

	"github.com/katanomi/pkg/plugin/path"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	kerrors "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/plugin/client"
)

type gitRepositoryTagLister struct {
	impl client.GitRepositoryTagLister
	tags []string
}

// NewGitRepositoryTagLister get a git repository tag route with plugin client
func NewGitRepositoryTagLister(impl client.GitRepositoryTagLister) Route {
	return &gitRepositoryTagLister{
		tags: []string{"git", "repositories", "tag"},
		impl: impl,
	}
}

// Register route
func (a *gitRepositoryTagLister) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "repository belong to project").DataType("string")
	repositoryParam := ws.PathParameter("repository", "tag belong to repository")
	ws.Route(
		ws.GET("/projects/{project:*}/coderepositories/{repository}/tags").To(a.ListGitRepositoryTag).
			Doc("ListGitRepositoryTag").Param(projectParam).Param(repositoryParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.GitRepositoryTagList{}),
	)
}

// ListGitRepositoryTag list repo tags
func (a *gitRepositoryTagLister) ListGitRepositoryTag(request *restful.Request, response *restful.Response) {
	repository := path.Parameter(request, "repository")
	project := path.Parameter(request, "project")
	listOption := GetListOptionsFromRequest(request)
	tagList, err := a.impl.ListGitRepositoryTag(request.Request.Context(),
		metav1alpha1.GitRepositoryTagListOption{
			GitRepo: metav1alpha1.GitRepo{Repository: repository, Project: project},
		},
		listOption)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, tagList)
}

type gitRepositoryTagGetter struct {
	impl client.GitRepositoryTagGetter
	tags []string
}

// NewGitRepositoryTagGetter get a git repo route with plugin client
func NewGitRepositoryTagGetter(impl client.GitRepositoryTagGetter) Route {
	return &gitRepositoryTagGetter{
		tags: []string{"git", "repositories", "tag"},
		impl: impl,
	}
}

// Register route
func (a *gitRepositoryTagGetter) Register(ws *restful.WebService) {
	repositoryParam := ws.PathParameter("repository", "repository name")
	projectParam := ws.PathParameter("project", "repository belong to project").DataType("string")
	tagParam := ws.PathParameter("tag", "the name of the tag")
	ws.Route(
		ws.GET("/projects/{project:*}/coderepositories/{repository}/tags/{tag}").To(a.GetGitRepositoryTag).
			Doc("GetGitRepositoryTag").Param(projectParam).Param(repositoryParam).Param(tagParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.GitRepositoryTag{}),
	)
}

// GetGitRepositoryTag get repo tag
func (a *gitRepositoryTagGetter) GetGitRepositoryTag(request *restful.Request, response *restful.Response) {
	project := path.Parameter(request, "project")
	repo := path.Parameter(request, "repository")
	tag := path.Parameter(request, "tag")
	repoInfo, err := a.impl.GetGitRepositoryTag(request.Request.Context(),
		metav1alpha1.GitRepositoryTagOption{
			GitRepo: metav1alpha1.GitRepo{Repository: repo, Project: project},
			Tag:     tag,
		},
	)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, repoInfo)
}
