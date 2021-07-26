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
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/emicklei/go-restful/v3"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/rest"
	"knative.dev/pkg/injection"
)

func TestManagerFilter(t *testing.T) {
	os.Setenv("KUBERNETES_MASTER", "127.0.0.1:16003")
	target := func(req *restful.Request, resp *restful.Response) {}
	chain := &restful.FilterChain{Target: target}

	t.Run("should succeed", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		mgr := NewManager(ctx, FromBearerToken, func() (*rest.Config, error) {
			return &rest.Config{
				Host:     "https://127.0.0.1:6443",
				Username: "abc",
				Password: "def",
			}, nil
		})
		req := restful.NewRequest(httptest.NewRequest(http.MethodGet, "http://example.com", nil))
		req.Request.Header.Set("Authorization", "Bearer 0123456789")
		req.Request = req.Request.WithContext(ctx)
		resp := restful.NewResponse(httptest.NewRecorder())

		ManagerFilter(mgr)(req, resp, chain)

		config := injection.GetConfig(req.Request.Context())
		g.Expect(config).ToNot(BeNil())
		g.Expect(resp.StatusCode()).ToNot(Equal(http.StatusInternalServerError))
		g.Expect(config.BearerToken).To(Equal("0123456789"))
	})

	t.Run("should return error", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		mgr := NewManager(ctx, FromBearerToken, func() (*rest.Config, error) {
			// will return an error
			return &rest.Config{}, nil
		})
		req := restful.NewRequest(httptest.NewRequest(http.MethodGet, "http://example.com", nil))
		req.Request = req.Request.WithContext(ctx)
		resp := restful.NewResponse(httptest.NewRecorder())
		ManagerFilter(mgr)(req, resp, chain)
		g.Expect(resp.StatusCode()).To(Equal(http.StatusInternalServerError))
		config := injection.GetConfig(req.Request.Context())
		g.Expect(config).To(BeNil())
	})

}
