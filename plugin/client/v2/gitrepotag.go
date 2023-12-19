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

// GetGitRepositoryTag get git repository tag
func (p *PluginClientV2) GetGitRepositoryTag(ctx context.Context, option metav1alpha1.GitTag) (metav1alpha1.GitRepositoryTag, error) {
	tag := metav1alpha1.GitRepositoryTag{}
	options := []base.OptionFunc{base.ResultOpts(&tag)}
	if option.Repository == "" {
		return tag, errors.NewBadRequest("repo is empty string")
	} else if option.Tag == "" {
		return tag, errors.NewBadRequest("tag is null")
	}
	uri := path.Format("projects/%s/coderepositories/%s/tags/%s", option.Project, option.Repository, option.Tag)
	if err := p.Get(ctx, p.ClassAddress, uri, options...); err != nil {
		return tag, err
	}
	return tag, nil
}

// ListGitRepositoryTag list git repository tags
func (p *PluginClientV2) ListGitRepositoryTag(ctx context.Context, option metav1alpha1.GitRepositoryTagListOption, listOption metav1alpha1.ListOptions) (metav1alpha1.GitRepositoryTagList, error) {
	tagList := metav1alpha1.GitRepositoryTagList{}
	if option.Repository == "" {
		return tagList, errors.NewBadRequest("repo is empty string")
	}
	options := []base.OptionFunc{base.ResultOpts(&tagList), base.ListOpts(listOption)}
	uri := path.Format("projects/%s/coderepositories/%s/tags", option.Project, option.Repository)
	if err := p.Get(ctx, p.ClassAddress, uri, options...); err != nil {
		return tagList, err
	}
	return tagList, nil
}
