/*
Copyright 2024 The Katanomi Authors.

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

package types

//go:generate mockgen -package=types -destination=../../testing/mock/github.com/katanomi/pkg/plugin/types/git_plugin_client_set.go github.com/katanomi/pkg/plugin/types GitPluginClientSet

// GitPluginClientSet is a set of interfaces that a Git plugin client should implement
type GitPluginClientSet interface {
	GitPullRequestCommentCreator
	GitPullRequestCommentUpdater
	GitPullRequestCommentLister
	GitPullRequestHandler
	GitPullRequestLister
	GitPullRequestGetter
	GitPullRequestCreator

	GitCommitStatusLister
	GitCommitStatusCreator
	GitCommitCommentLister
	GitCommitCommentCreator

	// GitRepositoryTagCreator
	GitRepositoryTagGetter
	GitRepositoryTagLister

	GitCommitGetter
	GitCommitCreator
	GitCommitLister

	GitRepoFileGetter
	GitRepoFileCreator
	GitRepositoryFileTreeGetter

	GitRepositoryLister
	GitRepositoryGetter

	GitBranchLister
	GitBranchGetter
	GitBranchCreator
}
