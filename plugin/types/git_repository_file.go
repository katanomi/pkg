/*
Copyright 2024 The Katanomi Authors.

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

package types

import (
	"context"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
)

//go:generate mockgen -package=types -destination=../../testing/mock/github.com/katanomi/pkg/plugin/types/git_repository_file.go github.com/katanomi/pkg/plugin/types GitRepoFileGetter,GitRepoFileCreator,GitRepositoryFileTreeGetter

// GitRepoFileGetter used to get a file content
type GitRepoFileGetter interface {
	Interface
	GetGitRepoFile(ctx context.Context, option metav1alpha1.GitRepoFileOption) (metav1alpha1.GitRepoFile, error)
}

// GitRepoFileCreator used to create a file, gogs don't support
type GitRepoFileCreator interface {
	Interface
	CreateGitRepoFile(ctx context.Context, payload metav1alpha1.CreateRepoFilePayload) (metav1alpha1.GitCommit, error)
}

// GitRepositoryFileTreeGetter get git repository file tree
type GitRepositoryFileTreeGetter interface {
	Interface
	GetGitRepositoryFileTree(
		ctx context.Context,
		repoOption metav1alpha1.GitRepoFileTreeOption,
		listOption metav1alpha1.ListOptions,
	) (metav1alpha1.GitRepositoryFileTree, error)
}
