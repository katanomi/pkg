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

// ArchiveCapable defines methods of archive capability
// TODO: will be updated later
type ArchiveCapable interface {
	// Foo is temporary for demo
	Foo()
	// Upsert(ctx context.Context, record *v1alpha1.Record) error
	// Delete(ctx context.Context, cluster string, uid string, opts *DeleteOption) error
	// DeleteBatch(ctx context.Context, conditions []Condition, opts *DeleteOption) error

	// ListRecords(ctx context.Context, query Query, opts *ListOptions) (*v1alpha1.RecordList, error)
	// ListRelatedRecords(ctx context.Context, query Query, opts *ListOptions) (*v1alpha1.RecordList, error)
	// Aggregate(ctx context.Context, aggs AggregateQuery, opts *ListOptions) ([]map[string]interface{}, error)
	//
}
