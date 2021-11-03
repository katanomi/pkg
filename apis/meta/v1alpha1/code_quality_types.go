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
	Name              string            `json:"name"`
	IsMain            bool              `json:"isMain"`
	QualityGateStatus string            `json:"qualityGateStatus"`
	AnalysisDate      metav1.Time       `json:"analysisDate"`
	Metrics           map[string]string `json:"metrics"`
}
