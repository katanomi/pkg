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
	"regexp"
	"strings"

	authv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var TestCaseExecutionGVK = GroupVersion.WithKind("TestCaseExecution")
var TestCaseExecutionListGVK = GroupVersion.WithKind("TestCaseExecutionList")

// TestCaseExecutionStatus covers possible values of TestcaseExecutionStatus
type TestCaseExecutionStatus string

// Possible test case execution status below
const (
	TestcaseExecutionStatusPassed  TestCaseExecutionStatus = "passed"
	TestcaseExecutionStatusFailed  TestCaseExecutionStatus = "failed"
	TestcaseExecutionStatusBlocked TestCaseExecutionStatus = "blocked"
	TestcaseExecutionStatusWaiting TestCaseExecutionStatus = "waiting"
)

// TestCaseExecution object for plugins
type TestCaseExecution struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec TestCaseExecutionSpec `json:"spec"`
}

// TestCaseExecutionSpec spec for TestCaseExecution
type TestCaseExecutionSpec struct {
	// TestPlanID refers to the test plan including current test case
	TestPlanID string `json:"testPlanId"`

	// BuildRef refers to the build related to current test case
	// +optional
	BuildRef *TestObjectRef `json:"buildRef"`

	// Status is the execution result status
	Status TestCaseExecutionStatus `json:"status"`

	// CreatedAt is the time when test case was executed
	CreatedAt metav1.Time `json:"createdAt"`

	// CreatedBy is the user who created the TestCaseExecution
	CreatedBy UserSpec `json:"createdBy,omitempty"`
}

// TestCaseExecutionList list of TestCaseExecutions
type TestCaseExecutionList struct {
	metav1.TypeMeta `json:",inline"`
	ListMeta        `json:"metadata,omitempty"`

	Items []TestCaseExecution `json:"items"`
}

// TestCaseExecutionResourceAttributes returns a ResourceAttribute object to be used in a filter
func TestCaseExecutionResourceAttributes(verb string) authv1.ResourceAttributes {
	return authv1.ResourceAttributes{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: "testcaseexecutions",
		Verb:     verb,
	}
}

func UserSpecFromNote(note string) *UserSpec {
	if note == "" {
		return nil
	}

	reg, _ := regexp.Compile("\\[createdBy: ([\\w@.]*\\|[\\w.-]+@[\\w.-]+\\.\\w+)]")
	matches := reg.FindStringSubmatch(note)
	if len(matches) > 1 {
		splits := strings.Split(matches[1], "|")
		if len(splits) == 2 {
			return &UserSpec{
				Name:  splits[0],
				Email: splits[1],
			}
		}
	}
	return nil
}
