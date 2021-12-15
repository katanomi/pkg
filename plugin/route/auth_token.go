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
	// "context"
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	kerrors "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/plugin/client"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	// corev1 "k8s.io/api/core/v1"
)

type authToken struct {
	impl client.AuthTokenGenerator
	tags []string
}

// NewAuthToken new route for auth checking
func NewAuthToken(impl client.AuthTokenGenerator) Route {
	return &authToken{
		impl: impl,
		tags: []string{"auth", "token"},
	}
}

func (a *authToken) Register(ws *restful.WebService) {
	ws.Route(
		ws.POST("/auth/token").To(a.AuthToken).
			Doc("AuthToken").
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.AuthToken{}),
	)
}

func (a *authToken) AuthToken(req *restful.Request, resp *restful.Response) {
	result, err := a.impl.AuthToken(req.Request.Context())
	if err != nil {
		kerrors.HandleError(req, resp, err)
		return
	}
	resp.WriteHeaderAndEntity(http.StatusOK, result)
}
