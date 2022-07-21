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
	"encoding/json"
	"strconv"
)

// AssignByGitBranch is used to assign the infomation from GitBranch
func (b *BuildGitBranchStatus) AssignByGitBranch(gitBranch *GitBranch) *BuildGitBranchStatus {
	if gitBranch == nil {
		return b
	}
	if b == nil {
		b = &BuildGitBranchStatus{}
	}
	b.Name = gitBranch.Name
	b.Default = *gitBranch.Spec.Default
	b.Protected = *gitBranch.Spec.Protected
	if gitBranch.Spec.Properties.Raw != nil {
		var content map[string]string
		json.Unmarshal(gitBranch.Spec.Properties.Raw, &content)
		b.WebURL = content["webURL"]
	}
	return b
}

// CommitProperties commit properties info
// +kubebuilder:object:generate=false
type CommitProperties struct {
	// ShortID commit short id
	ShortID string `json:"shortID"`
	// Title commit title
	Title string `json:"title"`
}

// AssignByGitCommit is used to assign the infomation from GitCommit
func (b *BuildGitCommitStatus) AssignByGitCommit(gitCommit *GitCommit) *BuildGitCommitStatus {
	if gitCommit == nil {
		return b
	}
	if b == nil {
		b = &BuildGitCommitStatus{}
	}
	if gitCommit.Spec.SHA != nil {
		b.ID = *gitCommit.Spec.SHA
	}
	if gitCommit.Spec.Message != nil {
		b.Message = *gitCommit.Spec.Message
	}
	if gitCommit.Spec.Author != nil {
		b.AuthorEmail = gitCommit.Spec.Author.Email
	}
	if gitCommit.Spec.Address != nil && gitCommit.Spec.Address.URL != nil {
		b.WebURL = gitCommit.Spec.Address.URL.String()
	}

	if gitCommit.Spec.Properties.Raw != nil {
		propertiesInfo := &CommitProperties{}
		if err := json.Unmarshal(gitCommit.Spec.Properties.Raw, propertiesInfo); err == nil {
			b.ShortID = propertiesInfo.ShortID
			b.Title = propertiesInfo.Title
		}
	}

	return b
}

// AssignByGitPullRequest is used to assign the infomation from GitPullRequest
func (b *BuildGitPullRequestStatus) AssignByGitPullRequest(gitPullRequest *GitPullRequest) *BuildGitPullRequestStatus {
	if gitPullRequest == nil {
		return b
	}
	if b == nil {
		b = &BuildGitPullRequestStatus{}
	}
	b.ID = strconv.FormatInt(gitPullRequest.Spec.Number, 10)
	b.Title = gitPullRequest.Spec.Title
	b.HasConflicts = (gitPullRequest.Spec.MergeStatus == MergeStatusCannotBeMerged)

	b.Target = gitPullRequest.Spec.Target.Name
	b.Source = gitPullRequest.Spec.Source.Name
	b.AuthorEmail = gitPullRequest.Spec.Author.Email
	if gitPullRequest.Spec.Properties.Raw != nil {
		var content map[string]string
		json.Unmarshal(gitPullRequest.Spec.Properties.Raw, &content)
		b.WebURL = content["webURL"]
	}
	return b
}
