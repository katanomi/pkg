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

package args_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/katanomi/pkg/command/args"
)

var _ = Describe("GetKeyValues", func() {
	var (
		ctx     context.Context
		cliArgs []string
		flag    string
		result  map[string]string
		ok      bool
	)
	BeforeEach(func() {
		ctx = context.Background()
		cliArgs = []string{}
		flag = ""
	})
	JustBeforeEach(func() {
		result, ok = args.GetKeyValues(ctx, cliArgs, flag)
	})
	When("args are nil", func() {
		BeforeEach(func() {
			cliArgs = nil
		})
		It("should return empty map with false ok", func() {
			Expect(result).To(Equal(map[string]string{}))
			Expect(ok).To(BeFalse())
		})
	})
	When("args does not have the flag", func() {
		BeforeEach(func() {
			cliArgs = []string{"--flag1", "key=value", "--other-flag", "another-flag-key=123"}
			flag = "non-existing-flag"
		})
		It("should return empty map with false ok", func() {
			Expect(result).To(Equal(map[string]string{}))
			Expect(ok).To(BeFalse())
		})
	})
	When("flag does not have items", func() {
		BeforeEach(func() {
			cliArgs = []string{"--flag1", "--other-flag", "another-flag-key=123"}
			flag = "flag1"
		})
		It("should return empty map with false ok", func() {
			Expect(result).To(Equal(map[string]string{}))
			Expect(ok).To(BeTrue())
		})
	})
	When("flag does not have key-value pairs", func() {
		BeforeEach(func() {
			cliArgs = []string{"--flag1", "this-is-a-value", "this-is-another-value", "--other-flag", "another-flag-key=123"}
			flag = "flag1"
		})
		It("should return empty map with true ok", func() {
			Expect(result).To(Equal(map[string]string{}))
			Expect(ok).To(BeTrue())
		})
	})
	When("flag has key-value pairs", func() {
		BeforeEach(func() {
			cliArgs = []string{
				"--flag1", "this-is-key=this-is-value", "this-is-another-key=this-is-another-value", "key=with=multiple=equalsigns",
				"--other-flag", "another-flag-key=123"}
			flag = "flag1"
		})
		It("should return map with keys and values and true ok", func() {
			Expect(result).To(Equal(map[string]string{
				"this-is-key":         "this-is-value",
				"this-is-another-key": "this-is-another-value",
				"key":                 "with=multiple=equalsigns",
			}))
			Expect(ok).To(BeTrue())
		})
	})
})
