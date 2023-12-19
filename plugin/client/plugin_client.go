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
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client/base"
	corev1 "k8s.io/api/core/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// PluginClient client for plugins
type PluginClient struct {
	*base.PluginClient

	// IntegrationClassName is the name of the integration class
	// +optional
	IntegrationClassName string
}

// NewPluginClient creates a new plugin client
func NewPluginClient(opts ...BuildOptions) *PluginClient {
	return &PluginClient{
		PluginClient: base.NewPluginClient(opts...),
	}
}

// Make sure that PluginClient implements the Client interface
var _ Client = &PluginClient{}

// Clone shallow clone the plugin client
// used to update some fields without changing the original
func (p *PluginClient) Clone() *PluginClient {
	if p == nil {
		return nil
	}
	clone := *p
	clone.PluginClient = p.PluginClient.Clone()
	return &clone
}

func (p *PluginClient) WithMeta(meta Meta) *PluginClient {
	p.PluginClient = p.PluginClient.WithMeta(meta)
	return p
}

func (p *PluginClient) WithSecret(secret corev1.Secret) *PluginClient {
	p.PluginClient = p.PluginClient.WithSecret(secret)
	return p
}

func (p *PluginClient) WithClassAddress(classAddress *duckv1.Addressable) *PluginClient {
	p.PluginClient = p.PluginClient.WithClassAddress(classAddress)
	return p
}

func (p *PluginClient) WithRequestOptions(opts ...OptionFunc) *PluginClient {
	clone := p.Clone()
	clone.PluginClient = clone.PluginClient.WithRequestOptions(opts...)
	return clone
}

func (p *PluginClient) WithIntegrationClassName(integrationClassName string) *PluginClient {
	p.IntegrationClassName = integrationClassName
	return p
}

// Secret provides a secret to be assigned to the request in the header
func (p *PluginClient) Secret(secret corev1.Secret) OptionFunc {
	return SecretOpts(secret)
}

// Meta provides metadata for the request
func (p *PluginClient) Meta(meta Meta) OptionFunc {
	return MetaOpts(meta)
}

// ListOptions options for lists
func (p *PluginClient) ListOptions(opts metav1alpha1.ListOptions) OptionFunc {
	return ListOpts(opts)
}

// Query query parameters for the request
func (p *PluginClient) Query(params map[string]string) OptionFunc {
	return QueryOpts(params)
}

// Body request body
func (p *PluginClient) Body(body interface{}) OptionFunc {
	return BodyOpts(body)
}

// Dest request result automatically marshalled into object
func (p *PluginClient) Dest(dest interface{}) OptionFunc {
	return ResultOpts(dest)
}

// Error error response object
func (p *PluginClient) Error(err interface{}) OptionFunc {
	return ErrorOpts(err)
}

// Header sets a header
func (p *PluginClient) Header(key, value string) OptionFunc {
	return HeaderOpts(key, value)
}

// GitPluginClient convert PluginClient to GitPluginClient
func (p *PluginClient) GitPluginClient() *GitPluginClient {
	return &GitPluginClient{PluginClient: p}
}

// Auth provides an auth methods for clients
func (p *PluginClient) Auth(meta Meta, secret corev1.Secret) ClientAuth {
	clone := p.Clone().WithMeta(meta).WithSecret(secret)

	return newAuthClient(clone)
}

// NewAuth provides an auth methods for clients
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewAuth() ClientAuth {
	return newAuthClient(p)
}

// Project get project client
func (p *PluginClient) Project(meta Meta, secret corev1.Secret) ClientProject {
	clone := p.Clone().WithMeta(meta).WithSecret(secret)

	return newProject(clone)
}

// NewProject get project client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewProject() ClientProject {
	return newProject(p)
}

// Repository get Repository client
func (p *PluginClient) Repository(meta Meta, secret corev1.Secret) ClientRepository {
	clone := p.Clone().WithMeta(meta).WithSecret(secret)

	return newRepository(clone)
}

// NewRepository get Repository client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewRepository() ClientRepository {
	return newRepository(p)
}

// Artifact get Artifact client
func (p *PluginClient) Artifact(meta Meta, secret corev1.Secret) ClientArtifact {
	clone := p.Clone().WithMeta(meta).WithSecret(secret)

	return newArtifact(clone)
}

// ProjectArtifact get ProjectArtifact client
func (p *PluginClient) ProjectArtifact(meta Meta, secret corev1.Secret) ClientProjectArtifact {
	clone := p.Clone().WithMeta(meta).WithSecret(secret)

	return newProjectArtifact(clone)
}

// NewProjectArtifact get ProjectArtifact client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewProjectArtifact() ClientProjectArtifact {
	return newProjectArtifact(p)
}

// NewArtifact get Artifact client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewArtifact() ClientArtifact {
	return newArtifact(p)
}

// GitBranch get branch client
func (p *PluginClient) GitBranch(meta Meta, secret corev1.Secret) ClientGitBranch {
	clone := p.Clone().WithMeta(meta).WithSecret(secret)

	return newGitBranch(clone)
}

// NewGitBranch get branch client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewGitBranch() ClientGitBranch {
	return newGitBranch(p)
}

