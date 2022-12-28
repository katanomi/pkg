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

	"github.com/katanomi/pkg/command/io"
	"github.com/katanomi/pkg/encoding"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// UnitTestReportOption
type UnitTestReportOption struct {
	ReportPathOption
	ReportTypeOption

	ResultPathOption
}

// AddFlags add flags to options
func (m *UnitTestReportOption) AddFlags(flags *pflag.FlagSet) {
	m.ReportPathOption.AddFlags(flags)
	m.ReportTypeOption.AddFlags(flags)
	m.ResultPathOption.AddFlags(flags)
}

// Setup init quality gate rules from args
func (m *UnitTestReportOption) Setup(ctx context.Context, cmd *cobra.Command, args []string) (err error) {
	return m.ReportTypeOption.Setup(ctx, cmd, args)
}

func (m *UnitTestReportOption) Validate(path *field.Path) (errs field.ErrorList) {
	errs = append(errs, m.ReportTypeOption.Validate(path.Child("report-type"))...)
	if m.ReportType != "" {
		errs = append(errs, m.ReportPathOption.Validate(path.Child("report-path"))...)
	}
	return
}

// ParseReport parse report
func (m *UnitTestReportOption) ParseReport() (interface{}, error) {
	return m.ReportTypeOption.Parse(m.ReportPath)
}

// WriteResult save data to result path
func (c *UnitTestReportOption) WriteResult(obj interface{}) error {
	if c.ResultPath == "" {
		return nil
	}

	data := encoding.Encode(obj)
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return io.WriteFile(c.ResultPath, content, 0644)
}
