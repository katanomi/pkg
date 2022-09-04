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
	"k8s.io/apimachinery/pkg/util/validation/field"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test.Check.Validate", func() {
	var (
		ctx   context.Context
		check *Check
		path  *field.Path
		errs  field.ErrorList
	)

	BeforeEach(func() {
		ctx = context.TODO()
		path = field.NewPath("spec")
		check = &Check{}
	})

	JustBeforeEach(func() {
		errs = check.Validate(ctx, path)
	})

	Context("Lots of validation errors", func() {
		BeforeEach(func() {
			Expect(pkgtesting.LoadYAML("testdata/check_validation.InvalidData.original.yaml", check)).To(Succeed())
		})

		It("should return validation error", func() {
			Expect(errs).ToNot(BeNil(), "should return an error")
			// Filter removes items from the ErrorList that match the provided fns.
			Expect(errs.Filter(field.NewErrorTypeMatcher(field.ErrorTypeDuplicate))).To(HaveLen(0))
			Expect(errs).To(HaveLen(1))
		})
	})

	Context("Valid", func() {
		BeforeEach(func() {
			Expect(pkgtesting.LoadYAML("testdata/check_validation.Valid.original.yaml", check)).To(Succeed())
		})

		It("should return validation error", func() {
			Expect(errs).To(BeNil(), "should NOT return an error")
		})
	})

})

var _ = Describe("Test.Check.ValidateChange", func() {
	var (
		ctx      context.Context
		old, new *Check
		path     *field.Path
		errs     field.ErrorList
	)

	BeforeEach(func() {
		ctx = context.TODO()
		path = field.NewPath("spec")
		old, new = &Check{}, &Check{}
	})

	JustBeforeEach(func() {
		errs = new.ValidateChange(ctx, old, path)
	})

	Context("old is nil and new is not nil", func() {
		BeforeEach(func() {
			old = nil
		})

		It("should return validation error", func() {
			Expect(errs).ToNot(BeNil(), "should return an error")
			// Filter removes items from the ErrorList that match the provided fns.
			Expect(errs.Filter(field.NewErrorTypeMatcher(field.ErrorTypeForbidden))).To(HaveLen(0))
			Expect(errs).To(HaveLen(1))
		})
	})

	Context("old is not nil and new is nil", func() {
		BeforeEach(func() {
			new = nil
		})

		It("should return validation error", func() {
			Expect(errs).ToNot(BeNil(), "should return an error")
			// Filter removes items from the ErrorList that match the provided fns.
			Expect(errs.Filter(field.NewErrorTypeMatcher(field.ErrorTypeForbidden))).To(HaveLen(0))
			Expect(errs).To(HaveLen(1))
		})
	})

	Context("Valid", func() {
		It("should return validation error", func() {
			Expect(errs).To(BeNil(), "should NOT return an error")
		})
	})

})
