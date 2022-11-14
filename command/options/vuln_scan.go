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
	"encoding/json"
	"errors"
	"fmt"
	"io"

	securityv1alpha1 "github.com/katanomi/pkg/apis/security/v1alpha1"
	"github.com/katanomi/pkg/command/logger"
	"github.com/katanomi/pkg/command/qualitygate"
	"github.com/katanomi/pkg/command/validators"
	"github.com/katanomi/pkg/encoding"
	"github.com/katanomi/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

const (
	gateRuleVulnSeverity = "severity"
	gateRuleVulnScore    = "score"
)

// VulnScanOption quality gate options for tasks of vulnerability type
type VulnScanOption struct {
	securityv1alpha1.VulnScanResult `json:",inline"`

	QualityGateOption
	QualityGateRulesOption
}

// AddFlags add flags to options
func (c *VulnScanOption) AddFlags(flags *pflag.FlagSet) {
	c.QualityGateOption.AddFlags(flags)
}

// Setup init quality gate rules from args
func (c *VulnScanOption) Setup(ctx context.Context, cmd *cobra.Command, args []string) (err error) {
	c.Result = "Successed"
	c.Targets = make([]securityv1alpha1.VulnScanTarget, 0)
	return c.QualityGateRulesOption.Setup(ctx, cmd, args)
}

// Validate validate params
func (c *VulnScanOption) Validate(path *field.Path) (errs field.ErrorList) {
	if c.QualityGate {
		validator := validators.NewMetric(c.QualityGateRules)
		base := path.Child("quality-gate")
		errs = append(errs, validator.ValidateFloat(base, gateRuleVulnScore, pointer.Float64(0), nil)...)

		availableSeverities := make([]string, 0, len(securityv1alpha1.AvailableVulnSeverities))
		for _, item := range securityv1alpha1.AvailableVulnSeverities {
			availableSeverities = append(availableSeverities, string(item))
		}
		errs = append(errs, validator.ValidateStringEnums(base, gateRuleVulnSeverity, availableSeverities...)...)
	}
	return errs
}

// WriteResult save quality gate result
func (c *VulnScanOption) WriteResult(err error, w io.Writer) {
	if err != nil && !errors.Is(err, qualitygate.QualityGateCheckFailedErr) {
		// Some exceptions occurred, skipping write the results
		return
	}
	if err != nil {
		c.Result = "Failed"
	}
	resultData, _ := json.Marshal(map[string]string{
		"result": c.Result,
	})
	w.Write(resultData)
}

// WriteMetricsResult save metrics result
func (c *VulnScanOption) WriteMetricsResult(err error, w io.Writer) {
	if err != nil && !errors.Is(err, qualitygate.QualityGateCheckFailedErr) {
		// Some exceptions occurred, skipping write the results
		return
	}
	result := c.VulnScanResult
	if len(result.Targets) > 3 {
		result.Targets = result.Targets[:3]
	}
	data := encoding.NewJsonPath().Encode(map[string]interface{}{
		"targets": result.ToVulnScanResultShadow().Targets,
	})
	resultData, _ := json.Marshal(data)
	w.Write(resultData)
}

func (c *VulnScanOption) ValidateQualityGate(ctx context.Context) (errs field.ErrorList) {
	logger := logger.NewLoggerFromContext(ctx)

	if !c.QualityGate {
		logger.Infow("==> 📢  quality gate disabled, skip checking.")
		return
	}

	base := field.NewPath("quality-gate")
	scoreGate, existScoreGate, _ := c.GetRuleValueFloat(gateRuleVulnScore)
	severityGate, existSeverityGate := c.GetRuleValue(gateRuleVulnSeverity)
	if !existScoreGate && !existSeverityGate {
		logger.Infow("==> 📢  no rules are set, skip quality gate checking.")
		return
	}

	maxScore := 0.0
	maxSeverity := ""
	for _, target := range c.Targets {
		if target.Cvss.Score > maxScore {
			maxScore = target.Cvss.Score
			maxSeverity = target.Cvss.Severity
		}
	}

	if existScoreGate && maxScore >= scoreGate {
		errs = append(errs, field.Forbidden(base.Child("score"), fmt.Sprintf("vuln score should be less than %g(current %g)", scoreGate, maxScore)))
	}

	if existSeverityGate {
		for _, item := range securityv1alpha1.AvailableVulnSeverities {
			if string(item) == maxSeverity {
				errs = append(errs, field.Forbidden(base.Child("score"), fmt.Sprintf("vuln severity should be less than %s(current %s)", severityGate, maxSeverity)))
				break
			}
			if string(item) == severityGate {
				break
			}
		}
	}

	if len(errs) == 0 {
		logger.Infow("==> ✅  quality gate check passed.")
	}

	return errs
}
