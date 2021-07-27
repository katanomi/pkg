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
	kerrors "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/plugin/client"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
)

type resourceList struct {
	impl client.ResourceLister
	tags []string
}

// NewResourceList create a list resource route with plugin client
func NewResourceList(impl client.ResourceLister) Route {
	return &resourceList{
		tags: []string{"resources"},
		impl: impl,
	}
}

func (r *resourceList) Register(ws *restful.WebService) {
	ws.Route(
		ListOptionsDocs(ws.GET("/resources").To(r.ResourceList).
			// docs
			// keep the name the same as the method to store in IntegrationClass
			Doc("ListResources").
			Metadata(restfulspec.KeyOpenAPITags, r.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.ResourceList{})))
}

// ResourceList http handler for list resource
func (r *resourceList) ResourceList(request *restful.Request, response *restful.Response) {
	option := GetListOptionsFromRequest(request)
	resources, err := r.impl.ListResources(request.Request.Context(), option)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, resources)
}
