/*
Copyright 2021 The AlaudaDevops Authors.

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
	"net/http/httptest"

	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	authnv1 "k8s.io/api/authentication/v1"

	"knative.dev/pkg/injection"

	"k8s.io/client-go/rest"

	"github.com/emicklei/go-restful/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ImpersonateFilter", func() {

	var (
		ctx                    context.Context
		config                 *rest.Config
		fakeClientBeforeFilter ctrlclient.Client
		fakeClientAfterFilter  ctrlclient.Client

		req  *restful.Request
		resp *restful.Response

		chain *restful.FilterChain
	)

	BeforeEach(func() {
		ctx = context.TODO()
		config = &rest.Config{}
		fakeClientBeforeFilter = fake.NewClientBuilder().Build()

		ctx = injection.WithConfig(ctx, config)
		ctx = WithClient(ctx, fakeClientBeforeFilter)
		//ctx = apiserverrequest.WithUser(ctx, &user.DefaultInfo{Name: "dev"})

		request := httptest.NewRequest("GET", "http://localhost", nil)
		req = restful.NewRequest(request)
		req.Request = req.Request.WithContext(ctx)

		response := httptest.NewRecorder()
		resp = restful.NewResponse(response)
	})

	JustBeforeEach(func() {
		filter := ImpersonateFilter(ctx)
		filter(req, resp, chain)
	})

	When("without impersonate in request header", func() {
		BeforeEach(func() {

			chain = &restful.FilterChain{
				Target: func(request *restful.Request, response *restful.Response) {
					response.WriteHeader(204)
					return
				},
			}
		})
		It("should do nothing in filter", func() {
			Expect(resp.StatusCode()).Should(BeEquivalentTo(204))
			Expect(fakeClientBeforeFilter).Should(BeEquivalentTo(fakeClientBeforeFilter))
		})
	})

	When("with impersonate in request header", func() {
		BeforeEach(func() {
			req.Request.Header.Set(authnv1.ImpersonateUserHeader, "dev")
			chain = &restful.FilterChain{
				Target: func(request *restful.Request, response *restful.Response) {
					fakeClientAfterFilter = Client(request.Request.Context())
					response.WriteHeader(200)
					return
				},
			}
		})
		It("should inject impersonate user and overwrite client", func() {
			Expect(resp.StatusCode()).Should(BeEquivalentTo(200))
			u := User(req.Request.Context())
			Expect(u.GetName()).Should(BeEquivalentTo("dev"))
			Expect(fakeClientAfterFilter).ShouldNot(BeEquivalentTo(fakeClientBeforeFilter))
		})
	})

})
