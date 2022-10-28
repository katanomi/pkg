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
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestCommandOption_AddFlags(t *testing.T) {
	g := NewGomegaWithT(t)
	obj := struct {
		CommandOption
	}{}
	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
	RegisterFlags(&obj, flagSet)

	err := flagSet.Parse([]string{
		"--command", "xx",
	})
	g.Expect(err).Should(Succeed())
	g.Expect(obj.Command).Should(Equal("xx"))
}

var _ = Describe("Test.CommandOption.Validate", func() {
	var (
		obj *CommandOption
		err field.ErrorList
	)
	BeforeEach(func() {
		obj = &CommandOption{}
		err = field.ErrorList{}
	})
	JustBeforeEach(func() {
		err = obj.Validate(field.NewPath("test"))
	})
	Context("when command is empty", func() {
		BeforeEach(func() {
			obj.Command = ""
		})
		It("should return error", func() {
			Expect(err).ShouldNot(BeEmpty())
		})
	})
	Context("when command is not empty", func() {
		BeforeEach(func() {
			obj.Command = "xx"
		})
		It("should not return error", func() {
			Expect(err).Should(BeEmpty())
		})
	})
})
