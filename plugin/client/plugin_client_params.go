/*
Copyright 2022 The Katanomi Authors.

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
	corev1 "k8s.io/api/core/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
)

// PluginClientParams used to store the parameters of the plugin client call.
// Avoid repeated generation some parameters
type PluginClientParams struct {
	// Meta Plugin meta with base url and version info, for calling plugin api
	Meta Meta

	// ClassAddress is the address of the integration class
	ClassAddress *duckv1.Addressable

	// Secret is the secret to use for the plugin client
	Secret *corev1.Secret

	// IntegrationClassName is the name of the integration class
	IntegrationClassName string

	// GitRepo Repo base info, such as project, repository
	GitRepo metav1alpha1.GitRepo
}

// NewPluginClientParams generate plugin client params
func NewPluginClientParams() *PluginClientParams {
	return &PluginClientParams{}
}

func (p *PluginClientParams) WithMeta(meta Meta) *PluginClientParams {
	p.Meta = meta
	return p
}

func (p *PluginClientParams) WithClassAddress(classAddress *duckv1.Addressable) *PluginClientParams {
	p.ClassAddress = classAddress
	return p
}

func (p *PluginClientParams) WithSecret(secret *corev1.Secret) *PluginClientParams {
	p.Secret = secret
	return p
}

func (p *PluginClientParams) WithIntegrationClassName(integrationClassName string) *PluginClientParams {
	p.IntegrationClassName = integrationClassName
	return p
}

func (p *PluginClientParams) WithGitRepo(gitRepo metav1alpha1.GitRepo) *PluginClientParams {
	p.GitRepo = gitRepo
	return p
}
