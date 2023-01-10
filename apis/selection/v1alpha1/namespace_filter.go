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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// NamespaceFilter filter options for namespaces
// +k8s:deepcopy-gen=true
type NamespaceFilter struct {
	// selector is a label query over namespaces that match the filter.
	// It must match the namespaces's labels.
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`

	// a match rule to filter namespace based on Selector or Refs
	// +optional
	Filter *NamespaceFilterRule `json:"filter,omitempty"`

	// Refs is a slice of specific references for namespaces
	// +optional
	Refs []corev1.ObjectReference `json:"refs,omitempty"`
}

// NamespaceFilterRule is alias of BaseFilterRule
// +k8s:deepcopy-gen=true
type NamespaceFilterRule BaseFilterRule

// Filter is a filter for namespaces
func (n NamespaceFilterRule) Filter(namespaces []corev1.Namespace) []corev1.Namespace {
	return FilterGenericResources(BaseFilterRule(n), namespaces)
}

// Validate namespaceFilter validation method
func (n *NamespaceFilter) Validate(fld *field.Path) field.ErrorList {
	bFilter := BaseFilter{
		Selector: n.Selector,
		Filter:   (*BaseFilterRule)(n.Filter),
		Refs:     n.Refs,
	}
	return bFilter.Validate(fld)
}
