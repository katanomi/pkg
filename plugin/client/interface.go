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
	"context"

	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"

	cloudevent "github.com/cloudevents/sdk-go/v2"
	"go.uber.org/zap"
)

// Interface base interface for plugins
type Interface interface {
	Path() string
	Setup(context.Context, *zap.SugaredLogger) error
}

// PluginRegister plugin registration methods to update IntegrationClass status
type PluginRegister interface {
	Interface
	GetIntegrationClassName() string
	// GetAddressURL Returns its own plugin access URL
	GetAddressURL() *apis.URL
	// GetWebhookURL Returns a Webhook accessible URL for external tools
	// If not supported return nil, false
	GetWebhookURL() (*apis.URL, bool)
	// GetSupportedVersions Returns a list of supported versions by the plugin
	// For SaaS platform plugins use a "online" version.
	GetSupportedVersions() []string
	// GetSecretTypes Returns all secret types supported by the plugin
	GetSecretTypes() []string
	// GetReplicationPolicyTypes return replication policy types for ClusterIntegration
	GetReplicationPolicyTypes() []string
	// GetResourceTypes Returns a list of Resource types that can be used in ClusterIntegration and Integration
	GetResourceTypes() []string
	// GetAllowEmptySecret Returns if an empty secret is allowed with IntegrationClass
	GetAllowEmptySecret() []string
}

// AuthCheck implements an authorization check method for plugins
type AuthChecker interface {
	AuthCheck(ctx context.Context, option metav1alpha1.AuthCheckOptions) (*metav1alpha1.AuthCheck, error)
}

// AuthTokenGenerator implements token generation/refresh API method
type AuthTokenGenerator interface {
	AuthToken(ctx context.Context) (*metav1alpha1.AuthToken, error)
}

// ProjectLister list project api
type ProjectLister interface {
	Interface
	ListProjects(ctx context.Context, option metav1alpha1.ListOptions) (*metav1alpha1.ProjectList, error)
}

type PluginAttributes interface {
	SetAttribute(k string, values ...string)
	GetAttribute(k string) []string
	Attributes() map[string][]string
}

// ProjectGetter list project api
type ProjectGetter interface {
	Interface
	GetProject(ctx context.Context, id string) (*metav1alpha1.Project, error)
}

// ProjectCreator create project api
type ProjectCreator interface {
	Interface
	CreateProject(ctx context.Context, project *metav1alpha1.Project) (*metav1alpha1.Project, error)
}

// RepositoryLister list repository
type RepositoryLister interface {
	Interface
	ListRepositories(ctx context.Context, params metav1alpha1.RepositoryOptions, option metav1alpha1.ListOptions) (*metav1alpha1.RepositoryList, error)
}

// ArtifactLister list artifact
type ArtifactLister interface {
	Interface
	ListArtifacts(ctx context.Context, params metav1alpha1.ArtifactOptions, option metav1alpha1.ListOptions) (*metav1alpha1.ArtifactList, error)
}

// ArtifactGetter get artifact detail
type ArtifactGetter interface {
	Interface
	GetArtifact(ctx context.Context, params metav1alpha1.ArtifactOptions) (*metav1alpha1.Artifact, error)
}

// ArtifactDeleter delete artifact
type ArtifactDeleter interface {
	Interface
	DeleteArtifact(ctx context.Context, params metav1alpha1.ArtifactOptions) error
}

// ScanImage scan image
type ScanImage interface {
	Interface
	ScanImage(ctx context.Context, params metav1alpha1.ArtifactOptions) error
}

// WebhookRegister used to register and manage webhooks
type WebhookRegister interface {
	// Use the methods below to manage webhooks in the target platform
	CreateWebhook(ctx context.Context, spec metav1alpha1.WebhookRegisterSpec, secret corev1.Secret) (metav1alpha1.WebhookRegisterStatus, error)
	UpdateWebhook(ctx context.Context, spec metav1alpha1.WebhookRegisterSpec, secret corev1.Secret) (metav1alpha1.WebhookRegisterStatus, error)
	DeleteWebhook(ctx context.Context, spec metav1alpha1.WebhookRegisterSpec, secret corev1.Secret) error
	ListWebhooks(ctx context.Context, uri apis.URL, secret corev1.Secret) ([]metav1alpha1.WebhookRegisterStatus, error)
}

// GitTriggerRegister used to register GitTrigger
// TODO: need refactor: maybe integration plugin should decided how to generate cloudevents filters
// up to now, it is not a better solution that relying on plugins to give some events type to GitTriggerReconcile.
//
//   PullRequestCloudEventFilter() CloudEventFilters
//   BranchCloudEventFilter() CloudEventFilters
//   TagCloudEventFilter() CloudEventFilters
//   WebHook() WebHook
type GitTriggerRegister interface {
	GetIntegrationClassName() string

	// cloud event type of pull request hook that will match
	PullRequestEventType() string

	// cloud event type of push hook that will match
	PushEventType() string

	// cloud event type of push hook that will match
	TagEventType() string
}

// WebhookResourceDiffer used to compare different webhook resources in order to provide
// a way to merge webhook registration requests. If not provided, the resource's URI will be directly compared
type WebhookResourceDiffer interface {
	// IsSameResource will provide two ResourceURI
	// the plugin should discern if they are the same.
	// If this method is not implemented a standard comparisons will be used
	IsSameResource(ctx context.Context, i, j metav1alpha1.ResourceURI) bool
}

// WebhookReceiver receives a webhook request with validation and transform it into a cloud event
type WebhookReceiver interface {
	Interface
	ReceiveWebhook(ctx context.Context, req *restful.Request, secret string) (cloudevent.Event, error)
}

