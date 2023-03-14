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
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/emicklei/go-restful/v3"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	storagev1alpha1 "github.com/katanomi/pkg/apis/storage/v1alpha1"
	pkgclient "github.com/katanomi/pkg/client"
	filestoreroute "github.com/katanomi/pkg/plugin/storage/route/filestore/v1alpha1"
	ktesting "github.com/katanomi/pkg/testing"
	"github.com/katanomi/pkg/testing/mock/github.com/katanomi/pkg/plugin/storage/capabilities/filestore/v1alpha1"
	mockClient "github.com/katanomi/pkg/testing/mock/sigs.k8s.io/controller-runtime/pkg/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/authorization/v1"
	apiserveruser "k8s.io/apiserver/pkg/authentication/user"
	apiserverrequest "k8s.io/apiserver/pkg/endpoints/request"
)

var _ = Describe("Filemeta", func() {

	var (
		ctx       context.Context
		recorder  *httptest.ResponseRecorder
		fileMetas []storagev1alpha1.FileMeta
		container *restful.Container
		resp      *restful.Response
		mClient   *mockClient.MockClient
	)

	BeforeEach(func() {
		ctx = context.Background()
		ktesting.MustLoadJSON("testdata/filemetas.stroage.route.response.golden.json", &fileMetas)
		recorder = httptest.NewRecorder()

		mockCtl := gomock.NewController(GinkgoT())
		fileMetaCapability := v1alpha1.NewMockFileMetaInterface(mockCtl)
		fileMetaCapability.EXPECT().ListFileMetas(gomock.Any(),
			gomock.AssignableToTypeOf(storagev1alpha1.FileMetaListOptions{})).Return(fileMetas, nil)

		mClient = mockClient.NewMockClient(mockCtl)
		ctx = pkgclient.WithClient(ctx, mClient)

		ctx = apiserverrequest.WithUser(ctx, &apiserveruser.DefaultInfo{Name: "system:serviceaccount:devops:foo"})

		resourceAttributes := storagev1alpha1.FileMetaResourceAttributes("list")
		mClient.EXPECT().Create(gomock.AssignableToTypeOf(ctx),
			gomock.AssignableToTypeOf(&v1.SelfSubjectAccessReview{})).
			SetArg(1, v1.SelfSubjectAccessReview{
				Spec: v1.SelfSubjectAccessReviewSpec{
					ResourceAttributes: &resourceAttributes,
				},
				Status: v1.SubjectAccessReviewStatus{
					Allowed: true,
				},
			}).Return(nil)

		container = restful.NewContainer()
		container.Router(restful.RouterJSR311{})
		ws := &restful.WebService{}
		restful.DefaultResponseContentType(restful.MIME_JSON)
		fileMeta := filestoreroute.NewFileMeta(fileMetaCapability)
		Expect(fileMeta.Register(ctx, ws)).Should(Succeed())
		ws.Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
		container.Add(ws)
	})

	JustBeforeEach(func() {
		req := httptest.NewRequest(http.MethodGet, "/storageplugins/minio.default/filemetas", nil)
		req = req.WithContext(ctx)
		req.Header.Add("Accept", "*/*")

		recorder = httptest.NewRecorder()
		resp = restful.NewResponse(recorder)
		resp.SetRequestAccepts("application/json")
		container.Dispatch(resp, req)
	})

	Context("put file with meta", func() {
		It("returns file meta response", func() {
			Expect(resp.StatusCode()).Should(Equal(http.StatusOK))
			var respFileMetas []storagev1alpha1.FileMeta
			json.Unmarshal(recorder.Body.Bytes(), &respFileMetas)
			Expect(cmp.Diff(respFileMetas, fileMetas)).To(BeEmpty())
		})
	})
})
