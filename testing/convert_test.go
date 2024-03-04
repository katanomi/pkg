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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ = Describe("Test.ConvertTypeMetaToGroupVersionResource", func() {
	DescribeTable("Converts TypeMeta to GroupVersionResource",
		func(apiVersion, kind string, expectedGVR schema.GroupVersionResource) {
			typeMeta := metav1.TypeMeta{
				APIVersion: apiVersion,
				Kind:       kind,
			}
			Expect(ConvertTypeMetaToGroupVersionResource(typeMeta)).To(Equal(expectedGVR))
		},
		Entry("should handle normal resources",
			"batch/v1", "Job", schema.GroupVersionResource{
				Group:    "batch",
				Version:  "v1",
				Resource: "jobs",
			}),
		Entry("should handle namespaced API group",
			"tekton.dev/v1beta1", "Task", schema.GroupVersionResource{
				Group:    "tekton.dev",
				Version:  "v1beta1",
				Resource: "tasks",
			}),
		Entry("should handle core API group",
			"v1", "Pod", schema.GroupVersionResource{
				Group:    "",
				Version:  "v1",
				Resource: "pods",
			}),
		Entry("should handle special resources",
			"group/version", "Jenkins", schema.GroupVersionResource{
				Group:    "group",
				Version:  "version",
				Resource: "jenkinses",
			}),
		Entry("should handle special resources",
			"group/version", "Policy", schema.GroupVersionResource{
				Group:    "group",
				Version:  "version",
				Resource: "policies",
			}),
	)
})

var _ = Describe("Test.SliceToInterfaceSlice", func() {
	Context("with different slice types", func() {
		It("should convert empty slice", func() {
			input := []string{}
			expected := []interface{}{}
			result := SliceToInterfaceSlice[string](input)
			Expect(result).To(Equal(expected))
		})

		It("should convert slice of int", func() {
			input := []int{1, 2, 3}
			expected := []interface{}{1, 2, 3}
			result := SliceToInterfaceSlice[int](input)
			Expect(result).To(Equal(expected))
		})

		It("should convert slice of float", func() {
			input := []float64{1.1, 2.2, 3.3}
			expected := []interface{}{1.1, 2.2, 3.3}
			result := SliceToInterfaceSlice[float64](input)
			Expect(result).To(Equal(expected))
		})

		It("should convert slice of string", func() {
			input := []string{"a", "b", "c"}
			expected := []interface{}{"a", "b", "c"}
			result := SliceToInterfaceSlice[string](input)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Test.SliceToRuntimeOjbect", func() {
	Context("when given a slice of any type", func() {
		It("should convert the slice to runtime.Object", func() {
			// Test case 1: Empty slice
			Expect(SliceToRuntimeOjbect[int]([]int{})).To(BeEmpty())

			// Test case 2: All elements can be converted to runtime.Object
			Expect(SliceToRuntimeOjbect[interface{}]([]interface{}{
				&corev1.Pod{},
				&corev1.Pod{},
				&corev1.Pod{},
			})).To(Equal([]runtime.Object{
				&corev1.Pod{},
				&corev1.Pod{},
				&corev1.Pod{},
			}))

			// Test case 3: Some elements can be converted to runtime.Object
			Expect(SliceToRuntimeOjbect([]interface{}{
				&corev1.Pod{},
				false,
				&corev1.Pod{},
			})).To(Equal([]runtime.Object{
				&corev1.Pod{},
				&corev1.Pod{},
			}))

			// Test case 4: No element can be converted to runtime.Object
			Expect(SliceToRuntimeOjbect([]interface{}{
				42,
				false,
				"string",
			})).To(BeEmpty())
		})
	})
})
