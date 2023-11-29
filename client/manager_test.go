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
	"github.com/katanomi/pkg/multicluster"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/rest"
	"knative.dev/pkg/injection"
)

func EmptyHandler(rq *restful.Request, rp *restful.Response) {
}

func TestManagerFilter(t *testing.T) {
	// TODO: Find a better and more reliable way to do these
	t.Skip()
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
		}, func(config *rest.Config) (multicluster.Interface, error) { return nil, nil })
		ws := new(restful.WebService)
		ws.Consumes(restful.MIME_JSON)
		ws.Route(ws.GET("/config").Filter(ManagerFilter(ctx, mgr)).To(EmptyHandler))
		restful.Add(ws)
		testReq := httptest.NewRequest(http.MethodGet, "/config", nil)
		req := restful.NewRequest(testReq)
		req.Request.Header.Set("Authorization", "Bearer 0123456789")
		req.Request = req.Request.WithContext(ctx)
		resp := restful.NewResponse(httptest.NewRecorder())

		restful.DefaultContainer.ServeHTTP(resp, testReq)
		// ManagerFilter(ctx, mgr)(req, resp, chain)

		config := injection.GetConfig(testReq.Context())
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
		}, nil)
		req := restful.NewRequest(httptest.NewRequest(http.MethodGet, "http://example.com", nil))
		req.Request = req.Request.WithContext(ctx)
		resp := restful.NewResponse(httptest.NewRecorder())
		ManagerFilter(ctx, mgr)(req, resp, chain)
		g.Expect(resp.StatusCode()).To(Equal(http.StatusUnauthorized))
		config := injection.GetConfig(req.Request.Context())
		g.Expect(config).To(BeNil())
	})
}

func TestUserFromBearerToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJlbWFpbCI6ImRldiJ9.v5leOJQ8mxkOzWW-dWWFfPGPn__0eYUGtDCdwx1LWkM"
	info, err := UserFromBearerToken(token)
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if info.GetName() != "dev" {
		t.Error("username should be dev")
	}
}

func TestWithCtxManagerFilters(t *testing.T) {
	g := NewGomegaWithT(t)

	ws := restful.WebService{}
	ctx := context.TODO()
	ctx = WithManager(ctx, &Manager{})

	err := WithCtxManagerFilters(ctx, &ws)
	g.Expect(err).To(Succeed(), "should return nil")
}
