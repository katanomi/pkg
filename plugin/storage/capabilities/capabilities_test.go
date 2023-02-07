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

package capabilities

import (
	"context"
	"io"
	"testing"

	"github.com/katanomi/pkg/apis/storage/v1alpha1"

	archivev1alpha1 "github.com/katanomi/pkg/apis/archive/v1alpha1"
	archivecapv1alpha1 "github.com/katanomi/pkg/plugin/storage/capabilities/archive/v1alpha1"
	filestorev1alpha1 "github.com/katanomi/pkg/plugin/storage/capabilities/filestore/v1alpha1"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type fakeFileStoreImp struct{}

func (f fakeFileStoreImp) GetFileObject(ctx context.Context,
	key string) (*filestorev1alpha1.FileObject, error) {
	// TODO implement me
	panic("implement me")
}

func (f fakeFileStoreImp) PutFileObject(ctx context.Context, fileReadCloser io.ReadCloser,
	meta v1alpha1.FileMeta) (*v1alpha1.FileMeta, error) {
	// TODO implement me
	panic("implement me")
}

func (f fakeFileStoreImp) DeleteFileObject(ctx context.Context, key string) error {
	// TODO implement me
	panic("implement me")
}

func (f fakeFileStoreImp) ListFileMetas(ctx context.Context,
	opt *metav1.ListOptions) ([]v1alpha1.FileMeta, error) {
	// TODO implement me
	panic("implement me")
}

func (f fakeFileStoreImp) GetFileMeta(ctx context.Context, key string) (*v1alpha1.FileMeta, error) {
	// TODO implement me
	panic("implement me")
}

type fakeArchiveImp struct{}

func (f *fakeArchiveImp) Upsert(ctx context.Context, record *archivev1alpha1.Record) error {
	// TODO implement me
	panic("implement me")
}

func (f fakeArchiveImp) Delete(ctx context.Context, cluster string, uid string, opts *archivev1alpha1.DeleteOption) error {
	// TODO implement me
	panic("implement me")
}

func (f fakeArchiveImp) DeleteBatch(ctx context.Context, conditions []archivev1alpha1.Condition, opts *archivev1alpha1.DeleteOption) error {
	// TODO implement me
	panic("implement me")
}

func (f fakeArchiveImp) ListRecords(ctx context.Context, query archivev1alpha1.Query, opts *archivev1alpha1.ListOptions) (*archivev1alpha1.RecordList, error) {
	// TODO implement me
	panic("implement me")
}

func (f fakeArchiveImp) ListRelatedRecords(ctx context.Context, query archivev1alpha1.Query, opts *archivev1alpha1.ListOptions) (*archivev1alpha1.RecordList, error) {
	// TODO implement me
	panic("implement me")
}

func (f fakeArchiveImp) Aggregate(ctx context.Context, aggs archivev1alpha1.AggregateQuery, opts *archivev1alpha1.ListOptions) (*archivev1alpha1.AggregateResult, error) {
	// TODO implement me
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
		want []string
	}{
		{
			name: "nil object returns nil",
			obj:  nil,
			want: nil,
		},
		{
			name: "file-store capability",
			obj:  fakeFileStoreImp{},
			want: []string{},
		},
		{
			name: "file-store capability pointer",
			obj:  &fakeFileStoreImp{},
			want: []string{filestorev1alpha1.FileStoreV1alpha1GV.String()},
		},
		{
			name: "archive capability",
			obj:  &fakeArchiveImp{},
			want: []string{archivecapv1alpha1.ArchiveV1alpha1.String()},
		},
		{
			name: "multiple capabilities",
			obj:  &fakeMultipleImp{},
			want: []string{filestorev1alpha1.FileStoreV1alpha1GV.String(), archivecapv1alpha1.ArchiveV1alpha1.String()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g.Expect(GetImplementedCapabilities(tt.obj)).To(ContainElements(tt.want))
		})
	}
}
