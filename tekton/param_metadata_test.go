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

package tekton_test

import (
	"context"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pipelinev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"

	"github.com/katanomi/pkg/tekton"
	pkgtesting "github.com/katanomi/pkg/testing"
)

// ParamSpecWrapper wraps ParamSpec to implement NamedItem interface
type ParamSpecWrapper struct {
	*pipelinev1beta1.ParamSpec
}

// GetName implements NamedItem interface
func (p *ParamSpecWrapper) GetName() string {
	return p.Name
}

var _ = Describe("NewSpecMerger", func() {
	var (
		ctx context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	DescribeTable("parameter spec merging scenarios",
		func(testDir string, description string) {
			testDataPath := filepath.Join("testdata", "param_metadata", testDir)

			// Load test data files
			var destination []pipelinev1beta1.ParamSpec
			pkgtesting.MustLoadYaml(filepath.Join(testDataPath, "destination.yaml"), &destination)
			var source []pipelinev1beta1.ParamSpec
			pkgtesting.MustLoadYaml(filepath.Join(testDataPath, "source.yaml"), &source)
			var defaultSource []pipelinev1beta1.ParamSpec
			if _, err := os.Stat(filepath.Join(testDataPath, "default_source.yaml")); err == nil {
				pkgtesting.MustLoadYaml(filepath.Join(testDataPath, "default_source.yaml"), &defaultSource)
			}
			var expected []pipelinev1beta1.ParamSpec
			pkgtesting.MustLoadYaml(filepath.Join(testDataPath, "golden.yaml"), &expected)

			// Create spec merger processor
			processor := tekton.NewSpecMerger(&destination, source, defaultSource)

			// Process destination specs - the processor modifies destination in place
			for i := range destination {
				wrapper := &ParamSpecWrapper{ParamSpec: &destination[i]}
				processor(ctx, destination[i].Name, wrapper)
			}

			// The destination array is modified in place, so we compare it directly
			Expect(destination).To(Equal(expected), "Merged parameter specs should match golden file")
		},
		Entry("basic parameter merging with string, array, and object types",
			"spec_merger_basic",
			"Tests basic parameter merging where source descriptions and defaults take precedence over destination values"),
		Entry("object type parameter merging with property cleanup",
			"spec_merger_object_type",
			"Tests object parameter merging, default source merging, and schema-based property cleanup for object types"),
	)

	DescribeTable("edge case scenarios",
		func(testDir string, description string) {
			testDataPath := filepath.Join("testdata", "param_metadata", "spec_merger_edge_cases", testDir)

			// Load test data files
			var destination []pipelinev1beta1.ParamSpec
			pkgtesting.MustLoadYaml(filepath.Join(testDataPath, "destination.yaml"), &destination)
			var source []pipelinev1beta1.ParamSpec
			if _, err := os.Stat(filepath.Join(testDataPath, "source.yaml")); err == nil {
				pkgtesting.MustLoadYaml(filepath.Join(testDataPath, "source.yaml"), &source)
			}
			var expected []pipelinev1beta1.ParamSpec
			pkgtesting.MustLoadYaml(filepath.Join(testDataPath, "golden.yaml"), &expected)

			// Create spec merger processor
			processor := tekton.NewSpecMerger(&destination, source, nil)

			// Process destination specs - the processor modifies destination in place
			for i := range destination {
				wrapper := &ParamSpecWrapper{ParamSpec: &destination[i]}
				processor(ctx, destination[i].Name, wrapper)
			}

			// The destination array is modified in place, so we compare it directly
			Expect(destination).To(Equal(expected), "Edge case handling should match golden file")
		},
		Entry("empty destination array handling",
			"empty_destination",
			"Tests behavior with empty destination parameter specs - should handle gracefully without errors"),
		Entry("missing source parameter handling",
			"missing_source",
			"Tests behavior when source parameters don't match destination - original values should remain unchanged"),
		Entry("nil value scenarios",
			"nil_scenarios",
			"Tests behavior with nil default values and various edge conditions - should handle gracefully"),
		Entry("type mismatch between source and destination parameters is ignored",
			"type_mismatch",
			"Tests that merge skips metadata updates when source and destination parameter types differ"),
	)

	Context("processor function behavior", func() {
		It("should return false for parameters not found in source", func() {
			source := []pipelinev1beta1.ParamSpec{
				{Name: "existing-param", Type: pipelinev1beta1.ParamTypeString},
			}
			destination := []pipelinev1beta1.ParamSpec{
				{Name: "non-existent-param", Type: pipelinev1beta1.ParamTypeString},
			}
			processor := tekton.NewSpecMerger(&destination, source, nil)

			nonExistentParam := pipelinev1beta1.ParamSpec{
				Name: "non-existent-param",
				Type: pipelinev1beta1.ParamTypeString,
			}

			wrapper := &ParamSpecWrapper{ParamSpec: &nonExistentParam}
			processed := processor(ctx, nonExistentParam.Name, wrapper)
			Expect(processed).To(BeFalse(), "Processor should return false for parameters not found in source")
		})

		It("should return true for parameters found in source and modify destination", func() {
			source := []pipelinev1beta1.ParamSpec{
				{
					Name:        "test-param",
					Type:        pipelinev1beta1.ParamTypeString,
					Description: "Updated description",
					Default: &pipelinev1beta1.ParamValue{
						Type:      pipelinev1beta1.ParamTypeString,
						StringVal: "updated-value",
					},
				},
			}
			destination := []pipelinev1beta1.ParamSpec{
				{
					Name:        "test-param",
					Type:        pipelinev1beta1.ParamTypeString,
					Description: "Original description",
					Default: &pipelinev1beta1.ParamValue{
						Type:      pipelinev1beta1.ParamTypeString,
						StringVal: "original-value",
					},
				},
			}
			processor := tekton.NewSpecMerger(&destination, source, nil)

			wrapper := &ParamSpecWrapper{ParamSpec: &destination[0]}
			processed := processor(ctx, destination[0].Name, wrapper)
			Expect(processed).To(BeTrue(), "Processor should return true for parameters found in source")
			Expect(destination[0].Description).To(Equal("Updated description"), "Description should be updated from source")
			Expect(destination[0].Default.StringVal).To(Equal("updated-value"), "Default value should be updated from source")
		})
	})
})
