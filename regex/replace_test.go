/*
Copyright 2023 The Katanomi Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and
limitations under the License.
*/

package regex

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"

	kvalidation "github.com/katanomi/pkg/apis/validation"
	. "github.com/katanomi/pkg/testing"
)

var _ = Describe("Test.Replaces.ReplaceAllString", func() {

	// nameAsTagReplaces used to generate the tag name from the branch name.
	// 1. replacing `/` and `_` to `-`
	// 2. remove the ending non [0-9a-zA-Z] characters
	// 3. maximum length limit is 30 (extra characters in the prefix will be removed.)
	nameAsTagReplaces := &Replaces{
		// replacing `/` and `_` to `-`
		{Regex: `[/_]`, Replacement: "-"},
		// remove the ending non [0-9a-zA-Z] characters
		{Regex: `[^0-9a-zA-Z]*$`, Replacement: ""},
		// maximum length limit is 30 (extra characters in the prefix will be removed.)
		// (?U) indicates that the regex is non-greedy regex.
		{Regex: `^(?U)(.*)(.{0,30})$`, Replacement: "${2}"},
		// remove the starting non [0-9a-zA-Z] characters
		{Regex: `^[^0-9a-zA-Z]*`, Replacement: ""},
	}

	DescribeTable("ReplaceAllString",
		func(rs *Replaces, original, expected string) {
			actual := rs.ReplaceAllString(original)
			Expect(actual).To(Equal(expected))
		},
		Entry("replaces is empty", nil, "original", "original"),
		Entry("contains / and _", nameAsTagReplaces, "feat/awesome_feature", "feat-awesome-feature"),
		Entry("ending contains / and _", nameAsTagReplaces, "original_/-", "original"),
		Entry("starting contains / and _", nameAsTagReplaces, "_/-original", "original"),
		Entry("length is 1", nameAsTagReplaces, "a", "a"),
		Entry("length greater than 30", nameAsTagReplaces, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", "wxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"),
		Entry("length greater than 30", nameAsTagReplaces, "0123456789abcdefghijklmnopqrstuvwxyz/-_+ABCDEFGHIJKLMNOPQRSTUVWXYZ", "ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
	)
})

var _ = Describe("Test.Replaces.Validate", func() {

	var (
		rs   Replaces
		path *field.Path
		errs field.ErrorList
		err  error
	)

	BeforeEach(func() {
		path = field.NewPath("prefix")
		rs = Replaces{}
	})

	JustBeforeEach(func() {
		errs = rs.Validate(path)
		err = kvalidation.ReturnInvalidError(schema.GroupKind{}, "kind", errs)
	})

	Context("empty struct", func() {
		It("should NOT return validation error", func() {
			Expect(err).To(BeNil(), "should NOT return an error")
		})
	})

	Context("Lots of validation errors", func() {
		BeforeEach(func() {
			MustLoadYaml("testdata/replaces_validation.InvalidData.original.yaml", &rs)
		})

		It("should return validation error", func() {
			Expect(err).ToNot(BeNil(), "should return an error")
			Expect(errors.IsInvalid(err)).To(BeTrue(), "should return an invalid error")

			statusErr, _ := err.(*errors.StatusError)
			Expect(statusErr.ErrStatus.Details.Causes).To(ContainElements(
				metav1.StatusCause{
					Type:    "FieldValueInvalid",
					Message: "Invalid value: \"[abcd$\": error parsing regexp: missing closing ]: `[abcd$`",
					Field:   "prefix[0].regex",
				},
				metav1.StatusCause{
					Type:    "FieldValueInvalid",
					Message: "Invalid value: \"^.0,128[.*{$\": error parsing regexp: missing closing ]: `[.*{$`",
					Field:   "prefix[1].regex",
				},
			))
			Expect(statusErr.ErrStatus.Details.Causes).To(HaveLen(2))
		})
	})

	Context("Valid", func() {
		BeforeEach(func() {
			MustLoadYaml("testdata/replaces_validation.Valid.original.yaml", &rs)
		})

		It("should not return validation error", func() {
			Expect(err).To(BeNil(), "should NOT return an error")
		})
	})

})
