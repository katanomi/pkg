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

package fake

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/emicklei/go-restful/v3"
	"github.com/golang/mock/gomock"
	"github.com/katanomi/pkg/testing/fake/opa"
	"github.com/katanomi/pkg/testing/mock/testing/fake"
	. "github.com/onsi/gomega"
)

func TestPolicyHandlerFilter(t *testing.T) {
	tests := map[string]struct {
		method    string
		query     string
		body      string
		expected  interface{}
		getPolicy func(builder *Builder) *opa.Policy
	}{
		"empty": {
			method: "GET",
			query:  "?abc=def",
			body:   "{}",
			getPolicy: func(builder *Builder) *opa.Policy {
				return nil
			},
			expected: map[string]string{"result": "simple"},
		},
		"query": {
			method: "GET",
			query:  "?abc=def",
			body:   "{}",
			getPolicy: func(builder *Builder) *opa.Policy {
				builder.When(InputQuery.Field("abc"), "def").Result(map[string]string{"a": "b"})
				p, _ := builder.Complete()
				return p
			},
			expected: map[string]string{"a": "b"},
		},
		"body": {
			method: "POST",
			body:   `{"qwer": "1234"}`,
			getPolicy: func(builder *Builder) *opa.Policy {
				builder.When(InputBody.Field("qwer"), "1234").Result(map[string]string{"c": "d"})
				p, _ := builder.Complete()
				return p
			},
			expected: map[string]string{"c": "d"},
		},
		"not-matched": {
			method: "POST",
			body:   `{"qwer": "1234"}`,
			getPolicy: func(builder *Builder) *opa.Policy {
				builder.When(InputBody.Field("qwer"), "2345").Result(map[string]string{"c": "d"})
				p, _ := builder.Complete()
				return p
			},
			expected: map[string]string{"result": "simple"},
		},
	}

	for name, item := range tests {
		t.Run(name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			bodyReader := strings.NewReader(item.body)
			url := fmt.Sprintf("http://example.com/%s%s", name, item.query)
			httpRequest, _ := http.NewRequest(item.method, url, bodyReader)
			httpRequest.Header.Set("Accept", "*/*")
			httpRequest.Header.Set("Content-Type", "application/json")
			httpWriter := httptest.NewRecorder()

			ctrl := gomock.NewController(t)
			store := fake.NewMockStore(ctrl)

			builder := NewPolicyBuilder(item.method, name)
			policy := item.getPolicy(builder)
			id := IDFromMethodPath(item.method, name)
			store.EXPECT().Get(gomock.Any(), id).Return(policy, nil)

			h := PolicyHandler{Store: store}
			ws := new(restful.WebService).Filter(h.Filter)
			ws.Route(ws.Produces(restful.MIME_JSON).POST(name).To(simple(g, item.body)))
			ws.Route(ws.Produces(restful.MIME_JSON).GET(name).To(simple(g, item.body)))
			c := restful.NewContainer().Add(ws)
			c.Dispatch(httpWriter, httpRequest)

			var result map[string]string
			_ = json.Unmarshal([]byte(httpWriter.Body.String()), &result)

			g.Expect(result).To(Equal(item.expected))
		})
	}
}
func simple(g Gomega, body string) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		bytes, err := io.ReadAll(req.Request.Body)
		g.Expect(err).To(BeNil())
		g.Expect(string(bytes)).To(Equal(body))

		_, _ = io.WriteString(resp.ResponseWriter, `{"result":"simple"}`)
	}
}
