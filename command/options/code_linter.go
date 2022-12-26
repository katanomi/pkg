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
	"strconv"

	"github.com/katanomi/pkg/apis/codequality/v1alpha1"
	"github.com/katanomi/pkg/command/logger"
	"github.com/katanomi/pkg/command/qualitygate"
	"github.com/katanomi/pkg/command/validators"
	"github.com/katanomi/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

const (
	gateRuleIssuesCount = "issues-count"
)

// CodeLinterOption quality gate options for tasks of codeLinter type
type CodeLinterOption struct {
	v1alpha1.CodeLintResult `json:",inline"`

	QualityGateOption
	QualityGateRulesOption
}

// AddFlags add flags to options
func (c *CodeLinterOption) AddFlags(flags *pflag.FlagSet) {
	c.QualityGateOption.AddFlags(flags)
}

// Setup init quality gate rules from args
func (c *CodeLinterOption) Setup(ctx context.Context, cmd *cobra.Command, args []string) (err error) {
	c.Result = v1alpha1.Succeeded
	c.Issues = &v1alpha1.CodeLintIssues{}
	return c.QualityGateRulesOption.Setup(ctx, cmd, args)
}

// Validate validate params
func (c *CodeLinterOption) Validate(path *field.Path) (errs field.ErrorList) {
	if c.QualityGate {
		validator := validators.NewMetric(c.QualityGateRules)
		base := path.Child("quality-gate")
		errs = append(errs, validator.ValidateInt(base, gateRuleIssuesCount, pointer.Int(0), nil)...)
	}
	return errs
}

// WriteResult save quality gate result
func (c *CodeLinterOption) WriteResult(err error, w io.Writer) {
	if err != nil && !errors.Is(err, qualitygate.QualityGateCheckFailedErr) {
		// Some exceptions occurred, skipping write the results
		return
	}
	if err != nil {
		c.Result = v1alpha1.Failed
	}
	data := map[string]string{
		"result":       c.Result,
		"issues.count": strconv.Itoa(c.Issues.Count),
	}
	resultData, _ := json.Marshal(data)
	w.Write(resultData)
}

func (c *CodeLinterOption) ValidateQualityGate(ctx context.Context) (errs field.ErrorList) {
	logger := logger.NewLoggerFromContext(ctx)

	if !c.QualityGate {
		logger.Infow("==> 📢  quality gate disabled, skip checking.")
		return
	}

	base := field.NewPath("quality-gate")
	count, exist, _ := c.GetRuleValueInt(gateRuleIssuesCount)
	if !exist {
		logger.Infow("==> 📢  no rules are set, skip quality gate checking.")
		return
	}

	if c.Issues.Count > count {
		errs = append(errs, field.Forbidden(base.Child("issues-count"), fmt.Sprintf("issues count %d is greater than %d", c.Issues.Count, count)))
	} else {
		logger.Infow("==> ✅  quality gate check passed.")
	}

	return errs
}
