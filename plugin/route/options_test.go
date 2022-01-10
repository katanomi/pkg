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
	"testing"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHandleSortQuery(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t,
		"test for get sort params",
		[]Reporter{printer.NewlineReporter{}})
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
