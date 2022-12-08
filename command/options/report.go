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

	pkgargs "github.com/katanomi/pkg/command/args"
	"github.com/katanomi/pkg/report"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ReportPathOption describe report path option
type ReportPathOption struct {
	ReportPath string
}

// AddFlags add flags to options
func (m *ReportPathOption) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&m.ReportPath, "report-path", "", `the path contains report`)
}

// Validate check if command is empty
func (m *ReportPathOption) Validate(path *field.Path) (errs field.ErrorList) {
	if m.ReportPath == "" {
		errs = append(errs, field.Required(path.Child("report-path"), "report-path is required"))
	}
	return errs
}

// ReportPathsByTypesOption is the option for multiple report types with its report paths
type ReportPathsByTypesOption struct {
	ReportPathByTypes map[report.ReportType]string
}

// Setup defines how to start up with report-configs option
func (r *ReportPathsByTypesOption) Setup(ctx context.Context, cmd *cobra.Command, args []string) (err error) {
	if r.ReportPathByTypes == nil {
		r.ReportPathByTypes = make(map[report.ReportType]string)
	}
	pathByTypes, _ := pkgargs.GetKeyValues(ctx, args, "report-configs")
	// report type validation
	var errs field.ErrorList
	base := field.NewPath("report-configs")
	for t, p := range pathByTypes {
		reportType := report.ReportType(t)
		if !lo.Contains(report.SupportedTypes, reportType) {
			errs = append(errs, field.TypeInvalid(base, t, "Not support report type"))
			continue
		}
		r.ReportPathByTypes[reportType] = p
	}
	return errs.ToAggregate()
}

// TestSummariesByType gets test summaries by report type
func (r *ReportPathsByTypesOption) TestSummariesByType(parentPath string) (summaries report.SummariesByType,
	err error) {
	summaries = map[report.ReportType]report.TestSummary{}
	var errs field.ErrorList
	base := field.NewPath("report-configs")
	for reportType, reportPath := range r.ReportPathByTypes {
		switch reportType {
		// Add more type in the future...
		case report.TypeJunitXml:
			summary, err := report.GetSummaryFromJunitXml(path.Join(parentPath, reportPath))
			if err != nil {
				errs = append(errs, field.InternalError(base.Child(string(reportType)),
					fmt.Errorf("GetSummaryFromJunitXml err: %s",
						err.Error())))
				continue
			}
			summaries[report.TypeJunitXml] = *summary
		default:
			// do nothing...
		}
	}
	return summaries, errs.ToAggregate()
}
