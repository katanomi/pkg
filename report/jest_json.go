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
	"encoding/json"
	"os"

	"github.com/katanomi/pkg/apis/codequality/v1alpha1"
)

const (
	// TypeJestJson is the type of jest-json
	TypeJestJson ReportType = "jest-json"
)

// JestJsonParser jest json parser
type JestJsonParser struct {
	// failed testsuite number.
	NumFailedTestSuites int `json:"numFailedTestSuites"`
	// failed test number.
	NumFailedTests int `json:"numFailedTests"`
	// passed testsuite number.
	NumPassedTestSuites int `json:"numPassedTestSuites"`
	// passed test number.
	NumPassedTests int `json:"numPassedTests"`
	// pinding testsuite number
	NumPendingTestSuites int `json:"numPendingTestSuites"`
	// pinding test number
	NumPendingTests int `json:"numPendingTests"`
	// error testsuite number
	NumRuntimeErrorTestSuites int `json:"numRuntimeErrorTestSuites"`
	// todo test number
	NumTodoTests int `json:"numTodoTests"`
	// total testsuite number
	NumTotalTestSuites int `json:"numTotalTestSuites"`
	// total test number
	NumTotalTests int `json:"numTotalTests"`
	// testsuites run status
	Success bool `json:"success"`
}

// Parse pase jest json report.
func (m *JestJsonParser) Parse(path string) (result interface{}, err error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// ConvertToTestResult convert to TestResult
func (m *JestJsonParser) ConvertToTestResult() v1alpha1.TestResult {
	testResult := v1alpha1.TestResult{}

	testResult.Failed += m.NumFailedTests
	testResult.Passed += m.NumPassedTests + m.NumTodoTests
	testResult.Skipped += m.NumPendingTests
	testResult.PassedTestsRate = v1alpha1.PassedTestsRate(&testResult)

	return testResult
}
