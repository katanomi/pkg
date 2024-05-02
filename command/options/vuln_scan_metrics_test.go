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

func TestVulnScanMetricsOption_AddFlags(t *testing.T) {
	g := NewGomegaWithT(t)
	obj := struct {
		VulnScanMetricsOption
	}{}
	flagSet := pflag.NewFlagSet("metrics-result-limit", pflag.ContinueOnError)
	RegisterFlags(&obj, flagSet)

	err := flagSet.Parse([]string{
		"--metrics-result-limit", "10",
	})
	g.Expect(err).Should(Succeed())
	g.Expect(obj.ResultLimit).Should(Equal(10))
}

var _ = Describe("Test.VulnScanMetricsOption.Validate", func() {
	var (
		obj *VulnScanMetricsOption
		err field.ErrorList
	)
	BeforeEach(func() {
		obj = &VulnScanMetricsOption{}
		err = field.ErrorList{}
	})
	JustBeforeEach(func() {
		err = obj.Validate(field.NewPath("test"))
	})
	Context("when result limit is equal 0", func() {
		BeforeEach(func() {
			obj.ResultLimit = -1
		})
		It("should return error", func() {
			Expect(err).ShouldNot(BeEmpty())
		})
	})
	Context("when result limit is bigger than 0", func() {
		BeforeEach(func() {
			obj.ResultLimit = 1
		})
		It("should not return error", func() {
			Expect(err).Should(BeEmpty())
		})
	})
})
