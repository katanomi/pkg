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

	// Revision code revision used. uses a git clone format
	// refs/head/main or refs/pulls/1/head etc
	// +optional
	Revision *GitRevision `json:"revision,omitempty"`

	// LastCommit means last commit status of current build
	// +optional
	LastCommit *BuildGitCommitStatus `json:"lastCommit,omitempty"`
	// PullRequest means pull request status of current build
	// +optional
	PullRequest *BuildGitPullRequestStatus `json:"pullRequest,omitempty"`
	// Branch status of current build
	// +optional
	Branch *BuildGitBranchStatus `json:"branch,omitempty"`

	// Target branch status of current build for Pull requests
	// +optional
	Target *BuildGitBranchStatus `json:"target,omitempty"`

	// Version is the version generated for this git revision
	// +optional
	Version string `json:"version,omitempty"`

	// VersionVariants are different variants generated based on version
	// key is the name of the variant, value is the value after the variant.
	// +optional
	VersionVariants map[string]string `json:"versionVariants,omitempty"`
}

// BuildGitBranchStatus represent branch status of build run
type BuildGitBranchStatus struct {
	// Name of git branch
	// +optional
	Name string `json:"name,omitempty"`
	// Protected represent if is the protected branch
	// +optional
	Protected bool `json:"protected"`
	// Default represent if is the protected branch
	// +optional
	Default bool `json:"default"`
	// WebURL to access the branch
	// +optional
	WebURL string `json:"webURL,omitempty"`
}

type BuildGitCommitStatus struct {
	// ShortID means last commit short id
	// +optional
	ShortID string `json:"shortID,omitempty"`
	// ID represent last commit id
	// +optional
	ID string `json:"id,omitempty"`
	// Title represent last commit title
	// +optional
	Title string `json:"title,omitempty"`
	// Message of last commit
	// +optional
	Message string `json:"message,omitempty"`
	// AuthorEmail of last commit
	// +optional
	AuthorEmail string `json:"authorEmail,omitempty"`
	// PushedAt means push time of last commit
	// +optional
	PushedAt *metav1.Time `json:"pushedAt,omitempty"`
	// webURL access link of the commit
	// +optional
	WebURL string `json:"webURL,omitempty"`
}

type BuildGitPullRequestStatus struct {
	// ID is identity of pull request
	// +optional
	ID string `json:"id,omitempty"`
	// Title of pullrequest if current build is building a pull request
	// +optional
	Title string `json:"title,omitempty"`
	// Source of pullrequest if current build is building a pull request
	// +optional
	Source string `json:"source,omitempty"`
	// Target of pullrequest if current build is building a pull request
	// +optional
	Target string `json:"target,omitempty"`
	// AuthorEmail of pull request
	// +optional
	AuthorEmail string `json:"authorEmail,omitempty"`
	// WebURL to access pull request
	// +optional
	WebURL string `json:"webURL,omitempty"`
	// HasConflicts represent if has conflicts in pull request
	// +optional
	HasConflicts bool `json:"hasConflicts"`
}
