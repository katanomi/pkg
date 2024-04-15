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
	"slices"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	GitPullRequestsGVK        = GroupVersion.WithKind("GitPullRequest")
	GitPullRequestsListGVK    = GroupVersion.WithKind("GitPullRequestList")
	GitPullRequestNotesGVK    = GroupVersion.WithKind("GitPullRequestNote")
	GitPullRequestNoteListGVK = GroupVersion.WithKind("GitPullRequestNoteList")
)

// MergeStatus is the status of a merge request
type MergeStatus string

// IsValid returns true if the MergeStatus is one of the possible values.
func (t MergeStatus) IsValid() bool {
	return slices.Contains(possibleMergeStatus, t)
}

// possibleMergeStatus is a list of valid MergeStatus values.
// This list is used to validate the MergeStatus type.
var possibleMergeStatus = []MergeStatus{
	MergeStatusChecking,
	MergeStatusUnknown,
	MergeStatusCanBeMerged,
	MergeStatusCannotBeMerged,
}

const (
	// MergeStatusChecking indicates that the merge request is being checked
	MergeStatusChecking MergeStatus = "checking"
	// MergeStatusUnknown is the unknown status of the merge request
	MergeStatusUnknown MergeStatus = "unknown"
	// MergeStatusCanBeMerged indicates that the merge request can be merged
	MergeStatusCanBeMerged MergeStatus = "can_be_merged"
	// MergeStatusCannotBeMerged indicates that the merge request cannot be merged
	MergeStatusCannotBeMerged MergeStatus = "cannot_be_merged"
)

// GitPullRequest object for plugins
type GitPullRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GitPullRequestSpec `json:"spec"`
}

// PullRequestState defines the state type for a git pull request
type PullRequestState string

const (
	// PullRequestOpenedState indicates that the pull request is open and under consideration.
	PullRequestOpenedState PullRequestState = "opened"

	// PullRequestClosedState indicates that the pull request has been closed without being merged.
	// Note: This state is not supported for filtering in pull request lists.
	PullRequestClosedState PullRequestState = "closed"

	// PullRequestMergedState indicates that the pull request has been successfully merged into the target branch.
	// Note: This state is not supported for filtering in pull request lists.
	PullRequestMergedState PullRequestState = "merged"

	// PullRequestAllState is not an actual state of a pull request.
	// It is used to select pull requests of all states, including open, closed.
	PullRequestAllState PullRequestState = "all"
)

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
	State PullRequestState `json:"state"`
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
	// It is dependent on the merge_status.
	HasConflicts bool `json:"hasConflicts,omitempty"`
	// MergeStatus indicates if there is a merge conflict
	MergeStatus MergeStatus `json:"mergeStatus,omitempty"`
	// OriginMergeStatus used to store origin merge status
	OriginMergeStatus string `json:"originMergeStatus,omitempty"`
	// MergedBy indicates pr was merged by user use email
	MergedBy GitUserBaseInfo `json:"mergedBy,omitempty"`
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

// GitPullRequestNoteList note list for pr
type GitPullRequestNoteList struct {
	metav1.TypeMeta `json:",inline"`
	ListMeta        `json:"metadata,omitempty"`

	Items []GitPullRequestNote `json:"items"`
}
