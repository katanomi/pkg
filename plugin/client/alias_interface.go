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

package client

import (
	"github.com/katanomi/pkg/plugin/types"
)

// Interface base interface for plugins
type Interface = types.Interface

// PluginRegister plugin registration methods to update IntegrationClass status
type PluginRegister = types.PluginRegister

// PluginAddressable provides methods to get plugin address url
type PluginAddressable = types.PluginAddressable

// DependentResourceGetter checks and returns dependent resource references
type DependentResourceGetter = types.DependentResourceGetter

type AdditionalWebhookRegister = types.AdditionalWebhookRegister

// ResourcePathFormatter implements a formatter for resource path base on different scene
type ResourcePathFormatter = types.ResourcePathFormatter

// AuthChecker implements an authorization check method for plugins
type AuthChecker = types.AuthChecker

// AuthTokenGenerator implements token generation/refresh API method
type AuthTokenGenerator = types.AuthTokenGenerator

// ProjectLister list project api
type ProjectLister = types.ProjectLister

type PluginAttributes = types.PluginAttributes

// PluginDisplayColumns implements display columns manager.
type PluginDisplayColumns = types.PluginDisplayColumns

// PluginVersionAttributes get diff configurations for different versions.
type PluginVersionAttributes = types.PluginVersionAttributes

// ProjectGetter list project api
type ProjectGetter = types.ProjectGetter

// ProjectCreator create project api
type ProjectCreator = types.ProjectCreator

// ProjectDeleter create project api
type ProjectDeleter = types.ProjectDeleter

// RepositoryLister list repository
type RepositoryLister = types.RepositoryLister

// RepositoryGetter get repository
type RepositoryGetter = types.RepositoryGetter

// ArtifactLister list artifact
type ArtifactLister = types.ArtifactLister

// ProjectArtifactLister list project-level artifacts
type ProjectArtifactLister = types.ProjectArtifactLister

// ArtifactGetter get artifact detail
type ArtifactGetter = types.ArtifactGetter

// ProjectArtifactGetter get artifact detail
type ProjectArtifactGetter = types.ProjectArtifactGetter

// ProjectArtifactFileGetter download artifact within a project
type ProjectArtifactFileGetter = types.ProjectArtifactFileGetter

// ProjectArtifactDeleter delete artifact
type ProjectArtifactDeleter = types.ProjectArtifactDeleter

// ProjectArtifactUploader upload artifact
type ProjectArtifactUploader = types.ProjectArtifactUploader

// ArtifactDeleter delete artifact
type ArtifactDeleter = types.ArtifactDeleter

// ArtifactTagDeleter delete a specific tag of the artifact.
type ArtifactTagDeleter = types.ArtifactTagDeleter

// ScanImage scan image
type ScanImage = types.ScanImage

// ImageConfigGetter get image config
type ImageConfigGetter = types.ImageConfigGetter

// WebhookRegister used to register and manage webhooks
type WebhookRegister = types.WebhookRegister

// GitTriggerRegister used to register GitTrigger
// TODO: need refactor: maybe integration plugin should decided how to generate cloudevents filters
// up to now, it is not a better solution that relying on plugins to give some events type to GitTriggerReconcile.
//
// PullRequestCloudEventFilter() CloudEventFilters
// BranchCloudEventFilter() CloudEventFilters
// TagCloudEventFilter() CloudEventFilters
// WebHook() WebHook
type GitTriggerRegister = types.GitTriggerRegister

// WebhookResourceDiffer used to compare different webhook resources in order to provide
// a way to merge webhook registration requests. If not provided, the resource's URI will be directly compared
type WebhookResourceDiffer = types.WebhookResourceDiffer

// WebhookReceiver receives a webhook request with validation and transform it into a cloud event
type WebhookReceiver = types.WebhookReceiver

// GitPullRequestCommentCreator create pull request comment functions
type GitPullRequestCommentCreator = types.GitPullRequestCommentCreator

// GitPullRequestCommentUpdater updates pull request comment
type GitPullRequestCommentUpdater = types.GitPullRequestCommentUpdater

