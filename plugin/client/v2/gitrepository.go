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

// ListGitRepository list git repositories
func (p *PluginClientV2) ListGitRepository(ctx context.Context, id, keyword string, subtype metav1alpha1.ProjectSubType, listOption metav1alpha1.ListOptions) (metav1alpha1.GitRepositoryList, error) {
	repoList := metav1alpha1.GitRepositoryList{}
	options := []base.OptionFunc{
		base.QueryOpts(map[string]string{
			"keyword": keyword,
			"subtype": subtype.String(),
		}),
		base.ResultOpts(&repoList),
		base.ListOpts(listOption),
	}
	if id == "" {
		return repoList, errors.NewBadRequest("project is empty string")
	}
	uri := path.Format("projects/%s/coderepositories", id)
	if err := p.Get(ctx, p.ClassAddress, uri, options...); err != nil {
		return repoList, err
	}
	return repoList, nil
}

// GetGitRepository get git repository
func (p *PluginClientV2) GetGitRepository(ctx context.Context, repoOption metav1alpha1.GitRepo) (metav1alpha1.GitRepository, error) {
	repo := metav1alpha1.GitRepository{}
	options := []base.OptionFunc{base.ResultOpts(&repo)}
	if repoOption.Project == "" {
		return repo, errors.NewBadRequest("project is empty string")
	}
	if repoOption.Repository == "" {
		return repo, errors.NewBadRequest("repo is empty string")
	}
	uri := path.Format("projects/%s/coderepositories/%s", repoOption.Project, repoOption.Repository)
	if err := p.Get(ctx, p.ClassAddress, uri, options...); err != nil {
		return repo, err
	}
	return repo, nil
}
