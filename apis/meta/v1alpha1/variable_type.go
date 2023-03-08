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
	authv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// VariableGVK is the GVK for Variable
	VariableGVK = GroupVersion.WithKind("Variable")

	// VariableListGVK is the GVK for VariableList
	VariableListGVK = GroupVersion.WithKind("VariableList")
)

// Variable description of custom environment variables.
type Variable struct {
	// Name variable name, usually a variable JsonPath.
	Name string `json:"name"`

	// Example variable example. equal signs and semicolons cannot be included in the strength.
	Example string `json:"example,omitempty"`

	// Label variable labels, used to distinguish different types of variables, multiple labels are separated by commas.
	Label string `json:"label,omitempty"`

	// Group is used to mark display groups. Unlike labels, labels represent business groups.
	Group string `json:"group,omitempty"`
}

// VariableList variable list.
type VariableList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items contains variable list.
	Items []Variable `json:"items"`
}

// VariableResourceAttributes returns a ResourceAttribute object to be used in a filter
func VariableResourceAttributes(verb string) authv1.ResourceAttributes {
	return authv1.ResourceAttributes{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: "variables",
		Verb:     verb,
	}
}

func (v *VariableList) Filter(filters ...func(*Variable) bool) {
	if v == nil {
		return
	}

	vars := []Variable{}
	for i := range v.Items {
		if filtVariable(&v.Items[i], filters...) {
			vars = append(vars, v.Items[i])
		}
	}
	v.Items = vars
}
