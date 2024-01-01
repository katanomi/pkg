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
	"github.com/katanomi/pkg/plugin/client"
	"github.com/katanomi/pkg/plugin/client/base"
	corev1 "k8s.io/api/core/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// NewPluginClientV2 construct a new plugin client
func NewPluginClientV2(pluginAddress *duckv1.Addressable, meta base.Meta, secret corev1.Secret, opts ...base.BuildOptions) *PluginClientV2 {
	return &PluginClientV2{
		PluginClient: base.NewPluginClient(opts...).
			WithMeta(meta).WithSecret(secret).WithClassAddress(pluginAddress),
	}
}

// PluginClientV2 describe a plugin client with version 2
type PluginClientV2 struct {
	client.Interface

	*base.PluginClient
}

// WithMeta set metadata of target tool
func (p *PluginClientV2) WithMeta(meta base.Meta) *PluginClientV2 {
	p.PluginClient = p.PluginClient.WithMeta(meta)
	return p
}

// WithSecret set authorization secret of target tool
func (p *PluginClientV2) WithSecret(secret corev1.Secret) *PluginClientV2 {
	p.PluginClient = p.PluginClient.WithSecret(secret)
	return p
}

// WithClassAddress set class address of target plugin
func (p *PluginClientV2) WithClassAddress(classAddress *duckv1.Addressable) *PluginClientV2 {
	p.PluginClient = p.PluginClient.WithClassAddress(classAddress)
	return p
}

// WithRequestOptions set options for the next request
func (p *PluginClientV2) WithRequestOptions(opts ...base.OptionFunc) *PluginClientV2 {
	return &PluginClientV2{
		PluginClient: p.PluginClient.WithRequestOptions(opts...),
	}
}