// GitPullRequestCommentCreator create pull request comment functions
type GitPullRequestCommentCreator interface {
	Interface
	CreatePullRequestComment(ctx context.Context, option metav1alpha1.CreatePullRequestCommentPayload) (metav1alpha1.GitPullRequestNote, error)
}

// GitPullRequestCommentLister list pull request comment functions
type GitPullRequestCommentLister interface {
	Interface
	ListPullRequestComment(ctx context.Context, option metav1alpha1.GitPullRequestOption, listOption metav1alpha1.ListOptions) (metav1alpha1.GitPullRequestNoteList, error)
}

// GitPullRequestHandler list, get and create pr function
type GitPullRequestHandler interface {
	Interface
	ListGitPullRequest(ctx context.Context, option metav1alpha1.GitPullRequestListOption, listOption metav1alpha1.ListOptions) (metav1alpha1.GitPullRequestList, error)
	GetGitPullRequest(ctx context.Context, option metav1alpha1.GitPullRequestOption) (metav1alpha1.GitPullRequest, error)
	CreatePullRequest(ctx context.Context, payload metav1alpha1.CreatePullRequestPayload) (metav1alpha1.GitPullRequest, error)
}

// GitCommitGetter get git commit
type GitCommitGetter interface {
	Interface
	GetGitCommit(ctx context.Context, option metav1alpha1.GitCommitOption) (metav1alpha1.GitCommit, error)
}

// GitBranchLister List git branch
type GitBranchLister interface {
	Interface
	ListGitBranch(ctx context.Context, branchOption metav1alpha1.GitBranchOption, option metav1alpha1.ListOptions) (metav1alpha1.GitBranchList, error)
}

// GitBranchGetter get git branch
type GitBranchGetter interface {
	Interface
	GetGitBranch(ctx context.Context, repoOption metav1alpha1.GitRepo, branch string) (metav1alpha1.GitBranch, error)
}

// GitBranchCreator create git branch,github, gogs don't support create branch
type GitBranchCreator interface {
	Interface
	CreateGitBranch(ctx context.Context, payload metav1alpha1.CreateBranchPayload) (metav1alpha1.GitBranch, error)
}

// GitRepoFileGetter used to get a file content
type GitRepoFileGetter interface {
	Interface
	GetGitRepoFile(ctx context.Context, option metav1alpha1.GitRepoFileOption) (metav1alpha1.GitRepoFile, error)
}

// GitRepoFileCreator used to create a file, gogs don't support
type GitRepoFileCreator interface {
	Interface
	CreateGitRepoFile(ctx context.Context, payload metav1alpha1.CreateRepoFilePayload) (metav1alpha1.GitCommit, error)
}

// GitRepositoryLister list git repository
type GitRepositoryLister interface {
	Interface
	ListGitRepository(ctx context.Context, id, keyword string, subtype metav1alpha1.ProjectSubType, listOption metav1alpha1.ListOptions) (metav1alpha1.GitRepositoryList, error)
}

// GitRepositoryGetter get git repository
type GitRepositoryGetter interface {
	Interface
	GetGitRepository(ctx context.Context, repoOption metav1alpha1.GitRepo) (metav1alpha1.GitRepository, error)
}

// GitCommitStatusLister list git commit status
type GitCommitStatusLister interface {
	Interface
	ListGitCommitStatus(ctx context.Context, option metav1alpha1.GitCommitOption, listOption metav1alpha1.ListOptions) (metav1alpha1.GitCommitStatusList, error)
}

// GitCommitStatusCreator create git commit status
type GitCommitStatusCreator interface {
	Interface
	CreateGitCommitStatus(ctx context.Context, payload metav1alpha1.CreateCommitStatusPayload) (metav1alpha1.GitCommitStatus, error)
}

// GitRepositoryLister list git commit comment
type GitCommitCommentLister interface {
	Interface
	ListGitCommitComment(ctx context.Context, option metav1alpha1.GitCommitOption, listOption metav1alpha1.ListOptions) (metav1alpha1.GitCommitCommentList, error)
}

// GitRepositoryLister create git commit comment
type GitCommitCommentCreator interface {
	Interface
	CreateGitCommitComment(ctx context.Context, payload metav1alpha1.CreateCommitCommentPayload) (metav1alpha1.GitCommitComment, error)
}

type CodeQualityGetter interface {
	Interface
	GetCodeQuality(ctx context.Context, projectKey string) (*metav1alpha1.CodeQuality, error)
	GetCodeQualityOverviewByBranch(ctx context.Context, opt metav1alpha1.CodeQualityBaseOption) (*metav1alpha1.CodeQuality, error)
	GetCodeQualityLineCharts(ctx context.Context, opt metav1alpha1.CodeQualityLineChartOption) (*metav1alpha1.CodeQualityLineChart, error)
}

type BlobStoreLister interface {
	Interface
	ListBlobStores(ctx context.Context, listOption metav1alpha1.ListOptions) (*metav1alpha1.BlobStoreList, error)
}

// Client inteface for PluginClient, client code shoud use the interface
// as dependency
type Client interface {
	Get(ctx context.Context, baseURL *duckv1.Addressable, uri string, options ...OptionFunc) error
	Post(ctx context.Context, baseURL *duckv1.Addressable, uri string, options ...OptionFunc) error
	Put(ctx context.Context, baseURL *duckv1.Addressable, uri string, options ...OptionFunc) error
	Delete(ctx context.Context, baseURL *duckv1.Addressable, uri string, options ...OptionFunc) error
}

type ClientProjectGetter interface {
	Project(meta Meta, secret corev1.Secret) ClientProject
}
