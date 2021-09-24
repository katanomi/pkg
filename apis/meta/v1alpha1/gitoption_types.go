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

// GitRepoFileOption option for get repo's file
type GitRepoFileOption struct {
	GitRepo
	// Ref commit/branch/tag name
	Ref  string `json:"ref"`
	Path string `json:"path"`
}

// GitCommitOption option for get one commit by sha
type GitCommitOption struct {
	GitRepo
	GitCommitBasicInfo
}

// GitPullRequestOption option for one pr by id
type GitPullRequestOption struct {
	GitRepo
	Index int `json:"Index"`
}

type PullRequestState string

const (
	PullRequestOpenedState PullRequestState = "opened"
	PullRequestClosedState PullRequestState = "closed"
	PullRequestMergedState PullRequestState = "merged"
	PullRequestAllState    PullRequestState = "all"
)

type GitPullRequestListOption struct {
	GitRepo
	State *PullRequestState `json:"state,omitempty"`
}

func String2PullRequestState(state string) *PullRequestState {
	if state == "" {
		return nil
	}
	return (*PullRequestState)(&state)
}
