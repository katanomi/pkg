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

package v1alpha1_test

import (
	"context"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/katanomi/pkg/apis/storage/v1alpha1"
	v1alpha13 "github.com/katanomi/pkg/plugin/storage/capabilities/filestore/v1alpha1"
	v1alpha12 "github.com/katanomi/pkg/plugin/storage/client/versioned/filestore/v1alpha1"
	"github.com/katanomi/pkg/testing"
	filestorev1alpha1 "github.com/katanomi/pkg/testing/mock/github.com/katanomi/pkg/plugin/storage/client/versioned/filestore/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Test.Fileobject", func() {

	var (
		fileObject   v1alpha12.FileObjectInterface
		fakeFileMeta v1alpha1.FileMeta
		httpHeaders  http.Header
	)

	BeforeEach(func() {
		fileObject = mockFileStoreClient.FileObject("foo")
		testing.MustLoadYaml("testdata/filemeta.normal.yaml", &fakeFileMeta)
		httpHeaders = make(http.Header)
		httpHeaders.Set(v1alpha1.HeaderFileMeta, fakeFileMeta.Encode())
	})

	Context("Test.FileObject.GET", func() {
		It("returns fileobject", func() {
			mockStoragePluginClient.EXPECT().
				GetResponse(ctx, "storageplugins/foo/fileobjects/dir1/file1", gomock.Any()).
				Return(&resty.Response{
					RawResponse: &http.Response{
						StatusCode: http.StatusOK,
						Header:     httpHeaders,
					},
				}, nil)

			expectedFileMeta := v1alpha1.FileMeta{}
			testing.MustLoadYaml("testdata/filemeta.normal.golden.yaml", &expectedFileMeta)

			fileObject, err := fileObject.GET(ctx, "dir1/file1")
			Expect(err).ShouldNot(HaveOccurred())

			Expect(cmp.Diff(fileObject.FileMeta, expectedFileMeta)).To(BeEmpty())
		})
	})

	Context("Test.FileObject.PUT", func() {
		It("returns fileobject with meta", func() {
			mockStoragePluginClient.EXPECT().
				Put(ctx, "storageplugins/foo/fileobjects/file", gomock.Any()).
				Return(nil)

			expectedFileMeta := v1alpha1.FileMeta{}
			testing.MustLoadYaml("testdata/filemeta.normal.golden.yaml", &expectedFileMeta)

			mockFileObject := v1alpha13.FileObject{
				FileMeta: fakeFileMeta,
			}
			fileMeta, err := fileObject.PUT(ctx, mockFileObject)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(fileMeta).ShouldNot(BeNil())
		})
	})

	Context("Test.FileObject.DELETE", func() {
		It("returns ok", func() {
			mockStoragePluginClient.EXPECT().
				Delete(ctx, "storageplugins/foo/fileobjects/file").
				Return(nil)

			Expect(fileObject.DELETE(ctx, "file")).To(Succeed())
		})
	})
})

var _ = Describe("Test.FileObject.Context", func() {
	var (
		ctx        context.Context
		fileObjCli v1alpha12.FileObjectInterface
	)

	BeforeEach(func() {
		ctx = context.Background()
		fileObjCli = &filestorev1alpha1.MockFileObjectInterface{}
	})

	Context("ContextWithFileObjectClient", func() {
		It("can extract FileObjectClient from context", func() {
			ctx = v1alpha12.WithFileObjectClient(ctx, fileObjCli)
			fileObjCliFromCtx := v1alpha12.FileObjectClientFrom(ctx)
			Expect(fileObjCliFromCtx).To(Equal(fileObjCli))
		})
	})
})
