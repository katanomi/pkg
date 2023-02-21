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
	"github.com/katanomi/pkg/plugin/path"
	"github.com/katanomi/pkg/plugin/storage"
	filestorev1alpha1 "github.com/katanomi/pkg/plugin/storage/capabilities/filestore/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type fileMeta struct {
	impl filestorev1alpha1.FileMetaInterface
	tags []string
}

func (a *fileMeta) GroupVersion() schema.GroupVersion {
	return filestorev1alpha1.FileStoreV1alpha1GV
}

// NewFileMeta new route for auth checking
func NewFileMeta(impl filestorev1alpha1.FileMetaInterface) storage.VersionedRouter {
	return &fileMeta{
		impl: impl,
		tags: []string{"auth"},
	}
}

func (a *fileMeta) Register(ws *restful.WebService) {
	storagePluginParam := ws.PathParameter("storagePlugin", "storage plugin to be used")
	objectNameParam := ws.PathParameter("objectName", "file object name in storage plugin")
	ws.Route(
		ws.GET("/storageplugin/{storagePlugin}/filemetas/{objectName:*}").To(a.GetFileMeta).
			Doc("Storage plugin put raw file").
			Param(objectNameParam).Param(storagePluginParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", v1alpha1.FileMeta{}),
	)
	// TODO: add list filemetas route
}

// GetFileMeta is handler of auth check route
func (a *fileMeta) GetFileMeta(req *restful.Request, resp *restful.Response) {
	pluginName := path.Parameter(req, "storagePlugin")
	objectName := path.Parameter(req, "objectName")

	ctx := req.Request.Context()
	meta, err := a.impl.GetFileMeta(storage.CtxWithPluginName(ctx, pluginName), objectName)
	if err != nil {
		kerrors.HandleError(req, resp, err)
		return
	}
	resp.WriteHeaderAndEntity(http.StatusOK, meta)
}
