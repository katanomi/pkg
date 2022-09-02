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
	"knative.dev/pkg/apis"
)

const DefaultTestPlansPerPage = 20

var TestPlanGVK = GroupVersion.WithKind("TestPlan")
var TestPlanListGVK = GroupVersion.WithKind("TestPlanList")

// TestPlan object for plugins
type TestPlan struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec TestPlanSpec `json:"spec"`

	Status TestPlanStatus `json:"status,omitempty"`
}

// TestPlanSpec spec for TestPlan
type TestPlanSpec struct {
	// ID is the test plan id
	ID string `json:"id"`
	// Assignee is the user assigned to execute the TestPlan
	Assignee UserSpec `json:"assignee"`
	// CreatedBy is the user who created the TestPlan
	CreatedBy UserSpec `json:"createdBy"`
	// BuildRefs are the build references related to the TestPlan
	BuildRefs []TestObjectRef `json:"buildRefs"`
}

// TestObjectRef refers to a test object
type TestObjectRef struct {
	// Name is the name of the build
	Name string `json:"name"`
	// ID is the id of the build
	ID string `json:"id"`
}

// TestPlanStatus for test plan status
type TestPlanStatus struct {
	// Conditions indicates the latest status of TestPlan within a build
	apis.Conditions `json:"conditions"`
	// BuildRef is the ref of the last build of a TestPlan
	BuildRef *TestObjectRef `json:"buildRef"`
	// Executable indicates test cases under the TestPlan could be executed
	Executable bool `json:"executable"`
	// StartDate is the start date of the TestPlan within a build
	StartDate *metav1.Time `json:"startDate,omitempty"`
	// EndDate is the deadline of the TestPlan within a build
	EndDate *metav1.Time `json:"endDate,omitempty"`
	// TestBuildStatus shows the latest test plan status of a build
	TestBuildStatus TestBuildStatusInfo `json:"testBuildStatus,omitempty"`
}

// TestBuildStatusInfo for test build status
type TestBuildStatusInfo struct {
	// TotalCases is the total number of cases
	TotalCases int `json:"totalCases"`
	// Passed is the number of passed cases
	Passed int `json:"passed"`
	// Failed is the number of failed cases
	Failed int `json:"failed"`
	// BLocked is the number of blocked cases
	Blocked int `json:"blocked"`
	// Waiting is the number of waiting cases
	Waiting int `json:"waiting"`
	// PassRate is the percentage of passed cases among all cases
	PassRate float64 `json:"passRate"`
}

// TestPlanList list of TestPlans
type TestPlanList struct {
	metav1.TypeMeta `json:",inline"`
	ListMeta        `json:"metadata,omitempty"`

	Items []TestPlan `json:"items"`
}

type TestProjectOptions struct {
	// Project identity
	Project string `json:"project"`
	// TestPlanID identity
	TestPlanID string `json:"testPlanID"`
	// TestCaseID identity
	TestCaseID string `json:"testCaseID"`
	// BuildID query param
	BuildID string `json:"buildID"`
	// Search query param for listing
	Search string `json:"search"`
}

func RefFromMap(m map[string]*TestObjectRef, ID string) *TestObjectRef {
	if m == nil {
		return nil
	}
	ret, ok := m[ID]
	if ok {
		return ret
	}
	return nil
}

// TestPlanResourceAttributes returns a ResourceAttribute object to be used in a filter
func TestPlanResourceAttributes(verb string) authv1.ResourceAttributes {
	return authv1.ResourceAttributes{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: "testplans",
		Verb:     verb,
	}
}
