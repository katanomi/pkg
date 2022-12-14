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
	"context"

	cmdArgs "github.com/katanomi/pkg/command/args"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test.EnvFlagsOption.Setup", func() {
	var (
		ctx context.Context
		obj struct {
			EnvFlagsOption
		}
		args []string
		err  error
	)

	BeforeEach(func() {
		ctx = context.Background()
		obj = struct{ EnvFlagsOption }{EnvFlagsOption: EnvFlagsOption{}}
	})

	JustBeforeEach(func() {
		err = RegisterSetup(&obj, ctx, nil, args)
	})

	When("args is empty", func() {
		BeforeEach(func() {
			args = []string{}
		})
		It("returns nil err", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	When("invalid env items", func() {
		BeforeEach(func() {
			args = []string{
				"--env-flags", "foo", "bar",
			}
		})
		It("return error message with invalid env items", func() {
			Expect(err.Error()).Should(Equal("invalid env-flags: invalid env items: [foo bar]"))
		})
	})

	When("env items with duplicated keys", func() {
		BeforeEach(func() {
			args = []string{
				"--env-flags", "foo=1", "foo=2",
			}
		})
		It("return error message with duplicated key env items", func() {
			Expect(err.Error()).Should(Equal("invalid env-flags: duplicated env keys: [foo]"))
		})
	})

	When("with both invalid items and duplicated keys", func() {
		BeforeEach(func() {
			args = []string{
				"--env-flags", "foo=1", "foo=2", "bar",
			}
		})
		It("return error message with duplicated key and invalid env items", func() {
			Expect(err.Error()).Should(Equal("invalid env-flags: duplicated env keys: [foo],invalid env items: [bar]"))
		})
	})

	When("with required validation option and empty env-flags", func() {
		BeforeEach(func() {
			args = []string{
				"--env-flags",
			}
			ctx = WithValuesValidationOpts(ctx, []cmdArgs.ValuesValidateOption{cmdArgs.ValuesValidationOptRequired})

		})
		It("return error message with duplicated key and invalid env items", func() {
			Expect(err.Error()).Should(Equal("empty values"))
		})
	})

	When("with valid env-flags", func() {
		BeforeEach(func() {
			args = []string{
				"--env-flags", "FOO=BAR",
			}
		})

		It("should return valid key pairs", func() {
			Expect(err).Should(Succeed())
			Expect(obj.EnvFlags).To(Equal(map[string]string{
				"FOO": "BAR",
			}))
		})
	})

})
