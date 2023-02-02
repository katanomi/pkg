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
	"encoding/json"
	"regexp"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/katanomi/pkg/apis/validation"
	"github.com/tidwall/gjson"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1validation "k8s.io/apimachinery/pkg/apis/meta/v1/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// BaseFilter is the base filter struct
// Provide some general methods
// +k8s:deepcopy-gen=true
type BaseFilter struct {
	// selector is a label query over resources that match the filter.
	// It must match the resource's labels.
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`

	// a match rule to filter resources based on Selector or Refs
	// +optional
	Filter *BaseFilterRule `json:"filter,omitempty"`

	// Refs is a slice of specific references for resources
	// +optional
	Refs []corev1.ObjectReference `json:"refs,omitempty"`
}

func (p *BaseFilter) MatchObject(obj client.Object) bool {
	if p == nil {
		return true
	}
	isExcept := false
	if p.Selector != nil {
		labelSelector, _ := metav1.LabelSelectorAsSelector(p.Selector)
		if labelSelector.Matches(labels.Set(obj.GetLabels())) {
			isExcept = true
		}
	}
	objRef := metav1alpha1.GetObjectReferenceFromObject(obj,
		metav1alpha1.ObjectRefWithNamespace(),
		metav1alpha1.ObjectRefWithTypeMeta(),
	)
	if p.Refs != nil {
		for _, ref := range p.Refs {
			if metav1alpha1.IsTheSameObject(ref, objRef) {
				isExcept = true
				break
			}
		}
	}
	if !isExcept {
		return false
	}
	if p.Filter != nil && !p.Filter.MatchExact(obj) {
		return false
	}
	return true
}

// BaseFilterRule is the base filter rule
// +k8s:deepcopy-gen=true
type BaseFilterRule struct {
	// Exact filter objects by attributes
	Exact map[string]string `json:"exact"`

	//TODO: add more filter rules
}

// MatchExact match exact filter rule
func (n *BaseFilterRule) MatchExact(obj interface{}) bool {
	data, _ := json.Marshal(obj)

	for k, v := range n.Exact {
		k := getAttribute(string(data), k)
		v := getAttribute(string(data), v)

		if k == "" || v == "" || k != v {
			return false
		}
	}

	return true
}

// Validate BaseFilter validation method
func (n *BaseFilter) Validate(fld *field.Path) field.ErrorList {
	errs := field.ErrorList{}

	if n.Selector == nil && len(n.Refs) == 0 {
		errs = append(errs, field.Required(fld, "one of selector OR refs is required"))
	}
	if n.Selector != nil {
		errs = append(errs, v1validation.ValidateLabelSelector(n.Selector, fld.Child("selector"))...)
	}
	if len(n.Refs) > 0 {
		fld := fld.Child("refs")
		for i, ref := range n.Refs {
			currRef := ref
			fld := fld.Index(i)
			errs = append(errs, validation.ValidateObjectReference(&currRef, false, false, fld)...)
		}
	}

	return errs
}

// FilterGenericResources filter generic resources by BaseFilterRule
func FilterGenericResources[T any](n BaseFilterRule, objs []T) []T {
	if n.Exact == nil {
		return nil
	}
	result := make([]T, 0, len(objs))
	for _, obj := range objs {
		if n.MatchExact(obj) {
			result = append(result, obj)
		}
	}

	return result
}

// baseFilterExactRegex used to match `$(.+)` and get the content in parentheses
var baseFilterExactRegex = regexp.MustCompile(`^\$\((.+)\)$`)

// getAttribute get attribute from json data
// If the attribute not found, return attribute string
func getAttribute(data string, k string) string {
	restoreKey := restoreEscapedCharacters(k)
	result := baseFilterExactRegex.FindStringSubmatch(restoreKey)
	if len(result) == 2 {
		path := result[1]
		return gjson.Get(data, path).String()
	}

	return k
}
