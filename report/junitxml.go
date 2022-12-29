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

package report

import (
	"github.com/joshdk/go-junit"
	"github.com/katanomi/pkg/apis/codequality/v1alpha1"
)

const (
	// TypeJunitXml is the type of junit-xml
	TypeJunitXml ReportType = "junit-xml"
)

// SummariesByType stores TestSummary by ReportType
type SummariesByType map[ReportType]v1alpha1.AutomatedTestResult

// JunitParser junit parser
type JunitParser struct {
	// save junit testsuites
	Suites []junit.Suite
}

// Parse pase junit report.
func (p *JunitParser) Parse(path string) (testResult interface{}, err error) {
	p.Suites, err = junit.IngestFile(path)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// ConvertToTestResult convert to TestResult
func (m *JunitParser) ConvertToTestResult() v1alpha1.TestResult {
	result := v1alpha1.TestResult{}

	for _, suite := range m.Suites {
		result.Failed += suite.Totals.Error + suite.Totals.Failed
		result.Passed += suite.Totals.Passed
		result.Skipped += suite.Totals.Skipped
	}

	result.PassedTestsRate = v1alpha1.PassedTestsRate(&result)
	return result
}

// ConvertToAutomatedTestResult convert to AutomatedTestResult
func (m *JunitParser) ConvertToAutomatedTestResult() v1alpha1.AutomatedTestResult {
	summary := v1alpha1.AutomatedTestResult{}

	for _, suite := range m.Suites {
		summary.Total += suite.Totals.Tests
		summary.Error += suite.Totals.Error
		summary.Failed += suite.Totals.Failed
		summary.Skipped += suite.Totals.Skipped
		summary.Passed += suite.Totals.Passed
	}

	divider := summary.Passed + summary.Failed + summary.Error

	if divider > 0 {
		summary.PassedTestsRate = float64(summary.Passed) / float64(divider)
	}
	return summary
}
