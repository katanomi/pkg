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
	authv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var AttributeGVK = GroupVersion.WithKind("Attribute")

// Attribute object for project management issue attribute
type Attribute struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AttributeSpec `json:"spec"`
}

// AttributeSpec for issue
type AttributeSpec struct {
	// Issue type of user self define in project
	Types []IssueType `json:"types"`

	// Issue priority of user self define in project
	Priorities []IssuePriority `json:"priorities"`

	// Issue status of user self define in project
	Statuses []AttributeStatus `json:"statuses"`
}

// AttributeStatus for issue status type
type AttributeStatus struct {
	// project attribute status name
	Name string `json:"name"`
}

// AttributeResourceAttributes returns a ResourceAttribute object to be used in a filter
func AttributeResourceAttributes(verb string) authv1.ResourceAttributes {
	return authv1.ResourceAttributes{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: "attributes",
		Verb:     verb,
	}
}
