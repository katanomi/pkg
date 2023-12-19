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

package base

import (
	"context"
	"fmt"
	"strings"

	"github.com/emicklei/go-restful/v3"
	"github.com/go-resty/resty/v2"
	corev1 "k8s.io/api/core/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/client"
	perrors "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/tracing"
)

// Client interface for PluginClient, client code should use the interface
// as dependency
type Client interface {
	Get(ctx context.Context, baseURL *duckv1.Addressable, uri string, options ...OptionFunc) error
	GetResponse(ctx context.Context, baseURL *duckv1.Addressable, uri string, options ...OptionFunc) (*resty.Response, error)
	Post(ctx context.Context, baseURL *duckv1.Addressable, uri string, options ...OptionFunc) error
	Put(ctx context.Context, baseURL *duckv1.Addressable, uri string, options ...OptionFunc) error
	Delete(ctx context.Context, baseURL *duckv1.Addressable, uri string, options ...OptionFunc) error
}

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

	// requestOptions options to wrap resty request
	requestOptions []OptionFunc
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

// Make sure that PluginClient implements the PluginClient interface
var _ Client = &PluginClient{}

// OptionFunc options for requests
type OptionFunc func(request *resty.Request)

func (p *PluginClient) GetMeta() Meta {
	return p.meta
}

func (p *PluginClient) GetSecret() corev1.Secret {
	return p.secret
}

// Clone shallow clone the plugin client
// used to update some fields without changing the original
func (p *PluginClient) Clone() *PluginClient {
	if p == nil {
		return nil
	}
	clone := *p
	return &clone
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

// WithRequestOptions set request options
func (p *PluginClient) WithRequestOptions(opts ...OptionFunc) *PluginClient {
	clone := p.Clone()
	clone.requestOptions = opts
	return clone
}

func (p *PluginClient) builtinOptions() []OptionFunc {
	options := append(DefaultOptions(), MetaOpts(p.meta), SecretOpts(p.secret))
	return append(options, p.requestOptions...)
}

// Get performs a GET request using defined options
func (p *PluginClient) Get(ctx context.Context, baseURL *duckv1.Addressable, path string, options ...OptionFunc) error {
	options = append(p.builtinOptions(), options...)

	request := p.R(ctx, options...)
	response, err := request.Get(p.FullUrl(baseURL, path))

	return p.HandleError(response, err)
}

// GetResponse performs a GET request using defined options and return response
func (p *PluginClient) GetResponse(ctx context.Context, baseURL *duckv1.Addressable, path string, options ...OptionFunc) (*resty.Response, error) {
	options = append(p.builtinOptions(), options...)

	request := p.R(ctx, options...)
	return request.Get(p.FullUrl(baseURL, path))
}

// Post performs a POST request with the given parameters
func (p *PluginClient) Post(ctx context.Context, baseURL *duckv1.Addressable, path string, options ...OptionFunc) error {
	options = append(p.builtinOptions(), options...)

	request := p.R(ctx, options...)
	response, err := request.Post(p.FullUrl(baseURL, path))

	return p.HandleError(response, err)
}

// Put performs a PUT request with the given parameters
func (p *PluginClient) Put(ctx context.Context, baseURL *duckv1.Addressable, path string, options ...OptionFunc) error {
	options = append(p.builtinOptions(), options...)

	request := p.R(ctx, options...)
	response, err := request.Put(p.FullUrl(baseURL, path))

	return p.HandleError(response, err)
}

// Delete performs a DELETE request with the given parameters
func (p *PluginClient) Delete(ctx context.Context, baseURL *duckv1.Addressable, path string, options ...OptionFunc) error {
	options = append(p.builtinOptions(), options...)

	request := p.R(ctx, options...)
	response, err := request.Delete(p.FullUrl(baseURL, path))

	return p.HandleError(response, err)
}

// R prepares a request based on the given information
func (p *PluginClient) R(ctx context.Context, options ...OptionFunc) *resty.Request {
	request := p.client.R()

	request.SetContext(ctx)

	for _, fn := range options {
		fn(request)
	}
	return request
}

func (p *PluginClient) FullUrl(address *duckv1.Addressable, uri string) string {
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

// DefaultOptions for default plugin client options
func DefaultOptions() []OptionFunc {
	return []OptionFunc{
		ErrorOpts(&ResponseStatusErr{}),
		HeaderOpts("Content-Type", "application/json"),
	}
}

// GetSubResourcesOptionsFromRequest returns SubResourcesOptions based on a request
func GetSubResourcesOptionsFromRequest(req *restful.Request) (opts metav1alpha1.SubResourcesOptions) {
	subResourcesHeader := req.HeaderParameter(PluginSubresourcesHeader)
	if strings.TrimSpace(subResourcesHeader) != "" {
		opts.SubResources = strings.Split(subResourcesHeader, ",")
	}
	return
}
