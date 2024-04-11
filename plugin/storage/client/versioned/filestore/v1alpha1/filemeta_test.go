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
	"net/http"

	"github.com/katanomi/pkg/apis/storage/v1alpha1"
	v1alpha12 "github.com/katanomi/pkg/plugin/storage/client/versioned/filestore/v1alpha1"
	"github.com/katanomi/pkg/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Test.Fileobject", func() {

	var (
		fileMeta     v1alpha12.FileMetaInterface
		fakeFileMeta v1alpha1.FileMeta
		httpHeaders  http.Header
	)

	BeforeEach(func() {
		fileMeta = mockFileStoreClient.FileMeta("foo")
		testing.MustLoadYaml("testdata/filemeta.normal.yaml", &fakeFileMeta)
		httpHeaders = make(http.Header)
		httpHeaders.Set(v1alpha1.HeaderFileMeta, fakeFileMeta.Encode())
	})

	Context("Test.FileMeta.GET", func() {
		It("returns filemeta", func() {
			mockStoragePluginClient.EXPECT().
				Get(ctx, "storageplugins/foo/filemetas/dir1/file1", gomock.Any()).
				Return(nil)

			expectedFileMeta := v1alpha1.FileMeta{}
			testing.MustLoadYaml("testdata/filemeta.normal.golden.yaml", &expectedFileMeta)

			_, err := fileMeta.GET(ctx, "dir1/file1")
			Expect(err).ShouldNot(HaveOccurred())

		})
	})

	Context("Test.FileMeta.List", func() {
		It("returns filemetas", func() {

			expectedFileMetas := []v1alpha1.FileMeta{}
			testing.MustLoadJSON("testdata/filemetas.normal.golden.json", &expectedFileMetas)

			mockStoragePluginClient.EXPECT().
				Get(ctx, "storageplugins/foo/filemetas", gomock.Any()).
				Return(nil)

			_, err := fileMeta.List(ctx, v1alpha1.FileMetaListOptions{})
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

})
