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
	duckv1 "knative.dev/pkg/apis/duck/v1"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
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
		client: resty.New(),
	}

	for _, op := range opts {
		op(pluginClient)
	}
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

// Get performs a GET request using defined options
func (p *PluginClient) Get(ctx context.Context, baseURL *duckv1.Addressable, path string, options ...OptionFunc) error {
	request := p.R(ctx, baseURL, options...)
	_, err := request.Get(p.fullUrl(baseURL, path))
	return err
}

// Post performs a POST request with the given parameters
func (p *PluginClient) Post(ctx context.Context, baseURL *duckv1.Addressable, path string, options ...OptionFunc) error {
	request := p.R(ctx, baseURL, options...)
	_, err := request.Post(p.fullUrl(baseURL, path))
	return err
}

// Put performs a PUT request with the given parameters
func (p *PluginClient) Put(ctx context.Context, baseURL *duckv1.Addressable, path string, options ...OptionFunc) error {
	request := p.R(ctx, baseURL, options...)
	_, err := request.Put(p.fullUrl(baseURL, path))
	return err
}

// Delete performs a DELETE request with the given parameters
func (p *PluginClient) Delete(ctx context.Context, baseURL *duckv1.Addressable, path string, options ...OptionFunc) error {
	request := p.R(ctx, baseURL, options...)
	_, err := request.Put(p.fullUrl(baseURL, path))
	return err
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

// Project get project client
func (p *PluginClient) Project(meta Meta, secret corev1.Secret) ClientProject {
	return newProject(p, meta, secret)
}
