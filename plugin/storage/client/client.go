/*
Copyright 2023 The Katanomi Authors.

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
	"path"

	"github.com/go-resty/resty/v2"
	pkgClient "github.com/katanomi/pkg/client"
	perrors "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/plugin/client"
	"github.com/katanomi/pkg/tracing"
	"k8s.io/apimachinery/pkg/runtime/schema"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// Interface captures the set of operations for generically interacting with Kubernetes REST apis.
type Interface interface {
	Get(ctx context.Context, path string,
		options ...client.OptionFunc) error
	Put(ctx context.Context, path string,
		options ...client.OptionFunc) error
	Post(ctx context.Context, path string,
		options ...client.OptionFunc) error
	Delete(ctx context.Context, path string,
		options ...client.OptionFunc) error
	APIVersion() *schema.GroupVersion
}

// StoragePluginClient is the client for storage client.
type StoragePluginClient struct {
	client *resty.Client

	// groupVersion stands for group(core or capabilities' name) and its version
	// used for generating client request path
	groupVersion *schema.GroupVersion

	// classAddress is the address of the integration class
	// +optional
	classAddress *duckv1.Addressable
}

// NewStoragePluginClient creates a new plugin client
func NewStoragePluginClient(baseURL *duckv1.Addressable, opts ...BuildOptions) *StoragePluginClient {
	restyClient := resty.NewWithClient(pkgClient.NewHTTPClient())
	restyClient.SetDisableWarn(true)

	pluginClient := &StoragePluginClient{
		client:       restyClient,
		classAddress: baseURL,
	}

	for _, op := range opts {
		op(pluginClient)
	}

	tracing.WrapTransportForRestyClient(pluginClient.client)
	return pluginClient
}

// Get performs a GET request using defined options
func (p *StoragePluginClient) Get(ctx context.Context, path string,
	options ...client.OptionFunc) error {
	options = append(client.DefaultOptions, options...)
	request := p.R(ctx, options...)
	response, err := request.Get(p.FullUrl(path))

	return p.HandleError(response, err)
}

// Post performs a POST request with the given parameters
func (p *StoragePluginClient) Post(ctx context.Context, path string,
	options ...client.OptionFunc) error {
	clientOptions := append(client.DefaultOptions)
	options = append(clientOptions, options...)

	request := p.R(ctx, options...)
	response, err := request.Post(p.FullUrl(path))

	return p.HandleError(response, err)
}

// Put performs a PUT request with the given parameters
func (p *StoragePluginClient) Put(ctx context.Context, path string,
	options ...client.OptionFunc) error {
	clientOptions := append(client.DefaultOptions)
	options = append(clientOptions, options...)

	request := p.R(ctx, options...)
	response, err := request.Put(p.FullUrl(path))

	return p.HandleError(response, err)
}

// Delete performs a DELETE request with the given parameters
func (p *StoragePluginClient) Delete(ctx context.Context, path string,
	options ...client.OptionFunc) error {
	clientOptions := append(client.DefaultOptions)
	options = append(clientOptions, options...)

	request := p.R(ctx, options...)
	response, err := request.Delete(p.FullUrl(path))

	return p.HandleError(response, err)
}

// R prepares a request based on the given information
func (p *StoragePluginClient) R(ctx context.Context, options ...client.OptionFunc) *resty.Request {
	request := p.client.R()
	request.SetContext(ctx)
	for _, fn := range options {
		fn(request)
	}
	return request
}

func (p *StoragePluginClient) APIVersion() *schema.GroupVersion {
	return p.groupVersion
}

// FullUrl returns actual url
func (p *StoragePluginClient) FullUrl(uri string) string {
	if p.classAddress == nil {
		return uri
	}
	url := p.classAddress.URL.DeepCopy()
	if p.groupVersion != nil {
		url.Path = path.Join(url.Path, p.groupVersion.Identifier())
	}
	url.Path = path.Join(url.Path, uri)
	return url.String()
}

// HandleError assigns error as http response
func (p *StoragePluginClient) HandleError(response *resty.Response, err error) error {
	if err != nil {
		return err
	}

	if response.IsError() {
		return perrors.AsStatusError(response)
	}

	return nil
}

// Clone shallow clone the plugin client
// used to update some fields without changing the original
func (p *StoragePluginClient) Clone() *StoragePluginClient {
	if p == nil {
		return nil
	}
	newP := *p
	return &newP
}

func (p *StoragePluginClient) ForGroupVersion(gv *schema.GroupVersion) *StoragePluginClient {
	newClient := p.Clone()
	newClient.groupVersion = gv
	return newClient
}
