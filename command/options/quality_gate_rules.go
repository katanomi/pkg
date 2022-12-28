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
	"strconv"

	pkgargs "github.com/katanomi/pkg/command/args"
	"github.com/spf13/cobra"
)

// QualityGateRulesOption describe quality gate rules option
type QualityGateRulesOption struct {
	QualityGateRules map[string]string
	FlagName         string
}

// Setup init quality gate rules from args
func (m *QualityGateRulesOption) Setup(ctx context.Context, _ *cobra.Command, args []string) (err error) {
	if m.QualityGateRules == nil {
		m.QualityGateRules = make(map[string]string)
	}

	if m.FlagName == "" {
		m.FlagName = "quality-gate-rules"
	}

	m.QualityGateRules, err = pkgargs.GetKeyValues(ctx, args, m.FlagName)
	return err
}

// GetRuleValue get rule value
func (m *QualityGateRulesOption) GetRuleValue(key string) (value string, exist bool) {
	if len(m.QualityGateRules) == 0 {
		return "", false
	}
	value, exist = m.QualityGateRules[key]
	return
}

// GetRuleValueInt get rule value as int
func (m *QualityGateRulesOption) GetRuleValueInt(key string) (value int, exist bool, err error) {
	v, exist := m.GetRuleValue(key)
	if !exist {
		return 0, false, nil
	}
	value, err = strconv.Atoi(v)
	return value, true, err
}

// GetRuleValueFloat get rule value as float64
func (m *QualityGateRulesOption) GetRuleValueFloat(key string) (value float64, exist bool, err error) {
	v, exist := m.GetRuleValue(key)
	if !exist {
		return 0, false, nil
	}
	value, err = strconv.ParseFloat(v, 64)
	return value, true, err
}
