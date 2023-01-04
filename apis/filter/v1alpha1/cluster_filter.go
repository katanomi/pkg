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

package filter

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/validation/field"

	kvalidation "github.com/katanomi/pkg/apis/validation"
)

// ClusterFilter filter options for clusters
// +k8s:deepcopy-gen=true
type ClusterFilter struct {
	// The namespace where the clusters.clusterregistry.k8s.io exist, empty indicates the current namespace.
	// optional
	Namespace string `json:"namespace,omitempty"`

	// selector is a label query over clusters that match the filter.
	// It must match the clusters's labels.
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`

	// a match rule to filter cluster based on Selector or Refs
	// +optional
	Filter *ClusterFilterRule `json:"filter,omitempty"`

	// Refs is a slice of specific references for clusters
	// +optional
	Refs []corev1.ObjectReference `json:"refs,omitempty"`
}

// ClusterFilterRule is alias of baseFilterRule
// +k8s:deepcopy-gen=true
type ClusterFilterRule baseFilterRule

// Filter is a filter for clusters
func (n ClusterFilterRule) Filter(uClusters []unstructured.Unstructured) []corev1.ObjectReference {
	// To avoid the duplicate path prefix `object.`, the contents of the object are extracted separately.
	clustersObject := make([]map[string]interface{}, len(uClusters))
	for i, cluster := range uClusters {
		clustersObject[i] = cluster.Object
	}
	clustersObject = filterGenericResources(baseFilterRule(n), clustersObject)
	if len(clustersObject) == 0 {
		return nil
	}
	refs := make([]corev1.ObjectReference, 0)
	for _, obj := range clustersObject {
		cluster := unstructured.Unstructured{Object: obj}
		refs = append(refs, corev1.ObjectReference{
			APIVersion: cluster.GetAPIVersion(),
			Kind:       cluster.GetKind(),
			Namespace:  cluster.GetNamespace(),
			Name:       cluster.GetName(),
		})
	}
	return refs
}

// Validate clusterFilter validation method
func (n *ClusterFilter) Validate(fld *field.Path) field.ErrorList {
	errs := field.ErrorList{}

	errs = append(errs, kvalidation.ValidateItemName(n.Namespace, false, field.NewPath("namespace"))...)

	bFilter := baseFilter{
		Selector: n.Selector,
		Filter:   (*baseFilterRule)(n.Filter),
		Refs:     n.Refs,
	}
	errs = append(errs, bFilter.Validate(fld)...)

	return errs
}
