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

import (
	// "github.com/katanomi/pkg/apis/meta/v1alpha1"
	// "github.com/katanomi/builds"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	CodeQualityTaskMetricsGVK = GroupVersion.WithKind("CodeQualityTaskMetrics")
)

type CodeQualityTaskMetrics struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CodeQualityTaskMetricsSpec `json:"spec,omitempty"`
}

type CodeQualityTaskMetricsSpec struct {
	Summary  CodeScanStatusSummary `json:"summary"`
	TaskInfo CodeQualityTaskInfo   `json:"codeQualityProjects"`
}

type CodeScanStatusSummary struct {
	NewBugs                   string `json:"newBugs"`
	NewDuplicatedLinesDensity string `json:"newDuplicatedLinesDensity"`
	NewVulnerabilities        string `json:"newVulnerabilities"`
	NewCodeSmells             string `json:"newCodeSmells"`
	Bugs                      string `json:"Bugs"`
	DuplicatedLinesDensity    string `json:"DuplicatedLinesDensity"`
	Vulnerabilities           string `json:"Vulnerabilities"`
	CodeSmells                string `json:"CodeSmells"`
}

type CodeQualityTaskInfo struct {
	TaskID          string            `json:"taskID"`
	Status          string            `json:"status"`
	AnalysisId      string            `json:"analysisId"`
	ComponentId     string            `json:"componentId"`
	ComponentKey    string            `json:"componentKey"`
	ComponentName   string            `json:"componentName"`
	StartedAt       string            `json:"startedAt"`
	ExecutedAt      string            `json:"executedAt"`
	ExecutionTimeMs string            `json:"executionTimeMs"`
	Branch          string            `json:"branch"`
	BranchType      string            `json:"branchType"`
	Metrics         map[string]string `json:"metrics"`
}
