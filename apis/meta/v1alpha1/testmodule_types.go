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
	authv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var TestModuleGVK = GroupVersion.WithKind("TestModule")
var TestModuleListGVK = GroupVersion.WithKind("TestModuleList")

// TestModule object for plugins
type TestModule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec TestModuleSpec `json:"spec"`
}

// TestModuleSpec spec for TestModule
type TestModuleSpec struct {
	// ID is the test module id
	ID string `json:"id"`
	// Order is used to sort modules by ASC order
	Order int `json:"order"`
	// ParentID is the parent module ID
	ParentID string `json:"parentID"`
	// TestCases are the cases included by a module
	TestCases []TestModuleCaseRef `json:"testCases"`
}

type TestModuleCaseRef struct {
	// TestObjectRef refers to a test case
	TestObjectRef `json:"ref"`
	// Order indicates the ASC order of the object at same level
	Order int `json:"order"`
}

// TestModuleList list of TestModules
type TestModuleList struct {
	metav1.TypeMeta `json:",inline"`
	ListMeta        `json:"metadata,omitempty"`

	Items []TestModule `json:"items"`
}

// TestModuleResourceAttributes returns a ResourceAttribute object to be used in a filter
func TestModuleResourceAttributes(verb string) authv1.ResourceAttributes {
	return authv1.ResourceAttributes{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: "testmodules",
		Verb:     verb,
	}
}

func (tm *TestModule) ContainsTestCaseID(caseID string) bool {
	if tm == nil || tm.Spec.TestCases == nil {
		return false
	}

	for _, tc := range tm.Spec.TestCases {
		if tc.ID == caseID {
			return true
		}
	}
	return false
}
