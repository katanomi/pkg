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

// Package v1alpha1 for archive v1alpha1 routes
package v1alpha1

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/katanomi/pkg/apis/archive/v1alpha1"
	kerrors "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/plugin/path"
	"github.com/katanomi/pkg/plugin/storage"
	archivev1alpha1 "github.com/katanomi/pkg/plugin/storage/capabilities/archive/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type archive struct {
	impl archivev1alpha1.ArchiveCapable
	//impl archivev1alpha1.ArchiveCapable
	tags []string
}

func (a *archive) GroupVersion() schema.GroupVersion {
	return archivev1alpha1.ArchiveV1alpha1GV
}

// NewArchive new route for auth checking
func NewArchive(impl archivev1alpha1.ArchiveCapable) storage.VersionedRouter {
	return &archive{
		impl: impl,
		tags: []string{"archive"},
	}
}

func (a *archive) Register(ws *restful.WebService) {
	storagePluginParam := ws.PathParameter("storageplugin", "storage plugin to be used")
	ws.Route(
		ws.POST("storageplugin/{storageplugin}/records").To(a.ListRecords).
			Doc("List archive records").
			Param(storagePluginParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", &v1alpha1.RecordList{}),
	)

	ws.Route(
		ws.POST("storageplugin/{storageplugin}/relatedRecords").To(a.ListRelatedRecords).
			Doc("List archive related records").
			Param(storagePluginParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", &v1alpha1.RecordList{}),
	)

	ws.Route(
		ws.POST("storageplugin/{storageplugin}/record").To(a.UpsertRecord).
			Doc("Create or update archive record").
			Param(storagePluginParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", nil),
	)

	ws.Route(
		ws.DELETE("storageplugin/{storageplugin}/clusters/{cluster}/uids/{uid}").To(a.DeleteRecord).
			Doc("Delete archive record").
			Param(storagePluginParam).
			Param(ws.PathParameter("cluster", "cluster name")).
			Param(ws.PathParameter("uid", "resource uid")).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", nil),
	)

	ws.Route(
		ws.DELETE("storageplugin/{storageplugin}/records").To(a.BatchDeleteRecord).
			Doc("Delete archive record batch").
			Param(storagePluginParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", nil),
	)

	ws.Route(
		ws.POST("storageplugin/{storageplugin}/aggregate").To(a.Aggregate).
			Doc("Aggregate query").
			Param(storagePluginParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", &v1alpha1.AggregateResult{}),
	)
}

// ListRecords to get archive record list
func (a *archive) ListRecords(req *restful.Request, resp *restful.Response) {
	pluginName := path.Parameter(req, "storageplugin")
	ctx := req.Request.Context()
	ctx = storage.CtxWithPluginName(ctx, pluginName)

	params := v1alpha1.ListParams{}
	err := req.ReadEntity(&params)
	if err != nil {
		kerrors.HandleError(req, resp, errors.NewBadRequest(err.Error()))
		return
	}

	list, err := a.impl.ListRecords(ctx, params.Query, params.Options)
	if err != nil {
		kerrors.HandleError(req, resp, err)
		return
	}
	resp.WriteHeaderAndEntity(http.StatusOK, list)
}

// ListRelatedRecords to get archive related record list
func (a *archive) ListRelatedRecords(req *restful.Request, resp *restful.Response) {
	pluginName := path.Parameter(req, "storageplugin")
	ctx := req.Request.Context()
	ctx = storage.CtxWithPluginName(ctx, pluginName)

	params := v1alpha1.ListParams{}
	err := req.ReadEntity(&params)
	if err != nil {
		kerrors.HandleError(req, resp, errors.NewBadRequest(err.Error()))
		return
	}

	list, err := a.impl.ListRelatedRecords(ctx, params.Query, params.Options)
	if err != nil {
		kerrors.HandleError(req, resp, err)
		return
	}
	resp.WriteHeaderAndEntity(http.StatusOK, list)
}

// UpsertRecord to create or update archive record
func (a *archive) UpsertRecord(req *restful.Request, resp *restful.Response) {
	pluginName := path.Parameter(req, "storageplugin")
	ctx := req.Request.Context()
	ctx = storage.CtxWithPluginName(ctx, pluginName)

	params := v1alpha1.Record{}
	err := req.ReadEntity(&params)
	if err != nil {
		kerrors.HandleError(req, resp, errors.NewBadRequest(err.Error()))
		return
	}

	err = a.impl.Upsert(ctx, &params)
	if err != nil {
		kerrors.HandleError(req, resp, err)
		return
	}
	resp.WriteHeader(http.StatusOK)
}

// DeleteRecord to delete archive record
func (a *archive) DeleteRecord(req *restful.Request, resp *restful.Response) {
	pluginName := path.Parameter(req, "storageplugin")
	cluster := path.Parameter(req, "cluster")
	uid := path.Parameter(req, "uid")
	ctx := req.Request.Context()
	ctx = storage.CtxWithPluginName(ctx, pluginName)

	params := v1alpha1.DeleteParams{}
	err := req.ReadEntity(&params)
	if err != nil {
		kerrors.HandleError(req, resp, errors.NewBadRequest(err.Error()))
		return
	}

	err = a.impl.Delete(ctx, cluster, uid, params.Options)
	if err != nil {
		kerrors.HandleError(req, resp, err)
		return
	}
	resp.WriteHeader(http.StatusOK)
}

// BatchDeleteRecord to delete archive record batch
func (a *archive) BatchDeleteRecord(req *restful.Request, resp *restful.Response) {
	pluginName := path.Parameter(req, "storageplugin")
	ctx := req.Request.Context()
	ctx = storage.CtxWithPluginName(ctx, pluginName)

	params := v1alpha1.DeleteParams{}
	err := req.ReadEntity(&params)
	if err != nil {
		kerrors.HandleError(req, resp, errors.NewBadRequest(err.Error()))
		return
	}

	err = a.impl.DeleteBatch(ctx, params.Conditions, params.Options)
	if err != nil {
		kerrors.HandleError(req, resp, err)
		return
	}
	resp.WriteHeader(http.StatusOK)
}

// Aggregate to delete archive record batch
func (a *archive) Aggregate(req *restful.Request, resp *restful.Response) {
	pluginName := path.Parameter(req, "storageplugin")
	ctx := req.Request.Context()
	ctx = storage.CtxWithPluginName(ctx, pluginName)

	params := v1alpha1.AggregateParams{}
	err := req.ReadEntity(&params)
	if err != nil {
		kerrors.HandleError(req, resp, errors.NewBadRequest(err.Error()))
		return
	}

	ret, err := a.impl.Aggregate(ctx, params.Query, params.Options)
	if err != nil {
		kerrors.HandleError(req, resp, err)
		return
	}
	resp.WriteHeaderAndEntity(http.StatusOK, ret)
}
