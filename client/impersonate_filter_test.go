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
	"net/http/httptest"
	"testing"

	authnv1 "k8s.io/api/authentication/v1"

	"knative.dev/pkg/injection"

	"k8s.io/client-go/rest"

	"github.com/emicklei/go-restful/v3"
	. "github.com/onsi/gomega"
	"k8s.io/apiserver/pkg/authentication/user"
	apiserverrequest "k8s.io/apiserver/pkg/endpoints/request"
)

func TestImpersonateFilter(t *testing.T) {

	t.Run("with impersonate in rest.Config", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctx := context.TODO()

		config := &rest.Config{}
		ctx = injection.WithConfig(ctx, config)
		ctx = apiserverrequest.WithUser(ctx, &user.DefaultInfo{Name: "dev"})

		request := httptest.NewRequest("GET", "http://localhost", nil)
		request.Header.Set(authnv1.ImpersonateUserHeader, "dev")
		req := restful.NewRequest(request)
		req.Request = req.Request.WithContext(ctx)

		chain := &restful.FilterChain{
			Target: func(request *restful.Request, response *restful.Response) {
				return
			},
		}

		response := httptest.NewRecorder()
		resp := restful.NewResponse(response)

		filter := ImpersonateFilter(ctx)

		filter(req, resp, chain)

		u := User(req.Request.Context())
		g.Expect(u.GetName()).Should(BeEquivalentTo("dev"))
	})
}
