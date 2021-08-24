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
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// ClientGitRepository client for repo
type ClientGitRepository interface {
	List(ctx context.Context, baseURL *duckv1.Addressable, project, keyword string, options ...OptionFunc) (*metav1alpha1.GitRepositoryList, error)
}

type gitRepository struct {
	client Client
	meta   Meta
	secret corev1.Secret
}

func newGitRepository(client Client, meta Meta, secret corev1.Secret) ClientGitRepository {
	return &gitRepository{
		client: client,
		meta:   meta,
		secret: secret,
	}
}

func (g *gitRepository) List(ctx context.Context, baseURL *duckv1.Addressable, project, keyword string, options ...OptionFunc) (*metav1alpha1.GitRepositoryList, error) {
	list := &metav1alpha1.GitRepositoryList{}
	options = append(options, MetaOpts(g.meta), SecretOpts(g.secret), QueryOpts(map[string]string{"keyword": keyword}), ResultOpts(list))
	if keyword == "" {
		return nil, errors.New("project is empty string")
	}
	uri := fmt.Sprintf("projects/%s/coderepositories", project)
	if err := g.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}
	return list, nil
}
