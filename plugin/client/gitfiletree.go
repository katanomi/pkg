/*
Copyright 2022 The Katanomi Authors.

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
	"strconv"

	"github.com/katanomi/pkg/plugin/path"

	"k8s.io/apimachinery/pkg/api/errors"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// ClientGitRepositoryFileTree defines the request interface for the file tree
type ClientGitRepositoryFileTree interface {
	GetGitRepositoryFileTree(
		ctx context.Context,
		baseURL *duckv1.Addressable,
		option metav1alpha1.GitRepoFileTreeOption,
		options ...OptionFunc,
	) (*metav1alpha1.GitRepositoryFileTree, error)
}

type gitRepositoryFileTree struct {
	client Client
}

// init gitRepositoryFileTree
func newGitRepositoryFileTree(client Client) ClientGitRepositoryFileTree {
	return &gitRepositoryFileTree{
		client: client,
	}
}

// GetGitRepositoryFileTree call the integrations api
func (g *gitRepositoryFileTree) GetGitRepositoryFileTree(
	ctx context.Context,
	baseURL *duckv1.Addressable,
	option metav1alpha1.GitRepoFileTreeOption,
	options ...OptionFunc,
) (*metav1alpha1.GitRepositoryFileTree, error) {
	fileTree := &metav1alpha1.GitRepositoryFileTree{}
	recursiveValue := strconv.FormatBool(option.Recursive)

	options = append(
		options,
		QueryOpts(map[string]string{"path": option.Path, "tree_sha": option.TreeSha, "recursive": recursiveValue}),
		ResultOpts(fileTree),
	)
	if option.Repository == "" {
		return nil, errors.NewBadRequest("repo is empty string")
	} else if option.Path == "" {
		return nil, errors.NewBadRequest("file path is empty string")
	}
	uri := path.Format("projects/%s/coderepositories/%s/tree", option.Project, option.Repository)
	if err := g.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}
	return fileTree, nil
}
