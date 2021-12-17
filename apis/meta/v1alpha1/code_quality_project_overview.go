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

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

var (
	CodeQualityProjectOverviewGVK = GroupVersion.WithKind("CodeQualityProjectOverview")
)

type CodeQualityProjectOverview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CodeQualityProjectOverviewSpec `json:"spec"`
}

type CodeQualityProjectOverviewSpec struct {
	Summary             CodeQualityProjectOverviewSpecSummary `json:"summary"`
	CodeQualityProjects []CodeQualityProjectOverviewComponent `json:"codeQualityProjects"`
}

type CodeQualityProjectOverviewSpecSummary struct {
	Error int `json:"error"`
	OK    int `json:"ok"`
	Warn  int `json:"warn"`
	Total int `json:"total"`
}

type CodeQualityProjectOverviewComponent struct {
	Branch                 string                              `json:"branch"`
	CodeRepositoryFullName string                              `json:"codeRepositoryFullName"`
	CodeRepositoryLink     string                              `json:"codeRepositoryLink"`
	LastAnalysisTime       int                                 `json:"lastAnalysisTime"`
	Name                   string                              `json:"name"`
	QualityGateStatus      string                              `json:"qualityGateStatus"`
	SonarQubeLink          string                              `json:"sonarQubeLink"`
	Metrics                map[string]CodeQualityAnalyzeMetric `json:"metrics"`
}
