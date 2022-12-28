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
	// Lcov is the type of lcov
	MochaJson ReportType = "mocha-json"
)

// MochaJsonParser mocha json report parser
type MochaJsonParser struct {
	Stats MochaJsonStats `json:"stats"`
}

// MochaJsonStats unit test data structures for Mocha json.
type MochaJsonStats struct {
	// total testsuite
	Suites int `json:"suites"`
	// total test
	Tests int `json:"tests"`
	// passed test number
	Passes int `json:"passes"`
	// pending test number
	Pending int `json:"pending"`
	// failed test number
	Failures int `json:"failures"`
	// other test number
	Other int `json:"other"`
	// skipped test number
	Skipped int `json:"skipped"`
}

// Parse parse mocha json report.
func (m *MochaJsonParser) Parse(path string) (interface{}, error) {
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
func (m *MochaJsonParser) ConvertToTestResult() v1alpha1.TestResult {
	testResult := v1alpha1.TestResult{}
	testResult.Failed += m.Stats.Failures
	testResult.Passed += m.Stats.Passes
	testResult.Skipped += m.Stats.Pending + m.Stats.Skipped
	testResult.PassedTestsRate = v1alpha1.PassedTestsRate(&testResult)
	return testResult
}
