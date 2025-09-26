/*
Copyright 2025 The Katanomi Authors.

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

package tekton

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/katanomi/pkg/testing"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

var _ = Describe("CleanupObjectParamByProperties", func() {
	type testCase struct {
		description string
		paramFile   string
		goldenFile  string
	}

	DescribeTable("object parameter cleanup scenarios",
		func(tc testCase) {
			// Load test parameter from file
			var param v1beta1.ParamSpec
			Expect(LoadYAML(tc.paramFile, &param)).To(Succeed())

			// Load expected result
			var expected v1beta1.ParamSpec
			Expect(LoadYAML(tc.goldenFile, &expected)).To(Succeed())

			// Execute the function
			CleanupObjectParamByProperties(&param)

			// Verify the result
			Expect(param).To(Equal(expected))
		},

		Entry("non-object parameter should remain unchanged", testCase{
			description: "should not modify string parameters",
			paramFile:   "testdata/cleanup_object_param/string_param.yaml",
			goldenFile:  "testdata/cleanup_object_param/string_param.yaml",
		}),

		Entry("object parameter with no default should remain unchanged", testCase{
			description: "should not modify object parameter without default value",
			paramFile:   "testdata/cleanup_object_param/object_no_default.yaml",
			goldenFile:  "testdata/cleanup_object_param/object_no_default.yaml",
		}),

		Entry("object parameter with empty properties should clear non-empty default", testCase{
			description: "should clear ObjectVal when Properties is empty and ObjectVal is non-empty",
			paramFile:   "testdata/cleanup_object_param/empty_properties_with_default.yaml",
			goldenFile:  "testdata/cleanup_object_param/empty_properties_cleared.yaml",
		}),

		Entry("object parameter with empty properties and empty default should remain unchanged", testCase{
			description: "should preserve empty ObjectVal when Properties is empty",
			paramFile:   "testdata/cleanup_object_param/empty_properties_empty_default.yaml",
			goldenFile:  "testdata/cleanup_object_param/empty_properties_empty_default.yaml",
		}),

		Entry("object parameter with defined properties should filter default values", testCase{
			description: "should only retain properties that exist in Properties schema",
			paramFile:   "testdata/cleanup_object_param/defined_properties_with_extra_defaults.yaml",
			goldenFile:  "testdata/cleanup_object_param/defined_properties_filtered.yaml",
		}),

		Entry("object parameter with defined properties and no matching defaults", testCase{
			description: "should result in empty ObjectVal when no defaults match properties",
			paramFile:   "testdata/cleanup_object_param/defined_properties_no_matching_defaults.yaml",
			goldenFile:  "testdata/cleanup_object_param/defined_properties_empty_result.yaml",
		}),

		Entry("object parameter with partial matching properties", testCase{
			description: "should retain only properties that match schema definition",
			paramFile:   "testdata/cleanup_object_param/partial_matching_properties.yaml",
			goldenFile:  "testdata/cleanup_object_param/partial_matching_filtered.yaml",
		}),
	)

	Context("edge cases", func() {
		It("should handle nil parameter without panic", func() {
			Expect(func() {
				CleanupObjectParamByProperties(nil)
			}).ToNot(Panic())
		})

		It("should handle array parameter correctly", func() {
			var param v1beta1.ParamSpec
			Expect(LoadYAML("testdata/cleanup_object_param/array_param.yaml", &param)).To(Succeed())

			original := param.DeepCopy()
			CleanupObjectParamByProperties(&param)

			Expect(param).To(Equal(*original))
		})
	})
})
