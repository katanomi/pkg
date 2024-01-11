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

package multicluster

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/emicklei/go-restful/v3"
	"github.com/golang/mock/gomock"
	"github.com/katanomi/pkg/testing/mock/github.com/katanomi/pkg/multicluster"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/authentication/user"
	apiserverrequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/client-go/rest"
	"knative.dev/pkg/system"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("TestCrossClusterSubjectReview_GetClient", func() {
	var (
		crossReview *CrossClusterSubjectReview
		req         *restful.Request
		clt         client.Client
		err         error
	)

	BeforeEach(func() {
		crossReview = &CrossClusterSubjectReview{
			ClusterParameter: "cluster",
			ClusterNamespace: "test",
			restMapper:       meta.NewDefaultRESTMapper([]schema.GroupVersion{}),
		}
		req = restful.NewRequest(&http.Request{
			URL: &url.URL{},
		})
		clt = nil
		err = nil
	})

	JustBeforeEach(func() {
		clt, err = crossReview.GetClient(context.Background(), req)
	})

	Context("cluster parameter is empty", func() {
		It("client should be nil", func() {
			Expect(clt).Should(BeNil())
			Expect(err).Should(BeNil())
		})
	})
	Context("cluster parameter is not empty", func() {
		BeforeEach(func() {
			request := httptest.NewRequest("GET", "http://localhost:8080?cluster=test", nil)
			req = restful.NewRequest(request)
		})

		When("when get multi cluster client failed", func() {
			BeforeEach(func() {
				mockCtl := gomock.NewController(GinkgoT())
				mockClient := multicluster.NewMockInterface(mockCtl)
				mockClient.EXPECT().GetConfig(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("test"))
				crossReview.multiClusterClient = mockClient
			})

			It("client should return error", func() {
				Expect(clt).Should(BeNil())
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("when get multi cluster client successfully", func() {
			BeforeEach(func() {
				mockCtl := gomock.NewController(GinkgoT())
				mockClient := multicluster.NewMockInterface(mockCtl)
				mockClient.EXPECT().GetConfig(gomock.Any(), gomock.Any()).Return(&rest.Config{
					Username: "test",
				}, nil)
				crossReview.multiClusterClient = mockClient
				ctx := apiserverrequest.WithUser(context.Background(), &user.DefaultInfo{})
				req.Request = req.Request.WithContext(ctx)
			})

			It("client should return a client successfully", func() {
				Expect(clt).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
			})
		})
	})
})

func TestCrossClusterSubjectReview_SetClusterParameter(t *testing.T) {
	g := NewGomegaWithT(t)
	os.Setenv(system.NamespaceEnvKey, "test")
	review := NewCrossClusterSubjectReview(nil, nil, nil)
	g.Expect(review.ClusterNamespace).Should(Equal("test"))
	g.Expect(review.ClusterParameter).Should(Equal("cluster"))

	review.SetClusterNamespace("new-ns")
	g.Expect(review.ClusterNamespace).Should(Equal("new-ns"))

	review.SetClusterParameter("new-param")
	g.Expect(review.ClusterParameter).Should(Equal("new-param"))
}
