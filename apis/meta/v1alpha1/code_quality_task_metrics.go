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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	CodeQualityTaskMetricsGVK = GroupVersion.WithKind("CodeQualityTaskMetrics")
)

type CodeQualityTaskMetrics struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CodeQualityTaskMetricsSpec   `json:"spec"`
	Status            CodeQualityTaskMetricsStatus `json:"status,omitempty"`
}

type CodeQualityTaskMetricsSpec struct {
	Summary   CodeQualityTaskMetricsSpecSummary   `json:"summary,omitempty"`
	Task      CodeQualityTaskMetricsSpecTask      `json:"task"`
	Component CodeQualityTaskMetricsSpecComponent `json:"component"`
	Metrics   map[string]string                   `json:"metrics"`
}

type CodeQualityTaskMetricsSpecSummary struct {
	New   *CodeQualityTaskMetricsSpecSummaryOverview `json:"new,omitempty"`
	Total *CodeQualityTaskMetricsSpecSummaryOverview `json:"total,omitempty"`
}

type CodeQualityTaskMetricsSpecSummaryOverview struct {
	Bugs                   string `json:"bugs"`
	DuplicatedLinesDensity string `json:"duplicatedLinesDensity"`
	Vulnerabilities        string `json:"vulnerabilities"`
	CodeSmells             string `json:"codeSmells"`
}

type CodeQualityTaskMetricsSpecComponent struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

type CodeQualityTaskMetricsSpecTask struct {
	StartedAt       string                   `json:"startedAt"`
	CompletedAt     string                   `json:"executedAt"`
	ExecutionTimeMs string                   `json:"executionTimeMs"`
	ID              string                   `json:"id"`
	Status          CodeScanReportSyncReason `json:"status"`
	AnalysisId      string                   `json:"analysisId"`
}

type CodeQualityTaskMetricsStatus struct {
	Reason CodeScanReportSyncReason `json:"reason,omitempty"`
	Status corev1.ConditionStatus   `json:"status,omitempty"`
}
