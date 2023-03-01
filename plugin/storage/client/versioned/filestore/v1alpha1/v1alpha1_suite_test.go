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
	"testing"

	"github.com/golang/mock/gomock"
	filestorev1alpha1 "github.com/katanomi/pkg/plugin/storage/capabilities/filestore/v1alpha1"
	v1alpha13 "github.com/katanomi/pkg/plugin/storage/client/versioned/filestore/v1alpha1"
	v1alpha1 "github.com/katanomi/pkg/testing/mock/github.com/katanomi/pkg/plugin/storage/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	ctx           context.Context
	mockSuiteCtrl *gomock.Controller

	mockStoragePluginClient *v1alpha1.MockInterface
	mockFileStoreClient     *v1alpha13.FileStoreV1alpha1Client
)

var _ = BeforeSuite(func() {
	ctx = context.Background()
	mockSuiteCtrl = gomock.NewController(GinkgoT())
	mockStoragePluginClient = v1alpha1.NewMockInterface(mockSuiteCtrl)

	mockStoragePluginClient.EXPECT().ForGroupVersion(&filestorev1alpha1.FileStoreV1alpha1GV).
		Return(mockStoragePluginClient)
	mockStoragePluginClient.EXPECT().APIVersion().Return(&filestorev1alpha1.FileStoreV1alpha1GV)

	mockFileStoreClient = v1alpha13.NewForClient(mockStoragePluginClient)
	Expect(mockFileStoreClient.RESTClient()).To(Equal(mockStoragePluginClient))
	Expect(mockFileStoreClient.RESTClient().APIVersion()).To(Equal(&filestorev1alpha1.FileStoreV1alpha1GV))
})

var _ = AfterSuite(func() {
	mockSuiteCtrl.Finish()
})

func TestV1alpha1(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "V1alpha1 Suite")
}
