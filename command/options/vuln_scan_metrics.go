/*
Copyright 2024 The Katanomi Authors.

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
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// VulnScanMetricsOption describe quality gate option
type VulnScanMetricsOption struct {
	// FlagName is the name of the flag
	FlagName string

	// ResultLimit is the maximum number of metrics results to write
	ResultLimit int
}

// AddFlags add flags to options
func (m *VulnScanMetricsOption) AddFlags(flags *pflag.FlagSet) {
	if m.FlagName == "" {
		m.FlagName = "metrics-result-limit"
	}
	flags.IntVar(&m.ResultLimit, m.FlagName, 3, `The maximum number of metrics results to write.`)
}

// Validate check if result limit is greater than 0
func (m *VulnScanMetricsOption) Validate(path *field.Path) (errs field.ErrorList) {
	if m.ResultLimit <= 0 {
		base := path.Child("metrics-result-limit")
		errs = append(errs, field.Invalid(base, m.ResultLimit, "result limit should be greater than 0"))
	}
	return
}
