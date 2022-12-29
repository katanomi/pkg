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
	"path"

	"github.com/katanomi/pkg/apis/codequality/v1alpha1"
	pkgargs "github.com/katanomi/pkg/command/args"
	"github.com/katanomi/pkg/report"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ReportPathOption describe report path option
type ReportPathOption struct {
	ReportPath string
	FlagName   string
}

// AddFlags add flags to options
func (m *ReportPathOption) AddFlags(flags *pflag.FlagSet) {
	if m.FlagName == "" {
		m.FlagName = "report-path"
	}
	flags.StringVar(&m.ReportPath, m.FlagName, "", `the path contains report`)
}

// Validate check if command is empty
func (m *ReportPathOption) Validate(path *field.Path) (errs field.ErrorList) {
	if m.ReportPath == "" {
		errs = append(errs, field.Required(path.Child("report-path"), "report-path is required"))
	}
	return errs
}

// AutomatedTestResultDefaultParser AutomatedTestResult default parser
var AutomatedTestResultDefaultParser = map[report.ReportType]report.ReportParser{report.TypeJunitXml: &report.JunitParser{}}

// ReportPathsByTypesOption is the option for multiple report types with its report paths
type ReportPathsByTypesOption struct {
	ReportPathByTypes map[report.ReportType]string
	ReportParsers     map[report.ReportType]report.ReportParser
}

// Validate verify that the type is supported.
func (m *ReportPathsByTypesOption) Validate(path *field.Path) (errs field.ErrorList) {
	for t := range m.ReportPathByTypes {
		if _, ok := m.ReportParsers[t]; !ok {
			errs = append(errs, field.TypeInvalid(path, t, "Not supported report type"))
		}
	}
	return
}

// Setup defines how to start up with report-configs option
func (r *ReportPathsByTypesOption) Setup(ctx context.Context, cmd *cobra.Command, args []string) (err error) {
	if r.ReportPathByTypes == nil {
		r.ReportPathByTypes = make(map[report.ReportType]string)
	}

	if r.ReportParsers == nil {
		r.ReportParsers = AutomatedTestResultDefaultParser
	}

	pathByTypes, _ := pkgargs.GetKeyValues(ctx, args, "report-configs")
	// report type validation
	var errs field.ErrorList
	base := field.NewPath("report-configs")
	for t, p := range pathByTypes {
		reportType := report.ReportType(t)
		if _, ok := r.ReportParsers[reportType]; !ok {
			errs = append(errs, field.TypeInvalid(base, t, "Not supported report type"))
			continue
		}
		r.ReportPathByTypes[reportType] = p
	}
	return errs.ToAggregate()
}

// TestSummariesByType gets test summaries by report type
func (r *ReportPathsByTypesOption) TestSummariesByType(parentPath string) (summaries report.SummariesByType,
	err error) {
	// Avoid using the current function when no setup is called.
	if r.ReportParsers == nil {
		r.ReportParsers = AutomatedTestResultDefaultParser
	}

	summaries = map[report.ReportType]v1alpha1.AutomatedTestResult{}
	var errs field.ErrorList
	base := field.NewPath("reportType")
	for reportType, reportPath := range r.ReportPathByTypes {
		parser, ok := r.ReportParsers[reportType]
		if !ok {
			errs = append(errs, field.Invalid(base, reportType, "parser for report type not found"))
			continue
		}

		result, err := parser.Parse(path.Join(parentPath, reportPath))
		if err != nil {
			errs = append(errs, field.InternalError(base.Child(string(reportType)),
				fmt.Errorf("failed to parse report. err: %s", err.Error())))
			continue
		}

		converter, ok := result.(report.ConvertToAutomatedTestResult)
		if !ok {
			errs = append(errs, field.TypeInvalid(base.Child(string(reportType)), converter, "ConvertToAutomatedTestResult interface is not implemented"))
		}
		summaries[reportType] = converter.ConvertToAutomatedTestResult()
	}
	return summaries, errs.ToAggregate()
}
