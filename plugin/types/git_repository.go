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

//go:generate mockgen -package=types -destination=../../testing/mock/github.com/katanomi/pkg/plugin/types/git_repository.go github.com/katanomi/pkg/plugin/types GitRepositoryCreator,GitRepositoryDeleter,GitRepositoryLister,GitRepositoryGetter

// GitRepositoryCreator create a git repository
type GitRepositoryCreator interface {
	Interface
	CreateGitRepository(ctx context.Context, payload metav1alpha1.CreateGitRepositoryPayload) (metav1alpha1.GitRepository, error)
}

// GitRepositoryDeleter delete a git repository
type GitRepositoryDeleter interface {
	Interface
	DeleteGitRepository(ctx context.Context, gitRepo metav1alpha1.GitRepo) error
}

// GitRepositoryLister list git repository
type GitRepositoryLister interface {
	Interface
	ListGitRepository(
		ctx context.Context,
		id, keyword string,
		subtype metav1alpha1.ProjectSubType,
		listOption metav1alpha1.ListOptions,
	) (metav1alpha1.GitRepositoryList, error)
}

// GitRepositoryGetter get git repository
type GitRepositoryGetter interface {
	Interface
	GetGitRepository(ctx context.Context, repoOption metav1alpha1.GitRepo) (metav1alpha1.GitRepository, error)
}
