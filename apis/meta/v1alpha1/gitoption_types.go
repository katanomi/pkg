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
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GitRepoFileOption option for get repo's file
type GitRepoFileOption struct {
	GitRepo
	// Ref commit/branch/tag name
	Ref  string `json:"ref"`
	Path string `json:"path"`
}

// GitCommitListOption option for list commit
type GitCommitListOption struct {
	GitRepo
	// Ref source branch name
	Ref string `json:"ref"`
	// Since Time query parameter, the lower bound of the time range
	Since *v1.Time `json:"since,omitempty"`
	// Until Time query parameter, the upper limit of the time range
	Until *v1.Time `json:"util,omitempty"`
}

// GitCommitOption option for get one commit by sha
type GitCommitOption struct {
	GitRepo
	GitCommitBasicInfo
}

// GitBranchOption option for list branch
type GitBranchOption struct {
	GitRepo
	Keyword string `json:"keyword"`
}

// GitPullRequestOption option for one pr by id
type GitPullRequestOption struct {
	GitRepo
	Index int `json:"Index"`
}

type GitPullRequestListOption struct {
	GitRepo
	// State indicattes pullrequest state.
	// Note that only opened and all are supported here, as enum values may vary across different tools.
	State *PullRequestState `json:"state,omitempty"`
	// Commit will filter pullrequest that just associate to this commit
	Commit string `json:"commit,omitempty"`
}

func String2PullRequestState(state string) *PullRequestState {
	if state == "" {
		return nil
	}
	return (*PullRequestState)(&state)
}
