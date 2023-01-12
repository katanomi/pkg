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
	"fmt"
	"strconv"

	"github.com/katanomi/pkg/apis/codequality/v1alpha1"
	"github.com/katanomi/pkg/command/logger"
	"github.com/katanomi/pkg/command/validators"
	"github.com/katanomi/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// MetricType define quaity gate rules
type MetricType string

const (
	// LinesCoverageMetric coverage lines rule.
	LinesCoverageMetric = "lines-coverage"
	// BranchesCoverageMetric coverage branches rule.
	BranchesCoverageMetric = "branches-coverage"
	// PassedTestsRateMetric test result passed rate rule.
	PassedTestsRateMetric = "passed-tests-rate"
)

// UnitTestQuaityGateOption unittest quaity gate option
type UnitTestQuaityGateOption struct {
	QualityGateOption
	QualityGateRulesOption
}

// AddFlags add flags to options
func (m *UnitTestQuaityGateOption) AddFlags(flags *pflag.FlagSet) {
	m.QualityGateOption.AddFlags(flags)
}

// Setup init quality gate rules from args
func (m *UnitTestQuaityGateOption) Setup(ctx context.Context, cmd *cobra.Command, args []string) (err error) {
	return m.QualityGateRulesOption.Setup(ctx, cmd, args)
}

// Validate verify that the input rules are legal.
func (m *UnitTestQuaityGateOption) Validate(path *field.Path) (errs field.ErrorList) {
	validator := validators.NewMetric(m.QualityGateRules)
	base := path.Child("quality-gate")
	errs = append(errs, validator.ValidateFloat(base, LinesCoverageMetric, pointer.Float64(0), pointer.Float64(100))...)
	errs = append(errs, validator.ValidateFloat(base, BranchesCoverageMetric, pointer.Float64(0), pointer.Float64(100))...)
	errs = append(errs, validator.ValidateFloat(base, PassedTestsRateMetric, pointer.Float64(0), pointer.Float64(100))...)

	return
}

// ValidateQualityGate verify that the UnitTestsResult satisfy the quality gate.
func (m *UnitTestQuaityGateOption) ValidateQualityGate(ctx context.Context, testResults *v1alpha1.UnitTestsResult) (errs field.ErrorList) {
	if testResults == nil {
		errs = append(errs, field.Forbidden(field.NewPath("unittest-quality-gate"), "no data was found for unittest result."))
		return
	}

	logger := logger.NewLoggerFromContext(ctx)
	if !m.QualityGate {
		logger.Infof("==> 📢  %s quality gate disabled, skip checking.", m.QualityGateOption.FlagName)
		return
	}

	base := field.NewPath("quality-gate")
	if testResults.TestResult != nil {
		errs = append(errs, m.validateQualityGate(ctx, base, PassedTestsRateMetric, testResults.TestResult.PassedTestsRate)...)
	}

	if testResults.Coverage != nil {
		errs = append(errs, m.validateQualityGate(ctx, base, LinesCoverageMetric, testResults.Coverage.Lines)...)
		errs = append(errs, m.validateQualityGate(ctx, base, BranchesCoverageMetric, testResults.Coverage.Branches)...)
	}
	return
}

// validateQualityGate verify that the metric satisfy the quality gate.
func (m *UnitTestQuaityGateOption) validateQualityGate(ctx context.Context, base *field.Path, metric, value string) (errs field.ErrorList) {
	logger := logger.NewLoggerFromContext(ctx)

	rate, exist, err := m.GetRuleValueFloat(metric)
	if err != nil {
		errs = append(errs, field.InternalError(base.Child(metric), fmt.Errorf("parse metirc[%s] value failed. error: %s", metric, err.Error())))
		return
	}

	if !exist {
		return
	}

	testRate, err := strconv.ParseFloat(value, 64)
	if err != nil {
		errs = append(errs, field.InternalError(base.Child(metric), fmt.Errorf("parsing metric %q value %q failed: error: %s", metric, value, err.Error())))
		return
	}

	if testRate < rate {
		errs = append(errs, field.Forbidden(base.Child(metric), fmt.Sprintf("%s quality gate failed: current %.2f%% < expected %.2f%%", metric, testRate, rate)))
	} else {
		logger.Infof("==> ✅  %s quality gate passed: current %.2f%% >= expected %.2f%%", metric, testRate, rate)
	}

	return errs
}