// GitPullRequestCommentLister list pull request comment functions
type GitPullRequestCommentLister = types.GitPullRequestCommentLister

// GitPullRequestHandler list, get and create pr function
type GitPullRequestHandler = types.GitPullRequestHandler

// GitCommitGetter get git commit
type GitCommitGetter = types.GitCommitGetter

// GitCommitCreator create git commit
type GitCommitCreator = types.GitCommitCreator

// GitCommitLister List git commit
type GitCommitLister = types.GitCommitLister

// GitBranchLister List git branch
type GitBranchLister = types.GitBranchLister

// GitBranchGetter get git branch
type GitBranchGetter = types.GitBranchGetter

// GitBranchCreator create git branch,github, gogs don't support create branch
type GitBranchCreator = types.GitBranchCreator

// GitRepoFileGetter used to get a file content
type GitRepoFileGetter = types.GitRepoFileGetter

// GitRepoFileCreator used to create a file, gogs don't support
type GitRepoFileCreator = types.GitRepoFileCreator

// GitRepositoryCreator create a git repository
type GitRepositoryCreator = types.GitRepositoryCreator

// GitRepositoryDeleter delete a git repository
type GitRepositoryDeleter = types.GitRepositoryDeleter

// GitRepositoryLister list git repository
type GitRepositoryLister = types.GitRepositoryLister

// GitRepositoryGetter get git repository
type GitRepositoryGetter = types.GitRepositoryGetter

// GitRepositoryFileTreeGetter get git repository file tree
type GitRepositoryFileTreeGetter types.GitRepositoryFileTreeGetter

// GitCommitStatusLister list git commit status
type GitCommitStatusLister = types.GitCommitStatusLister

// GitCommitStatusCreator create git commit status
type GitCommitStatusCreator = types.GitCommitStatusCreator

// GitCommitCommentLister list git commit comment
type GitCommitCommentLister = types.GitCommitCommentLister

// GitCommitCommentCreator create git commit comment
type GitCommitCommentCreator = types.GitCommitCommentCreator

// GitRepositoryTagCreator create git repository tag
type GitRepositoryTagCreator = types.GitRepositoryTagCreator

// GitRepositoryTagGetter get git repository Tag
type GitRepositoryTagGetter = types.GitRepositoryTagGetter

// GitRepositoryTagLister list git repository Tag
type GitRepositoryTagLister = types.GitRepositoryTagLister

type CodeQualityGetter = types.CodeQualityGetter

type BlobStoreLister = types.BlobStoreLister

// ArtifactTriggerRegister used to register ArtifactTrigger
type ArtifactTriggerRegister = types.ArtifactTriggerRegister

// IssueLister issue lister
type IssueLister = types.IssueLister

type IssueGetter = types.IssueGetter

type IssueBranchLister = types.IssueBranchLister

type IssueBranchCreator = types.IssueBranchCreator

type IssueBranchDeleter = types.IssueBranchDeleter

type IssueAttributeGetter = types.IssueAttributeGetter

type ProjectUserLister = types.ProjectUserLister

// TestPlanLister list test plans
type TestPlanLister = types.TestPlanLister

// TestPlanGetter get a test plan
type TestPlanGetter = types.TestPlanGetter

// TestCaseLister list test cases
type TestCaseLister = types.TestCaseLister

// TestCaseGetter get a test case
type TestCaseGetter = types.TestCaseGetter

// TestModuleLister list a test module
type TestModuleLister = types.TestModuleLister

// TestCaseExecutionLister list test case executions
type TestCaseExecutionLister = types.TestCaseExecutionLister

// TestCaseExecutionCreator create a new test case execution
type TestCaseExecutionCreator = types.TestCaseExecutionCreator

// LivenessChecker check the tool service is alive
type LivenessChecker = types.LivenessChecker

// Initializer initialize the tool service
type Initializer = types.Initializer

// ToolMetadataGetter get the version information corresponding to the address.
type ToolMetadataGetter = types.ToolMetadataGetter
