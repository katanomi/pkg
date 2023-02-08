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
	storagePluginParam := ws.PathParameter("storagePlugin", "storage plugin to be used")
	objectNameParam := ws.PathParameter("objectName", "file object name in storage tools")
	fileTypeParam := ws.QueryParameter("fileType", "business type of file")
	ws.Route(
		ws.PUT("storageplugin/{storagePlugin}/fileobjects/{objectName:*}").To(a.PutFileObject).
			AllowedMethodsWithoutContentType([]string{http.MethodPut}).
			Produces(restful.MIME_OCTET).
			Doc("Storage plugin put raw file").
			Param(storagePluginParam).Param(objectNameParam).Param(fileTypeParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", v1alpha1.FileMeta{}),
	)

	ws.Route(
		ws.GET("storageplugin/{storagePlugin}/fileobjects/{objectName:*}").To(a.GetFileObject).
			AllowedMethodsWithoutContentType([]string{http.MethodGet}).
			Doc("Storage plugin get raw file").
			Param(storagePluginParam).Param(objectNameParam).Param(fileTypeParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", v1alpha1.FileMeta{}),
	)

	ws.Route(
		ws.DELETE("storageplugin/{storagePlugin}/fileobjects/{objectName:*}").To(a.DeleteFileObject).
			Doc("Storage plugin delete file by key").
			Param(storagePluginParam).Param(objectNameParam).Param(fileTypeParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", v1alpha1.FileMeta{}),
	)

}

// PutFileObject is handler of put file object
func (a *fileObject) PutFileObject(req *restful.Request, resp *restful.Response) {
	log := logging.FromContext(req.Request.Context())
	objectName := path.Parameter(req, "objectName")
	pluginName := path.Parameter(req, "storagePlugin")
	fileType := req.QueryParameter("fileType")

	fileSizeHeader := req.Request.Header.Get("Content-Length")
	fileSize, _ := strconv.Atoi(fileSizeHeader)

	// TODO: get filemeta and annotation from header x-katanomi-meta and x-katanomi-annotation-

	fileObject := filestorev1alpha1.FileObject{
		FileMeta: v1alpha1.FileMeta{
			Spec: v1alpha1.FileMetaSpec{
				ContentType:   req.Request.Header.Get(restful.HEADER_ContentType),
				FileType:      v1alpha1.FileType(fileType),
				ContentLength: int64(fileSize),
			},
		},
		FileReadCloser: req.Request.Body,
	}

	ctx := req.Request.Context()
	newMeta, err := a.impl.PutFileObject(storage.CtxWithPluginName(ctx, pluginName), objectName, &fileObject)
	if err != nil {
		log.Errorw("PutFileObject err", "err", err)
		kerrors.HandleError(req, resp, err)
		return
	}
	resp.WriteHeaderAndEntity(http.StatusCreated, newMeta)
}

// GetFileObject is handler of put file object
func (a *fileObject) GetFileObject(req *restful.Request, resp *restful.Response) {
	objectName := path.Parameter(req, "objectName")
	pluginName := path.Parameter(req, "storagePlugin")

	ctx := req.Request.Context()
	fileObject, err := a.impl.GetFileObject(storage.CtxWithPluginName(ctx, pluginName), objectName)

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
	objectName := path.Parameter(req, "objectName")
	pluginName := path.Parameter(req, "storagePlugin")

	ctx := req.Request.Context()
	err := a.impl.DeleteFileObject(storage.CtxWithPluginName(ctx, pluginName), objectName)
	if err != nil {
		kerrors.HandleError(req, resp, err)
		return
	}
	resp.WriteHeader(http.StatusOK)
}
