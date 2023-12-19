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
	"time"

	coderepositoryv1alpha1 "github.com/katanomi/pkg/apis/coderepository/v1alpha1"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client/base"
	"github.com/katanomi/pkg/plugin/path"
	"k8s.io/apimachinery/pkg/api/errors"
)

// GetGitCommit get git commit detail
func (p *PluginClientV2) GetGitCommit(ctx context.Context, option metav1alpha1.GitCommitOption) (metav1alpha1.GitCommit, error) {
	gitCommit := metav1alpha1.GitCommit{}
	options := []base.OptionFunc{base.ResultOpts(&gitCommit)}
	if option.Repository == "" {
		return gitCommit, errors.NewBadRequest("repo is empty string")
	} else if option.SHA == nil {
		return gitCommit, errors.NewBadRequest("sha is null")
	} else if *option.SHA == "" {
		return gitCommit, errors.NewBadRequest("sha is empty string")
	}
	sha := *option.SHA
	uri := path.Format("projects/%s/coderepositories/%s/commit/%s", option.Project, option.Repository, sha)
	if err := p.Get(ctx, p.ClassAddress, uri, options...); err != nil {
		return gitCommit, err
	}
	return gitCommit, nil
}

// CreateGitCommit create git commit
func (p *PluginClientV2) CreateGitCommit(ctx context.Context, option coderepositoryv1alpha1.CreateGitCommitOption) (metav1alpha1.GitCommit, error) {
	gitCommit := metav1alpha1.GitCommit{}
	options := []base.OptionFunc{base.ResultOpts(&gitCommit), base.BodyOpts(option)}
	if err := option.GitRepo.Validate(); err != nil {
		return gitCommit, errors.NewBadRequest(err.Error())
	}
	if errs := option.GitCreateCommit.Validate(ctx); len(errs) != 0 {
		return gitCommit, errors.NewBadRequest(errs.ToAggregate().Error())
	}
	uri := path.Format("projects/%s/coderepositories/%s/commits", option.Project, option.Repository)
	if err := p.Post(ctx, p.ClassAddress, uri, options...); err != nil {
		return gitCommit, err
	}
	return gitCommit, nil
}

// ListGitCommit list git commits
func (p *PluginClientV2) ListGitCommit(ctx context.Context, option metav1alpha1.GitCommitListOption, listOption metav1alpha1.ListOptions) (metav1alpha1.GitCommitList, error) {
	commitList := metav1alpha1.GitCommitList{}
	if option.Repository == "" {
		return commitList, errors.NewBadRequest("repo is empty string")
	}

	query := map[string]string{"ref": option.Ref}
	if option.Since != nil {
		query["since"] = option.Since.Format(time.RFC3339)
	}
	if option.Until != nil {
		query["until"] = option.Until.Format(time.RFC3339)
	}
	options := []base.OptionFunc{base.ResultOpts(&commitList), base.ListOpts(listOption), base.QueryOpts(query)}
	uri := path.Format("projects/%s/coderepositories/%s/commits", option.Project, option.Repository)
	if err := p.Get(ctx, p.ClassAddress, uri, options...); err != nil {
		return commitList, err
	}
	return commitList, nil
}