// GitContent get content client
func (p *PluginClient) GitContent(meta Meta, secret corev1.Secret) ClientGitContent {
	clone := p.Clone().WithMeta(meta).WithSecret(secret)

	return newGitContent(clone)
}

// NewGitContent get content client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewGitContent() ClientGitContent {
	return newGitContent(p)
}

// GitPullRequest get pr client
func (p *PluginClient) GitPullRequest(meta Meta, secret corev1.Secret) GitPullRequestCRUClient {
	clone := p.Clone().WithMeta(meta).WithSecret(secret)

	return newGitPullRequest(clone)
}

// NewGitPullRequest get pr client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewGitPullRequest() GitPullRequestCRUClient {
	return newGitPullRequest(p)
}

// GitCommit get pr client
func (p *PluginClient) GitCommit(meta Meta, secret corev1.Secret) ClientGitCommit {
	clone := p.Clone().WithMeta(meta).WithSecret(secret)

	return newGitCommit(clone)
}

// NewGitCommit get pr client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewGitCommit() ClientGitCommit {
	return newGitCommit(p)
}

// GitRepository get repo client
func (p *PluginClient) GitRepository(meta Meta, secret corev1.Secret) ClientGitRepository {
	clone := p.Clone().WithMeta(meta).WithSecret(secret)

	return newGitRepository(clone)
}

// NewGitRepository get repo client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewGitRepository() ClientGitRepository {
	return newGitRepository(p)
}

// GitRepositoryFileTree get repo file tree client
func (p *PluginClient) GitRepositoryFileTree(meta Meta, secret corev1.Secret) ClientGitRepositoryFileTree {
	clone := p.Clone().WithMeta(meta).WithSecret(secret)

	return newGitRepositoryFileTree(clone)
}

// NewGitRepositoryFileTree get repo file tree client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewGitRepositoryFileTree() ClientGitRepositoryFileTree {
	return newGitRepositoryFileTree(p)
}

// GitCommitComment get commit comment client
func (p *PluginClient) GitCommitComment(meta Meta, secret corev1.Secret) ClientGitCommitComment {
	clone := p.Clone().WithMeta(meta).WithSecret(secret)

	return newGitCommitComment(clone)
}

// NewGitCommitComment get commit comment client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewGitCommitComment() ClientGitCommitComment {
	return newGitCommitComment(p)
}

// GitCommitStatus get commit comment client
func (p *PluginClient) GitCommitStatus(meta Meta, secret corev1.Secret) ClientGitCommitStatus {
	clone := p.Clone().WithMeta(meta).WithSecret(secret)

	return newGitCommitStatus(clone)
}

// NewGitCommitStatus get commit comment client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewGitCommitStatus() ClientGitCommitStatus {
	return newGitCommitStatus(p)
}

// CodeQuality get code quality client
func (p *PluginClient) CodeQuality(meta Meta, secret corev1.Secret) ClientCodeQuality {
	clone := p.Clone().WithMeta(meta).WithSecret(secret)

	return newCodeQuality(clone)
}

// NewCodeQuality get code quality client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewCodeQuality() ClientCodeQuality {
	return newCodeQuality(p)
}

// BlobStore get blob store client
func (p *PluginClient) BlobStore(meta Meta, secret corev1.Secret) ClientBlobStore {
	clone := p.Clone().WithMeta(meta).WithSecret(secret)

	return newBlobStore(clone)
}

// NewBlobStore get blob store client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewBlobStore() ClientBlobStore {
	return newBlobStore(p)
}

// GitRepositoryTag get repository tag client
func (p *PluginClient) GitRepositoryTag(meta Meta, secret corev1.Secret) ClientGitRepositoryTag {
	clone := p.Clone().WithMeta(meta).WithSecret(secret)

	return newGitRepositoryTag(clone)
}

// NewGitRepositoryTag get repository tag client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewGitRepositoryTag() ClientGitRepositoryTag {
	return newGitRepositoryTag(p)
}

// TestPlan get test plan client
func (p *PluginClient) TestPlan(meta Meta, secret corev1.Secret) ClientTestPlan {
	clone := p.Clone().WithMeta(meta).WithSecret(secret)

	return newTestPlan(clone)
}

// TestCase get test case client
func (p *PluginClient) TestCase(meta Meta, secret corev1.Secret) ClientTestCase {
	clone := p.Clone().WithMeta(meta).WithSecret(secret)

	return newTestCase(clone)
}

// TestModule get test module client
func (p *PluginClient) TestModule(meta Meta, secret corev1.Secret) ClientTestModule {
	clone := p.Clone().WithMeta(meta).WithSecret(secret)

	return newTestModule(clone)
}

// TestCaseExecution get test case execution client
func (p *PluginClient) TestCaseExecution(meta Meta, secret corev1.Secret) ClientTestCaseExecution {
	clone := p.Clone().WithMeta(meta).WithSecret(secret)

	return newTestCaseExecution(clone)
}

// NewToolService get tool service execution client
func (p *PluginClient) NewToolService() ClientToolService {
	return newToolService(p, p.ClassAddress)
}

// NewToolMetadata get tool metadata execution client
func (p *PluginClient) NewToolMetadata(meta Meta, secret corev1.Secret) ClientToolMetadata {
	clone := p.Clone().WithMeta(meta).WithSecret(secret)

	return newToolMetadata(clone)
}
