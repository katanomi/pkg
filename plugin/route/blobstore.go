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

type blobStoreLister struct {
	impl client.BlobStoreLister
	tags []string
}

//NewCodeQualityGetter create a get codeQuality route with plugin client
func NewBlobStoreLister(impl client.BlobStoreLister) Route {
	return &blobStoreLister{
		tags: []string{"blobStore"},
		impl: impl,
	}
}

func (c *blobStoreLister) Register(ws *restful.WebService) {
	ws.Route(
		ListOptionsDocs(
			ws.GET("/blobStores").To(c.ListBlobStores).
				// docs
				Doc("ListBlobStores").
				Metadata(restfulspec.KeyOpenAPITags, c.tags).
				Returns(http.StatusOK, "OK", metav1alpha1.BlobStoreList{}),
		),
	)
}

// ListBlobStores http handler for list blob stores
func (c *blobStoreLister) ListBlobStores(request *restful.Request, response *restful.Response) {
	option := GetListOptionsFromRequest(request)
	list, err := c.impl.ListBlobStores(request.Request.Context(), option)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, list)
}
