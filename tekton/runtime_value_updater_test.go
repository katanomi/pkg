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
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/katanomi/pkg/testing"
	pipelinev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

var _ = Describe("Test.MergeDefaultWithRuntimeValue", func() {
	logger := GetDefaultLogger()
	type testCase struct {
		description  string
		originalFile string
		runtimeFile  string
		expectedFile string
	}

	loadParamSpec := func(path string) *pipelinev1beta1.ParamSpec {
		if path == "" {
			return nil
		}

		paramSpec := &pipelinev1beta1.ParamSpec{}
		MustLoadYaml(path, &paramSpec)
		return paramSpec
	}

	loadParamValue := func(path string) *pipelinev1beta1.ParamValue {
		if path == "" {
			return nil
		}

		paramValue := &pipelinev1beta1.ParamValue{}
		MustLoadYaml(path, &paramValue)
		return paramValue
	}

	DescribeTable("should update default value based on runtime value",
		func(tc testCase) {
			original := loadParamSpec(tc.originalFile)
			runtimeValue := loadParamValue(tc.runtimeFile)
			expected := loadParamSpec(tc.expectedFile)
			logger.Debugw("test case", "description", tc.description, "original", original, "runtimeValue", runtimeValue, "expected", expected)

			MergeDefaultWithRuntimeValue(original, runtimeValue)

			if expected == nil {
				Expect(original).To(BeNil())
			} else {
				Expect(original).ToNot(BeNil())
				diff := cmp.Diff(expected, original)
				if diff != "" {
					logger.Errorw("result does not match expected", "original", original, "expected", expected)
					GinkgoT().Errorf("result does not match expected:\n%s", diff)
				}
				Expect(diff).To(BeEmpty())
			}
		},

		Entry("nil original should remain unchanged", testCase{
			description:  "should handle nil original without modification",
			originalFile: "",
			runtimeFile:  "testdata/update_default_with_runtime_value/nil_handling/runtime.yaml",
			expectedFile: "",
		}),

		Entry("nil runtimeValue should keep original unchanged", testCase{
			description:  "should preserve original when runtimeValue is nil",
			originalFile: "testdata/update_default_with_runtime_value/nil_handling/original.yaml",
			runtimeFile:  "",
			expectedFile: "testdata/update_default_with_runtime_value/nil_handling/original.yaml",
		}),

		Entry("both nil should handle gracefully", testCase{
			description:  "should handle both parameters being nil gracefully",
			originalFile: "",
			runtimeFile:  "",
			expectedFile: "",
		}),

		Entry("type mismatch should keep original unchanged", testCase{
			description:  "should not modify original when types don't match",
			originalFile: "testdata/update_default_with_runtime_value/type_mismatch/original.yaml",
			runtimeFile:  "testdata/update_default_with_runtime_value/type_mismatch/runtime.yaml",
			expectedFile: "testdata/update_default_with_runtime_value/type_mismatch/expected.yaml",
		}),

		Entry("string type should be replaced with runtime value", testCase{
			description:  "should replace string original with runtimeValue",
			originalFile: "testdata/update_default_with_runtime_value/string_update/original.yaml",
			runtimeFile:  "testdata/update_default_with_runtime_value/string_update/runtime.yaml",
			expectedFile: "testdata/update_default_with_runtime_value/string_update/expected.yaml",
		}),

		Entry("array type should be replaced with runtime value", testCase{
			description:  "should replace array original with runtimeValue",
			originalFile: "testdata/update_default_with_runtime_value/array_update/original.yaml",
			runtimeFile:  "testdata/update_default_with_runtime_value/array_update/runtime.yaml",
			expectedFile: "testdata/update_default_with_runtime_value/array_update/expected.yaml",
		}),

		Entry("object type with empty original should merge runtime values", testCase{
			description:  "should merge runtimeValue into empty original object",
			originalFile: "testdata/update_default_with_runtime_value/object_merge/empty_original.yaml",
			runtimeFile:  "testdata/update_default_with_runtime_value/object_merge/runtime.yaml",
			expectedFile: "testdata/update_default_with_runtime_value/object_merge/empty_expected.yaml",
		}),

		Entry("object type with overlapping keys should merge properly", testCase{
			description:  "should merge objects with runtime values overriding original values",
			originalFile: "testdata/update_default_with_runtime_value/object_merge/original.yaml",
			runtimeFile:  "testdata/update_default_with_runtime_value/object_merge/runtime.yaml",
			expectedFile: "testdata/update_default_with_runtime_value/object_merge/expected.yaml",
		}),

		Entry("object type with nil runtime objectVal should remain unchanged", testCase{
			description:  "should keep original unchanged when runtime objectVal is nil",
			originalFile: "testdata/update_default_with_runtime_value/object_merge/original.yaml",
			runtimeFile:  "testdata/update_default_with_runtime_value/object_merge/nil_runtime.yaml",
			expectedFile: "testdata/update_default_with_runtime_value/object_merge/original.yaml",
		}),
	)
})
