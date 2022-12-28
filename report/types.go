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

import "github.com/katanomi/pkg/apis/codequality/v1alpha1"

// ReportType defines report type
type ReportType string

// SupportedTypes consists of default supported report types
var SupportedTypes = []ReportType{
	TypeJunitXml,
}

// DefaultReportParsers consists of default supported report types
var DefaultReportParsers = map[ReportType]ReportParser{
	TypeJunitXml: &JunitParser{},
	MochaJson:    &MochaJsonParser{},
	JestJson:     &JestJsonParser{},
	Lcov:         &LcovParser{},
}

// ReportParser provides an interface for parsing reports.
type ReportParser interface {
	Parse(string) (interface{}, error)
}

// ConvertToTestResult provides an interface converted to TestResult.
type ConvertToTestResult interface {
	ConvertToTestResult() v1alpha1.TestResult
}

// ConvertToTestCoverage provides an interface converted to TestCoverage.
type ConvertToTestCoverage interface {
	ConvertToTestCoverage() v1alpha1.TestCoverage
}

// ConvertToAutomatedTestResult provides an interface converted to AutomatedTestResult.
type ConvertToAutomatedTestResult interface {
	ConvertToAutomatedTestResult() v1alpha1.AutomatedTestResult
}
