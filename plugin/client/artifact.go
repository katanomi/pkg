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

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

type ClientArtifact interface {
	List(ctx context.Context, baseURL *duckv1.Addressable, project string, repository string, options ...OptionFunc) (*metav1alpha1.ArtifactList, error)
	Delete(ctx context.Context, baseURL *duckv1.Addressable, project string, repository string, artifact string, options ...OptionFunc) error
	DeleteTag(ctx context.Context, baseURL *duckv1.Addressable, project string, repository string, artifact string, tag string, options ...OptionFunc) error
}

type artifact struct {
	client Client
	meta   Meta
	secret corev1.Secret
}

func newArtifact(client Client, meta Meta, secret corev1.Secret) ClientArtifact {
	return &artifact{
		client: client,
		meta:   meta,
		secret: secret,
	}
}

// List get artifact using plugin
func (p *artifact) List(ctx context.Context,
	baseURL *duckv1.Addressable,
	project, repsitory string,
	options ...OptionFunc) (*metav1alpha1.ArtifactList, error) {

	list := &metav1alpha1.ArtifactList{}

	uri := fmt.Sprintf("projects/%s/repositories/%s/artifacts", project, repsitory)
	options = append(options, MetaOpts(p.meta), SecretOpts(p.secret), ResultOpts(list))
	if err := p.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}

	return list, nil
}

// Delete artifact using plugin
func (p *artifact) Delete(ctx context.Context, baseURL *duckv1.Addressable, project string, repsitory string, artifact string, options ...OptionFunc) error {
	uri := fmt.Sprintf("projects/%s/repositories/%s/artifacts/%s", project, repsitory, artifact)
	options = append(options, MetaOpts(p.meta), SecretOpts(p.secret))
	if err := p.client.Delete(ctx, baseURL, uri, options...); err != nil {
		return err
	}

	return nil
}

// DeleteTag artifact's tag using plugin
func (p *artifact) DeleteTag(ctx context.Context,
	baseURL *duckv1.Addressable,
	project, repsitory, artifact, tag string,
	options ...OptionFunc) error {

	uri := fmt.Sprintf("projects/%s/repositories/%s/artifacts/%s/tags/%s", project, repsitory, artifact, tag)
	options = append(options, MetaOpts(p.meta), SecretOpts(p.secret))
	if err := p.client.Delete(ctx, baseURL, uri, options...); err != nil {
		return err
	}

	return nil
}
