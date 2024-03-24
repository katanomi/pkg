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

//go:generate mockgen -package=types -destination=../../testing/mock/github.com/katanomi/pkg/plugin/types/plugin_services.go github.com/katanomi/pkg/plugin/types PluginServices

// PluginServices contains all the services that a plugin can use
type PluginServices interface {
	TestCaseLister
	TestCaseGetter
	TestCaseExecutionLister
	TestCaseExecutionCreator

	TestModuleLister

	// ProjectUserLister

	GitPullRequestCommentCreator
	GitPullRequestCommentUpdater
	GitPullRequestCommentLister
	GitPullRequestHandler
	GitPullRequestLister
	GitPullRequestGetter
	GitPullRequestCreator

	TestPlanLister
	TestPlanGetter

	ProjectLister
	ProjectGetter
	SubtypeProjectGetter
	ProjectCreator
	// ProjectDeleter

	RepositoryLister
	// RepositoryGetter

	// WebhookRegister
	// WebhookCreator
	// WebhookUpdater
	// WebhookDeleter
	// WebhookLister
	// WebhookResourceDiffer
	// WebhookReceiver

	Interface
	// PluginRegister
	// PluginAddressable
	// DependentResourceGetter
	// AdditionalWebhookRegister
	// ResourcePathFormatter
	// PluginDisplayColumns
	// PluginAttributes
	// PluginVersionAttributes
	LivenessChecker
	Initializer
	ToolMetadataGetter

	// ImageConfigGetter

	// ScanImage

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

	// IssueLister
	// IssueGetter
	// IssueBranchLister
	// IssueBranchCreator
	// IssueBranchDeleter
	// IssueAttributeGetter

	// GitRepositoryCreator
	// GitRepositoryDeleter
	GitRepositoryLister
	GitRepositoryGetter

	CodeQualityGetter

	GitBranchLister
	GitBranchGetter
	GitBranchCreator

	ArtifactLister
	ArtifactGetter
	ArtifactDeleter
	ProjectArtifactLister
	ProjectArtifactGetter
	ProjectArtifactDeleter
	ProjectArtifactUploader
	ProjectArtifactFileGetter
	ArtifactTagDeleter
	// ArtifactTriggerRegister

	BlobStoreLister

	AuthChecker
	AuthTokenGenerator
}
