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
	"k8s.io/apimachinery/pkg/api/errors"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// ClientGitRepositoryTag client for repository tag
type ClientGitRepositoryTag interface {
	Get(ctx context.Context, baseURL *duckv1.Addressable, option metav1alpha1.GitRepositoryTagOption, options ...OptionFunc) (*metav1alpha1.GitRepositoryTag, error)
	List(ctx context.Context, baseURL *duckv1.Addressable, option metav1alpha1.GitRepositoryTagListOption, options ...OptionFunc) (*metav1alpha1.GitRepositoryTagList, error)
}

type gitRepositoryTag struct {
	client Client
	meta   Meta
	secret corev1.Secret
}

func newGitRepositoryTag(client Client, meta Meta, secret corev1.Secret) ClientGitRepositoryTag {
	return &gitRepositoryTag{
		client: client,
		meta:   meta,
		secret: secret,
	}
}

// Get repository tag info
func (g *gitRepositoryTag) Get(ctx context.Context, baseURL *duckv1.Addressable, option metav1alpha1.GitRepositoryTagOption, options ...OptionFunc) (*metav1alpha1.GitRepositoryTag, error) {
	tagObj := &metav1alpha1.GitRepositoryTag{}
	options = append(options, MetaOpts(g.meta), SecretOpts(g.secret), ResultOpts(tagObj))
	if option.Repository == "" {
		return nil, errors.NewBadRequest("repo is empty string")
	} else if option.Tag == "" {
		return nil, errors.NewBadRequest("tag is null")
	}
	tagName := option.Tag
	uri := fmt.Sprintf("projects/%s/coderepositories/%s/tags/%s", option.Project, handlePathParamHasSlash(option.Repository), tagName)
	if err := g.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}
	return tagObj, nil
}

// List repository tag info
func (g *gitRepositoryTag) List(ctx context.Context, baseURL *duckv1.Addressable, option metav1alpha1.GitRepositoryTagListOption, options ...OptionFunc) (*metav1alpha1.GitRepositoryTagList, error) {
	result := &metav1alpha1.GitRepositoryTagList{}
	if option.Repository == "" {
		return nil, errors.NewBadRequest("repo is empty string")
	}
	options = append(options, MetaOpts(g.meta), SecretOpts(g.secret), ResultOpts(result))
	uri := fmt.Sprintf("projects/%s/coderepositories/%s/tags", option.Project, option.Repository)
	if err := g.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}
	return result, nil
}
