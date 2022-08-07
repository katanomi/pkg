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

package substitution

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/sets"
	"knative.dev/pkg/apis"
)

var _ = Describe("ValidateVariable", func() {
	DescribeTable("Tekton.Validate.Variable",
		func(name, value, prefix, locationName, path string, vars sets.String, err error) {
			realError := ValidateVariable(name, value, prefix, locationName, path, vars,
				parameterSubstitutionSample, TektonValidateSingleVariable)
			if err == nil {
				Expect(realError).To(BeNil())
			} else {
				Expect(realError).To(Equal(err))
			}

		},
		Entry("invalid sub variable",
			"abc", "$(params.abc.def)", "params", "somewhere", "params", sets.NewString("abc", "def"),
			&apis.FieldError{
				Message: fmt.Sprintf("non-existent variable in %q for %s %s", `$(params.abc.def)`, "somewhere", "abc"),
				Paths:   []string{"params.abc"},
			},
		),
		Entry("invalid variable",
			"xyz", "$(params.xyz)", "params", "somewhere", "params", sets.NewString("abc", "def"),
			&apis.FieldError{
				Message: fmt.Sprintf("non-existent variable in %q for %s %s", `$(params.xyz)`, "somewhere", "xyz"),
				Paths:   []string{"params.xyz"},
			},
		),
		Entry("no match",
			"xyz", "$(aaaa)", "params", "somewhere", "params", sets.NewString("abc", "def"),
			nil,
		),
		Entry("valid variable",
			"abc", "$(params.abc)", "params", "somewhere", "params", sets.NewString("abc", "def"),
			nil,
		),
		Entry("valid sub variable",
			"abc", "$(params.abc.def)", "params", "somewhere", "params", sets.NewString("abc", "def", "abc.def"),
			nil,
		),
	)

	DescribeTable("Version.Validate.Variable",
		func(name, value, prefix, locationName, path string, vars sets.String, err error) {
			realError := ValidateVariable(name, value, prefix, locationName, path, vars,
				parameterSubstitutionSample, VersionValidateSingleVariable)
			if err == nil {
				Expect(realError).To(BeNil())
			} else {
				Expect(realError).To(Equal(err))
			}

		},
		Entry("invalid sub variable",
			"abc", "$(context.abc.def)", "context", "localtionName", "path", sets.NewString("abc", "def"),
			&apis.FieldError{
				Message: "non-existent variable in \"$(context.abc.def)\" for localtionName abc.def",
				Paths:   []string{"path"},
			},
		),
		Entry("valid sub variable",
			"abc", "$(context.abc.def)", "context", "localtionName", "path", sets.NewString("abc", "def", "abc.def"),
			nil,
		),
	)
})
