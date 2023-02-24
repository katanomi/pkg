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

// Package v1alpha1 defines versioned interfaces for archive capability
package v1alpha1

import (
	"context"

	archivev1alpha1 "github.com/katanomi/pkg/apis/archive/v1alpha1"
)

// ArchiveCapable defines methods of archive capability
//
//go:generate ../../../../../bin/mockgen -source=interface.go -destination=../../../../../testing/mock/github.com/katanomi/pkg/storage/capabilities/archive/v1alpha1/interface.go -package=v1alpha1 ArchiveCapable
type ArchiveCapable interface {
	// Upsert create or update a record
	Upsert(ctx context.Context, record *archivev1alpha1.Record) error
	// Delete delete a record
	Delete(ctx context.Context, cluster string, uid string, opts *archivev1alpha1.DeleteOption) error
	// DeleteBatch delete records by conditions
	DeleteBatch(ctx context.Context, conditions []archivev1alpha1.Condition, opts *archivev1alpha1.DeleteOption) error

	// ListRecords list records by conditions
	ListRecords(ctx context.Context, query archivev1alpha1.Query, opts *archivev1alpha1.ListOptions) (*archivev1alpha1.RecordList, error)
	// ListRelatedRecords list related records by conditions
	ListRelatedRecords(ctx context.Context, query archivev1alpha1.Query, opts *archivev1alpha1.ListOptions) (*archivev1alpha1.RecordList, error)
	// Aggregate aggregate records by conditions
	Aggregate(ctx context.Context, aggs archivev1alpha1.AggregateQuery, opts *archivev1alpha1.ListOptions) (*archivev1alpha1.AggregateResult, error)
}
