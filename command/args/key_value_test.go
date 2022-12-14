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
	"fmt"
	"testing"

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
		err     error
		opts    []args.ValuesValidateOption
	)
	BeforeEach(func() {
		ctx = context.Background()
		cliArgs = []string{}
		flag = ""
	})
	JustBeforeEach(func() {
		result, err = args.GetKeyValues(ctx, cliArgs, flag, opts...)
	})
	When("args are nil", func() {
		BeforeEach(func() {
			cliArgs = nil
			opts = []args.ValuesValidateOption{args.ValuesValidationOptRequired}
		})
		It("should return empty map with error", func() {
			Expect(result).To(Equal(map[string]string{}))
			Expect(err).Should(HaveOccurred())
		})
	})
	When("args does not have the flag", func() {
		BeforeEach(func() {
			cliArgs = []string{"--flag1", "key=value", "--other-flag", "another-flag-key=123"}
			flag = "non-existing-flag"
			opts = []args.ValuesValidateOption{args.ValuesValidationOptRequired}
		})
		It("should return empty map with error", func() {
			Expect(result).To(Equal(map[string]string{}))
			Expect(err).To(HaveOccurred())
		})
	})
	When("flag does not have items", func() {
		BeforeEach(func() {
			cliArgs = []string{"--flag1", "--other-flag", "another-flag-key=123"}
			flag = "flag1"
			opts = []args.ValuesValidateOption{args.ValuesValidationOptRequired}
		})
		It("should return empty map with error", func() {
			Expect(result).To(Equal(map[string]string{}))
			Expect(err).Should(HaveOccurred())
		})
	})
	When("flag does not have key-value pairs", func() {
		BeforeEach(func() {
			cliArgs = []string{"--flag1", "this-is-a-value", "this-is-another-value", "--other-flag", "another-flag-key=123"}
			flag = "flag1"
		})
		It("should return empty map without error", func() {
			Expect(result).To(Equal(map[string]string{}))
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
	When("flag has key-value pairs", func() {
		BeforeEach(func() {
			cliArgs = []string{
				"--flag1", "this-is-key=this-is-value", "this-is-another-key=this-is-another-value", "key=with=multiple=equalsigns",
				"--other-flag", "another-flag-key=123"}
			flag = "flag1"
		})
		It("should return map with keys and values and nil error", func() {
			Expect(result).To(Equal(map[string]string{
				"this-is-key":         "this-is-value",
				"this-is-another-key": "this-is-another-value",
				"key":                 "with=multiple=equalsigns",
			}))
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	When("flag with failed validation option", func() {
		BeforeEach(func() {
			cliArgs = []string{
				"--flag1", "this-is-key=this-is-value", "this-is-another-key=this-is-another-value", "key=with=multiple=equalsigns",
				"--other-flag", "another-flag-key=123"}
			flag = "flag1"
			opts = []args.ValuesValidateOption{
				func(values []string) error {
					return fmt.Errorf("invalid validation")
				},
			}
		})
		It("should return map with keys and values and validation error", func() {
			Expect(result).To(Equal(map[string]string{
				"this-is-key":         "this-is-value",
				"this-is-another-key": "this-is-another-value",
				"key":                 "with=multiple=equalsigns",
			}))
			Expect(err).Should(HaveOccurred())
		})
	})

	When("flag with successful validation option", func() {
		BeforeEach(func() {
			cliArgs = []string{
				"--flag1", "this-is-key=this-is-value", "this-is-another-key=this-is-another-value", "key=with=multiple=equalsigns",
				"--other-flag", "another-flag-key=123"}
			flag = "flag1"
			opts = []args.ValuesValidateOption{
				func(values []string) error {
					return nil
				},
			}
		})
		It("should return map with keys and values and validation error", func() {
			Expect(result).To(Equal(map[string]string{
				"this-is-key":         "this-is-value",
				"this-is-another-key": "this-is-another-value",
				"key":                 "with=multiple=equalsigns",
			}))
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

})

func Test_ValuesValidationOptDuplicatedKeys(t *testing.T) {
	tests := []struct {
		name   string
		pairs  []string
		errMsg string
	}{
		{
			name:   "nil pairs",
			pairs:  nil,
			errMsg: "",
		},
		{
			name:   "empty pairs",
			pairs:  []string{},
			errMsg: "",
		},
		{
			name:   "duplicated pairs",
			pairs:  []string{"foo=1", "foo=2", "foo=3"},
			errMsg: "invalid env-flags: duplicated env keys: [foo]",
		},
		{
			name:   "invalid pairs",
			pairs:  []string{"foo", "foo", "bar"},
			errMsg: "invalid env-flags: invalid env items: [foo bar]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := args.ValuesValidationOptDuplicatedKeys(tt.pairs)
			if (err != nil) != (tt.errMsg != "") {
				t.Errorf("ValuesValidationOptDuplicatedKeys() error = %v, want errmsg: %v", err, tt.errMsg)
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("ValuesValidationOptDuplicatedKeys() error.Error = %v, want errmsg: %v", err.Error(), tt.errMsg)
			}
		})
	}
}
