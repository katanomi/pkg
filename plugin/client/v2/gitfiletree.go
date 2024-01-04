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

// GetGitRepositoryFileTree get git repository file tree
func (p *PluginClient) GetGitRepositoryFileTree(ctx context.Context, option metav1alpha1.GitRepoFileTreeOption, listOption metav1alpha1.ListOptions) (metav1alpha1.GitRepositoryFileTree, error) {
	fileTree := metav1alpha1.GitRepositoryFileTree{}
	if option.Repository == "" {
		return fileTree, errors.NewBadRequest("repo is empty string")
	} else if option.Path == "" {
		return fileTree, errors.NewBadRequest("file path is empty string")
	}

	recursiveValue := strconv.FormatBool(option.Recursive)
	options := []base.OptionFunc{
		base.QueryOpts(map[string]string{
			"path":      option.Path,
			"tree_sha":  option.TreeSha,
			"recursive": recursiveValue,
		}),
		base.ResultOpts(&fileTree),
		base.ListOpts(listOption),
	}

	uri := path.Format("projects/%s/coderepositories/%s/tree", option.Project, option.Repository)
	if err := p.Get(ctx, p.ClassAddress, uri, options...); err != nil {
		return fileTree, err
	}
	return fileTree, nil
}
