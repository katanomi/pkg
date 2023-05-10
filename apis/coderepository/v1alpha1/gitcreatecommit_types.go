/*
Copyright 2023 The Katanomi Authors.

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
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// GitCreateCommitGVK is the GroupVersionKind for GitCreateCommit
	GitCreateCommitGVK = GroupVersion.WithKind("GitCreateCommit")
)

// GitCreateCommit object for plugins
type GitCreateCommit struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GitCreateCommitSpec `json:"spec"`
}

// CreateGitCommitOption option for create commit
type CreateGitCommitOption struct {
	metav1alpha1.GitRepo
	GitCreateCommit
}

// GitCreateCommitSpec spec for GitCreateCommit
type GitCreateCommitSpec struct {
	// Name of the branch to commit into.
	// To create a new branch, also provide either startBranch or startSHA or startTag.
	Branch string `json:"branch"`
	// Commit message
	Message string `json:"message"`
	// An array of action to commit as a batch.
	Actions []CreateCommitAction `json:"actions"`
	// Author commit author
	Author *metav1alpha1.GitUserBaseInfo `json:"author,omitempty"`
	// StartBranch start branch
	StartBranch string `json:"startBranch,omitempty"`
	// StartSHA start sha
	StartSHA string `json:"startSHA,omitempty"`
	// StartTag start tag
	StartTag string `json:"startTag,omitempty"`
	// Create a pull request after creating the commit.
	CreatePullRequest bool `json:"createPullRequest"`
	// Delete the source branch after merging the pull request.
	RemoveSourceBranch bool `json:"removeSourceBranch"`
}

// CreateCommitAction action for create commit
type CreateCommitAction struct {
	Action       string `json:"action,omitempty"`
	FilePath     string `json:"filePath"`
	PreviousPath string `json:"previousPath,omitempty"`
	Encoding     string `json:"encoding,omitempty"`
	Content      string `json:"content,omitempty"`
}
