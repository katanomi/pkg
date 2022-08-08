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
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	duckv1 "knative.dev/pkg/apis/duck/v1"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/client"
	perrors "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/tracing"
)

var (
	defaultOptions = []OptionFunc{
		ErrorOpts(&errors.StatusError{}),
		HeaderOpts("Content-Type", "application/json"),
	}
)

// PluginClient client for plugins
type PluginClient struct {
	client *resty.Client

	// meta plugin meta with base url and version info, for calling plugin api
	// +optional
	meta Meta

	// secret is the secret to use for the plugin client
	// +optional
	secret corev1.Secret

	// ClassAddress is the address of the integration class
	// +optional
	ClassAddress *duckv1.Addressable

	// IntegrationClassName is the name of the integration class
	// +optional
	IntegrationClassName string
}

// BuildOptions Options to build the plugin client
type BuildOptions func(client *PluginClient)

// NewPluginClient creates a new plugin client
func NewPluginClient(opts ...BuildOptions) *PluginClient {
	restyClient := resty.NewWithClient(client.NewHTTPClient())
	restyClient.SetDisableWarn(true)

	pluginClient := &PluginClient{
		client: restyClient,
	}

	for _, op := range opts {
		op(pluginClient)
	}

	tracing.WrapTransportForRestyClient(pluginClient.client)
	return pluginClient
}

// ClientOpts adds a custom client build options for plugin client
func ClientOpts(clt *resty.Client) BuildOptions {
	return func(client *PluginClient) {
		client.client = clt
	}
}

// Make sure that PluginClient implements the Client interface
var _ Client = &PluginClient{}

// OptionFunc options for requests
type OptionFunc func(request *resty.Request)

// Clone shallow clone the plugin client
// used to update some fields without changing the original
func (p *PluginClient) Clone() *PluginClient {
	if p == nil {
		return nil
	}
	p = &(*p)
	return p
}

func (p *PluginClient) WithMeta(meta Meta) *PluginClient {
	p.meta = meta
	return p
}

func (p *PluginClient) WithSecret(secret corev1.Secret) *PluginClient {
	p.secret = secret
	return p
}

func (p *PluginClient) WithClassAddress(classAddress *duckv1.Addressable) *PluginClient {
	p.ClassAddress = classAddress
	return p
}

func (p *PluginClient) WithIntegrationClassName(integrationClassName string) *PluginClient {
	p.IntegrationClassName = integrationClassName
	return p
}

// Get performs a GET request using defined options
func (p *PluginClient) Get(ctx context.Context, baseURL *duckv1.Addressable, path string, options ...OptionFunc) error {
	options = append(defaultOptions, options...)

	request := p.R(ctx, baseURL, options...)
	response, err := request.Get(p.fullUrl(baseURL, path))

	return p.HandleError(response, err)
}

// Post performs a POST request with the given parameters
func (p *PluginClient) Post(ctx context.Context, baseURL *duckv1.Addressable, path string, options ...OptionFunc) error {
	options = append(defaultOptions, options...)

	request := p.R(ctx, baseURL, options...)
	response, err := request.Post(p.fullUrl(baseURL, path))

	return p.HandleError(response, err)
}

// Put performs a PUT request with the given parameters
func (p *PluginClient) Put(ctx context.Context, baseURL *duckv1.Addressable, path string, options ...OptionFunc) error {
	options = append(defaultOptions, options...)

	request := p.R(ctx, baseURL, options...)
	response, err := request.Put(p.fullUrl(baseURL, path))

	return p.HandleError(response, err)
}

// Delete performs a DELETE request with the given parameters
func (p *PluginClient) Delete(ctx context.Context, baseURL *duckv1.Addressable, path string, options ...OptionFunc) error {
	options = append(defaultOptions, options...)

	request := p.R(ctx, baseURL, options...)
	response, err := request.Delete(p.fullUrl(baseURL, path))

	return p.HandleError(response, err)
}

// R prepares a request based on the given information
func (p *PluginClient) R(ctx context.Context, baseURL *duckv1.Addressable, options ...OptionFunc) *resty.Request {
	request := p.client.R()

	request.SetContext(ctx)

	for _, fn := range options {
		fn(request)
	}
	return request
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

func (p *PluginClient) fullUrl(address *duckv1.Addressable, uri string) string {
	url := address.URL.String()
	return fmt.Sprintf("%s/%s", strings.TrimSuffix(url, "/"), strings.TrimPrefix(uri, "/"))
}

func (p *PluginClient) HandleError(response *resty.Response, err error) error {
	if err != nil {
		return err
	}

	if response.IsError() {
		return perrors.AsStatusError(response)
	}

	return nil
}

// GitPluginClient convert PluginClient to GitPluginClient
func (p *PluginClient) GitPluginClient() *GitPluginClient {
	return &GitPluginClient{PluginClient: p}
}

// Auth provides an auth methods for clients
func (p *PluginClient) Auth(meta Meta, secret corev1.Secret) ClientAuth {
	return newAuthClient(p, meta, secret)
}

// NewAuth provides an auth methods for clients
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewAuth() ClientAuth {
	return newAuthClient(p, p.meta, p.secret)
}

// Project get project client
func (p *PluginClient) Project(meta Meta, secret corev1.Secret) ClientProject {
	return newProject(p, meta, secret)
}

// NewProject get project client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewProject() ClientProject {
	return newProject(p, p.meta, p.secret)
}

