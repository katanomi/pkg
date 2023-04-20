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
	"context"

	kvalidation "github.com/katanomi/pkg/apis/validation"
	. "github.com/katanomi/pkg/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var _ = Describe("Test.GitCreateCommitSpec.Validate", func() {

	var (
		r    *GitCreateCommit
		errs field.ErrorList
		ctx  context.Context
		err  error
	)

	BeforeEach(func() {
		ctx = context.Background()
		r = &GitCreateCommit{}
	})

	JustBeforeEach(func() {
		errs = r.Validate(ctx)
		err = kvalidation.ReturnInvalidError(schema.GroupKind{}, "kind", errs)
	})

	Context("Lots of validation errors", func() {
		BeforeEach(func() {
			MustLoadYaml("testdata/gitcreatecommit_validation.InvalidData.original.yaml", &r)
		})

		It("should return validation error", func() {
			Expect(err).ToNot(BeNil(), "should return an error")
			Expect(errors.IsInvalid(err)).To(BeTrue(), "should return an invalid error")

			statusErr, _ := err.(*errors.StatusError)
			Expect(statusErr.ErrStatus.Details.Causes).To(ContainElements(
				metav1.StatusCause{
					Type:    "FieldValueInvalid",
					Message: "Invalid value: \"\": commit message is required",
					Field:   "spec.message",
				},
				metav1.StatusCause{
					Type:    "FieldValueInvalid",
					Message: "Invalid value: \"[]\": actions is required",
					Field:   "spec.actions",
				},
				metav1.StatusCause{
					Type:    "FieldValueForbidden",
					Message: "Forbidden: only one of startBranch, startSHA OR startTag can be used, not all at the same time.",
					Field:   "spec",
				},
			))
			Expect(statusErr.ErrStatus.Details.Causes).To(HaveLen(3))
		})
	})

	Context("Valid", func() {
		BeforeEach(func() {
			MustLoadYaml("testdata/gitcreatecommit_validation.Valid.original.yaml", &r)
		})

		It("should not return validation error", func() {
			Expect(err).To(BeNil(), "should NOT return an error")
		})
	})
})
