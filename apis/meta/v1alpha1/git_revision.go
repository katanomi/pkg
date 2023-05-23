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
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var (
	// GitRevisionGVK is the GVK for GitRevision
	GitRevisionGVK = GroupVersion.WithKind("GitRevision")
	// GitRevisionListGVK is the GVK for GitRevisionList
	GitRevisionListGVK = GroupVersion.WithKind("GitRevisionList")
)

// GitRevisionList is a list of GitRevision
type GitRevisionList struct {
	metav1.TypeMeta `json:",inline"`
	ListMeta        `json:"metadata,omitempty"`

	Items []GitRevisionDetail `json:"items"`
}

// GitRevisionDetail is the detail of a GitRevision
type GitRevisionDetail struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GitRevisionSpec `json:"spec"`
}

type GitRevisionSpec struct {
	BaseGitStatus `json:",inline"`

	// Attributes stores additional attributes for the revision
	// +optional
	Attributes map[string]string `json:"attribute,omitempty"`
}

// GitRevisionType revision types used in scm
type GitRevisionType string

const (
	// GitRevisionTypePullRequest pull request type. Uses `refs/<type>/<pr id>/head` as reference.
	// <type> depends on the vendor. Gihub uses `pulls`, Gitlab uses `merge-requests`, etc.
	GitRevisionTypePullRequest = "PullRequest"
	// GitRevisionTypeBranch branch type. Generally uses `refs/head/<branch>`
	GitRevisionTypeBranch = "Branch"
	// GitRevisionTypeTag specific tag name. Uses `refs/tag/<tag>` as reference.
	GitRevisionTypeTag = "Tag"
	// GitRevisionTypeCommit specific commit
	GitRevisionTypeCommit = "Commit"
)

// GitRevision stores revision data from git provider
type GitRevision struct {
	// Raw revision in git clone format
	// refs/head/main or refs/pulls/1/head etc
	// +optional
	Raw string `json:"raw,omitempty" variable:"example=refs/head/main"`

	// Type stores the type of revision:
	// Branch, PullRequest, Tag, or Commit
	// +optional
	Type GitRevisionType `json:"type,omitempty" variable:"example=Branch"`

	// ID for the specific revision type:
	// Branch: branch name
	// PullRequest: Pull request ID
	// Tag: tag name
	// Commit: commit short ID
	// +optional
	ID string `json:"id,omitempty" variable:"example=1"`

	// RevisionPusher for record user email of revision push
	// +optional
	RevisionPusher string `json:"revisionPusher,omitempty"`
}

// GetValWithKey returns the list of keys and values to support variable substitution
func (rev *GitRevision) GetValWithKey(ctx context.Context, path *field.Path) (values map[string]string) {
	if rev == nil {
		rev = &GitRevision{}
	}
	values = map[string]string{
		path.String():               rev.Raw,
		path.Child("raw").String():  rev.Raw,
		path.Child("type").String(): string(rev.Type),
		path.Child("id").String():   rev.ID,
	}
	return
}
