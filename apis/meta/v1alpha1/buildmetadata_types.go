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
)

const (
	// Build metadata key.
	BuildMetadataKey = "builds.katanomi.dev/buildrun"
)

// BuildMetaData this structure is a derivative of buildrun and is used for artifacts to record build information.
type BuildMetaData struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Status BuildMetaDataStatus `json:"status,omitempty"`
}

type BuildMetaDataStatus struct {
	// Git represent code repository status of buildrun
	// +optional
	Git *BuildRunGitStatus `json:"git,omitempty"`

	// TriggeredBy is the reason for the event trigger
	// +optional
	TriggeredBy *TriggeredBy `json:"triggeredBy,omitempty"`
}

// BuildRunGitStatus represent code repository status
type BuildRunGitStatus struct {
	// URL means git repository url of current buildrun
	// +optional
	URL string `json:"url,omitempty"`
	// LastCommit means last commit status of current build
	// +optional
	LastCommit *BuildGitCommitStatus `json:"lastCommit,omitempty"`
	// PullRequest means pull request status of current build
	// +optional
	PullRequest *BuildGitPullRequestStatus `json:"pullRequest,omitempty"`
	// Branch status of current build
	// +optional
	Branch *BuildGitBranchStatus `json:"branch,omitempty"`
}

// BuildGitBranchStatus represent branch status of build run
type BuildGitBranchStatus struct {
	// Name of git branch
	Name string `json:"name"`
	// Protected represent if is the protected branch
	Protected bool `json:"protected"`
	// Default represent if is the protected branch
	Default bool `json:"default"`
	// WebURL to access the branch
	WebURL string `json:"webURL"`
}

type BuildGitCommitStatus struct {
	// ShortID means last commit short id
	ShortID string `json:"shortID"`
	// ID represent last commit id
	ID string `json:"id"`
	// Title represent last commit title
	Title string `json:"title"`
	// Message of last commit
	Message string `json:"message"`
	// AuthorEmail of last commit
	AuthorEmail string `json:"authorEmail"`
	// PushedAt means push time of last commit
	// +optional
	PushedAt *metav1.Time `json:"pushedAt,omitempty"`
}

type BuildGitPullRequestStatus struct {
	// ID is identity of pull request
	ID string `json:"id"`
	// Title of pullrequest if current build is building a pull request
	Title string `json:"title"`
	// Source of pullrequest if current build is building a pull request
	Source string `json:"source"`
	// Target of pullrequest if current build is building a pull request
	Target string `json:"target"`
	// AuthorEmail of pull request
	AuthorEmail string `json:"authorEmail"`
	// WebURL to access pull request
	WebURL string `json:"webURL"`
	// HasConflicts represent if has conflicts in pull request
	HasConflicts bool `json:"hasConflicts"`
}
