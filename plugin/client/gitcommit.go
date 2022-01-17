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
	"time"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// ClientGitCommit client for commit
type ClientGitCommit interface {
	Get(ctx context.Context, baseURL *duckv1.Addressable, option metav1alpha1.GitCommitOption, options ...OptionFunc) (*metav1alpha1.GitCommit, error)
	List(ctx context.Context, baseURL *duckv1.Addressable, option metav1alpha1.GitCommitListOption, options ...OptionFunc) (*metav1alpha1.GitCommitList, error)
}

type gitCommit struct {
	client Client
	meta   Meta
	secret corev1.Secret
}

func newGitCommit(client Client, meta Meta, secret corev1.Secret) ClientGitCommit {
	return &gitCommit{
		client: client,
		meta:   meta,
		secret: secret,
	}
}

// Get commit info
func (g *gitCommit) Get(ctx context.Context, baseURL *duckv1.Addressable, option metav1alpha1.GitCommitOption, options ...OptionFunc) (*metav1alpha1.GitCommit, error) {
	commitObj := &metav1alpha1.GitCommit{}
	options = append(options, MetaOpts(g.meta), SecretOpts(g.secret), ResultOpts(commitObj))
	if option.Repository == "" {
		return nil, errors.NewBadRequest("repo is empty string")
	} else if option.SHA == nil {
		return nil, errors.NewBadRequest("sha is null")
	} else if *option.SHA == "" {
		return nil, errors.NewBadRequest("sha is empty string")
	}
	sha := *option.SHA
	uri := fmt.Sprintf("projects/%s/coderepositories/%s/commit/%s", option.Project, handlePathParamHasSlash(option.Repository), sha)
	if err := g.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}
	return commitObj, nil
}

// List commit info
func (g *gitCommit) List(ctx context.Context, baseURL *duckv1.Addressable, option metav1alpha1.GitCommitListOption, options ...OptionFunc) (*metav1alpha1.GitCommitList, error) {
	result := &metav1alpha1.GitCommitList{}
	if option.Repository == "" {
		return nil, errors.NewBadRequest("repo is empty string")
	}
	options = append(options, MetaOpts(g.meta), SecretOpts(g.secret), ResultOpts(result))
	query := map[string]string{"ref": option.Ref}
	if option.Since != nil {
		query["since"] = option.Since.Format(time.RFC3339)
	}
	if option.Until != nil {
		query["until"] = option.Until.Format(time.RFC3339)
	}
	options = append(options, QueryOpts(query))
	uri := fmt.Sprintf("projects/%s/coderepositories/%s/commits", option.Project, option.Repository)
	if err := g.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}
	return result, nil
}
