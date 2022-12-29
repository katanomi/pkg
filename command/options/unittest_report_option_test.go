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

	. "github.com/onsi/gomega"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestUnitTestReportOption(t *testing.T) {
	t.Run("default flag and setup", func(t *testing.T) {
		g := NewGomegaWithT(t)
		base := field.NewPath("base")

		obj := struct {
			UnitTestReportOption
		}{}
		args := []string{
			"--report-type", "junit-xml",
			"--report-path", "./junit.xml",
			"--result-path", "/tmp/result",
		}

		flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
		RegisterFlags(&obj, flagSet)
		err := flagSet.Parse(args)
		g.Expect(err).Should(Succeed(), "parse flag succeed.")
		g.Expect(obj.ReportType).To(Equal("junit-xml"), "parse report type succeed.")
		g.Expect(obj.ReportPath).To(Equal("./junit.xml"), "parse report path succeed.")
		g.Expect(obj.ResultPath).To(Equal("/tmp/result"), "parse result path succeed.")
		err = RegisterSetup(&obj, context.Background(), nil, args)

		g.Expect(err).Should(Succeed(), "step succeed.")
		g.Expect(obj.Validate(base)).To(HaveLen(0), "validate succeed")
	})

	t.Run("custom flag and setup", func(t *testing.T) {
		g := NewGomegaWithT(t)
		base := field.NewPath("base")

		obj := struct {
			UnitTestReportOption
		}{}
		obj.ReportPathOption.FlagName = "custom-path"
		obj.ReportTypeOption.FlagName = "custom-type"
		obj.ResultPathOption.FlagName = "custom-result"
		args := []string{
			"--custom-type", "junit-xml",
			"--custom-path", "./junit.xml",
			"--custom-result", "/tmp/result",
		}

		flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
		RegisterFlags(&obj, flagSet)
		err := flagSet.Parse(args)
		g.Expect(err).Should(Succeed(), "parse flag succeed.")
		g.Expect(obj.ReportType).To(Equal("junit-xml"), "parse report type succeed.")
		g.Expect(obj.ReportPath).To(Equal("./junit.xml"), "parse report path succeed.")
		g.Expect(obj.ResultPath).To(Equal("/tmp/result"), "parse result path succeed.")
		err = RegisterSetup(&obj, context.Background(), nil, args)

		g.Expect(err).Should(Succeed(), "step succeed.")
		g.Expect(obj.Validate(base)).To(HaveLen(0), "validate succeed")
	})

	t.Run("Validate failed", func(t *testing.T) {
		g := NewGomegaWithT(t)
		base := field.NewPath("base")

		obj := struct {
			UnitTestReportOption
		}{}
		args := []string{
			"--report-type", "junit-xml",
			"--result-path", "/tmp/result",
		}
		flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
		RegisterFlags(&obj, flagSet)
		err := flagSet.Parse(args)
		g.Expect(err).Should(Succeed(), "parse flag succeed.")

		err = RegisterSetup(&obj, context.Background(), nil, args)
		g.Expect(err).Should(Succeed(), "step succeed.")
		g.Expect(obj.ReportType).To(Equal("junit-xml"), "parse report type succeed.")
		g.Expect(obj.ReportPath).To(Equal(""), "parse report path succeed.")
		g.Expect(obj.ResultPath).To(Equal("/tmp/result"), "parse result path succeed.")
		g.Expect(obj.Validate(base)).To(HaveLen(1), "validate failed")
	})

}
