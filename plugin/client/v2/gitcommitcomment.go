/*
Copyright 2023 The Katanomi Authors.

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

package v2

import (
	"context"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client/base"
	"github.com/katanomi/pkg/plugin/path"
	"k8s.io/apimachinery/pkg/api/errors"
)

// ListGitCommitComment list git commit comment
func (p *PluginClient) ListGitCommitComment(ctx context.Context, option metav1alpha1.GitCommitOption, listOption metav1alpha1.ListOptions) (metav1alpha1.GitCommitCommentList, error) {
	commentList := metav1alpha1.GitCommitCommentList{}
	options := []base.OptionFunc{base.ResultOpts(&commentList), base.ListOpts(listOption)}
	if option.Repository == "" {
		return commentList, errors.NewBadRequest("repo is empty string")
	} else if option.Project == "" {
		return commentList, errors.NewBadRequest("project is empty string")
	} else if option.SHA == nil {
		return commentList, errors.NewBadRequest("unknown sha for commit")
	}
	sha := *option.SHA
	uri := path.Format("/projects/%s/coderepositories/%s/commit/%s/comments", option.Project, option.Repository, sha)
	if err := p.Get(ctx, p.ClassAddress, uri, options...); err != nil {
		return commentList, err
	}
	return commentList, nil
}

// CreateGitCommitComment create a git commit comment
func (p *PluginClient) CreateGitCommitComment(ctx context.Context, payload metav1alpha1.CreateCommitCommentPayload) (metav1alpha1.GitCommitComment, error) {
	comment := metav1alpha1.GitCommitComment{}
	options := []base.OptionFunc{base.BodyOpts(payload.CreateCommitCommentParam), base.ResultOpts(&comment)}
	if payload.Repository == "" {
		return comment, errors.NewBadRequest("repo is empty string")
	} else if payload.Project == "" {
		return comment, errors.NewBadRequest("project is empty string")
	} else if payload.SHA == nil {
		return comment, errors.NewBadRequest("unknown sha for commit")
	}
	sha := *payload.SHA
	uri := path.Format("/projects/%s/coderepositories/%s/commit/%s/comments", payload.Project, payload.Repository, sha)
	if err := p.Post(ctx, p.ClassAddress, uri, options...); err != nil {
		return comment, err
	}

	return comment, nil
}
