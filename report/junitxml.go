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
)

const (
	// TypeJunitXml is the type of junit-xml
	TypeJunitXml ReportType = "junit-xml"
)

// SummariesByType stores TestSummary by ReportType
type SummariesByType map[ReportType]TestSummary

// TestSummary defines summary of test report
type TestSummary struct {
	Total           int     `json:"total"`
	Passed          int     `json:"passed"`
	Failed          int     `json:"failed"`
	Error           int     `json:"error"`
	Skipped         int     `json:"skipped"`
	PassedTestsRate float64 `json:"passedTestsRate,omitempty"`
}

// GetSummaryFromJunitXml gets summary from provided junit-xml file path
func GetSummaryFromJunitXml(reportPath string) (summary *TestSummary, err error) {
	var suites []junit.Suite
	suites, err = junit.IngestDir(reportPath)
	if err != nil {
		return nil, err
	}

	summary = &TestSummary{}
	for _, suite := range suites {
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
	return summary, nil
}
