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
	"k8s.io/apimachinery/pkg/api/errors"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

type ClientGitCommitComment interface {
	List(ctx context.Context, baseURL *duckv1.Addressable, option metav1alpha1.GitCommitOption, options ...OptionFunc) (*metav1alpha1.GitCommitCommentList, error)
	Create(ctx context.Context, baseURL *duckv1.Addressable, payload metav1alpha1.CreateCommitCommentPayload, options ...OptionFunc) (*metav1alpha1.GitCommitComment, error)
}

type gitCommitComment struct {
	client Client
}

func newGitCommitComment(client Client) ClientGitCommitComment {
	return &gitCommitComment{
		client: client,
	}
}

func (g *gitCommitComment) List(ctx context.Context, baseURL *duckv1.Addressable, option metav1alpha1.GitCommitOption, options ...OptionFunc) (*metav1alpha1.GitCommitCommentList, error) {
	commitCommentList := &metav1alpha1.GitCommitCommentList{}
	options = append(options, ResultOpts(commitCommentList))
	if option.Repository == "" {
		return nil, errors.NewBadRequest("repo is empty string")
	} else if option.Project == "" {
		return nil, errors.NewBadRequest("project is empty string")
	} else if option.SHA == nil {
		return nil, errors.NewBadRequest("unknown sha for commit")
	}
	sha := *option.SHA
	uri := path.Format("/projects/%s/coderepositories/%s/commit/%s/comments", option.Project, option.Repository, sha)
	if err := g.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}
	return commitCommentList, nil
}

func (g *gitCommitComment) Create(ctx context.Context, baseURL *duckv1.Addressable, payload metav1alpha1.CreateCommitCommentPayload, options ...OptionFunc) (*metav1alpha1.GitCommitComment, error) {
	commentInfo := &metav1alpha1.GitCommitComment{}
	options = append(options, BodyOpts(payload.CreateCommitCommentParam), ResultOpts(commentInfo))
	if payload.Repository == "" {
		return nil, errors.NewBadRequest("repo is empty string")
	} else if payload.Project == "" {
		return nil, errors.NewBadRequest("project is empty string")
	} else if payload.SHA == nil {
		return nil, errors.NewBadRequest("unknown sha for commit")
	}
	sha := *payload.SHA
	uri := path.Format("/projects/%s/coderepositories/%s/commit/%s/comments", payload.Project, payload.Repository, sha)
	if err := g.client.Post(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}

	return commentInfo, nil
}
