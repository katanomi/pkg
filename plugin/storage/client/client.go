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

	"github.com/go-resty/resty/v2"
	pkgClient "github.com/katanomi/pkg/client"
	"github.com/katanomi/pkg/plugin/client"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// HeaderFileMeta is the header name of file meta
const HeaderFileMeta = "X-Katanomi-Meta"

// HeaderFileAnnotationPrefix is the annotation header prefix
const HeaderFileAnnotationPrefix = "x-katanomi-annotation-"

// StoragePluginClient is the client for storage client, we simply embed client.PluginClient now.
// TODO: refactor to another client with different authorization mechanism
type StoragePluginClient struct {
	*client.PluginClient
}

// NewStoragePluginClient creates a new plugin client
func NewStoragePluginClient(opts ...client.BuildOptions) *StoragePluginClient {
	restyClient := resty.NewWithClient(pkgClient.NewHTTPClient())
	restyClient.SetDisableWarn(true)

	pluginClient := &StoragePluginClient{
		PluginClient: client.NewPluginClient(opts...),
	}
	return pluginClient
}

// Get performs a GET request using defined options
func (p *StoragePluginClient) Get(ctx context.Context, baseURL *duckv1.Addressable, path string,
	options ...client.OptionFunc) error {
	options = append(client.DefaultOptions, options...)
	request := p.R(ctx, baseURL, options...)
	response, err := request.Get(p.FullUrl(baseURL, path))

	return p.HandleError(response, err)
}

// Post performs a POST request with the given parameters
func (p *StoragePluginClient) Post(ctx context.Context, baseURL *duckv1.Addressable, path string,
	options ...client.OptionFunc) error {
	clientOptions := append(client.DefaultOptions)
	options = append(clientOptions, options...)

	request := p.R(ctx, baseURL, options...)
	response, err := request.Post(p.FullUrl(baseURL, path))

	return p.HandleError(response, err)
}

// Put performs a PUT request with the given parameters
func (p *StoragePluginClient) Put(ctx context.Context, baseURL *duckv1.Addressable, path string,
	options ...client.OptionFunc) error {
	clientOptions := append(client.DefaultOptions)
	options = append(clientOptions, options...)

	request := p.R(ctx, baseURL, options...)
	response, err := request.Put(p.FullUrl(baseURL, path))

	return p.HandleError(response, err)
}

// Delete performs a DELETE request with the given parameters
func (p *StoragePluginClient) Delete(ctx context.Context, baseURL *duckv1.Addressable, path string,
	options ...client.OptionFunc) error {
	clientOptions := append(client.DefaultOptions)
	options = append(clientOptions, options...)

	request := p.R(ctx, baseURL, options...)
	response, err := request.Delete(p.FullUrl(baseURL, path))

	return p.HandleError(response, err)
}
