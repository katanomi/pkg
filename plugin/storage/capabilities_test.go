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

package storage

import (
	"context"
	"io"
	"testing"

	archivev1alpha1 "github.com/katanomi/pkg/apis/archive/v1alpha1"
	apistoragev1alpha1 "github.com/katanomi/pkg/apis/storage/v1alpha1"
	"github.com/katanomi/pkg/plugin/storage/capabilities/filestore/v1alpha1"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type fakeFileStoreImp struct{}

func (f fakeFileStoreImp) GetFileObject(ctx context.Context, key string) (v1alpha1.FileObject, error) {
	// TODO implement me
	panic("implement me")
}

func (f fakeFileStoreImp) PutFileObject(ctx context.Context, key string, fileReader io.ReadCloser,
	meta apistoragev1alpha1.FileMeta) (apistoragev1alpha1.FileMeta, error) {
	// TODO implement me
	panic("implement me")
}

func (f fakeFileStoreImp) DeleteFileObject(ctx context.Context) error {
	// TODO implement me
	panic("implement me")
}

func (f fakeFileStoreImp) ListFileMetas(ctx context.Context, opt *metav1.ListOptions) ([]apistoragev1alpha1.FileMeta, error) {
	// TODO implement me
	panic("implement me")
}

func (f fakeFileStoreImp) GetFileMeta(ctx context.Context, key string) (apistoragev1alpha1.FileMeta, error) {
	// TODO implement me
	panic("implement me")
}

type fakeArchiveImp struct{}

func (f fakeArchiveImp) Upsert(ctx context.Context, record *archivev1alpha1.Record) error {
	//TODO implement me
	panic("implement me")
}

func (f fakeArchiveImp) Delete(ctx context.Context, cluster string, uid string, opts *archivev1alpha1.DeleteOption) error {
	//TODO implement me
	panic("implement me")
}

func (f fakeArchiveImp) DeleteBatch(ctx context.Context, conditions []archivev1alpha1.Condition, opts *archivev1alpha1.DeleteOption) error {
	//TODO implement me
	panic("implement me")
}

func (f fakeArchiveImp) ListRecords(ctx context.Context, query archivev1alpha1.Query, opts *archivev1alpha1.ListOptions) (*archivev1alpha1.RecordList, error) {
	//TODO implement me
	panic("implement me")
}

func (f fakeArchiveImp) ListRelatedRecords(ctx context.Context, query archivev1alpha1.Query, opts *archivev1alpha1.ListOptions) (*archivev1alpha1.RecordList, error) {
	//TODO implement me
	panic("implement me")
}

func (f fakeArchiveImp) Aggregate(ctx context.Context, aggs archivev1alpha1.AggregateQuery, opts *archivev1alpha1.ListOptions) (*archivev1alpha1.AggregateResult, error) {
	//TODO implement me
	panic("implement me")
}

type fakeMultipleImp struct {
	fakeArchiveImp
	fakeFileStoreImp
}

func TestGetImplementedCapabilities(t *testing.T) {
	g := NewGomegaWithT(t)

	tests := []struct {
		name string
		obj  interface{}
		want Capabilities
	}{
		{
			name: "nil object returns nil",
			obj:  nil,
			want: nil,
		},
		{
			name: "file-store capability",
			obj:  fakeFileStoreImp{},
			want: Capabilities{CapabilityFileStore},
		},
		{
			name: "file-store capability pointer",
			obj:  &fakeFileStoreImp{},
			want: Capabilities{CapabilityFileStore},
		},
		{
			name: "archive capability",
			obj:  fakeArchiveImp{},
			want: Capabilities{CapabilityArchive},
		},
		{
			name: "multiple capabilities",
			obj:  fakeMultipleImp{},
			want: Capabilities{CapabilityFileStore, CapabilityArchive},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g.Expect(GetImplementedCapabilities(tt.obj)).To(ContainElements(tt.want))
		})
	}
}
