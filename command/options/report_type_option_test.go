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

	"github.com/katanomi/pkg/report"
	. "github.com/onsi/gomega"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestReportTypeOption(t *testing.T) {
	t.Run("default flag ande setup", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctx := context.Background()
		base := field.NewPath("base")

		obj := struct {
			ReportTypeOption
		}{}
		args := []string{
			"--report-type", "junit-xml",
		}

		flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
		RegisterFlags(&obj, flagSet)
		err := flagSet.Parse(args)
		g.Expect(err).Should(Succeed(), "parse flag succeed.")

		err = RegisterSetup(&obj, ctx, nil, args)
		g.Expect(err).Should(Succeed(), "step flag succeed.")
		g.Expect(obj.SupportedType).ToNot(BeNil())
		g.Expect(obj.SupportedType).To(Equal(report.DefaultReportParsers), "should equal default")
		g.Expect(obj.Validate(base)).To(HaveLen(0), "validate succeed")
	})

	t.Run("invalid type", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctx := context.Background()
		base := field.NewPath("base")

		obj := struct {
			ReportTypeOption
		}{}
		args := []string{
			"--report-type", "Junit",
		}

		flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
		RegisterFlags(&obj, flagSet)
		err := flagSet.Parse(args)
		g.Expect(err).Should(Succeed(), "parse flag succeed.")

		err = RegisterSetup(&obj, ctx, nil, args)
		g.Expect(err).Should(Succeed(), "step flag succeed.")
		g.Expect(obj.SupportedType).ToNot(BeNil())
		g.Expect(obj.SupportedType).To(Equal(report.DefaultReportParsers), "should equal default")
		g.Expect(obj.Validate(base)).To(HaveLen(1), "validate failed")
	})

	t.Run("custom type parser", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctx := context.Background()
		base := field.NewPath("base")

		obj := struct {
			ReportTypeOption
		}{}

		obj.SupportedType = map[report.ReportType]report.ReportParser{"Junit": &report.JunitParser{}}

		args := []string{
			"--report-type", "Junit",
		}

		flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
		RegisterFlags(&obj, flagSet)
		err := flagSet.Parse(args)
		g.Expect(err).Should(Succeed(), "parse flag succeed.")

		err = RegisterSetup(&obj, ctx, nil, args)
		g.Expect(err).Should(Succeed(), "step flag succeed.")
		g.Expect(obj.SupportedType).ToNot(BeNil())
		g.Expect(obj.Validate(base)).To(HaveLen(0), "validate success")
	})

}
