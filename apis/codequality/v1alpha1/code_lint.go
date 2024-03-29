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

package v1alpha1

// NamedCodeLintResults list of NamedCodeLintResult
type NamedCodeLintResults []NamedCodeLintResult

// NamedCodeLintResult adds name over integrated CodeLintResult
type NamedCodeLintResult struct {
	// Name of a specific lint result
	Name string `json:"name,omitempty"`

	CodeLintResult `json:",inline"`
}

// IsSameResult implements method for generic comparable usage and checking if
// lists have the same results
func (n NamedCodeLintResult) IsSameResult(y NamedCodeLintResult) bool {
	return n.Name == y.Name
}

// CodeLintResult stores code linting results
type CodeLintResult struct {
	// Result for the linting process
	//  - Succeeded: successful code linting with passing quality gates
	//  - Failed: failed code linting
	//  - Canceled: canceled code linting due to canceled task
	Result string `json:"result"`

	// Issues found during linting process
	// +optional
	Issues *CodeLintIssues `json:"issues,omitempty"`
}

// CodeLintIssues issues found during linting
// count stores the total number of issues detected
type CodeLintIssues struct {

	// Count the number of detected issues
	Count int `json:"count"`
}
