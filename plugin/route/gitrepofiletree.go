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

	kerrors "github.com/katanomi/pkg/errors"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
)

// gitRepositoryFileTreeGetter is impl ClientGitRepositoryFileTree interface
type gitRepositoryFileTreeGetter struct {
	impl client.GitRepositoryFileTreeGetter
	tags []string
}

// NewGitRepositoryFileTreeGetter get a git repo route with plugin client
func NewGitRepositoryFileTreeGetter(impl client.GitRepositoryFileTreeGetter) Route {
	return &gitRepositoryFileTreeGetter{
		tags: []string{"git", "repositories", "tree"},
		impl: impl,
	}
}

// Register route
func (g *gitRepositoryFileTreeGetter) Register(ws *restful.WebService) {
	repositoryParam := ws.PathParameter("repository", "file belong to repository")
	projectParam := ws.PathParameter("project", "repository belong to project")
	pathParam := ws.QueryParameter("path", "file path")
	treeShaParam := ws.QueryParameter("tree_sha", "sha for file tree")
	recursive := ws.QueryParameter("recursive", "recursive switch")
	ws.Route(
		ws.GET("/projects/{project:*}/coderepositories/{repository}/tree").To(g.GetGitRepositoryFileTree).
			Doc("GetGitRepositoryFileTree").
			Param(projectParam).
			Param(repositoryParam).
			Param(pathParam).
			Param(treeShaParam).
			Param(recursive).
			Metadata(restfulspec.KeyOpenAPITags, g.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.GitRepositoryFileTree{}),
	)
	return
}

// GetGitRepositoryFileTree get repo file tree
func (g *gitRepositoryFileTreeGetter) GetGitRepositoryFileTree(request *restful.Request, response *restful.Response) {
	repo := path.Parameter(request, "repository")
	project := path.Parameter(request, "project")
	path := request.QueryParameter("path")
	recursive := request.QueryParameter("recursive")
	recursiveValue := recursive == "true"
	treeSha := request.QueryParameter("tree_sha")

	ctx := request.Request.Context()
	option := metav1alpha1.GitRepoFileTreeOption{
		GitRepo:   metav1alpha1.GitRepo{Repository: repo, Project: project},
		Path:      path,
		TreeSha:   treeSha,
		Recursive: recursiveValue,
	}
	listOption := metav1alpha1.ListOptions{}
	fileTree, err := g.impl.GetGitRepositoryFileTree(ctx, option, listOption)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, fileTree)
}
