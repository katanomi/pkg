/*
Copyright 2023 The Katanomi Authors.

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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"

	kvalidation "github.com/katanomi/pkg/apis/validation"
	. "github.com/katanomi/pkg/testing"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Test.ClusterFilterRule", func() {

	DescribeTable("ClusterFilterRule.Filter",
		func(rule ClusterFilterRule, resourcesFilePath string, expectedLen int) {
			loaded, loadErr := LoadKubeResourcesAsUnstructured(resourcesFilePath)
			Expect(loadErr).To(BeNil())
			actual := rule.Filter(loaded)
			Expect(actual).To(HaveLen(expectedLen))
		},

		Entry("exact field is empty", ClusterFilterRule{
			Exact: nil,
		}, "testdata/clusterfilterrule_filter.empty.yaml", 0),

		Entry("exact one", ClusterFilterRule{Exact: map[string]string{
			"$(metadata.name)": "cluster-1",
		}}, "testdata/clusterfilterrule_filter.list.yaml", 1),

		Entry("attribute with escaped characters", ClusterFilterRule{Exact: map[string]string{
			"$(metadata.labels.label\\.key1)": "value-1",
		}}, "testdata/clusterfilterrule_filter.list.yaml", 1),

		Entry("matches to multiple", ClusterFilterRule{Exact: map[string]string{
			"$(metadata.labels.label/key1)": "value-1",
		}}, "testdata/clusterfilterrule_filter.list.yaml", 3),
	)

})

var _ = Describe("Test.ClusterFilter.Validate", func() {

	var (
		clusterFilter *ClusterFilter
		path          *field.Path
		errs          field.ErrorList
		err           error
	)

	BeforeEach(func() {
		path = field.NewPath("prefix")
		clusterFilter = &ClusterFilter{}
	})

	JustBeforeEach(func() {
		errs = clusterFilter.Validate(path)
		err = kvalidation.ReturnInvalidError(schema.GroupKind{}, "kind", errs)
	})

	Context("empty struct", func() {
		It("should return validation error", func() {
			Expect(err).ToNot(BeNil(), "should return an error")
			Expect(errors.IsInvalid(err)).To(BeTrue(), "should return an invalid error")

			statusErr, _ := err.(*errors.StatusError)
			Expect(statusErr.ErrStatus.Details.Causes).To(ContainElements(
				metav1.StatusCause{
					Type:    "FieldValueRequired",
					Message: "Required value: one of selector OR refs is required",
					Field:   "prefix",
				},
			))
		})
	})

	Context("Lots of validation errors", func() {
		BeforeEach(func() {
			Expect(LoadYAML("testdata/clusterfilter_validation.InvalidData.original.yaml", clusterFilter)).To(Succeed())
		})

		It("should return validation error", func() {
			Expect(err).ToNot(BeNil(), "should return an error")
			Expect(errors.IsInvalid(err)).To(BeTrue(), "should return an invalid error")

			statusErr, _ := err.(*errors.StatusError)
			Expect(statusErr.ErrStatus.Details.Causes).To(ContainElements(
				metav1.StatusCause{
					Type:    "FieldValueInvalid",
					Message: "Invalid value: \"default-\": a lowercase RFC 1123 label must consist of lower case alphanumeric characters or '-', and must start and end with an alphanumeric character (e.g. 'my-name',  or '123-abc', regex used for validation is '[a-z0-9]([-a-z0-9]*[a-z0-9])?')",
					Field:   "prefix.namespace",
				},
				metav1.StatusCause{
					Type:    "FieldValueInvalid",
					Message: "Invalid value: \"app-\": name part must consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character (e.g. 'MyName',  or 'my.name',  or '123-abc', regex used for validation is '([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]')",
					Field:   "prefix.selector.matchLabels",
				},
				metav1.StatusCause{
					Type:    "FieldValueInvalid",
					Message: "Invalid value: \"Unknown\": not a valid selector operator",
					Field:   "prefix.selector.matchExpressions[0].operator",
				},
				metav1.StatusCause{
					Type:    "FieldValueRequired",
					Message: "Required value: must be specified when `operator` is 'In' or 'NotIn'",
					Field:   "prefix.selector.matchExpressions[1].values",
				},
			))
		})
	})

	Context("Valid", func() {
		BeforeEach(func() {
			Expect(LoadYAML("testdata/clusterfilter_validation.Valid.original.yaml", clusterFilter)).To(Succeed())
		})

		It("should not return validation error", func() {
			Expect(err).To(BeNil(), "should NOT return an error")
		})
	})

})
