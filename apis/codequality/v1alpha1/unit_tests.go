/*
Copyright 2021 The Katanomi Authors.

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

package v1alpha1

// "fmt"

// UnitTestsResult unit tests results encapsulating
// coverage and number of passed, failed and skipped tests
type UnitTestsResult struct {
	// Coverage represent unit test coverage of current build
	Coverage *TestCoverage `json:"coverage"`
	// TODO: add BuildRunUnitTestStatus

	// TestResults stores a summary with the number of
	// test cases  that passed, where skipped or failed.
	// +optional
	TestResult *TestResult `json:"testResults,omitempty"`
}

// TestCoverage stores coverage data for unit tests
type TestCoverage struct {
	// Lines represent unit test coverage on lines
	Lines string `json:"lines,omitempty"`

	// Branches stores code branch coverage rate
	// valid value range is 0 to 100
	Branches string `json:"branches,omitempty"`
}

// TestResult test results aggregation
// stores the number of passed, skipped and failed test cases
// also stores an calculated passed tests rate value
type TestResult struct {
	// Passed test cases number
	// +optional
	Passed int `json:"passed"`

	// Skipped test cases number
	// +optional
	Skipped int `json:"skipped"`

	// Failed test cases number
	// adds on any errored test cases
	// +optional
	Failed int `json:"failed"`

	// PassedTestsRate rate of passed tests
	// calculated using  passed / (passed + failed) * 100
	// +optional
	PassedTestsRate string `json:"passedTestsRate"`
}
