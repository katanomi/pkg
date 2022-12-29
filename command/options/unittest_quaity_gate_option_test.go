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

	"github.com/katanomi/pkg/apis/codequality/v1alpha1"
	. "github.com/onsi/gomega"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestUnitTestQuaityGateOption(t *testing.T) {
	t.Run("open quality and pass", func(t *testing.T) {
		g := NewGomegaWithT(t)
		base := field.NewPath("base")

		obj := struct {
			UnitTestQuaityGateOption
		}{}
		args := []string{
			"--enable-quality-gate", "true",
		}

		flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
		RegisterFlags(&obj, flagSet)
		err := flagSet.Parse(args)
		g.Expect(err).Should(Succeed(), "parse flag succeed.")
		g.Expect(obj.QualityGate).To(Equal(true), "parse quality gate succeed.")

		args = append(args, "--quality-gate-rules", "passed-tests-rate=90", "lines-coverage=80", "branches-coverage=66.66")
		err = RegisterSetup(&obj, context.Background(), nil, args)
		g.Expect(err).Should(Succeed(), "setp succeed.")

		expectRules := map[string]string{PassedTestsRateMetric: "90", LinesCoverageMetric: "80", BranchesCoverageMetric: "66.66"}
		g.Expect(obj.QualityGateRules).To(Equal(expectRules), "parse quality rule succeed.")

		g.Expect(obj.Validate(base)).To(HaveLen(0), "validate succeed")

	})

	t.Run("disable quality and check rule passed", func(t *testing.T) {
		g := NewGomegaWithT(t)
		base := field.NewPath("base")

		obj := struct {
			UnitTestQuaityGateOption
		}{}
		args := []string{
			"--enable-quality-gate", "false",
		}

		flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
		RegisterFlags(&obj, flagSet)
		err := flagSet.Parse(args)
		g.Expect(err).Should(Succeed(), "parse flag succeed.")
		g.Expect(obj.QualityGate).To(Equal(true), "parse quality gate succeed.")

		args = append(args, "--quality-gate-rules", "passed-tests-rate=90", "lines-coverage=80", "branches-coverage=66.66")
		err = RegisterSetup(&obj, context.Background(), nil, args)
		g.Expect(err).Should(Succeed(), "setp succeed.")

		expectRules := map[string]string{PassedTestsRateMetric: "90", LinesCoverageMetric: "80", BranchesCoverageMetric: "66.66"}
		g.Expect(obj.QualityGateRules).To(Equal(expectRules), "parse quality rule succeed.")

		g.Expect(obj.Validate(base)).To(HaveLen(0), "validate succeed")

	})

	t.Run("empty rule validate pass", func(t *testing.T) {
		g := NewGomegaWithT(t)
		base := field.NewPath("base")

		obj := struct {
			UnitTestQuaityGateOption
		}{}
		args := []string{}

		flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
		RegisterFlags(&obj, flagSet)
		err := flagSet.Parse(args)
		g.Expect(err).Should(Succeed(), "parse flag succeed.")
		g.Expect(obj.QualityGate).To(Equal(false), "parse quality gate succeed.")

		err = RegisterSetup(&obj, context.Background(), nil, args)
		g.Expect(err).Should(Succeed(), "setp succeed.")
		g.Expect(obj.QualityGateRules).To(Equal(map[string]string{}), "parse quality rule succeed.")
		g.Expect(obj.Validate(base)).To(HaveLen(0), "validate succeed")
	})

	t.Run("validate rules failed", func(t *testing.T) {
		g := NewGomegaWithT(t)
		base := field.NewPath("base")

		obj := struct {
			UnitTestQuaityGateOption
		}{}
		args := []string{
			"--enable-quality-gate", "true",
		}

		flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
		RegisterFlags(&obj, flagSet)
		err := flagSet.Parse(args)
		g.Expect(err).Should(Succeed(), "parse flag succeed.")
		g.Expect(obj.QualityGate).To(Equal(true), "parse quality gate succeed.")

		args = append(args, "--quality-gate-rules", "passed-tests-rate=90", "lines-coverage=-10", "branches-coverage=166.66")
		err = RegisterSetup(&obj, context.Background(), nil, args)
		g.Expect(err).Should(Succeed(), "setp succeed.")

		expectRules := map[string]string{PassedTestsRateMetric: "90", LinesCoverageMetric: "-10", BranchesCoverageMetric: "166.66"}
		g.Expect(obj.QualityGateRules).To(Equal(expectRules), "parse quality rule succeed.")

		g.Expect(obj.Validate(base)).To(HaveLen(2), "validate succeed")
	})

}

