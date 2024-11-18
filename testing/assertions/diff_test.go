/*
Copyright 2024 The AlaudaDevops Authors.

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

package assertions

import (
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("DiffEqual", func() {
	When("asserting that nil is equivalent to nil", func() {
		It("should error", func() {
			success, err := (&DiffEqualMatcher{Expected: nil}).Match(nil)

			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
		})
	})

	When("asserting using clean function", func() {
		It("should use diff clean function", func() {
			data := map[string]string{"a": "b"}
			matcher := &DiffEqualMatcher{Expected: data, DiffCleanFunc: []func(obj interface{}) interface{}{
				func(obj interface{}) interface{} {
					if dict, ok := obj.(map[string]string); ok {
						dict["a"] = "c"
						return dict
					}
					return obj
				},
			}}
			success, err := matcher.Match(data)

			Expect(success).Should(BeTrue())
			Expect(err).ShouldNot(HaveOccurred())

			Expect(matcher.FailureMessage(data)).NotTo(BeEmpty())
			Expect(matcher.NegatedFailureMessage(data)).NotTo(BeEmpty())
		})
	})

	When("asserting basic different types", func() {
		It("should return non empty value", func() {
			Expect("foo").ShouldNot(DiffEqual(nil))
			Expect("foo").To(DiffEqual("foo"))
			Expect("foo").ToNot(DiffEqual("bar"))
			Expect(nil).ShouldNot(DiffEqual(3))
			Expect([]int{1, 2}).ShouldNot(DiffEqual(nil))
			Expect([]int{1, 2}).To(DiffEqual([]int{1, 2}))
			Expect([]int{1, 2}).ToNot(DiffEqual([]int{2, 1}))
			Expect([]int{1, 2}).ToNot(DiffEqual(map[string]string{"a": "b"}))
			Expect(map[string]string{"a": "b"}).To(DiffEqual(map[string]string{"a": "b"}))
		})
	})

	When("asserting on basic struct", func() {
		type BasicType struct {
			A string
		}
		It("should do the right thing", func() {
			Expect(BasicType{A: "a"}).To(DiffEqual(BasicType{A: "a"}))
			Expect(BasicType{A: "a"}).ToNot(DiffEqual(BasicType{A: "b"}))
			Expect(BasicType{A: "a"}).ToNot(DiffEqual(BasicType{}))
		})
	})

	When("asserting on kubernetes object", func() {
		It("should do the right thing", func() {
			pod := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					ResourceVersion:   "11",
					Generation:        1,
					CreationTimestamp: metav1.NewTime(time.Now()),
					Name:              "default",
					Namespace:         "deault",
				},
				Data: map[string]string{
					"a": "1",
				},
			}

			another := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "default",
					Namespace: "deault",
				},
				Data: map[string]string{
					"a": "12",
				},
			}

			Expect(pod).ShouldNot(DiffEqual(another))

			another.Data["a"] = "1"
			Expect(pod).ShouldNot(DiffEqual(another), "data is the same, but object meta is different")

			Expect(pod).To(DiffEqual(another, IgnoreObjectMetaFields()), "Data is the same, object meta is ignored")
		})

		It("should do the right thing", func() {
			pod := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					ResourceVersion:   "11",
					Generation:        1,
					CreationTimestamp: metav1.NewTime(time.Now()),
					Name:              "default",
					Namespace:         "deault",
				},
				Data: map[string]string{
					"a": "1",
				},
			}

			Expect(pod).ShouldNot(DiffEqual(&corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "default",
					Namespace: "deault",
				},
				Data: map[string]string{
					"a": "12",
				},
			}))
		})
	})

	When("asserting on kubernetes object with cmp.Options", func() {
		It("should not have difference", func() {
			pod := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					ResourceVersion:   "11",
					Generation:        1,
					CreationTimestamp: metav1.NewTime(time.Now()),
					Name:              "default",
					Namespace:         "deault",
				},
				Data: map[string]string{
					"a": "1",
				},
			}

			Expect(pod).Should(DiffEqual(&corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-1",
					Namespace: "foo",
				},
				Data: map[string]string{
					"a": "1",
				},
			}, cmpopts.IgnoreTypes(metav1.ObjectMeta{})))
		})
	})
})
