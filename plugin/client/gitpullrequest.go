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
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"

	duckv1 "knative.dev/pkg/apis/duck/v1"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
)

type ClientGitPullRequest interface {
	Create(ctx context.Context, baseURL *duckv1.Addressable, payload metav1alpha1.CreatePullRequestPayload, options ...OptionFunc) (*metav1alpha1.GitPullRequest, error)
	CreateNote(ctx context.Context, baseURL *duckv1.Addressable, payload metav1alpha1.CreatePullRequestCommentPayload, options ...OptionFunc) (*metav1alpha1.GitPullRequestNote, error)
	List(ctx context.Context, baseURL *duckv1.Addressable, option metav1alpha1.GitRepo, options ...OptionFunc) (*metav1alpha1.GitPullRequestList, error)
	Get(ctx context.Context, baseURL *duckv1.Addressable, option metav1alpha1.GitPullRequestOption, options ...OptionFunc) (*metav1alpha1.GitPullRequest, error)
}

type gitPullRequest struct {
	client Client
	meta   Meta
	secret corev1.Secret
}

func newGitPullRequest(client Client, meta Meta, secret corev1.Secret) ClientGitPullRequest {
	return &gitPullRequest{
		client: client,
		meta:   meta,
		secret: secret,
	}
}

// Create pr
func (g *gitPullRequest) Create(ctx context.Context, baseURL *duckv1.Addressable, payload metav1alpha1.CreatePullRequestPayload, options ...OptionFunc) (*metav1alpha1.GitPullRequest, error) {
	prObj := &metav1alpha1.GitPullRequest{}
	options = append(options, MetaOpts(g.meta), SecretOpts(g.secret), BodyOpts(payload), ResultOpts(prObj))
	repoInfo := strings.Split(payload.Source.Repository, "/")
	if len(repoInfo) != 2 {
		return nil, errors.New("repo info should include project and repo name")
	} else if repoInfo[1] == "" {
		return nil, errors.New("repo name is empty")
	}
	uri := fmt.Sprintf("projects/%s/coderepositories/%s/pulls", repoInfo[0], repoInfo[1])
	if err := g.client.Post(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}
	return prObj, nil
}

// List pr
func (g *gitPullRequest) List(ctx context.Context, baseURL *duckv1.Addressable, option metav1alpha1.GitRepo, options ...OptionFunc) (*metav1alpha1.GitPullRequestList, error) {
	prList := &metav1alpha1.GitPullRequestList{}
	options = append(options, MetaOpts(g.meta), SecretOpts(g.secret), ResultOpts(prList))
	if option.Repository == "" {
		return nil, errors.New("repo is empty string")
	}
	uri := fmt.Sprintf("projects/%s/coderepositories/%s/pulls", option.Project, option.Repository)
	if err := g.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}
	return prList, nil
}

// Get pr info
func (g *gitPullRequest) Get(ctx context.Context, baseURL *duckv1.Addressable, option metav1alpha1.GitPullRequestOption, options ...OptionFunc) (*metav1alpha1.GitPullRequest, error) {
	prObj := &metav1alpha1.GitPullRequest{}
	options = append(options, MetaOpts(g.meta), SecretOpts(g.secret), ResultOpts(prObj))
	if option.Repository == "" {
		return nil, errors.New("repo is empty string")
	}
	if option.Index < 1 {
		return nil, errors.New("pr's index is unknown")
	}
	index := strconv.Itoa(option.Index)
	uri := fmt.Sprintf("projects/%s/coderepositories/%s/pulls/%s", option.Project, option.Repository, index)
	if err := g.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}
	return prObj, nil
}

// CreateNote create pr note
func (g *gitPullRequest) CreateNote(ctx context.Context, baseURL *duckv1.Addressable, payload metav1alpha1.CreatePullRequestCommentPayload, options ...OptionFunc) (*metav1alpha1.GitPullRequestNote, error) {
	noteObj := &metav1alpha1.GitPullRequestNote{}
	options = append(options, MetaOpts(g.meta), SecretOpts(g.secret), BodyOpts(payload.CreatePullRequestCommentParam), ResultOpts(noteObj))
	if payload.Repository == "" {
		return nil, errors.New("repo is empty string")
	}
	if payload.Index < 1 {
		return nil, errors.New("pr's index is unknown")
	}
	index := strconv.Itoa(payload.Index)
	uri := fmt.Sprintf("projects/%s/coderepositories/%s/pulls/%s", payload.Project, payload.Repository, index)
	if err := g.client.Post(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}
	return noteObj, nil
}
