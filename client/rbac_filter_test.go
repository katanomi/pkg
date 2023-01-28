/*
Copyright 2021 The Katanomi Authors.

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

package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	apiserverrequest "k8s.io/apiserver/pkg/endpoints/request"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"knative.dev/pkg/injection"

	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/client-go/rest"

	"testing"

	"github.com/emicklei/go-restful/v3"
	gomock "github.com/golang/mock/gomock"
	mockfakeclient "github.com/katanomi/pkg/testing/mock/sigs.k8s.io/controller-runtime/pkg/client"
	. "github.com/onsi/gomega"
	authv1 "k8s.io/api/authorization/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestRBACFilter(t *testing.T) {

	scheme := runtime.NewScheme()
	authv1.AddToScheme(scheme)

	attr := authv1.ResourceAttributes{
		Namespace: "default",
		Verb:      "update",
		Group:     "meta.katanomi.dev",
		Version:   "v1alpha1",
		Resource:  "artifacts",
		Name:      "abc",
	}
	t.Run("no client in request", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctx := context.TODO()

		client := Client(ctx)
		review := makeSelfSubjectAccessReview("default", "def", attr)
		review.GetObject().SetName("default")
		err := postSubjectAccessReview(ctx, client, review)

		g.Expect(err).ToNot(BeNil())
		g.Expect(errors.IsUnauthorized(err)).To(BeTrue())

	})
	t.Run("adding fake client in ctx", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		clt := fake.NewClientBuilder().WithScheme(scheme).Build()
		ctx = WithClient(ctx, clt)

		client := Client(ctx)
		review := makeSelfSubjectAccessReview("default", "xyz", attr)
		review.GetObject().SetName("default")
		err := postSubjectAccessReview(ctx, client, review)

		g.Expect(err).ToNot(BeNil())
		fmt.Println(err.Error())
		g.Expect(errors.IsForbidden(err)).To(BeTrue())
	})
	t.Run("adding fake client in ctx", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		clt := fake.NewClientBuilder().WithScheme(scheme).Build()
		ctx = WithClient(ctx, clt)

		client := Client(ctx)
		user := &user.DefaultInfo{
			Name:   "system:serviceaccount:default:katanomi",
			Groups: []string{"system:authenticated"},
			UID:    "39f88a8e-9090-4495-830c-fabf0d0cc7a3",
		}
		review := makeSubjectAccessReview("default", "xyz", attr, user)
		review.GetObject().SetName("default")
		err := postSubjectAccessReview(ctx, client, review)

		g.Expect(err).ToNot(BeNil())
		g.Expect(errors.IsForbidden(err)).To(BeTrue())
	})

	t.Run("when is impersonate request", func(t *testing.T) {
		g := NewGomegaWithT(t)
		mockCtl := gomock.NewController(t)

		ctx := context.TODO()

		mockClient := mockfakeclient.NewMockClient(mockCtl)
		mockClient.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx interface{}, obj client.Object, opts ...client.CreateOption) error {
				subject := obj.(*authv1.SubjectAccessReview)
				subject.Status.Denied = true
				subject.Status.Allowed = false
				return nil
			}).Times(1)

		ctx = WithClient(ctx, mockClient)
		config := &rest.Config{Impersonate: rest.ImpersonationConfig{UserName: "dev"}}
		ctx = injection.WithConfig(ctx, config)
		ctx = apiserverrequest.WithUser(ctx, &user.DefaultInfo{Name: "dev"})

		filter := SubjectReviewFilterForResource(ctx, attr, "namespace", "name")
		req := httptest.NewRequest("GET", "http://localhost", nil)
		request := restful.NewRequest(req)
		request.Request = request.Request.WithContext(ctx)

		recorder := httptest.NewRecorder()
		response := restful.NewResponse(recorder)

		chain := &restful.FilterChain{Target: EmptyHandler}
		filter(request, response, chain)

		g.Expect(recorder.Code).Should(Not(BeEquivalentTo(http.StatusOK)))
	})
}
