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
	"testing"

	authnv1 "k8s.io/api/authentication/v1"

	"github.com/emicklei/go-restful/v3"
	mockfakeclient "github.com/katanomi/pkg/testing/mock/sigs.k8s.io/controller-runtime/pkg/client"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	authv1 "k8s.io/api/authorization/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/authentication/user"
	apiserverrequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/client-go/rest"
	"knative.dev/pkg/injection"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
		attr.Namespace = "default"
		attr.Name = "def"
		review := makeSelfSubjectAccessReview(attr)
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
		attr.Namespace = "default"
		attr.Name = "xyz"
		review := makeSelfSubjectAccessReview(attr)
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
		attr.Namespace = "default"
		attr.Name = "xyz"
		review := makeSubjectAccessReview(attr, user)
		review.GetObject().SetName("default")
		err := postSubjectAccessReview(ctx, client, review)

		g.Expect(err).ToNot(BeNil())
		g.Expect(errors.IsForbidden(err)).To(BeTrue())
	})

	t.Run("when is impersonated request", func(t *testing.T) {
		g := NewGomegaWithT(t)
		mockCtl := gomock.NewController(t)

		ctx := context.TODO()

		mockClient := mockfakeclient.NewMockClient(mockCtl)
		mockClient.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx interface{}, obj client.Object, opts ...client.CreateOption) error {
				subject := obj.(*authv1.SelfSubjectAccessReview)
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

func TestGetResourceAttributesFunc_GetResourceAttributes(t *testing.T) {
	g := NewGomegaWithT(t)
	getter := GetResourceAttributesFunc(func(ctx context.Context, req *restful.Request) (authv1.ResourceAttributes,
		error) {
		return authv1.ResourceAttributes{Name: "test"}, nil
	})
	got, _ := getter.GetResourceAttributes(context.Background(), &restful.Request{})
	g.Expect(got.Name).Should(BeEquivalentTo("test"))
}

func TestImpersonateUser(t *testing.T) {
	g := NewGomegaWithT(t)

	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	user := ImpersonateUser(req)
	g.Expect(user).Should(BeNil())

	req.Header[authnv1.ImpersonateUserHeader] = []string{"user-1"}
	req.Header[authnv1.ImpersonateUIDHeader] = []string{"user-id-1"}
	req.Header["Impersonate-Extra-Key1"] = []string{"value-1"}

	user = ImpersonateUser(req)
	g.Expect(user.GetName()).Should(BeEquivalentTo("user-1"))
	g.Expect(user.GetUID()).Should(BeEquivalentTo("user-id-1"))
	g.Expect(len(user.GetExtra())).Should(BeEquivalentTo(1))
	g.Expect(user.GetExtra()["Impersonate-Extra-Key1"][0]).Should(BeEquivalentTo("value-1"))

	req.Header.Del(authnv1.ImpersonateUserHeader)
	req.Header[authnv1.ImpersonateGroupHeader] = []string{"group-1"}
	req.Header[authnv1.ImpersonateUIDHeader] = []string{"group-id-1"}
	user = ImpersonateUser(req)
	g.Expect(user.GetName()).Should(BeEmpty())
	g.Expect(user.GetGroups()[0]).Should(BeEquivalentTo("group-1"))
	g.Expect(user.GetUID()).Should(BeEquivalentTo("group-id-1"))
	g.Expect(len(user.GetExtra())).Should(BeEquivalentTo(1))
	g.Expect(user.GetExtra()["Impersonate-Extra-Key1"][0]).Should(BeEquivalentTo("value-1"))

}
