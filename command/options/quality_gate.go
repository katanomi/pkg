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
	"github.com/spf13/pflag"
)

// QualityGateOption describe quality gate option
type QualityGateOption struct {
	QualityGate bool
	FlagName    string
}

// AddFlags add flags to options
func (m *QualityGateOption) AddFlags(flags *pflag.FlagSet) {
	if m.FlagName == "" {
		m.FlagName = "enable-quality-gate"
	}
	flags.BoolVar(&m.QualityGate, m.FlagName, false, `enables the quality gate`)
}
