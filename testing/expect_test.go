/*
Copyright 2024 The Katanomi Authors.

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

package testing

import (
	"time"

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

	When("asserting on nil", func() {
		It("should do the right thing", func() {
			Expect("foo").ShouldNot(DiffEqualTo(nil))
			Expect(nil).ShouldNot(DiffEqualTo(3))
			Expect([]int{1, 2}).ShouldNot(DiffEqualTo(nil))
		})
	})

	When("asserting on object", func() {
		It("should do the right thing", func() {
			Expect("foo").ShouldNot(DiffEqualTo(nil))
			Expect(nil).ShouldNot(DiffEqualTo(3))
			Expect([]int{1, 2}).ShouldNot(DiffEqualTo(nil))
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

			Expect(pod).ShouldNot(DiffEqualTo(&corev1.ConfigMap{
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

	When("asserting on kubernetes object with ignoreFuncs", func() {
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

			Expect(pod).Should(DiffEqualTo(&corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-1",
					Namespace: "foo",
				},
				Data: map[string]string{
					"a": "1",
				},
			}, KubeObjectDiffClean, func(object interface{}) interface{} {
				object.(*corev1.ConfigMap).ObjectMeta = metav1.ObjectMeta{}
				return object
			}))
		})
	})
})

var _ = Describe("ExpectDiff", func() {

	When("asserting on kubernetes object with no ignore funcs", func() {
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

			actual := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "default",
					Namespace: "deault",
				},
				Data: map[string]string{
					"a": "1",
				},
			}

			ExpectDiff(pod, actual).Should(BeEmpty())
		})
	})

	When("asserting on kubernetes object with ignoreFuncs", func() {
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

			expected := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-1",
					Namespace: "foo",
				},
				Data: map[string]string{
					"a": "1",
				},
			}

			ExpectDiff(pod, expected, KubeObjectDiffClean, func(object interface{}) interface{} {
				object.(*corev1.ConfigMap).ObjectMeta = metav1.ObjectMeta{}
				return object
			}).Should(BeEmpty())

		})
	})
})
