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

package route

import (
	"context"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/emicklei/go-restful/v3"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHandleSortQuery(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "test for get sort params")
}

var _ = Describe("test for get sort params", func() {
	var res []metav1alpha1.SortOptions
	var sortBy = ""
	var ctx = context.Background()

	Describe("test for get sort params", func() {

		JustBeforeEach(func() {
			res = HandleSortQuery(ctx, sortBy)
		})

		Context("format string", func() {
			BeforeEach(func() {
				sortBy = "asc,a,desc,b"
			})

			It("should succeed", func() {
				Expect(len(res)).To(Equal(2))
				Expect(res[0].Order).To(Equal(metav1alpha1.OrderAsc))
				Expect(res[0].SortBy).To(Equal(metav1alpha1.SortBy("a")))
				Expect(res[1].Order).To(Equal(metav1alpha1.OrderDesc))
				Expect(res[1].SortBy).To(Equal(metav1alpha1.SortBy("b")))
			})
		})

		Context("singular params string", func() {
			BeforeEach(func() {
				sortBy = "asc,a,desc"
			})

			It("should succeed", func() {
				Expect(len(res)).To(Equal(0))
			})
		})

		Context("empty string", func() {
			BeforeEach(func() {
				sortBy = ""
			})

			It("should succeed", func() {
				Expect(len(res)).To(Equal(0))
			})
		})

		Context("reversed string", func() {
			BeforeEach(func() {
				sortBy = "a,asc"
			})

			It("should succeed", func() {
				Expect(len(res)).To(Equal(0))
			})
		})

		Context("error separator", func() {
			BeforeEach(func() {
				sortBy = "asc,a;desc,b"
			})

			It("should succeed", func() {
				Expect(len(res)).To(Equal(0))
			})
		})
	})
})

var _ = Describe("TestParsePagerFromRequest", func() {
	var (
		req   *restful.Request
		pager metav1alpha1.Pager
	)

	BeforeEach(func() {
		req = &restful.Request{
			Request: &http.Request{
				URL: &url.URL{},
			},
		}
		pager = metav1alpha1.Pager{}
	})

	JustBeforeEach(func() {
		pager = ParsePagerFromRequest(req)
	})

	Context("no parameters were specified in the request", func() {
		It("should use default value", func() {
			Expect(pager.Page).To(Equal(1))
			Expect(pager.ItemsPerPage).To(Equal(20))
		})
	})

	Context("parameters were specified in the request", func() {
		BeforeEach(func() {
			req.Request.URL = &url.URL{
				RawQuery: "page=2&itemsPerPage=30",
			}
		})

		It("should use the value provided in the request", func() {
			Expect(pager.Page).To(Equal(2))
			Expect(pager.ItemsPerPage).To(Equal(30))
		})
	})
})

var _ = Describe("ParseTimeCursorFormRequest", func() {
	var (
		req       *restful.Request
		cursor    *metav1alpha1.TimeCursor
		startTime time.Time
	)

	BeforeEach(func() {
		req = &restful.Request{
			Request: &http.Request{
				URL: &url.URL{},
			},
		}
		cursor = &metav1alpha1.TimeCursor{}
		startTime = time.Now()
	})

	JustBeforeEach(func() {
		cursor = ParseTimeCursorFormRequest(req)
	})

	Context("no parameters were specified in the request", func() {
		It("should use default value", func() {
			Expect(cursor.QueryStartAt <= time.Now().Unix() && cursor.QueryStartAt >= startTime.Unix()).
				To(BeTrue())
			Expect(cursor.Page).To(Equal(1))
			Expect(cursor.ItemsPerPage).To(Equal(20))
		})
	})

	Context("continue param were specified in the request", func() {
		BeforeEach(func() {
			req.Request.URL = &url.URL{
				RawQuery: "continue=eyJpdGVtc1BlclBhZ2UiOjMwLCJwYWdlIjoyLCJxdWVyeVN0YXJ0QXQiOjE2NzU2NTAwODJ9",
			}
		})

		It("should use the value provided in the request", func() {
			Expect(cursor.QueryStartAt).To(Equal(int64(1675650082)))
			Expect(cursor.Page).To(Equal(2))
			Expect(cursor.ItemsPerPage).To(Equal(30))
		})
	})
})

func TestParsePagerFromRequest(t *testing.T) {
	type args struct {
		req *restful.Request
	}
	tests := []struct {
		name string
		args args
		want metav1alpha1.Pager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParsePagerFromRequest(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePagerFromRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
