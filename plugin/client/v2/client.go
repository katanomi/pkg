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

package v2

import (
	"context"

	"github.com/katanomi/pkg/plugin/client/base"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// NewPluginClient construct a new plugin client
func NewPluginClient(pluginAddress *duckv1.Addressable, meta base.Meta, secret corev1.Secret, opts ...base.BuildOptions) *PluginClient {
	return &PluginClient{
		PluginClient: base.NewPluginClient(opts...).
			WithMeta(meta).WithSecret(secret).WithClassAddress(pluginAddress),
	}
}

// PluginClient describe a plugin client with version 2
type PluginClient struct {
	*base.PluginClient
}

// Path empty implementation for `plugin Interface`
func (p *PluginClient) Path() string {
	return ""
}

// Setup empty implementation for `plugin Interface`
func (p *PluginClient) Setup(_ context.Context, _ *zap.SugaredLogger) error {
	return nil
}

// WithMeta set metadata of target tool
func (p *PluginClient) WithMeta(meta base.Meta) *PluginClient {
	p.PluginClient = p.PluginClient.WithMeta(meta)
	return p
}

// WithSecret set authorization secret of target tool
func (p *PluginClient) WithSecret(secret corev1.Secret) *PluginClient {
	p.PluginClient = p.PluginClient.WithSecret(secret)
	return p
}

// WithClassAddress set class address of target plugin
func (p *PluginClient) WithClassAddress(classAddress *duckv1.Addressable) *PluginClient {
	p.PluginClient = p.PluginClient.WithClassAddress(classAddress)
	return p
}

// WithRequestOptions set options for the next request
func (p *PluginClient) WithRequestOptions(opts ...base.OptionFunc) *PluginClient {
	return &PluginClient{
		PluginClient: p.PluginClient.WithRequestOptions(opts...),
	}
}
