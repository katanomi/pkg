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

func TestContextOption_AddFlags(t *testing.T) {
	g := NewGomegaWithT(t)
	obj := struct {
		ContextOption
	}{}
	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
	RegisterFlags(&obj, flagSet)

	err := flagSet.Parse([]string{
		"--context", ".",
	})
	g.Expect(err).Should(Succeed())
	g.Expect(obj.Context).Should(Equal("."))
}

var _ = Describe("Test.ContextOption.Validate", func() {
	var (
		obj *ContextOption
		err field.ErrorList
	)
	BeforeEach(func() {
		obj = &ContextOption{}
		err = field.ErrorList{}
	})
	JustBeforeEach(func() {
		err = obj.Validate(field.NewPath("test"))
	})
	Context("when context is empty", func() {
		BeforeEach(func() {
			obj.Context = ""
		})
		It("should return error", func() {
			Expect(err).ShouldNot(BeEmpty())
		})
	})
	Context("when context is not empty", func() {
		BeforeEach(func() {
			obj.Context = "xx"
		})
		It("should not return error", func() {
			Expect(err).Should(BeEmpty())
		})
	})
})
