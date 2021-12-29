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

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	kerrors "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/plugin/client"
)

type repositoryList struct {
	impl client.RepositoryLister
	tags []string
}

//NewRepositoryList create a list repository route with plugin client
func NewRepositoryList(impl client.RepositoryLister) Route {
	return &repositoryList{
		tags: []string{"projects", "repositories"},
		impl: impl,
	}
}

func (r *repositoryList) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "repository belong to integraion")
	ws.Route(
		ListOptionsDocs(
			ws.GET("/projects/{project:*}/repositories").To(r.ListRepositories).
				// docs
				Doc("ListRepositories").Param(projectParam).
				Metadata(restfulspec.KeyOpenAPITags, r.tags).
				Returns(http.StatusOK, "OK", metav1alpha1.RepositoryList{}),
		),
	)
}

// ListRepositories http handler for list repository
func (r *repositoryList) ListRepositories(request *restful.Request, response *restful.Response) {
	option := GetListOptionsFromRequest(request)

	subType := request.QueryParameter("subType")

	pathParams := metav1alpha1.RepositoryOptions{
		Project: request.PathParameter("project"),
		SubType: metav1alpha1.ProjectSubType(subType),
	}
	repositories, err := r.impl.ListRepositories(request.Request.Context(), pathParams, option)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, repositories)
}
