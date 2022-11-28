/*
Copyright 2022 The Katanomi Authors.

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

package options

import (
	// "context"
	// "testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

var _ = Describe("KeyValueListOption.Validate", func() {
	var (
		opts KeyValueListOption
		path *field.Path

		errs field.ErrorList
	)

	BeforeEach(func() {
		opts = KeyValueListOption{FlagName: "key-value", KeyValues: map[string]string{}}
		path = field.NewPath("params")
	})

	JustBeforeEach(func() {
		errs = opts.Validate(path)
	})

	When("execute basic key/value validations", func() {
		BeforeEach(func() {
			opts.AddValidation(
				RequiredKeyValueOptionValidation,
			)
			opts.KeyValues[""] = "some-value"
			opts.KeyValues["some-key"] = ""
		})
		It("should return validation errors", func() {
			Expect(errs).To(HaveLen(2))
		})
	})
	When("execute basic empty validations", func() {
		BeforeEach(func() {
			opts.AddValidation(
				RequiredKeyValueOptionValidation,
				NotEmptyKeyValueOptionValidation,
			)
		})
		It("should return validation errors", func() {
			Expect(errs).To(HaveLen(1))
		})
	})
})
