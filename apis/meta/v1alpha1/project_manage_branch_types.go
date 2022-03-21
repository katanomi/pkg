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

var BranchGVK = GroupVersion.WithKind("Branch")
var BranchListGVK = GroupVersion.WithKind("BranchList")

// Branch object for plugin
type Branch struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec BranchSpec `json:"spec"`
}

// BranchSpec for branch
type BranchSpec struct {
	// issue info
	Issue IssueInfo `json:"issue"`

	// Branch author
	Author UserSpec `json:"author"`

	// Address stores the branch address
	Address *duckv1.Addressable `json:"address"`

	// CodeInfo stores the branch info in code repo
	CodeInfo CodeInfo `json:"codeInfo"`
}

type IssueInfo struct {
	// issue id
	Id string `json:"id"`

	// issue type
	Type string `json:"type"`
}

type CodeInfo struct {
	// Address stores the repo address
	Address *duckv1.Addressable `json:"address"`

	// code repo integration name
	IntegrationName string `json:"integrationName"`

	// code project name
	Project string `json:"project"`

	// code repo name
	Repository string `json:"repository"`

	// issue relate branch
	Branch string `json:"branch"`

	// issue relate base branch
	// +optional
	BaseBranch string `json:"baseBranch"`
}

// BranchList list of branch
type BranchList struct {
	metav1.TypeMeta `json:",inline"`
	ListMeta        `json:"metadata,omitempty"`

	Items []Branch `json:"items"`
}

// BranchResourceAttributes returns a ResourceAttribute object to be used in a filter
func BranchResourceAttributes(verb string) authv1.ResourceAttributes {
	return authv1.ResourceAttributes{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: "virtualissuebranches",
		Verb:     verb,
	}
}
