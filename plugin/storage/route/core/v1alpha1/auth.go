/*
Copyright 2023 The Katanomi Authors.

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

package v1alpha1

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/katanomi/pkg/apis/storage/v1alpha1"
	kerrors "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/plugin/storage"
	corev1alpha1 "github.com/katanomi/pkg/plugin/storage/core/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type authCheck struct {
	impl corev1alpha1.AuthChecker
	tags []string
}

func (a *authCheck) GroupVersion() schema.GroupVersion {
	return corev1alpha1.CoreV1alpha1GV
}

// NewAuthCheck new route for auth checking
func NewAuthCheck(impl corev1alpha1.AuthChecker) storage.VersionedRouter {
	return &authCheck{
		impl: impl,
		tags: []string{"auth"},
	}
}

func (a *authCheck) Register(ws *restful.WebService) {
	ws.Route(
		ws.POST("/auth/check").To(a.AuthCheck).
			Doc("Storage plugin auth check").
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Reads(v1alpha1.StorageAuthCheckRequest{}, "request storage plugin for auth check").
			Returns(http.StatusOK, "OK", v1alpha1.StorageAuthCheck{}),
	)
}

// AuthCheck is handler of auth check route
func (a *authCheck) AuthCheck(req *restful.Request, resp *restful.Response) {
	authReq := &v1alpha1.StorageAuthCheckRequest{}
	err := req.ReadEntity(authReq)
	if err != nil {
		kerrors.HandleError(req, resp, err)
		return
	}

	authCheck, err := a.impl.CheckAuth(req.Request.Context(), authReq.Params)

	if err != nil {
		kerrors.HandleError(req, resp, err)
		return
	}

	resp.WriteHeaderAndEntity(http.StatusOK, authCheck)
}
