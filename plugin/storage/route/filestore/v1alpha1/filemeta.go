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
	"net/http"
	"strconv"

	kclient "github.com/katanomi/pkg/client"
	"github.com/katanomi/pkg/plugin/route"
	"go.uber.org/zap"
	"knative.dev/pkg/logging"

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
	*zap.SugaredLogger
}

func (a *fileMeta) GroupVersion() schema.GroupVersion {
	return filestorev1alpha1.FileStoreV1alpha1GV
}

// NewFileMeta new route for auth checking
func NewFileMeta(impl filestorev1alpha1.FileMetaInterface) storage.VersionedRouter {
	return &fileMeta{
		impl: impl,
		tags: []string{"filemeta"},
	}
}

func (a *fileMeta) Register(ctx context.Context, ws *restful.WebService) error {
	storagePluginParam := ws.PathParameter("storagePlugin", "storage plugin to be used")
	objectNameParam := ws.PathParameter("objectName", "file object name in storage plugin")

	a.SugaredLogger = logging.FromContext(ctx).With("resource", "filemata")

	ws.Route(
		ws.GET("/storageplugins/{storagePlugin}/filemetas/{objectName:*}").To(a.GetFileMeta).
			Filter(kclient.SubjectReviewFilterForResource(ctx, v1alpha1.FileMetaResourceAttributes("get"), "", "")).
			Doc("Storage plugin get file meta").
			Param(objectNameParam).Param(storagePluginParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", v1alpha1.FileMeta{}),
	)

	ws.Route(
		ws.GET("/storageplugins/{storagePlugin}/filemetas").To(a.ListFileMetas).
			Filter(kclient.SubjectReviewFilterForResource(ctx, v1alpha1.FileMetaResourceAttributes("get"), "", "")).
			Doc("Storage plugin list file metas").
			Param(storagePluginParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", []v1alpha1.FileMeta{}),
	)

	return nil
}

// GetFileMeta is handler of getting file meta
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

// ListFileMetas is handler of listing file meta
func (a *fileMeta) ListFileMetas(req *restful.Request, resp *restful.Response) {
	pluginName := path.Parameter(req, "storagePlugin")
	listOptsFromReq := route.GetListOptionsFromRequest(req)
	a.Debugw("request received for file meta", "opts", listOptsFromReq)

	// trim plugin name if exists
	prefix := listOptsFromReq.GetSearchFirstElement("prefix")

	limitSearch := listOptsFromReq.GetSearchFirstElement("limit")

	limit, _ := strconv.Atoi(limitSearch)

	recursiveSearch := listOptsFromReq.GetSearchFirstElement("recursive")
	startAfter := listOptsFromReq.GetSearchFirstElement("startAfter")
	recursive := true
	if parsedRecursiveSearch, err := strconv.ParseBool(recursiveSearch); err == nil {
		recursive = parsedRecursiveSearch
	}

	ctx := req.Request.Context()
	metas, err := a.impl.ListFileMetas(storage.CtxWithPluginName(ctx, pluginName), v1alpha1.FileMetaListOptions{
		Prefix:     prefix,
		Recursive:  recursive,
		StartAfter: startAfter,
		Limit:      limit,
	})
	if err != nil {
		kerrors.HandleError(req, resp, err)
		return
	}
	resp.WriteHeaderAndEntity(http.StatusOK, metas)
}
