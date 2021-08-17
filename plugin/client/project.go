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

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

type ClientProject interface {
	List(ctx context.Context, baseURL *duckv1.Addressable, options ...OptionFunc) (*metav1alpha1.ProjectList, error)
	Create(ctx context.Context, baseURL *duckv1.Addressable, project *metav1alpha1.Project, options ...OptionFunc) (*metav1alpha1.Project, error)
	Get(ctx context.Context, baseURL *duckv1.Addressable, id string, options ...OptionFunc) (*metav1alpha1.Project, error)
}

type project struct {
	client Client
	meta   Meta
	secret corev1.Secret
}

func newProject(client Client, meta Meta, secret corev1.Secret) ClientProject {
	return &project{
		client: client,
		meta:   meta,
		secret: secret,
	}
}

// List get project using plugin
func (p *project) List(ctx context.Context, baseURL *duckv1.Addressable, options ...OptionFunc) (*metav1alpha1.ProjectList, error) {
	list := &metav1alpha1.ProjectList{}

	options = append(options, MetaOpts(p.meta), SecretOpts(p.secret), ResultOpts(list))
	if err := p.client.Get(ctx, baseURL, "projects", options...); err != nil {
		return nil, err
	}

	return list, nil
}

// Create create project using plugin
func (p *project) Create(ctx context.Context, baseURL *duckv1.Addressable, project *metav1alpha1.Project, options ...OptionFunc) (*metav1alpha1.Project, error) {
	options = append(options, MetaOpts(p.meta), SecretOpts(p.secret), BodyOpts(project))
	if err := p.client.Post(ctx, baseURL, "projects", options...); err != nil {
		return nil, err
	}

	return project, nil
}

// Get get project using plugin
func (p *project) Get(ctx context.Context, baseURL *duckv1.Addressable, id string, options ...OptionFunc) (*metav1alpha1.Project, error) {
	resp := &metav1alpha1.Project{}
	options = append(options, MetaOpts(p.meta), SecretOpts(p.secret), ResultOpts(resp))
	if err := p.client.Get(ctx, baseURL, "projects/"+id, options...); err != nil {
		return nil, err
	}

	return resp, nil
}
