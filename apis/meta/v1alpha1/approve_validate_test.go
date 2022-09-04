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

package v1alpha1

import (
	"context"

	pkgtesting "github.com/katanomi/pkg/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var _ = Describe("Test.ApprovalPolicy.Validate", func() {
	DescribeTable("ApprovalPolicy.Validate",
		func(policy ApprovalPolicy, expected bool) {
			actual := policy.Validate(context.TODO(), field.NewPath(""))
			Expect(actual != nil).To(Equal(expected))
		},
		Entry("ApprovalPolicyAny", ApprovalPolicyAny, false),
		Entry("ApprovalPolicyAll", ApprovalPolicyAll, false),
		Entry("ApprovalPolicyInOrder", ApprovalPolicyInOrder, false),
		Entry("invalid ApprovalPolicy", ApprovalPolicy("invalid"), true),
	)
})

var _ = Describe("Test.ApprovalSpec.Validate", func() {
	var (
		ctx          context.Context
		approvalSpec *ApprovalSpec
		path         *field.Path
		errs         field.ErrorList
	)

	BeforeEach(func() {
		ctx = context.TODO()
		path = field.NewPath("spec")
		approvalSpec = &ApprovalSpec{}
	})

	JustBeforeEach(func() {
		errs = approvalSpec.Validate(ctx, path)
	})

	Context("Lots of validation errors", func() {
		BeforeEach(func() {
			Expect(pkgtesting.LoadYAML("testdata/approvalspec_validation.InvalidData.original.yaml", approvalSpec)).To(Succeed())
		})

		It("should return validation error", func() {
			Expect(errs).ToNot(BeNil(), "should return an error")
			// Filter removes items from the ErrorList that match the provided fns.
			Expect(errs.Filter(field.NewErrorTypeMatcher(field.ErrorTypeNotSupported))).To(HaveLen(5))
			Expect(errs.Filter(field.NewErrorTypeMatcher(field.ErrorTypeRequired))).To(HaveLen(6))
			Expect(errs.Filter(field.NewErrorTypeMatcher(field.ErrorTypeInvalid))).To(HaveLen(3))
			Expect(errs).To(HaveLen(7))
		})
	})

	Context("users is empty", func() {
		It("should return validation error", func() {
			Expect(errs).ToNot(BeNil(), "should return an error")
			// Filter removes items from the ErrorList that match the provided fns.
			Expect(errs.Filter(field.NewErrorTypeMatcher(field.ErrorTypeNotSupported))).To(HaveLen(1))
			Expect(errs.Filter(field.NewErrorTypeMatcher(field.ErrorTypeInvalid))).To(HaveLen(1))
			Expect(errs).To(HaveLen(2))
		})
	})

	Context("Valid", func() {
		BeforeEach(func() {
			Expect(pkgtesting.LoadYAML("testdata/approvalspec_validation.Valid.original.yaml", approvalSpec)).To(Succeed())
		})

		It("should not return validation error", func() {
			Expect(errs).To(BeNil(), "should NOT return an error")
		})
	})

	Context("approval spec is nil", func() {
		BeforeEach(func() {
			approvalSpec = nil
		})

		It("should not return validation error", func() {
			Expect(errs).To(BeNil(), "should NOT return an error")
		})
	})

})
