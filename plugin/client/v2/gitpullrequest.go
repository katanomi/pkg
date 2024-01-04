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
	"strconv"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client/base"
	"github.com/katanomi/pkg/plugin/path"
	"k8s.io/apimachinery/pkg/api/errors"
)

// CreatePullRequest create a pull request
func (p *PluginClient) CreatePullRequest(ctx context.Context, payload metav1alpha1.CreatePullRequestPayload) (metav1alpha1.GitPullRequest, error) {
	pr := metav1alpha1.GitPullRequest{}
	options := []base.OptionFunc{base.BodyOpts(payload), base.ResultOpts(&pr)}
	uri := path.Format("projects/%s/coderepositories/%s/pulls", payload.Source.Project, payload.Source.Repository)
	if err := p.Post(ctx, p.ClassAddress, uri, options...); err != nil {
		return pr, err
	}
	return pr, nil
}

// ListGitPullRequest list pull requests
func (p *PluginClient) ListGitPullRequest(ctx context.Context, option metav1alpha1.GitPullRequestListOption, listOption metav1alpha1.ListOptions) (metav1alpha1.GitPullRequestList, error) {
	prList := metav1alpha1.GitPullRequestList{}
	if option.Repository == "" {
		return prList, errors.NewBadRequest("repo is empty string")
	}

	query := make(map[string]string)
	if option.State != nil {
		query["state"] = (string)(*option.State)
	}
	if len(option.Commit) > 0 {
		query["commit"] = option.Commit
	}

	options := []base.OptionFunc{base.ResultOpts(&prList), base.ListOpts(listOption), base.QueryOpts(query)}
	uri := path.Format("projects/%s/coderepositories/%s/pulls", option.Project, option.Repository)
	if err := p.Get(ctx, p.ClassAddress, uri, options...); err != nil {
		return prList, err
	}
	return prList, nil
}

// GetGitPullRequest get a pull request
func (p *PluginClient) GetGitPullRequest(ctx context.Context, option metav1alpha1.GitPullRequestOption) (metav1alpha1.GitPullRequest, error) {
	pr := metav1alpha1.GitPullRequest{}
	options := []base.OptionFunc{base.ResultOpts(&pr)}
	if option.Repository == "" {
		return pr, errors.NewBadRequest("repo is empty string")
	}
	if option.Index < 1 {
		return pr, errors.NewBadRequest("pr's index is unknown")
	}
	index := strconv.Itoa(option.Index)
	uri := path.Format("projects/%s/coderepositories/%s/pulls/%s", option.Project, option.Repository, index)
	if err := p.Get(ctx, p.ClassAddress, uri, options...); err != nil {
		return pr, err
	}
	return pr, nil
}

// CreatePullRequestComment create a pull request comment
func (p *PluginClient) CreatePullRequestComment(ctx context.Context, payload metav1alpha1.CreatePullRequestCommentPayload) (metav1alpha1.GitPullRequestNote, error) {
	note := metav1alpha1.GitPullRequestNote{}
	options := []base.OptionFunc{base.BodyOpts(payload.CreatePullRequestCommentParam), base.ResultOpts(&note)}
	if payload.Repository == "" {
		return note, errors.NewBadRequest("repo is empty string")
	}
	if payload.Index < 1 {
		return note, errors.NewBadRequest("pr's index is unknown")
	}
	index := strconv.Itoa(payload.Index)
	uri := path.Format("projects/%s/coderepositories/%s/pulls/%s/note", payload.Project, payload.Repository, index)
	if err := p.Post(ctx, p.ClassAddress, uri, options...); err != nil {
		return note, err
	}
	return note, nil
}

// UpdatePullRequestComment update a pull request comment
func (p *PluginClient) UpdatePullRequestComment(ctx context.Context, payload metav1alpha1.UpdatePullRequestCommentPayload) (metav1alpha1.GitPullRequestNote, error) {
	note := metav1alpha1.GitPullRequestNote{}
	options := []base.OptionFunc{base.BodyOpts(payload.CreatePullRequestCommentParam), base.ResultOpts(&note)}
	if payload.Repository == "" {
		return note, errors.NewBadRequest("repo is empty string")
	}
	if payload.Index < 1 {
		return note, errors.NewBadRequest("pr's index is unknown")
	}
	if payload.CommentID < 1 {
		return note, errors.NewBadRequest("pr's comment id is unknown")
	}
	index := strconv.Itoa(payload.Index)
	commentID := strconv.Itoa(payload.CommentID)
	uri := path.Format("projects/%s/coderepositories/%s/pulls/%s/note/%s", payload.Project, payload.Repository, index, commentID)
	if err := p.Put(ctx, p.ClassAddress, uri, options...); err != nil {
		return note, err
	}
	return note, nil
}

// ListPullRequestComment list pull request comments
func (p *PluginClient) ListPullRequestComment(ctx context.Context, option metav1alpha1.GitPullRequestOption, listOption metav1alpha1.ListOptions) (metav1alpha1.GitPullRequestNoteList, error) {
	noteList := metav1alpha1.GitPullRequestNoteList{}
	options := []base.OptionFunc{base.ResultOpts(&noteList), base.ListOpts(listOption)}
	if option.Repository == "" {
		return noteList, errors.NewBadRequest("repo is empty string")
	}
	if option.Index < 1 {
		return noteList, errors.NewBadRequest("pr's index is unknown")
	}
	index := strconv.Itoa(option.Index)
	uri := path.Format("projects/%s/coderepositories/%s/pulls/%s/note", option.Project, option.Repository, index)
	if err := p.Get(ctx, p.ClassAddress, uri, options...); err != nil {
		return noteList, err
	}
	return noteList, nil
}
