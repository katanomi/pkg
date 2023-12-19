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

// ListGitCommitStatus list git commit status
func (p *PluginClientV2) ListGitCommitStatus(ctx context.Context, option metav1alpha1.GitCommitOption, listOption metav1alpha1.ListOptions) (metav1alpha1.GitCommitStatusList, error) {
	statusList := metav1alpha1.GitCommitStatusList{}
	options := []base.OptionFunc{base.ResultOpts(&statusList), base.ListOpts(listOption)}
	if option.Repository == "" {
		return statusList, errors.NewBadRequest("repo is empty string")
	} else if option.Project == "" {
		return statusList, errors.NewBadRequest("project is empty string")
	} else if option.SHA == nil {
		return statusList, errors.NewBadRequest("unknown sha for commit")
	}
	sha := *option.SHA
	uri := path.Format("/projects/%s/coderepositories/%s/commit/%s/status", option.Project, option.Repository, sha)
	if err := p.Get(ctx, p.ClassAddress, uri, options...); err != nil {
		return statusList, err
	}
	return statusList, nil
}

// CreateGitCommitStatus create git commit status
func (p *PluginClientV2) CreateGitCommitStatus(ctx context.Context, payload metav1alpha1.CreateCommitStatusPayload) (metav1alpha1.GitCommitStatus, error) {
	commitStatus := metav1alpha1.GitCommitStatus{}
	options := []base.OptionFunc{base.BodyOpts(payload.CreateCommitStatusParam), base.ResultOpts(&commitStatus)}
	if payload.Repository == "" {
		return commitStatus, errors.NewBadRequest("repo is empty string")
	} else if payload.Project == "" {
		return commitStatus, errors.NewBadRequest("project is empty string")
	} else if payload.SHA == nil {
		return commitStatus, errors.NewBadRequest("unknown sha for commit")
	}
	sha := *payload.SHA
	uri := path.Format("/projects/%s/coderepositories/%s/commit/%s/status", payload.Project, payload.Repository, sha)
	if err := p.Post(ctx, p.ClassAddress, uri, options...); err != nil {
		return commitStatus, err
	}

	return commitStatus, nil
}
