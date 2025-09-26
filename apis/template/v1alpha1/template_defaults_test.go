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

package v1alpha1

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	pipev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

var _ = Describe("Template.SetDefaults", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
	})

	DescribeTable("template default value setting scenarios",
		func(description string, sourceTemplate Template, expectedTemplate Template) {
			By(description)

			// Apply SetDefaults
			sourceTemplate.SetDefaults(ctx)

			// Verify the result matches expected
			Expect(sourceTemplate).To(Equal(expectedTemplate),
				"SetDefaults result should match expected template")
		},
		Entry("template with parameters requiring defaults",
			"should apply ParamSlice.SetDefaults to template parameters while preserving other fields",
			Template{
				Name: "test-template",
				TemplateRef: pipev1beta1.ResolverRef{
					Resolver: "git",
					Params: []pipev1beta1.Param{
						{
							Name: "url",
							Value: pipev1beta1.ParamValue{
								StringVal: "https://github.com/example/repo",
							},
						},
						{
							Name: "revision",
							Value: pipev1beta1.ParamValue{
								StringVal: "main",
							},
						},
					},
				},
				Params: []pipev1beta1.Param{
					{
						Name: "string-param-no-type",
						Value: pipev1beta1.ParamValue{
							StringVal: "test-string-value",
						},
					},
					{
						Name: "array-param-no-type",
						Value: pipev1beta1.ParamValue{
							ArrayVal: []string{"item1", "item2"},
						},
					},
					{
						Name: "existing-type-param",
						Value: pipev1beta1.ParamValue{
							Type:      pipev1beta1.ParamTypeObject,
							StringVal: `{"key":"value"}`,
							ObjectVal: map[string]string{
								"key": "value",
							},
						},
					},
				},
				Metadata: v1alpha1.Metadata{
					Name: "template-metadata-name",
					PipelineTaskMetadata: pipev1beta1.PipelineTaskMetadata{
						Labels: map[string]string{
							"template-label": "test-value",
						},
						Annotations: map[string]string{
							"template-annotation": "test-annotation",
						},
					},
				},
			},
			Template{
				Name: "test-template",
				TemplateRef: pipev1beta1.ResolverRef{
					Resolver: "git",
					Params: []pipev1beta1.Param{
						{
							Name: "url",
							Value: pipev1beta1.ParamValue{
								Type:      pipev1beta1.ParamTypeString,
								StringVal: "https://github.com/example/repo",
							},
						},
						{
							Name: "revision",
							Value: pipev1beta1.ParamValue{
								Type:      pipev1beta1.ParamTypeString,
								StringVal: "main",
							},
						},
					},
				},
				Params: []pipev1beta1.Param{
					{
						Name: "string-param-no-type",
						Value: pipev1beta1.ParamValue{
							Type:      pipev1beta1.ParamTypeString,
							StringVal: "test-string-value",
						},
					},
					{
						Name: "array-param-no-type",
						Value: pipev1beta1.ParamValue{
							Type:     pipev1beta1.ParamTypeArray,
							ArrayVal: []string{"item1", "item2"},
						},
					},
					{
						Name: "existing-type-param",
						Value: pipev1beta1.ParamValue{
							Type:      pipev1beta1.ParamTypeObject,
							StringVal: `{"key":"value"}`,
							ObjectVal: map[string]string{
								"key": "value",
							},
						},
					},
				},
				Metadata: v1alpha1.Metadata{
					Name: "template-metadata-name",
					PipelineTaskMetadata: pipev1beta1.PipelineTaskMetadata{
						Labels: map[string]string{
							"template-label": "test-value",
						},
						Annotations: map[string]string{
							"template-annotation": "test-annotation",
						},
					},
				},
			},
		),
		Entry("template without parameters",
			"should handle templates without parameters gracefully without modification",
			Template{
				Name: "test-template-no-params",
				TemplateRef: pipev1beta1.ResolverRef{
					Resolver: "git",
					Params: []pipev1beta1.Param{
						{
							Name: "url",
							Value: pipev1beta1.ParamValue{
								StringVal: "https://github.com/example/repo",
							},
						},
					},
				},
				Metadata: v1alpha1.Metadata{
					Name: "template-metadata-name",
					PipelineTaskMetadata: pipev1beta1.PipelineTaskMetadata{
						Labels: map[string]string{
							"test-label": "test-value",
						},
					},
				},
			},
			Template{
				Name: "test-template-no-params",
				TemplateRef: pipev1beta1.ResolverRef{
					Resolver: "git",
					Params: []pipev1beta1.Param{
						{
							Name: "url",
							Value: pipev1beta1.ParamValue{
								Type:      pipev1beta1.ParamTypeString,
								StringVal: "https://github.com/example/repo",
							},
						},
					},
				},
				Metadata: v1alpha1.Metadata{
					Name: "template-metadata-name",
					PipelineTaskMetadata: pipev1beta1.PipelineTaskMetadata{
						Labels: map[string]string{
							"test-label": "test-value",
						},
					},
				},
			},
		),
		Entry("template with empty parameter slice",
			"should handle templates with empty parameter slice correctly without adding parameters",
			Template{
				Name: "test-template-empty-params",
				TemplateRef: pipev1beta1.ResolverRef{
					Resolver: "git",
					Params: []pipev1beta1.Param{
						{
							Name: "url",
							Value: pipev1beta1.ParamValue{
								StringVal: "https://github.com/example/repo",
							},
						},
					},
				},
				Params: []pipev1beta1.Param{},
				Metadata: v1alpha1.Metadata{
					Name: "template-metadata-name",
				},
			},
			Template{
				Name: "test-template-empty-params",
				TemplateRef: pipev1beta1.ResolverRef{
					Resolver: "git",
					Params: []pipev1beta1.Param{
						{
							Name: "url",
							Value: pipev1beta1.ParamValue{
								Type:      pipev1beta1.ParamTypeString,
								StringVal: "https://github.com/example/repo",
							},
						},
					},
				},
				Params: []pipev1beta1.Param{},
				Metadata: v1alpha1.Metadata{
					Name: "template-metadata-name",
				},
			},
		),
	)

	Describe("edge cases and implementation details", func() {
		Context("when template is nil", func() {
			It("should return immediately without errors", func() {
				var template *Template = nil
				Expect(func() {
					template.SetDefaults(ctx)
				}).ToNot(Panic(), "SetDefaults should not panic with nil template")
			})
		})

		Context("when context is nil", func() {
			It("should handle nil context gracefully", func() {
				template := &Template{
					Name: "test-template",
					Params: []pipev1beta1.Param{
						{
							Name: "test-param",
							Value: pipev1beta1.ParamValue{
								StringVal: "test-value",
							},
						},
					},
				}

				Expect(func() {
					template.SetDefaults(nil)
				}).ToNot(Panic(), "SetDefaults should not panic with nil context")

				// Verify that ParamSlice.SetDefaults was still called
				Expect(template.Params[0].Value.Type).To(Equal(pipev1beta1.ParamTypeString),
					"Should still set default parameter type even with nil context")
			})
		})

		Context("when template has params with mixed specifications", func() {
			It("should correctly delegate to ParamSlice.SetDefaults", func() {
				template := &Template{
					Name: "test-template",
					Params: []pipev1beta1.Param{
						{
							Name: "param-no-type",
							Value: pipev1beta1.ParamValue{
								StringVal: "value1",
							},
						},
						{
							Name: "param-with-type",
							Value: pipev1beta1.ParamValue{
								Type:     pipev1beta1.ParamTypeArray,
								ArrayVal: []string{"val1", "val2"},
							},
						},
					},
				}

				template.SetDefaults(ctx)

				// Verify first param got default type
				Expect(template.Params[0].Value.Type).To(Equal(pipev1beta1.ParamTypeString),
					"Parameter without type should get default string type")

				// Verify second param kept its type
				Expect(template.Params[1].Value.Type).To(Equal(pipev1beta1.ParamTypeArray),
					"Parameter with existing type should be preserved")
			})
		})
	})

	Describe("integration with ParamSlice.SetDefaults", func() {
		It("should preserve template fields while processing parameters", func() {
			template := &Template{
				Name: "integration-test-template",
				TemplateRef: pipev1beta1.ResolverRef{
					Resolver: "git",
					Params: []pipev1beta1.Param{
						{
							Name:  "url",
							Value: pipev1beta1.ParamValue{StringVal: "https://example.com/repo"},
						},
					},
				},
				Params: []pipev1beta1.Param{
					{
						Name: "template-param",
						Value: pipev1beta1.ParamValue{
							StringVal: "template-value",
						},
					},
				},
				Metadata: v1alpha1.Metadata{
					Name: "test-metadata",
					PipelineTaskMetadata: pipev1beta1.PipelineTaskMetadata{
						Labels: map[string]string{
							"test-label": "test-value",
						},
					},
				},
			}

			originalName := template.Name
			originalResolver := template.TemplateRef.Resolver
			originalMetadata := template.Metadata

			template.SetDefaults(ctx)

			// Verify template fields are preserved
			Expect(template.Name).To(Equal(originalName),
				"Template name should be preserved")
			Expect(template.TemplateRef.Resolver).To(Equal(originalResolver),
				"Template resolver should be preserved")
			Expect(template.Metadata).To(Equal(originalMetadata),
				"Template metadata should be preserved")

			// Verify parameter processing occurred
			Expect(template.Params[0].Value.Type).To(Equal(pipev1beta1.ParamTypeString),
				"Parameter should have default type set")
		})
	})
})
