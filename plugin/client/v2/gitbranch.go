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

package v2

import (
	"context"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client/base"
	"github.com/katanomi/pkg/plugin/path"
	"k8s.io/apimachinery/pkg/api/errors"
)

// ListGitBranch list git branches
func (p *PluginClient) ListGitBranch(ctx context.Context, branchOption metav1alpha1.GitBranchOption, option metav1alpha1.ListOptions) (metav1alpha1.GitBranchList, error) {
	list := metav1alpha1.GitBranchList{}
	repo := branchOption.GitRepo
	if repo.Repository == "" {
		return list, errors.NewBadRequest("repo is empty string")
	}
	options := []base.OptionFunc{
		base.ResultOpts(&list),
		base.ListOpts(option),
	}
	uri := path.Format("projects/%s/coderepositories/%s/branches", repo.Project, repo.Repository)
	if err := p.Get(ctx, p.ClassAddress, uri, options...); err != nil {
		return list, err
	}
	return list, nil
}

// GetGitBranch get specific git branch
func (p *PluginClient) GetGitBranch(ctx context.Context, repo metav1alpha1.GitRepo, branch string) (metav1alpha1.GitBranch, error) {
	branchObj := metav1alpha1.GitBranch{}
	options := []base.OptionFunc{base.ResultOpts(&branchObj)}
	if repo.Repository == "" {
		return branchObj, errors.NewBadRequest("repo is empty string")
	}
	uri := path.Format("projects/%s/coderepositories/%s/branches/%s", repo.Project, repo.Repository, branch)
	if err := p.Get(ctx, p.ClassAddress, uri, options...); err != nil {
		return branchObj, err
	}

	return branchObj, nil
}

// CreateGitBranch create git branch
func (p *PluginClient) CreateGitBranch(ctx context.Context, payload metav1alpha1.CreateBranchPayload) (metav1alpha1.GitBranch, error) {
	branchObj := metav1alpha1.GitBranch{}
	options := []base.OptionFunc{base.BodyOpts(payload.CreateBranchParams), base.ResultOpts(&branchObj)}
	if payload.Repository == "" {
		return branchObj, errors.NewBadRequest("repo is empty string")
	}
	uri := path.Format("projects/%s/coderepositories/%s/branches", payload.Project, payload.Repository)
	if err := p.Post(ctx, p.ClassAddress, uri, options...); err != nil {
		return branchObj, err
	}

	return branchObj, nil
}
