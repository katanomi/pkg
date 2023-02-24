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
	"testing"

	"github.com/katanomi/pkg/apis/archive/v1alpha1"
	"github.com/onsi/gomega"
)

type testStorage struct{}

func (t testStorage) Upsert(ctx context.Context, record *v1alpha1.Record) error {
	return nil
}

func (t testStorage) Delete(ctx context.Context, cluster string, uid string, opts *v1alpha1.DeleteOption) error {
	return nil
}

func (t testStorage) DeleteBatch(ctx context.Context, conditions []v1alpha1.Condition, opts *v1alpha1.DeleteOption) error {
	return nil
}

func (t testStorage) ListRecords(ctx context.Context, query v1alpha1.Query, opts *v1alpha1.ListOptions) (*v1alpha1.RecordList, error) {
	return nil, nil
}

func (t testStorage) ListRelatedRecords(ctx context.Context, query v1alpha1.Query, opts *v1alpha1.ListOptions) (*v1alpha1.RecordList, error) {
	return nil, nil
}

func (t testStorage) Aggregate(ctx context.Context, aggs v1alpha1.AggregateQuery, opts *v1alpha1.ListOptions) (*v1alpha1.AggregateResult, error) {
	return nil, nil
}

func TestWithStorage(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	ctx := context.Background()
	s := &testStorage{}
	gotStorage := GetStorage(WithStorage(ctx, s))
	g.Expect(gotStorage).To(gomega.Equal(s))
}
