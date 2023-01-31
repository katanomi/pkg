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

// NamedAnalysisResults list of NamedAnalysisResult
type NamedAnalysisResults []NamedAnalysisResult

// NamedAnalysisResult adds name over integrated AnalysisResult
type NamedAnalysisResult struct {
	// Name of a specific analsysis result
	Name string `json:"name,omitempty"`

	AnalysisResult `json:",inline"`
}

// IsSameResult implements method for generic comparable usage and checking if
// lists have the same results
func (n NamedAnalysisResult) IsSameResult(y NamedAnalysisResult) bool {
	return n.Name == y.Name
}

// AnalysisResult stores the result of a code analysis performed
// storing the specific Result, the remote report address,
// the specific task id used for the project id, with optional
// metrics
type AnalysisResult struct {
	// Result of the code analysis:
	//  - Succeeded: successful code analysis with passing quality gates
	//  - Failed: failed code analysis
	//  - Canceled: canceled code analysis due to canceled task
	// +optional
	Result string `json:"result,omitempty"`
	// ReportURL for analyzed code revision
	// +optional
	ReportURL string `json:"reportURL,omitempty"`
	// TaskID of the performed analysis
	// +optional
	TaskID string `json:"taskID,omitempty"`
	// ProjectID for code analysis tool
	// +optional
	ProjectID string `json:"projectID,omitempty"`

	// Metrics for code analysis
	Metrics *AnalisysMetrics `json:"metrics,omitempty"`
}

// AnalisysMetrics b
type AnalisysMetrics struct {
	// Branch store rate metrics for current
	Branch *CodeChangeMetrics `json:"branch"`

	// Target stores target branch data
	// only populated when using a Pull Request
	// +optional
	Target *CodeChangeMetrics `json:"target,omitempty"`

	// Ratings resulted from code scan analysis
	// stores one rating type as an item
	// using the key as unique identifier for the rating type
	Ratings map[string]AnalysisRating `json:"ratings,omitempty"`

	// Languages detected in code base
	// +optional
	Languages []string `json:"languages"`

	// CodeSize based on scanned code base
	CodeSize *CodeSize `json:"codeSize"`
}

// CodeChangeMetrics common metrics for code such as:
// CoverageRate: test coverage rate based on automated tests
// DuplicationRate: code duplication rate detected by code analysis tool
type CodeChangeMetrics struct {
	// CoverageRate by tests during code analysis
	// +optional
	CoverageRate CodeChangeRates `json:"coverage"`

	// DuplicationRate discovered during code analysis
	// +optional
	DuplicationRate CodeChangeRates `json:"duplications"`
}

// CodeChangeRates stores changes in a specific rate
// for new and already existing code
// used as code duplications, test coverage, etc.
type CodeChangeRates struct {
	// New code rate
	// calculated for new code only
	// +optional
	New string `json:"new"`
	// Total code rate
	// measured over existing code
	Total string `json:"total"`
}

// AnalysisRating ratings by type
// supports adding customized types
// Rate describe the analysis rate, such as:
// - A
// - B
// - C
// - D
// - E
// IssuesCount stores the related issues number
type AnalysisRating struct {
	// Rate describe the analysis rate, such as:
	// - A
	// - B
	// - C
	// - D
	// - E
	Rate string `json:"rate"`

	// IssuesCount stores the related issues number
	// for the analyzed metric
	IssuesCount int `json:"issues"`
}

// CodeSize metrics of code base
type CodeSize struct {
	// LinesOfCode inside the project
	LinesOfCode int `json:"linesOfCode"`
}
