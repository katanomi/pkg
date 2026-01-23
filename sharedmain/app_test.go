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

package sharedmain

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/AlaudaDevops/pkg/fieldindexer"
	"github.com/emicklei/go-restful/v3"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("AppBuilder", func() {
	var (
		appBuilder *AppBuilder
	)

	BeforeEach(func() {
		appBuilder = &AppBuilder{}
	})

	Context("BuiltInFilters", func() {
		It("should return filter that disables pretty print", func() {

			filters := appBuilder.BuiltInFilters()
			Expect(filters).To(HaveLen(1))

			req := restful.NewRequest(httptest.NewRequest(http.MethodGet, "/", nil))
			recorder := httptest.NewRecorder()
			resp := restful.NewResponse(recorder)
			chain := &restful.FilterChain{
				Filters: []restful.FilterFunction{},
				Target: func(req *restful.Request, resp *restful.Response) {
					resp.ResponseWriter.Write([]byte("test"))
				},
			}

			filters[0](req, resp, chain)

			respValue := reflect.ValueOf(resp).Elem()
			prettyPrintField := respValue.FieldByName("prettyPrint")

			Expect(prettyPrintField.Bool()).To(BeFalse())
			Expect(recorder.Body.String()).To(Equal("test"))
		})
	})
})

func TestAppWithFieldIndexer(t *testing.T) {
	t.Run("append one field indexer ", func(t *testing.T) {
		g := NewGomegaWithT(t)
		a := App("test").WithFieldIndexer(fieldindexer.FieldIndexer{
			Obj:   &corev1.ConfigMap{},
			Field: "data.key",
			ExtractValue: func(object client.Object) []string {
				return []string{object.(*corev1.ConfigMap).Data["key"]}
			},
		})
		g.Expect(a.fieldIndexeres).Should(HaveLen(1))
	})
	t.Run("append more than one field indexer", func(t *testing.T) {
		g := NewGomegaWithT(t)

		a := App("test").WithFieldIndexer(fieldindexer.FieldIndexer{
			Obj:   &corev1.ConfigMap{},
			Field: "data.key",
			ExtractValue: func(object client.Object) []string {
				return []string{object.(*corev1.ConfigMap).Data["key"]}
			},
		}).WithFieldIndexer(fieldindexer.FieldIndexer{
			Obj:   &corev1.ConfigMap{},
			Field: "data.name",
			ExtractValue: func(object client.Object) []string {
				return []string{object.(*corev1.ConfigMap).Data["name"]}
			},
		})
		g.Expect(a.fieldIndexeres).Should(HaveLen(2))
	})
}

func TestAppContextOpts(t *testing.T) {
	t.Run("empty function list should return directly", func(t *testing.T) {
		g := NewGomegaWithT(t)
		a := &AppBuilder{}

		result := a.ContextOpts()

		g.Expect(result).Should(BeIdenticalTo(a))
		g.Expect(a.Context).Should(BeNil())
	})

	t.Run("single context modification function", func(t *testing.T) {
		g := NewGomegaWithT(t)
		a := &AppBuilder{
			Context: context.Background(),
		}

		type testKey struct{}
		testValue := "test-value"

		result := a.ContextOpts(func(ctx context.Context) context.Context {
			return context.WithValue(ctx, testKey{}, testValue)
		})

		g.Expect(result).Should(BeIdenticalTo(a))
		g.Expect(a.Context.Value(testKey{})).Should(Equal(testValue))
	})

	t.Run("multiple context modification functions applied in order", func(t *testing.T) {
		g := NewGomegaWithT(t)
		a := &AppBuilder{
			Context: context.Background(),
		}

		type key1 struct{}
		type key2 struct{}
		type key3 struct{}

		result := a.ContextOpts(
			func(ctx context.Context) context.Context {
				return context.WithValue(ctx, key1{}, "value1")
			},
			func(ctx context.Context) context.Context {
				return context.WithValue(ctx, key2{}, "value2")
			},
			func(ctx context.Context) context.Context {
				return context.WithValue(ctx, key3{}, "value3")
			},
		)

		g.Expect(result).Should(BeIdenticalTo(a))
		g.Expect(a.Context.Value(key1{})).Should(Equal("value1"))
		g.Expect(a.Context.Value(key2{})).Should(Equal("value2"))
		g.Expect(a.Context.Value(key3{})).Should(Equal("value3"))
	})
}
