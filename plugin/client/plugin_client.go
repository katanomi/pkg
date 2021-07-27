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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	corev1 "k8s.io/api/core/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
)

// PluginClient client for plugins
type PluginClient struct {
	address *duckv1.Addressable
	client  *resty.Client
}

// BuildOptions Options to build the plugin client
type BuildOptions func(client *PluginClient)

// NewPluginClient creates a new plugin client
func NewPluginClient(address *duckv1.Addressable, opts ...BuildOptions) *PluginClient {
	pluginClient := &PluginClient{
		address: address,
		client:  resty.New(),
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

// OptionFunc options for requests
type OptionFunc func(request *resty.Request)

// Get performs a get request using defined options
func (p *PluginClient) Get(uri string, options ...OptionFunc) error {
	request := p.client.R()

	for _, fn := range options {
		fn(request)
	}

	_, err := request.Get(p.fullUrl(uri))

	return err
}

// Secret provides a secret to be assigned to the request in the header
func (p *PluginClient) Secret(secret corev1.Secret) OptionFunc {
	return func(request *resty.Request) {
		request.SetHeader(PluginAuthHeader, string(secret.Type))
		dataBytes, _ := json.Marshal(secret.Data)
		request.SetHeader(PluginSecretHeader, base64.StdEncoding.EncodeToString(dataBytes))
	}
}

// Meta provides metadata for the request
func (p *PluginClient) Meta(meta Meta) OptionFunc {
	return func(request *resty.Request) {
		dataBytes, _ := json.Marshal(meta)
		request.SetHeader(PluginMetaHeader, base64.StdEncoding.EncodeToString(dataBytes))
	}
}

// ListOptions options for lists
func (p *PluginClient) ListOptions(opts metav1alpha1.ListOptions) OptionFunc {
	return func(request *resty.Request) {
		params := make(map[string]string)
		if len(opts.Search) > 0 {
			for k, v := range opts.Search {
				// actually a request can be performed with the same key for params
				// multiple times, but in this case resty does not support
				// then we just fetch the last one
				if len(v) > 0 {
					params[k] = v[len(v)-1]
				}
			}
		}
		params["page"] = strconv.Itoa(opts.Page)
		params["itemsPerPage"] = strconv.Itoa(opts.ItemsPerPage)
		request.SetQueryParams(params)
	}
}

// Query query parameters for the request
func (p *PluginClient) Query(params map[string]string) OptionFunc {
	return func(request *resty.Request) {
		request.SetQueryParams(params)
	}
}

// Body request body
func (p *PluginClient) Body(body interface{}) OptionFunc {
	return func(request *resty.Request) {
		request.SetBody(body)
	}
}

// Dest request result automatically marshalled into object
func (p *PluginClient) Dest(dest interface{}) OptionFunc {
	return func(request *resty.Request) {
		request.SetResult(dest)
	}
}

// Error error response object
func (p *PluginClient) Error(err interface{}) OptionFunc {
	return func(request *resty.Request) {
		request.SetError(err)
	}
}

func (p *PluginClient) fullUrl(uri string) string {
	url := p.address.URL.String()
	return fmt.Sprintf("%s/%s", strings.TrimSuffix(url, "/"), strings.TrimPrefix(uri, "/"))
}
