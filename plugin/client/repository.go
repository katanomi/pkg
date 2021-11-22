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

type ClientRepository interface {
	List(ctx context.Context, baseURL *duckv1.Addressable, project string, options ...OptionFunc) (*metav1alpha1.RepositoryList, error)
}

type repository struct {
	client Client
	meta   Meta
	secret corev1.Secret
}

func newRepository(client Client, meta Meta, secret corev1.Secret) ClientRepository {
	return &repository{
		client: client,
		meta:   meta,
		secret: secret,
	}
}

// List get project using plugin
func (p *repository) List(ctx context.Context, baseURL *duckv1.Addressable, project string, options ...OptionFunc) (*metav1alpha1.RepositoryList, error) {
	list := &metav1alpha1.RepositoryList{}

	uri := fmt.Sprintf("projects/%s/repositories", project)
	options = append(options, MetaOpts(p.meta), SecretOpts(p.secret), ResultOpts(list))
	if err := p.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}

	return list, nil
}