func TestUnitTestQuaityGateOption_ValidateQualityGate(t *testing.T) {
	t.Run("ValidateQualityGate not testresult and coverage data passed", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctx := context.Background()

		obj := struct {
			UnitTestQuaityGateOption
		}{}
		obj.UnitTestQuaityGateOption.QualityGate = true
		obj.UnitTestQuaityGateOption.QualityGateRules = map[string]string{PassedTestsRateMetric: "90", LinesCoverageMetric: "80", BranchesCoverageMetric: "66.66"}
		g.Expect(obj.ValidateQualityGate(ctx, nil)).To(HaveLen(1), "unitest test data must be set.")

		testResults := &v1alpha1.UnitTestsResult{}
		g.Expect(obj.ValidateQualityGate(ctx, testResults)).To(HaveLen(0), "not testresult and coverage data.")
	})

	t.Run("ValidateQualityGate rules validate passed", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctx := context.Background()
		obj := struct {
			UnitTestQuaityGateOption
		}{}
		obj.UnitTestQuaityGateOption.QualityGate = true
		obj.UnitTestQuaityGateOption.QualityGateRules = map[string]string{PassedTestsRateMetric: "90", LinesCoverageMetric: "80", BranchesCoverageMetric: "66.66"}
		testResults := &v1alpha1.UnitTestsResult{
			TestResult: &v1alpha1.TestResult{
				PassedTestsRate: "90",
			},
			Coverage: &v1alpha1.TestCoverage{
				Lines:    "90",
				Branches: "70.15",
			},
		}
		g.Expect(obj.ValidateQualityGate(ctx, testResults)).To(HaveLen(0), "rules validate pass")
	})

	t.Run("ValidateQualityGate rules validate failed", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctx := context.Background()
		obj := struct {
			UnitTestQuaityGateOption
		}{}
		obj.UnitTestQuaityGateOption.QualityGate = true
		obj.UnitTestQuaityGateOption.QualityGateRules = map[string]string{PassedTestsRateMetric: "90", LinesCoverageMetric: "80", BranchesCoverageMetric: "66.66"}
		testResults := &v1alpha1.UnitTestsResult{
			TestResult: &v1alpha1.TestResult{
				PassedTestsRate: "50",
			},
			Coverage: &v1alpha1.TestCoverage{
				Lines:    "79.99",
				Branches: "70.15",
			},
		}
		g.Expect(obj.ValidateQualityGate(ctx, testResults)).To(HaveLen(2), "not testresult and coverage data.")
	})

	t.Run("ValidateQualityGate rules metric not set, and validate pass", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctx := context.Background()
		obj := struct {
			UnitTestQuaityGateOption
		}{}
		obj.UnitTestQuaityGateOption.QualityGate = true
		obj.UnitTestQuaityGateOption.QualityGateRules = map[string]string{}
		testResults := &v1alpha1.UnitTestsResult{
			TestResult: &v1alpha1.TestResult{
				PassedTestsRate: "90",
			},
			Coverage: &v1alpha1.TestCoverage{
				Lines:    "90",
				Branches: "70.15",
			},
		}
		g.Expect(obj.ValidateQualityGate(ctx, testResults)).To(HaveLen(0), "not testresult and coverage data.")
	})
}
