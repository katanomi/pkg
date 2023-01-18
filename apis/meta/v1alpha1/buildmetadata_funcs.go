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
	"encoding/json"
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/util/validation/field"

	ksubstitute "github.com/katanomi/pkg/substitution"
)

// AssignByGitBranch is used to assign the information from GitBranch
func (b *BuildGitBranchStatus) AssignByGitBranch(gitBranch *GitBranch) *BuildGitBranchStatus {
	if gitBranch == nil {
		return b
	}
	if b == nil {
		b = &BuildGitBranchStatus{}
	}
	b.Name = gitBranch.Name
	if gitBranch.Spec.Default != nil {
		b.Default = *gitBranch.Spec.Default
	}
	if gitBranch.Spec.Protected != nil {
		b.Protected = *gitBranch.Spec.Protected
	}
	if gitBranch.Spec.Properties != nil && gitBranch.Spec.Properties.Raw != nil {
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

// AssignByGitCommit is used to assign the information from GitCommit
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

// AssignByGitPullRequest is used to assign the information from GitPullRequest
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

func (b *BaseGitStatus) GetValWithKey(ctx context.Context, path *field.Path) map[string]string {
	if b == nil {
		b = &BaseGitStatus{}
	}
	stringReplacements := map[string]string{}
	// adds a blank to have it return empty value when referencing
	// may return a simplified value in the future
	stringReplacements[path.String()] = ""
	//
	stringReplacements[path.Child("url").String()] = b.URL
	//
	stringReplacements = ksubstitute.MergeMap(stringReplacements, b.Revision.GetValWithKey(ctx, path.Child("revision")))
	//
	stringReplacements = ksubstitute.MergeMap(stringReplacements, b.LastCommit.GetValWithKey(ctx, path.Child("lastCommit")))
	//
	stringReplacements = ksubstitute.MergeMap(stringReplacements, b.PullRequest.GetValWithKey(ctx, path.Child("pullRequest")))
	//
	stringReplacements = ksubstitute.MergeMap(stringReplacements, b.Branch.GetValWithKey(ctx, path.Child("branch")))
	//
	stringReplacements = ksubstitute.MergeMap(stringReplacements, b.Target.GetValWithKey(ctx, path.Child("target")))
	return stringReplacements
}

func (b *BuildRunGitStatus) GetValWithKey(ctx context.Context, path *field.Path) map[string]string {
	if b == nil {
		b = &BuildRunGitStatus{}
	}
	stringReplacements := map[string]string{}
	// adds a blank to have it return empty value when referencing
	// may return a simplified value in the future
	stringReplacements[path.String()] = ""
	//
	stringReplacements = ksubstitute.MergeMap(stringReplacements, b.BaseGitStatus.GetValWithKey(ctx, path))
	//
	stringReplacements[path.Child("version").String()] = b.Version
	//
	variantsMap := map[string]string{}
	for variant, version := range b.VersionVariants {
		// the key is `version` not `versionVariants`, convenient for users.
		variantsMap[path.Child("version").Child(variant).String()] = version
	}
	stringReplacements = ksubstitute.MergeMap(stringReplacements, variantsMap)
	return stringReplacements
}

func (b *BuildGitCommitStatus) GetValWithKey(ctx context.Context, path *field.Path) map[string]string {
	if b == nil {
		b = &BuildGitCommitStatus{}
	}
	stringVals := map[string]string{path.String(): b.ShortID}
	stringVals[path.Child("shortID").String()] = b.ShortID
	stringVals[path.Child("id").String()] = b.ID
	stringVals[path.Child("title").String()] = b.Title
	stringVals[path.Child("message").String()] = b.Message
	stringVals[path.Child("authorEmail").String()] = b.AuthorEmail
	stringVals[path.Child("pushedAt").String()] = ""
	if b.PushedAt != nil {
		stringVals[path.Child("pushedAt").String()] = b.PushedAt.UTC().Format(time.RFC3339)
	}
	stringVals[path.Child("webURL").String()] = b.WebURL
	return stringVals
}

func (b *BuildGitPullRequestStatus) GetValWithKey(ctx context.Context, path *field.Path) map[string]string {
	if b == nil {
		b = &BuildGitPullRequestStatus{}
	}
	stringVals := map[string]string{path.String(): b.ID}
	stringVals[path.Child("id").String()] = b.ID
	stringVals[path.Child("title").String()] = b.Title
	stringVals[path.Child("source").String()] = b.Source
	stringVals[path.Child("target").String()] = b.Target
	stringVals[path.Child("webURL").String()] = b.WebURL
	stringVals[path.Child("hasConflicts").String()] = strconv.FormatBool(b.HasConflicts)
	stringVals[path.Child("authorEmail").String()] = b.AuthorEmail
	return stringVals
}

func (b *BuildGitBranchStatus) GetValWithKey(ctx context.Context, path *field.Path) map[string]string {
	if b == nil {
		b = &BuildGitBranchStatus{}
	}
	stringVals := map[string]string{path.String(): b.Name}
	stringVals[path.Child("name").String()] = b.Name
	stringVals[path.Child("protected").String()] = strconv.FormatBool(b.Protected)
	stringVals[path.Child("default").String()] = strconv.FormatBool(b.Default)
	stringVals[path.Child("webURL").String()] = b.WebURL
	return stringVals
}
