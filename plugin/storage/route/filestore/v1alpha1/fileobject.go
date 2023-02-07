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
	"io"
	"net/http"
	"strconv"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/katanomi/pkg/apis/storage/v1alpha1"
	kerrors "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/plugin/path"
	"github.com/katanomi/pkg/plugin/storage"
	filestorev1alpha1 "github.com/katanomi/pkg/plugin/storage/capabilities/filestore/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/logging"
)

type fileObject struct {
	impl filestorev1alpha1.FileObjectInterface
	tags []string
}

// NewFileObject new route for auth checking
func NewFileObject(impl filestorev1alpha1.FileObjectInterface) storage.VersionedRouter {
	return &fileObject{
		impl: impl,
		tags: []string{"file-store"},
	}
}

func (a *fileObject) GroupVersion() schema.GroupVersion {
	return filestorev1alpha1.FileStoreV1alpha1GV
}

func (a *fileObject) Register(ws *restful.WebService) {
	storagePluginParam := ws.PathParameter("storageplugin", "storage plugin to be used")
	keyParam := ws.PathParameter("key", "file key for naming file in storage")
	fileTypeParam := ws.QueryParameter("filetype", "business type of file")
	ws.Route(
		ws.PUT("storageplugin/{storageplugin}/fileobjects/{key:*}").To(a.PutFileObject).
			Doc("Storage plugin put raw file").
			Param(storagePluginParam).Param(keyParam).Param(fileTypeParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", v1alpha1.FileMeta{}),
	)

	ws.Route(
		ws.GET("storageplugin/{storageplugin}/fileobjects/{key:*}").To(a.GetFileObject).
			Produces(restful.MIME_OCTET).
			Doc("Storage plugin get raw file").
			Param(storagePluginParam).Param(keyParam).Param(fileTypeParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", v1alpha1.FileMeta{}),
	)

	ws.Route(
		ws.DELETE("storageplugin/{storageplugin}/fileobjects/{key:*}").To(a.DeleteFileObject).
			Doc("Storage plugin delete file by key").
			Param(storagePluginParam).Param(keyParam).Param(fileTypeParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", v1alpha1.FileMeta{}),
	)

}

// PutFileObject is handler of put file object
func (a *fileObject) PutFileObject(req *restful.Request, resp *restful.Response) {
	log := logging.FromContext(req.Request.Context())
	key := path.Parameter(req, "key")
	pluginName := path.Parameter(req, "storageplugin")

	fileSizeHeader := req.Request.Header.Get("Content-Length")
	fileSize, _ := strconv.Atoi(fileSizeHeader)

	// TODO: get filemeta from header x-katanomi-meta
	filemeta := v1alpha1.FileMeta{
		Spec: v1alpha1.FileMetaSpec{
			Key:           key,
			ContentType:   req.Request.Header.Get(restful.HEADER_ContentType),
			FileType:      "",
			ContentLength: int64(fileSize),
		},
	}

	newMeta, err := a.impl.PutFileObject(req.Request.Context(), pluginName, req.Request.Body, filemeta)
	if err != nil {
		log.Errorw("PutFileObject err", "err", err)
		kerrors.HandleError(req, resp, err)
		return
	}
	resp.WriteHeaderAndEntity(http.StatusOK, newMeta)
}

// GetFileObject is handler of put file object
func (a *fileObject) GetFileObject(req *restful.Request, resp *restful.Response) {
	key := path.Parameter(req, "key")
	pluginName := path.Parameter(req, "storageplugin")

	fileObject, err := a.impl.GetFileObject(req.Request.Context(), pluginName, key)

	if err != nil {
		kerrors.HandleError(req, resp, err)
		return
	}

	_, err = io.Copy(resp.ResponseWriter, fileObject.FileReadCloser)
	if err != nil {
		kerrors.HandleError(req, resp, err)
		return
	}

	resp.AddHeader(restful.HEADER_ContentType, fileObject.Spec.ContentType)
	resp.WriteHeader(http.StatusOK)
}

// DeleteFileObject is handler of delete file object
func (a *fileObject) DeleteFileObject(req *restful.Request, resp *restful.Response) {
	key := path.Parameter(req, "key")
	pluginName := path.Parameter(req, "storageplugin")

	err := a.impl.DeleteFileObject(req.Request.Context(), pluginName, key)

	if err != nil {
		kerrors.HandleError(req, resp, err)
		return
	}
	resp.WriteHeader(http.StatusOK)
}
