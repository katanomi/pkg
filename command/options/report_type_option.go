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

	"github.com/katanomi/pkg/report"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ReportTypeOption describe report type option
type ReportTypeOption struct {
	ReportType string
	Required   bool

	FlagName      string
	SupportedType map[report.ReportType]report.ReportParser
}

// Setup perform the necessary initialization
func (m *ReportTypeOption) Setup(ctx context.Context, cmd *cobra.Command, args []string) (err error) {
	if m.SupportedType == nil {
		m.SupportedType = report.DefaultReportParsers
	}
	return nil
}

// AddFlags add flags to options
func (m *ReportTypeOption) AddFlags(flags *pflag.FlagSet) {
	if m.FlagName == "" {
		m.FlagName = "report-type"
	}
	flags.StringVar(&m.ReportType, m.FlagName, "", `the report type`)
}

// Validate verify that the type is supported.
func (m *ReportTypeOption) Validate(path *field.Path) (errs field.ErrorList) {
	if m.Required && m.ReportType == "" {
		errs = append(errs, field.Required(path, m.FlagName+` must be set`))
		return
	}

	if _, ok := m.SupportedType[report.ReportType(m.ReportType)]; !ok {
		errs = append(errs, field.TypeInvalid(path, m.ReportType, "Not supported report type"))
	}
	return
}

// Parse parse according to the message type value row, if the message type is empty, perform any operation.
func (m *ReportTypeOption) Parse(path string) (result interface{}, err error) {
	if m.ReportType == "" {
		return
	}

	parser, ok := m.SupportedType[report.ReportType(m.ReportType)]
	if !ok {
		return nil, fmt.Errorf("Not found parser by ReportType %s", m.ReportType)
	}
	return parser.Parse(path)
}
