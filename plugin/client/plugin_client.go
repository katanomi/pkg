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
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	perrors "github.com/katanomi/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	duckv1 "knative.dev/pkg/apis/duck/v1"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
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
}

// BuildOptions Options to build the plugin client
type BuildOptions func(client *PluginClient)

// NewPluginClient creates a new plugin client
func NewPluginClient(opts ...BuildOptions) *PluginClient {
	pluginClient := &PluginClient{
		client: resty.NewWithClient(NewHTTPClient()),
	}

	for _, op := range opts {
		op(pluginClient)
	}
	return pluginClient
}

func NewHTTPClient() *http.Client {

	client := &http.Client{
		Transport: GetDefaultTransport(),
		Timeout:   30 * time.Second,
	}

	return client
}

func GetDefaultTransport() http.RoundTripper {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		// #nosec
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
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

// Project get project client
func (p *PluginClient) Project(meta Meta, secret corev1.Secret) ClientProject {
	return newProject(p, meta, secret)
}

// GitBranch get branch client
func (p *PluginClient) GitBranch(meta Meta, secret corev1.Secret) ClientGitBranch {
	return newGitBranch(p, meta, secret)
}

// GitContent get content client
func (p *PluginClient) GitContent(meta Meta, secret corev1.Secret) ClientGitContent {
	return newGitContent(p, meta, secret)
}

// GitPullRequest get pr client
func (p *PluginClient) GitPullRequest(meta Meta, secret corev1.Secret) ClientGitPullRequest {
	return newGitPullRequest(p, meta, secret)
}

// GitPullRequest get pr client
func (p *PluginClient) GitCommit(meta Meta, secret corev1.Secret) ClientGitCommit {
	return newGitCommit(p, meta, secret)
}

// GitRepository get repo client
func (p *PluginClient) GitRepository(meta Meta, secret corev1.Secret) ClientGitRepository {
	return newGitRepository(p, meta, secret)
}

// GitCommitComment get commit comment client
func (p *PluginClient) GitCommitComment(meta Meta, secret corev1.Secret) ClientGitCommitComment {
	return newGitCommitComment(p, meta, secret)
}

// GitCommitStatus get commit comment client
func (p *PluginClient) GitCommitStatus(meta Meta, secret corev1.Secret) ClientGitCommitStatus {
	return newGitCommitStatus(p, meta, secret)
}
