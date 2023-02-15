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

// Package v1alpha1 for archive v1alpha1 client
package v1alpha1

import (
	"context"
	"fmt"

	"github.com/katanomi/pkg/apis/archive/v1alpha1"
	client2 "github.com/katanomi/pkg/plugin/client"
	archivev1alpha1 "github.com/katanomi/pkg/plugin/storage/capabilities/archive/v1alpha1"
	"github.com/katanomi/pkg/plugin/storage/client"
)

// RecordGetter returns Record getter object
type RecordGetter interface {
	Record(pluginName string) RecordInterface
}

//go:generate ../../../../../../bin/mockgen -source=archive.go -destination=../../../../../../testing/mock/github.com/katanomi/pkg/storage/client/versioned/archive/v1alpha1/interface.go -package=v1alpha1 RecordInterface
type RecordInterface interface {
	archivev1alpha1.ArchiveCapable
}

func newRecord(clt *ArchiveClient, pluginName string) *archive {
	return &archive{
		client:     clt.RESTClient(),
		pluginName: pluginName,
	}
}

type archive struct {
	client     client.Interface
	pluginName string
}

func (a archive) Upsert(ctx context.Context, record *v1alpha1.Record) error {
	path := fmt.Sprintf("storageplugin/%s/record", a.pluginName)
	err := a.client.Post(ctx, path, client2.BodyOpts(record))
	if err != nil {
		return err
	}
	return nil
}

func (a archive) Delete(ctx context.Context, cluster string, uid string, opts *v1alpha1.DeleteOption) error {
	path := fmt.Sprintf("storageplugin/%s/clusters/%s/uids/%s", a.pluginName, cluster, uid)
	err := a.client.Delete(ctx, path, client2.BodyOpts(opts))
	if err != nil {
		return err
	}
	return nil
}

func (a archive) DeleteBatch(ctx context.Context, conditions []v1alpha1.Condition, opts *v1alpha1.DeleteOption) error {
	path := fmt.Sprintf("storageplugin/%s/records", a.pluginName)
	params := v1alpha1.DeleteParams{
		Conditions: conditions,
		Options:    opts,
	}
	err := a.client.Delete(ctx, path, client2.BodyOpts(params))
	if err != nil {
		return err
	}
	return nil
}

func (a archive) ListRecords(ctx context.Context, query v1alpha1.Query, opts *v1alpha1.ListOptions) (*v1alpha1.RecordList, error) {
	path := fmt.Sprintf("storageplugin/%s/records", a.pluginName)
	params := v1alpha1.ListParams{
		Query:   query,
		Options: opts,
	}
	ret := &v1alpha1.RecordList{}
	err := a.client.Post(ctx, path, client2.ResultOpts(ret), client2.BodyOpts(params))
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (a archive) ListRelatedRecords(ctx context.Context, query v1alpha1.Query, opts *v1alpha1.ListOptions) (*v1alpha1.RecordList, error) {
	path := fmt.Sprintf("storageplugin/%s/relatedRecords", a.pluginName)
	params := v1alpha1.ListParams{
		Query:   query,
		Options: opts,
	}
	ret := &v1alpha1.RecordList{}
	err := a.client.Post(ctx, path, client2.ResultOpts(ret), client2.BodyOpts(params))
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (a archive) Aggregate(ctx context.Context, query v1alpha1.AggregateQuery, opts *v1alpha1.ListOptions) (*v1alpha1.AggregateResult, error) {
	path := fmt.Sprintf("storageplugin/%s/aggregate", a.pluginName)
	params := v1alpha1.AggregateParams{
		Query:   query,
		Options: opts,
	}
	ret := &v1alpha1.AggregateResult{}
	err := a.client.Post(ctx, path, client2.ResultOpts(ret), client2.BodyOpts(params))
	if err != nil {
		return nil, err
	}
	return ret, nil
}
