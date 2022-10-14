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

	"github.com/katanomi/pkg/plugin/path"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// ClientGitBranch client for branch
type ClientGitBranch interface {
	List(ctx context.Context, baseURL *duckv1.Addressable, repo metav1alpha1.GitBranchOption, options ...OptionFunc) (*metav1alpha1.GitBranchList, error)
	Create(ctx context.Context, baseURL *duckv1.Addressable, payload metav1alpha1.CreateBranchPayload, options ...OptionFunc) (*metav1alpha1.GitBranch, error)
	Get(ctx context.Context, baseURL *duckv1.Addressable, repo metav1alpha1.GitRepo, branch string, options ...OptionFunc) (*metav1alpha1.GitBranch, error)
}

type gitBranch struct {
	client Client
	meta   Meta
	secret corev1.Secret
}

func newGitBranch(client Client, meta Meta, secret corev1.Secret) ClientGitBranch {
	return &gitBranch{
		client: client,
		meta:   meta,
		secret: secret,
	}
}

// List list branch
func (g *gitBranch) List(ctx context.Context, baseURL *duckv1.Addressable, repo metav1alpha1.GitBranchOption, options ...OptionFunc) (*metav1alpha1.GitBranchList, error) {
	list := &metav1alpha1.GitBranchList{}
	options = append(options, MetaOpts(g.meta), SecretOpts(g.secret), ResultOpts(list))
	if repo.Repository == "" {
		return nil, errors.NewBadRequest("repo is empty string")
	}
	uri := path.Format("projects/%s/coderepositories/%s/branches", repo.Project, repo.Repository)
	if err := g.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}
	return list, nil
}

// Create create branch
func (g *gitBranch) Create(ctx context.Context, baseURL *duckv1.Addressable, payload metav1alpha1.CreateBranchPayload, options ...OptionFunc) (*metav1alpha1.GitBranch, error) {
	branchObj := &metav1alpha1.GitBranch{}
	options = append(options, MetaOpts(g.meta), SecretOpts(g.secret), BodyOpts(payload.CreateBranchParams), ResultOpts(branchObj))
	if payload.Repository == "" {
		return nil, errors.NewBadRequest("repo is empty string")
	}
	uri := path.Format("projects/%s/coderepositories/%s/branches", payload.Project, payload.Repository)
	if err := g.client.Post(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}

	return branchObj, nil
}

// Get branch info
func (g *gitBranch) Get(ctx context.Context, baseURL *duckv1.Addressable, repo metav1alpha1.GitRepo, branch string, options ...OptionFunc) (*metav1alpha1.GitBranch, error) {
	branchObj := &metav1alpha1.GitBranch{}
	options = append(options, MetaOpts(g.meta), SecretOpts(g.secret), ResultOpts(branchObj))
	if repo.Repository == "" {
		return nil, errors.NewBadRequest("repo is empty string")
	}
	uri := path.Format("projects/%s/coderepositories/%s/branches/%s", repo.Project, repo.Repository, branch)
	if err := g.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}

	return branchObj, nil
}
