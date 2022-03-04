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
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

var IssueGVK = GroupVersion.WithKind("Issue")
var IssueListGVK = GroupVersion.WithKind("IssueList")

// Issue object for plugin
type Issue struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec IssueSpec `json:"spec"`
}

// IssueSpec for Issue
type IssueSpec struct {
	// Issue id
	Id string `json:"id"`

	// Issue project info
	Project IssueProject `json:"project"`

	// Address stores the webconsole address if any
	Address *duckv1.Addressable `json:"address"`

	// Issue subject
	Subject string `json:"subject"`

	// Issue type
	Type string `json:"type"`

	// Issue subtype
	SubType string `json:"subType"`

	// Issue priority
	Priority IssuePriority `json:"priority"`

	// Issue status
	Status string `json:"status"`

	// Issue assgin to someone
	Assign UserSpec `json:"assign"`

	// Issue latest update time
	UpdatedTime metav1.Time

	// Issue created user info
	Author UserSpec `json:"author"`

	// Issue description
	Description string `json:"description"`

	// Issue relate other issues
	RelateIssues []RelateIssue `json:"relateIssues"`

	// Issue subtasks
	SubTasks []RelateIssue `json:"subTasks"`

	// Issue comments
	Comments []Comment `json:"comments"`
}

type IssueProject struct {
	// Issue of project id
	Id string `json:"id"`

	// Issue of project name
	Name string `json:"name"`
}

type RelateIssue struct {
	// Relate issue subject
	Subject string `json:"subject"`

	// Relate issue access
	Access *duckv1.Addressable `json:"access"`
}

type Comment struct {
	// Issue comment by user
	User UserSpec `json:"author"`

	// Issue comment create time
	CreatedTime metav1.Time `json:"createdTime"`

	// Issue comment message
	Detail string `json:"detail"`
}

type IssueType struct {
	// Issue type id
	Id string `json:"id"`

	// Issue type name
	Name string `json:"name"`
}

type IssuePriority struct {
	// Issue priority id
	Id string `json:"id"`

	// Issue priority level
	// +optional
	Level string `json:"level"`

	// Issue priority name
	Name string `json:"name"`
}

// IssueList list of Issue
type IssueList struct {
	metav1.TypeMeta `json:",inline"`
	ListMeta        `json:"metadata,omitempty"`

	Items []Issue `json:"items"`
}

// IssueResourceAttributes returns a ResourceAttribute object to be used in a filter
func IssueResourceAttributes(verb string) authv1.ResourceAttributes {
	return authv1.ResourceAttributes{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: "issues",
		Verb:     verb,
	}
}
