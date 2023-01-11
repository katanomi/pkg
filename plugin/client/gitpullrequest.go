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
	"strconv"

	"github.com/katanomi/pkg/plugin/path"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// ClientGitPullRequest is interface for Create and Read pull requests
// Deprecated: Integrations should implement GitPullRequestCRUClient instead (or additionally).
//
//go:generate ../../bin/mockgen -source=gitpullrequest.go -destination=../../testing/mock/github.com/katanomi/pkg/plugin/client/gitpullrequest.go -package=client ClientGitPullRequest
type ClientGitPullRequest interface {
	Create(
		ctx context.Context,
		baseURL *duckv1.Addressable,
		payload metav1alpha1.CreatePullRequestPayload,
		options ...OptionFunc,
	) (*metav1alpha1.GitPullRequest, error)
	CreateNote(
		ctx context.Context,
		baseURL *duckv1.Addressable,
		payload metav1alpha1.CreatePullRequestCommentPayload,
		options ...OptionFunc,
	) (*metav1alpha1.GitPullRequestNote, error)
	List(
		ctx context.Context,
		baseURL *duckv1.Addressable,
		option metav1alpha1.GitPullRequestListOption,
		options ...OptionFunc,
	) (*metav1alpha1.GitPullRequestList, error)
	ListNote(ctx context.Context,
		baseURL *duckv1.Addressable,
		option metav1alpha1.GitPullRequestOption,
		options ...OptionFunc,
	) (*metav1alpha1.GitPullRequestNoteList, error)
	Get(ctx context.Context,
		baseURL *duckv1.Addressable,
		option metav1alpha1.GitPullRequestOption,
		options ...OptionFunc,
	) (*metav1alpha1.GitPullRequest, error)
}

// GitPullRequestCRUClient is the client interface with Create Read and Update operations
type GitPullRequestCRUClient interface {
	ClientGitPullRequest
	UpdateNote(
		ctx context.Context,
		baseURL *duckv1.Addressable,
		payload metav1alpha1.UpdatePullRequestCommentPayload,
		options ...OptionFunc,
	) (*metav1alpha1.GitPullRequestNote, error)
}

type gitPullRequest struct {
	client Client
}

func newGitPullRequest(
	client Client,
) GitPullRequestCRUClient {
	return &gitPullRequest{
		client: client,
	}
}

// Create pr
func (g *gitPullRequest) Create(
	ctx context.Context,
	baseURL *duckv1.Addressable,
	payload metav1alpha1.CreatePullRequestPayload,
	options ...OptionFunc,
) (*metav1alpha1.GitPullRequest, error) {
	prObj := &metav1alpha1.GitPullRequest{}
	options = append(options, BodyOpts(payload), ResultOpts(prObj))
	uri := path.Format("projects/%s/coderepositories/%s/pulls", payload.Source.Project, payload.Source.Repository)
	if err := g.client.Post(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}
	return prObj, nil
}

// List pr
func (g *gitPullRequest) List(
	ctx context.Context,
	baseURL *duckv1.Addressable,
	option metav1alpha1.GitPullRequestListOption,
	options ...OptionFunc,
) (*metav1alpha1.GitPullRequestList, error) {
	prList := &metav1alpha1.GitPullRequestList{}
	options = append(options, ResultOpts(prList))
	if option.State != nil {
		stateFilter := make(map[string]string)
		stateFilter["state"] = (string)(*option.State)
		options = append(options, QueryOpts(stateFilter))
	}
	if option.Repository == "" {
		return nil, errors.NewBadRequest("repo is empty string")
	}
	uri := path.Format("projects/%s/coderepositories/%s/pulls", option.Project, option.Repository)
	if err := g.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}
	return prList, nil
}

// Get pr info
func (g *gitPullRequest) Get(
	ctx context.Context,
	baseURL *duckv1.Addressable,
	option metav1alpha1.GitPullRequestOption,
	options ...OptionFunc,
) (*metav1alpha1.GitPullRequest, error) {
	prObj := &metav1alpha1.GitPullRequest{}
	options = append(options, ResultOpts(prObj))
	if option.Repository == "" {
		return nil, errors.NewBadRequest("repo is empty string")
	}
	if option.Index < 1 {
		return nil, errors.NewBadRequest("pr's index is unknown")
	}
	index := strconv.Itoa(option.Index)
	uri := path.Format("projects/%s/coderepositories/%s/pulls/%s", option.Project, option.Repository, index)
	if err := g.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}
	return prObj, nil
}

// CreateNote create pr note
func (g *gitPullRequest) CreateNote(
	ctx context.Context,
	baseURL *duckv1.Addressable,
	payload metav1alpha1.CreatePullRequestCommentPayload,
	options ...OptionFunc,
) (*metav1alpha1.GitPullRequestNote, error) {
	noteObj := &metav1alpha1.GitPullRequestNote{}
	options = append(options, BodyOpts(payload.CreatePullRequestCommentParam), ResultOpts(noteObj))
	if payload.Repository == "" {
		return nil, errors.NewBadRequest("repo is empty string")
	}
	if payload.Index < 1 {
		return nil, errors.NewBadRequest("pr's index is unknown")
	}
	index := strconv.Itoa(payload.Index)
	uri := path.Format("projects/%s/coderepositories/%s/pulls/%s/note", payload.Project, payload.Repository, index)
	if err := g.client.Post(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}
	return noteObj, nil
}

// UpdateNote updates pr note
func (g *gitPullRequest) UpdateNote(
	ctx context.Context, baseURL *duckv1.Addressable,
	payload metav1alpha1.UpdatePullRequestCommentPayload,
	options ...OptionFunc,
) (*metav1alpha1.GitPullRequestNote, error) {
	noteObj := &metav1alpha1.GitPullRequestNote{}
	options = append(options, BodyOpts(payload.CreatePullRequestCommentParam), ResultOpts(noteObj))
	if payload.Repository == "" {
		return nil, errors.NewBadRequest("repo is empty string")
	}
	if payload.Index < 1 {
		return nil, errors.NewBadRequest("pr's index is unknown")
	}
	if payload.CommentID < 1 {
		return nil, errors.NewBadRequest("pr's comment id is unknown")
	}
	index := strconv.Itoa(payload.Index)
	commentID := strconv.Itoa(payload.CommentID)
	uri := path.Format("projects/%s/coderepositories/%s/pulls/%s/note/%s", payload.Project, payload.Repository, index, commentID)
	if err := g.client.Put(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}
	return noteObj, nil
}

// ListNote list pr note
func (g *gitPullRequest) ListNote(
	ctx context.Context,
	baseURL *duckv1.Addressable,
	option metav1alpha1.GitPullRequestOption,
	options ...OptionFunc,
) (*metav1alpha1.GitPullRequestNoteList, error) {
	noteList := &metav1alpha1.GitPullRequestNoteList{}
	options = append(options, ResultOpts(noteList))
	if option.Repository == "" {
		return nil, errors.NewBadRequest("repo is empty string")
	}
	if option.Index < 1 {
		return nil, errors.NewBadRequest("pr's index is unknown")
	}
	index := strconv.Itoa(option.Index)
	uri := path.Format("projects/%s/coderepositories/%s/pulls/%s/note", option.Project, option.Repository, index)
	if err := g.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}
	return noteList, nil
}
