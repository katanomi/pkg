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
	"k8s.io/apimachinery/pkg/util/validation/field"

	. "github.com/katanomi/pkg/testing"
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
	)

	BeforeEach(func() {
		path = field.NewPath("")
		clusterFilter = &ClusterFilter{}
	})

	JustBeforeEach(func() {
		errs = clusterFilter.Validate(path)
	})

	Context("empty struct", func() {
		It("should return validation error", func() {
			Expect(errs).ToNot(BeNil(), "should return an error")
			// Filter removes items from the ErrorList that match the provided fns.
			Expect(errs.Filter(field.NewErrorTypeMatcher(field.ErrorTypeRequired))).To(HaveLen(1))
			Expect(errs.Filter(field.NewErrorTypeMatcher(field.ErrorTypeInvalid))).To(HaveLen(1))
			Expect(errs).To(HaveLen(2))
		})
	})

	Context("Lots of validation errors", func() {
		BeforeEach(func() {
			Expect(LoadYAML("testdata/clusterfilter_validation.InvalidData.original.yaml", clusterFilter)).To(Succeed())
		})

		It("should return validation error", func() {
			Expect(errs).ToNot(BeNil(), "should return an error")
			// Filter removes items from the ErrorList that match the provided fns.
			Expect(errs.Filter(field.NewErrorTypeMatcher(field.ErrorTypeRequired))).To(HaveLen(3))
			Expect(errs.Filter(field.NewErrorTypeMatcher(field.ErrorTypeInvalid))).To(HaveLen(1))
			Expect(errs).To(HaveLen(4))
		})
	})

	Context("Valid", func() {
		BeforeEach(func() {
			Expect(LoadYAML("testdata/clusterfilter_validation.Valid.original.yaml", clusterFilter)).To(Succeed())
		})

		It("should not return validation error", func() {
			Expect(errs).To(HaveLen(0), "should NOT return an error")
		})
	})

})