// Repository get Repository client
func (p *PluginClient) Repository(meta Meta, secret corev1.Secret) ClientRepository {
	return newRepository(p, meta, secret)
}

// NewRepository get Repository client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewRepository() ClientRepository {
	return newRepository(p, p.meta, p.secret)
}

// Artifact get Artifact client
func (p *PluginClient) Artifact(meta Meta, secret corev1.Secret) ClientArtifact {
	return newArtifact(p, meta, secret)
}

// NewArtifact get Artifact client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewArtifact() ClientArtifact {
	return newArtifact(p, p.meta, p.secret)
}

// GitBranch get branch client
func (p *PluginClient) GitBranch(meta Meta, secret corev1.Secret) ClientGitBranch {
	return newGitBranch(p, meta, secret)
}

// NewGitBranch get branch client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewGitBranch() ClientGitBranch {
	return newGitBranch(p, p.meta, p.secret)
}

// GitContent get content client
func (p *PluginClient) GitContent(meta Meta, secret corev1.Secret) ClientGitContent {
	return newGitContent(p, meta, secret)
}

// NewGitContent get content client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewGitContent() ClientGitContent {
	return newGitContent(p, p.meta, p.secret)
}

// GitPullRequest get pr client
func (p *PluginClient) GitPullRequest(meta Meta, secret corev1.Secret) GitPullRequestCRUClient {
	return newGitPullRequest(p, meta, secret)
}

// NewGitPullRequest get pr client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewGitPullRequest() GitPullRequestCRUClient {
	return newGitPullRequest(p, p.meta, p.secret)
}

// GitCommit get pr client
func (p *PluginClient) GitCommit(meta Meta, secret corev1.Secret) ClientGitCommit {
	return newGitCommit(p, meta, secret)
}

// NewGitCommit get pr client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewGitCommit() ClientGitCommit {
	return newGitCommit(p, p.meta, p.secret)
}

// GitRepository get repo client
func (p *PluginClient) GitRepository(meta Meta, secret corev1.Secret) ClientGitRepository {
	return newGitRepository(p, meta, secret)
}

// NewGitRepository get repo client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewGitRepository() ClientGitRepository {
	return newGitRepository(p, p.meta, p.secret)
}

// GitRepositoryFileTree get repo file tree client
func (p *PluginClient) GitRepositoryFileTree(meta Meta, secret corev1.Secret) ClientGitRepositoryFileTree {
	return newGitRepositoryFileTree(p, meta, secret)
}

// NewGitRepositoryFileTree get repo file tree client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewGitRepositoryFileTree() ClientGitRepositoryFileTree {
	return newGitRepositoryFileTree(p, p.meta, p.secret)
}

// GitCommitComment get commit comment client
func (p *PluginClient) GitCommitComment(meta Meta, secret corev1.Secret) ClientGitCommitComment {
	return newGitCommitComment(p, meta, secret)
}

// NewGitCommitComment get commit comment client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewGitCommitComment() ClientGitCommitComment {
	return newGitCommitComment(p, p.meta, p.secret)
}

// GitCommitStatus get commit comment client
func (p *PluginClient) GitCommitStatus(meta Meta, secret corev1.Secret) ClientGitCommitStatus {
	return newGitCommitStatus(p, meta, secret)
}

// NewGitCommitStatus get commit comment client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewGitCommitStatus() ClientGitCommitStatus {
	return newGitCommitStatus(p, p.meta, p.secret)
}

// CodeQuality get code quality client
func (p *PluginClient) CodeQuality(meta Meta, secret corev1.Secret) ClientCodeQuality {
	return newCodeQuality(p, meta, secret)
}

// NewCodeQuality get code quality client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewCodeQuality() ClientCodeQuality {
	return newCodeQuality(p, p.meta, p.secret)
}

// BlobStore get blob store client
func (p *PluginClient) BlobStore(meta Meta, secret corev1.Secret) ClientBlobStore {
	return newBlobStore(p, meta, secret)
}

// NewBlobStore get blob store client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewBlobStore() ClientBlobStore {
	return newBlobStore(p, p.meta, p.secret)
}

// GitRepositoryTag get repository tag client
func (p *PluginClient) GitRepositoryTag(meta Meta, secret corev1.Secret) ClientGitRepositoryTag {
	return newGitRepositoryTag(p, meta, secret)
}

// NewGitRepositoryTag get repository tag client
// Use the internal meta and secret to generate the client, please assign in advance.
func (p *PluginClient) NewGitRepositoryTag() ClientGitRepositoryTag {
	return newGitRepositoryTag(p, p.meta, p.secret)
}
