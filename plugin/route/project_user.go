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

	"github.com/katanomi/pkg/plugin/path"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	kerrors "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/plugin/client"
)

type projectUserList struct {
	impl client.ProjectUserLister
	tags []string
}

// NewProjectUserList create a list user route with plugin client
func NewProjectUserList(impl client.ProjectUserLister) Route {
	return &projectUserList{
		tags: []string{"projects", "users"},
		impl: impl,
	}
}

func (u *projectUserList) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "user belong to integrate project")
	ws.Route(
		ListOptionsDocs(
			ws.GET("/projects/{project:*}/users").To(u.ListProjectUsers).
				Doc("ListProjectUsers").Param(projectParam).
				Metadata(restfulspec.KeyOpenAPITags, u.tags).
				Returns(http.StatusOK, "OK", metav1alpha1.UserList{}),
		),
	)
}

func (u *projectUserList) ListProjectUsers(request *restful.Request, response *restful.Response) {
	option := GetListOptionsFromRequest(request)
	pathParams := metav1alpha1.UserOptions{
		Project: path.Parameter(request, "project"),
	}
	users, err := u.impl.ListProjectUsers(request.Request.Context(), pathParams, option)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, users)
}
