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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var CodeQualityGVK = GroupVersion.WithKind("CodeQuality")

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

// CodeQualityAnalyzeResult present CodeQualityProject analyze result
type CodeQualityAnalyzeMetric struct {
	// Value defines the value of this metric
	Value string `json:"value"`
	// Level defines the level of the value
	// +optional
	Level *string `json:"level,omitempty"`
}

var CodeQualityLineChartGVK = GroupVersion.WithKind("CodeQualityLineChart")

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

// +k8s:deepcopy-gen=false
type CodeQualityLineChartOption struct {
	CodeQualityBaseOption
	StartTime      *time.Time `json:"startTime"`
	CompletionTime *time.Time `json:"completionTime"`
	Metrics        string     `json:"metrics"`
}
