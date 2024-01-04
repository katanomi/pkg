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
	"encoding/base64"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client/base"
	"github.com/katanomi/pkg/plugin/path"
	"k8s.io/apimachinery/pkg/api/errors"
)

// GetGitRepoFile get a repo file content
func (p *PluginClient) GetGitRepoFile(ctx context.Context, option metav1alpha1.GitRepoFileOption) (metav1alpha1.GitRepoFile, error) {
	repoFile := metav1alpha1.GitRepoFile{}
	options := []base.OptionFunc{base.QueryOpts(map[string]string{"ref": option.Ref}), base.ResultOpts(&repoFile)}
	if option.Repository == "" {
		return repoFile, errors.NewBadRequest("repo is empty string")
	}
	if option.Path == "" {
		return repoFile, errors.NewBadRequest("file path is empty string")
	}
	uri := path.Format("projects/%s/coderepositories/%s/content/%s", option.Project, option.Repository, option.Path)
	if err := p.Get(ctx, p.ClassAddress, uri, options...); err != nil {
		return repoFile, err
	}
	return repoFile, nil
}

// CreateGitRepoFile create a new repo file
func (p *PluginClient) CreateGitRepoFile(ctx context.Context, payload metav1alpha1.CreateRepoFilePayload) (metav1alpha1.GitCommit, error) {
	commit := metav1alpha1.GitCommit{}

	if payload.Repository == "" {
		return commit, errors.NewBadRequest("repo is empty string")
	}

	var encoded = make([]byte, base64.StdEncoding.EncodedLen(len(payload.Content)))
	base64.StdEncoding.Encode(encoded, payload.Content)
	payload.Content = encoded

	options := []base.OptionFunc{base.BodyOpts(payload.CreateRepoFileParams), base.ResultOpts(&commit)}
	uri := path.Format("projects/%s/coderepositories/%s/content/%s", payload.Project, payload.Repository, payload.FilePath)
	if err := p.Post(ctx, p.ClassAddress, uri, options...); err != nil {
		return commit, err
	}

	return commit, nil
}
