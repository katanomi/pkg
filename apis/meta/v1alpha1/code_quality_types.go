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

import (
	"time"

	authv1 "k8s.io/api/authorization/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	CodeQualityGVK            = GroupVersion.WithKind("CodeQuality")
	CodeQualityLineChartGVK   = GroupVersion.WithKind("CodeQualityLineChart")
	CodeQualityTaskMetricsGVK = GroupVersion.WithKind("CodeQualityTaskMetrics")
)

// CodeQuality object for plugin
type CodeQuality struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CodeQualitySpec `json:"spec"`
}

type CodeQualitySpec struct {
	Branches map[string]CodeQualityBranch `json:"branches"`
}

type CodeQualityBranch struct {
	HTML              string                              `json:"html"`
	Name              string                              `json:"name"`
	IsMain            bool                                `json:"isMain"`
	QualityGateStatus string                              `json:"qualityGateStatus"`
	AnalysisDate      metav1.Time                         `json:"analysisDate"`
	Metrics           map[string]CodeQualityAnalyzeMetric `json:"metrics"`
}

// CodeQualityAnalyzeMetric present CodeQualityProject analyze result
type CodeQualityAnalyzeMetric struct {
	// Value defines the value of this metric
	Value string `json:"value"`
	// Level defines the level of the value
	// +optional
	Level *string `json:"level,omitempty"`
}

// CodeQualityLineChart object for plugin
type CodeQualityLineChart struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CodeQualityLineChartSpec `json:"spec"`
}

type CodeQualityLineChartSpec struct {
	Index   []string            `json:"index"`
	Metrics map[string][]string `json:"metrics"`
}

type CodeQualityBaseOption struct {
	ProjectKey string `json:"projectKey"`
	BranchKey  string `json:"branchKey"`
}

type CodeQualityTaskOption struct {
	TaskID      string `json:"taskID"`
	ProjectKey  string `json:"projectKey"`
	Branch      string `json:"branch"`
	PullRequest string `json:"pullRequest"`
}

// CodeQualityLineChartOption code quality line chart option
// +k8s:deepcopy-gen=false
type CodeQualityLineChartOption struct {
	CodeQualityBaseOption
	StartTime      *time.Time `json:"startTime"`
	CompletionTime *time.Time `json:"completionTime"`
	Metrics        string     `json:"metrics"`
}

// CodeQualityTaskMetrics object obtains the code scan results obtained by the scan task
type CodeQualityTaskMetrics struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CodeQualityTaskMetricsSpec   `json:"spec"`
	Status            CodeQualityTaskMetricsStatus `json:"status,omitempty"`
}

// CodeQualityTaskMetricsSpec includes structured and raw data for code scans
type CodeQualityTaskMetricsSpec struct {
	// Summary is the panel data for code scans
	// +optional
	Summary CodeQualityTaskMetricsSpecSummary `json:"summary,omitempty"`
	// Task is scan task specific information
	Task CodeQualityTaskMetricsSpecTask `json:"task"`
	// Component is code scan component
	Component CodeQualityTaskMetricsSpecComponent `json:"component"`
	// Metrics is other scan result indicators are stored here
	Metrics map[string]string `json:"metrics"`
}

// CodeQualityTaskMetricsSpecSummary is the panel data for code scans
type CodeQualityTaskMetricsSpecSummary struct {
	// New Indicates the newly added data in this scan
	// +optional
	New *CodeQualityTaskMetricsSpecSummaryOverview `json:"new,omitempty"`
	// Total represents all the issues scanned
	// +optional
	Total *CodeQualityTaskMetricsSpecSummaryOverview `json:"total,omitempty"`
}

// CodeQualityTaskMetricsSpecSummaryOverview is the overview data for code scans
type CodeQualityTaskMetricsSpecSummaryOverview struct {
	// Bugs means number of bugs
	Bugs string `json:"bugs"`
	// DuplicatedLinesDensity means ratio of duplicate code
	DuplicatedLinesDensity string `json:"duplicatedLinesDensity"`
	// Vulnerabilities means number of vulnerabilities
	Vulnerabilities string `json:"vulnerabilities"`
	// CodeSmells means number of codeSmells
	CodeSmells string `json:"codeSmells"`
}

// CodeQualityTaskMetricsSpecComponent is code scan component
type CodeQualityTaskMetricsSpecComponent struct {
	// ID is component id
	ID string `json:"id"`
	// Key is component key
	Key string `json:"key"`
	// Name is component name
	Name string `json:"name"`
}

// CodeQualityTaskMetricsSpecTask is scan task specific information
type CodeQualityTaskMetricsSpecTask struct {
	// StartedAt is start time
	StartedAt string `json:"startedAt"`
	// CompletedAt is end time
	CompletedAt string `json:"executedAt"`
	// ExecutionTimeMs is consuming to execute
	ExecutionTimeMs string `json:"executionTimeMs"`
	// ID is code scan task id
	ID string `json:"id"`
	// Status is code scan task status
	Status CodeScanReportSyncReason `json:"status"`
	// AnalysisId is analysis report id
	AnalysisId string `json:"analysisId"`
}

// CodeQualityTaskMetricsStatus is a status attribute that specifies the need to integrate
type CodeQualityTaskMetricsStatus struct {
	// Reason is code scan reason
	// +optional
	Reason CodeScanReportSyncReason `json:"reason,omitempty"`
	// Status is code scan status
	// +optional
	Status corev1.ConditionStatus `json:"status,omitempty"`
}

// CodeQualityResourceAttributes returns a ResourceAttribute object to be used in a filter
func CodeQualityResourceAttributes(verb string) authv1.ResourceAttributes {
	return authv1.ResourceAttributes{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: "codequalities",
		Verb:     verb,
	}
}
