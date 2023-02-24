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

package route

import (
	"context"
	"testing"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/apis/storage/v1alpha1"
	filestorev1alpha1 "github.com/katanomi/pkg/plugin/storage/capabilities/filestore/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type fakeFileStorePlugin struct{}

func (f *fakeFileStorePlugin) Path() string {
	return "fake-filestore"
}

func (f *fakeFileStorePlugin) Setup(ctx context.Context, logger *zap.SugaredLogger) error {
	// TODO implement me
	panic("implement me")
}

func (f *fakeFileStorePlugin) GetFileObject(ctx context.Context, objectName string) (*filestorev1alpha1.FileObject, error) {
	// TODO implement me
	panic("implement me")
}

func (f *fakeFileStorePlugin) PutFileObject(ctx context.Context, objectName string, obj *filestorev1alpha1.FileObject) (*v1alpha1.FileMeta, error) {
	// TODO implement me
	panic("implement me")
}

func (f *fakeFileStorePlugin) DeleteFileObject(ctx context.Context, objectName string) error {
	// TODO implement me
	panic("implement me")
}

func (f *fakeFileStorePlugin) ListFileMetas(ctx context.Context, opt *metav1.ListOptions) ([]v1alpha1.FileMeta, error) {
	// TODO implement me
	panic("implement me")
}

func (f *fakeFileStorePlugin) GetFileMeta(ctx context.Context, objectName string) (*v1alpha1.FileMeta, error) {
	// TODO implement me
	panic("implement me")
}

type fakeCorePlugin struct{}

func (f *fakeCorePlugin) Path() string {
	return "fake-core"
}

func (f *fakeCorePlugin) Setup(ctx context.Context, logger *zap.SugaredLogger) error {
	// TODO implement me
	panic("implement me")
}

func (f *fakeCorePlugin) CheckAuth(ctx context.Context, params []metav1alpha1.Param) (*v1alpha1.StorageAuthCheck, error) {
	// TODO implement me
	panic("implement me")
}

func TestRoute(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Route Suite")
}
