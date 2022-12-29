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

// NamedAutomatedTestResults adds name over integrated UnitTestsResult
type NamedAutomatedTestResults struct {
	// Name of a specific unit tests result
	Name string `json:"name,omitempty"`

	AutomatedTestResult `json:",inline"`
}

// AutomatedTestResult automatic test result encapsulation
type AutomatedTestResult struct {
	// Total test cases number
	Total int `json:"total"`

	// Passed test cases number
	// +optional
	Passed int `json:"passed"`

	// Failed test cases number
	// +optional
	Failed int `json:"failed"`

	// Error test case number
	// +optional
	Error int `json:"error"`

	// Skipped test cases number
	// +optional
	Skipped int `json:"skipped"`

	// PassedTestsRate rate of passed tests
	// calculated using  passed / (passed + failed + error) * 100
	// +optional
	PassedTestsRate float64 `json:"passedTestsRate,omitempty"`
}
