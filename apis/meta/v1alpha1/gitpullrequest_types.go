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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	GitPullRequestsGVK     = GroupVersion.WithKind("GitPullRequest")
	GitPullrequestsListGVK = GroupVersion.WithKind("GitPullRequestList")
	GitPullRequestNotesGVK = GroupVersion.WithKind("GitPullRequestNote")
)

// GitPullRequest object for plugins
type GitPullRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GitPullRequestSpec `json:"spec"`
}

// GitPullRequestSpec spec for pull request
type GitPullRequestSpec struct {
	GitRepo
	// ID num for pr in platform
	ID int64 `json:"id"`
	// Number num for pr in repo
	Number int64 `json:"num"`
	// Title pr title
	Title string `json:"title"`
	// State pr state (different between platforms)
	State string `json:"state"`
	// CreatedAt pr create time
	CreatedAt metav1.Time `json:"createdAt"`
	// UpdateAt pr latest update time
	UpdateAt *metav1.Time `json:"updateAt,omitempty"`
	// ClosedAt pr close time
	ClosedAt *metav1.Time `json:"closedAt,omitempty"`
	// Target pr target branch and repo
	Target GitBranchBaseInfo `json:"target"`
	// Source pr source branch and repo
	Source GitBranchBaseInfo `json:"source"`
	// Author pr author
	Author GitUserBaseInfo `json:"author,omitempty"`
	// MergeLog pr merge info(user and time)
	MergeLog   *GitOperateLogBaseInfo `json:"mergeLog,omitempty"`
	Properties *runtime.RawExtension  `json:"properties,omitempty"`
	// HasConflicts means source and target branch has conflict change
	HasConflicts bool `json:"hasConflicts,omitempty"`
}

// GitPullRequestList list of pr
type GitPullRequestList struct {
	metav1.TypeMeta `json:",inline"`
	ListMeta        `json:"metadata,omitempty"`

	Items []GitPullRequest `json:"items"`
}

// GitPullRequestNote note for pr
type GitPullRequestNote struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GitPullRequestNoteSpec `json:"spec"`
}

// GitPullRequestNoteSpec note's spec for pr
type GitPullRequestNoteSpec struct {
	// ID note id
	ID int `json:"id"`
	// Body note content
	Body       string                `json:"body"`
	Properties *runtime.RawExtension `json:"properties,omitempty"`
}
