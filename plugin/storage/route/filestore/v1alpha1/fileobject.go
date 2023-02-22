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
	"context"
	"fmt"
	"io"
	"net/http"

	kclient "github.com/katanomi/pkg/client"

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

func (a *fileObject) Register(ctx context.Context, ws *restful.WebService) error {
	storagePluginParam := ws.PathParameter("storagePlugin", "storage plugin toe used")
	objectNameParam := ws.PathParameter("objectName", "file object name in storage tools")

	if manager := kclient.ManagerCtx(ctx); manager != nil {
		filters, err := manager.Filters(ctx)
		if err != nil {
			return err
		}
		for _, filter := range filters {
			ws = ws.Filter(filter)
		}
	}

	ws.Route(
		ws.PUT("storageplugins/{storagePlugin}/fileobjects/{objectName:*}").To(a.PutFileObject).
			Filter(kclient.SubjectReviewFilterForResource(ctx, v1alpha1.FileObjectResourceAttributes("update"), "", "")).
			AllowedMethodsWithoutContentType([]string{http.MethodPut}).
			Consumes(v1alpha1.SupportedContentTypeList...).
			Produces(restful.MIME_JSON).
			Doc("Storage plugin put raw file").
			Param(storagePluginParam).Param(objectNameParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", v1alpha1.FileMeta{}),
	)

	ws.Route(
		ws.GET("storageplugins/{storagePlugin}/fileobjects/{objectName:*}").To(a.GetFileObject).
			Filter(kclient.SubjectReviewFilterForResource(ctx, v1alpha1.FileObjectResourceAttributes("get"), "", "")).
			AllowedMethodsWithoutContentType([]string{http.MethodGet}).
			Doc("Storage plugin get raw file").
			Param(storagePluginParam).Param(objectNameParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", v1alpha1.FileMeta{}),
	)

	ws.Route(
		ws.DELETE("storageplugins/{storagePlugin}/fileobjects/{objectName:*}").To(a.DeleteFileObject).
			Filter(kclient.SubjectReviewFilterForResource(ctx, v1alpha1.FileObjectResourceAttributes("delete"), "", "")).
			Doc("Storage plugin delete file by key").
			Param(storagePluginParam).Param(objectNameParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", v1alpha1.FileMeta{}),
	)

	return nil
}

// PutFileObject is handler of put file object
func (a *fileObject) PutFileObject(req *restful.Request, resp *restful.Response) {
	log := logging.FromContext(req.Request.Context())
	pluginName := path.Parameter(req, "storagePlugin")
	objectName := path.Parameter(req, "objectName")

	fileMetaString := req.HeaderParameter(v1alpha1.HeaderFileMeta)
	if fileMetaString == "" {
		log.Errorw("get empty file meta from header",
			"objectName", objectName, "pluginName", pluginName,
		)
		kerrors.HandleError(req, resp, fmt.Errorf("empty filemeta from header"))
		return
	}
	fileMeta, err := v1alpha1.DecodeAsFileMeta(fileMetaString)
	if err != nil {
		log.Errorw("unmarshal file meta err",
			"err", err,
		)
		kerrors.HandleError(req, resp, err)
		return
	}

	fileMeta.Name = objectName
	fileObject := filestorev1alpha1.FileObject{
		FileMeta:       *fileMeta,
		FileReadCloser: req.Request.Body,
	}

	ctx := req.Request.Context()
	newMeta, err := a.impl.PutFileObject(storage.CtxWithPluginName(ctx, pluginName), &fileObject)
	if err != nil {
		log.Errorw("PutFileObject err", "err", err)
		kerrors.HandleError(req, resp, err)
		return
	}
	resp.WriteHeaderAndEntity(http.StatusCreated, newMeta)
}

// GetFileObject is handler of put file object
func (a *fileObject) GetFileObject(req *restful.Request, resp *restful.Response) {
	pluginName := path.Parameter(req, "storagePlugin")
	objectName := path.Parameter(req, "objectName")

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
	pluginName := path.Parameter(req, "storagePlugin")
	objectName := path.Parameter(req, "objectName")

	ctx := req.Request.Context()
	err := a.impl.DeleteFileObject(storage.CtxWithPluginName(ctx, pluginName), objectName)
	if err != nil {
		kerrors.HandleError(req, resp, err)
		return
	}
	resp.WriteHeader(http.StatusOK)
}
