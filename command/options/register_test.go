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
	"testing"

	"k8s.io/apimachinery/pkg/util/validation/field"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/pflag"
)

var _ = Describe("Test.RegisterFlags", func() {
	type testStruct struct {
		ToolImageOption
	}

	var (
		flags *pflag.FlagSet
		obj   testStruct
	)

	BeforeEach(func() {
		flags = pflag.NewFlagSet("test", pflag.ContinueOnError)
		obj = testStruct{}
	})

	JustBeforeEach(func() {
		RegisterFlags(&obj, flags)
	})

	When("provide expected flags", func() {
		It("should get the expected value", func() {
			err := flags.Parse([]string{"--tool-image", "test-image"})
			Expect(err).Should(Succeed())
			Expect(obj.ToolImage).Should(Equal("test-image"))
		})
	})

	When("not provide expected flags", func() {
		It("should get empty value", func() {
			err := flags.Parse([]string{})
			Expect(err).Should(Succeed())
			Expect(obj.ToolImage).Should(Equal(""))
		})
	})
})

var _ = Describe("Test.RegisterSetup", func() {
	type testStruct struct {
		QualityGateRulesOption
	}

	var (
		obj testStruct
		ctx context.Context
	)

	BeforeEach(func() {
		obj = testStruct{}
		ctx = context.Background()
	})

	When("provide expected flags", func() {
		JustBeforeEach(func() {
			err := RegisterSetup(&obj, ctx, nil, []string{"--quality-gate-rules", "a=b", "c=d"})
			Expect(err).Should(Succeed())
		})

		It("should get the expected value", func() {
			Expect(obj.QualityGateRules).Should(Equal(map[string]string{
				"a": "b",
				"c": "d",
			}))
		})
	})

	When("not provide expected flags", func() {
		JustBeforeEach(func() {
			err := RegisterSetup(&obj, ctx, nil, []string{})
			Expect(err).Should(Succeed())
		})
		It("should get empty value", func() {
			Expect(obj.QualityGateRules).Should(Equal(map[string]string{}))
		})
	})
})

func TestCommandOption(t *testing.T) {
	g := NewGomegaWithT(t)
	ctx := context.Background()
	obj := struct {
		QualityGateRulesOption
		ScanFlagsOption
	}{}
	args := []string{
		"--quality-gate-rules", "a=b", "c=d",
		"--scan-flags", "e=f", "g=h",
	}
	err := RegisterSetup(&obj, ctx, nil, args)
	g.Expect(err).Should(Succeed())
	g.Expect(obj.QualityGateRules).To(Equal(map[string]string{
		"a": "b",
		"c": "d",
	}))
	g.Expect(obj.ScanFlags).To(Equal(map[string]string{
		"e": "f",
		"g": "h",
	}))
}

func TestDependencyReposOption(t *testing.T) {
	g := NewGomegaWithT(t)
	obj := DependencyReposOption{}
	root := field.NewPath("root")
	err := obj.Setup(context.Background(), nil, []string{
		"--dependencies-repositories",
		"repo1",
		"repo2",
		"htt:// abc",
	})

	g.Expect(err).Should(Succeed())
	g.Expect(obj.DependencyRepos).To(Equal([]string{
		"repo1",
		"repo2",
		"htt:// abc",
	}))
	validateErr := obj.Validate(root)
	g.Expect(validateErr).To(HaveLen(1))
	g.Expect(validateErr.ToAggregate().Error()).To(ContainSubstring("dependency repository is not a valid url"))
}

func TestOption_AddFlags(t *testing.T) {
	g := NewGomegaWithT(t)
	obj := struct {
		CommandOption
		KatanomiPathOption
		SourcePathOption
		ResultPathOption
		QualityGateOption
		ReportPathOption
		ToolImageOption
		ContextOption
	}{}
	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
	RegisterFlags(&obj, flagSet)

	err := flagSet.Parse([]string{
		"--command", "command-value",
		"--share-path", "share-path-value",
		"--config-path", "config-path-value",
		"--bin-path", "bin-path-value",
		"--source-path", "source-path-value",
		"--result-path", "result-path-value",
		"--enable-quality-gate",
		"--report-path", "report-path-value",
		"--tool-image", "tool-image-value",
		"--context", "context-value",
	})
	g.Expect(err).Should(Succeed())
	g.Expect(obj.Command).To(Equal("command-value"))
	g.Expect(obj.SharePath).To(Equal("share-path-value"))
	g.Expect(obj.ConfigPath).To(Equal("config-path-value"))
	g.Expect(obj.BinPath).To(Equal("bin-path-value"))
	g.Expect(obj.SourcePath).To(Equal("source-path-value"))
	g.Expect(obj.ResultPath).To(Equal("result-path-value"))
	g.Expect(obj.QualityGate).To(Equal(true))
	g.Expect(obj.ReportPath).To(Equal("report-path-value"))
	g.Expect(obj.ToolImage).To(Equal("tool-image-value"))
	g.Expect(obj.Context).To(Equal("context-value"))
}
