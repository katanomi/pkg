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
	"k8s.io/apimachinery/pkg/runtime"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

var (
	GitCommitGVK     = GroupVersion.WithKind("GitCommit")
	GitCommitListGVK = GroupVersion.WithKind("GitCommitList")
)

// GitCommit object for plugin
type GitCommit struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GitCommitSpec `json:"spec"`
}

// GitCommitBasicInfo github support field is SHA & web_URL
type GitCommitBasicInfo struct {
	// SHA commit's sha
	SHA *string `json:"sha,omitempty"`
}

// GitCommitInfo github support field is SHA & web_URL
type GitCommitInfo struct {
	// SHA commit's sha
	SHA *string `json:"sha,omitempty"`
	// TODO remove this field, because is same as metadata.creationTimeStamp
	CreatedAt metav1.Time `json:"createdAt"`
	// Address for commit url for code repository web server
	Address *duckv1.Addressable `json:"address,omitempty"`
}

// GitCommitSpec spec for commit
type GitCommitSpec struct {
	GitCommitInfo
	// Coverage code coverage for test
	Coverage *float64 `json:"coverage,omitempty"`
	// Author commit author
	Author *GitUserBaseInfo `json:"author,omitempty"`
	// Committer commit committer
	Committer *GitUserBaseInfo `json:"committer,omitempty"`
	// Message commit message
	Message    *string               `json:"message,omitempty"`
	Properties *runtime.RawExtension `json:"properties,omitempty"`
}

// GitCommitList list of commits
type GitCommitList struct {
	metav1.TypeMeta `json:",inline"`
	ListMeta        `json:"metadata,omitempty"`

	Items []GitCommit `json:"items"`
}

// GitCommitResourceAttributes returns a ResourceAttribute object to be used in a filter
func GitCommitResourceAttributes(verb string) authv1.ResourceAttributes {
	return authv1.ResourceAttributes{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: "gitcommits",
		Verb:     verb,
	}
}
