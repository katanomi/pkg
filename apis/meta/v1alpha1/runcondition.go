/*
Copyright 2023 The Katanomi Authors.

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

// Package v1alpha1 contains stuct details for RunCondition
package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// RunConditionerType represent RunConditioner type
type RunConditionerType string

// RunCondition indicate the status of the conditioner
// +k8s:deepcopy-gen=true
type RunCondition struct {
	// Type is the type of RunConditioner
	Type RunConditionerType `json:"type"`
	// Status represent the status of this condition
	Status corev1.ConditionStatus `json:"status"`
	// Reason represent the reason for the current status
	// +optional
	Reason string `json:"reason,omitempty"`
	// Mesage is human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty"`
	// Properties contains specify information for this condition
	// +optional
	Properties *runtime.RawExtension `json:"properties,omitempty"`
	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime *metav1.Time `json:"lastTransitionTime,omitempty"`
}
