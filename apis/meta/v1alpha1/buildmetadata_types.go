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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Build metadata key.
	BuildMetadataKey = "builds.katanomi.dev/buildrun" // NOSONAR // ignore: "Key" detected here, make sure this is not a hard-coded credential
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

// BaseGitStatus is the base git status
type BaseGitStatus struct {
	// URL means git repository url of current buildrun
	// +optional
	URL string `json:"url,omitempty" variable:"example=https://github.com/repository/tree/main"`

	// Revision code revision used. uses a git clone format
	// refs/head/main or refs/pulls/1/head etc
	// +optional
	Revision *GitRevision `json:"revision,omitempty" variable:"example=refs/head/main"`

	// LastCommit means last commit status of current build
	// +optional
	LastCommit *BuildGitCommitStatus `json:"lastCommit,omitempty" variable:"example=3cb8901f"`
	// PullRequest means pull request status of current build
	// +optional
	PullRequest *BuildGitPullRequestStatus `json:"pullRequest,omitempty" variable:"example=1"`
	// Branch status of current build
	// +optional
	Branch *BuildGitBranchStatus `json:"branch,omitempty" variable:"example=main"`

	// Target branch status of current build for Pull requests
	// +optional
	Target *BuildGitBranchStatus `json:"target,omitempty" variable:"example=main"`
}

// BuildRunGitStatus represent code repository status
type BuildRunGitStatus struct {
	// BaseGitStatus is the base git status
	BaseGitStatus `json:",inline"`

	// VersionPhase is the phase on the versionscheme that matches this git revision.
	// +optional
	VersionPhase string `json:"versionPhase,omitempty" variable:"-"`

	// Version is the version generated for this git revision
	// +optional
	Version string `json:"version,omitempty" variable:"-"`

	// VersionVariants are different variants generated based on version
	// key is the name of the variant, value is the value after the variant.
	// +optional
	VersionVariants map[string]string `json:"versionVariants,omitempty"`
}

// BuildGitBranchStatus represent branch status of build run
type BuildGitBranchStatus struct {
	// Name of git branch
	// +optional
	Name string `json:"name,omitempty" variable:"label=default;example=main"`
	// Protected represent if is the protected branch
	// +optional
	Protected bool `json:"protected" variable:"example=true"`
	// Default represent if is the protected branch
	// +optional
	Default bool `json:"default" variable:"example=true"`
	// WebURL to access the branch
	// +optional
	WebURL string `json:"webURL,omitempty" variable:"example=https://github.com/repository/tree/main"`
}

type BuildGitCommitStatus struct {
	// ShortID means last commit short id
	// +optional
	ShortID string `json:"shortID,omitempty" variable:"label=default;example=3cb8901f"`
	// ID represent last commit id
	// +optional
	ID string `json:"id,omitempty" variable:"example=3cb8901fb325228ea27b751fcf0d6c0658a57f01"`
	// Title represent last commit title
	// +optional
	Title string `json:"title,omitempty" variable:"example=Update README.md"`
	// Message of last commit
	// +optional
	Message string `json:"message,omitempty" variable:"example=Author email when running build"`
	// AuthorEmail of last commit
	// +optional
	AuthorEmail string `json:"authorEmail,omitempty" variable:"example=joedoe@example.com"`
	// PushedAt means push time of last commit
	// +optional
	PushedAt *metav1.Time `json:"pushedAt,omitempty" variable:"example=2022-08-04T17:21:36Z"`
	// webURL access link of the commit
	// +optional
	WebURL string `json:"webURL,omitempty" variable:"example=https://github.com/repository/commit/3cb8901fb325228ea27b751fcf0d6c0658a57f01"`
	// // PullRequests means the pr list of last commit
	// // +optional
	PullRequests []BuildGitPullRequestStatus `json:"pullRequests,omitempty"`
}

type BuildGitPullRequestStatus struct {
	// ID is identity of pull request
	// +optional
	ID string `json:"id,omitempty" variable:"label=default;example=1"`
	// Title of pullrequest if current build is building a pull request
	// +optional
	Title string `json:"title,omitempty" variable:"example=Lets merge code"`
	// Source of pullrequest if current build is building a pull request
	// +optional
	Source string `json:"source,omitempty" variable:"label=default;example=branch"`
	// Target of pullrequest if current build is building a pull request
	// +optional
	Target string `json:"target,omitempty" variable:"label=default;example=main"`
	// AuthorEmail of pull request
	// +optional
	AuthorEmail string `json:"authorEmail,omitempty" variable:"example=joedoe@example.com"`
	// WebURL to access pull request
	// +optional
	WebURL string `json:"webURL,omitempty" variable:"example=https://github.com/repository/pull/1"`
	// HasConflicts represent if has conflicts in pull request
	// +optional
	HasConflicts bool `json:"hasConflicts" variable:"example=false"`
	// MergedBy indicates pr was merged by user use email
	MergedBy GitUserBaseInfo `json:"mergedBy,omitempty" variable:"example=joedoe@example.com"`
}
