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

	coderepositoryv1alpha1 "github.com/katanomi/pkg/apis/coderepository/v1alpha1"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
)

// GitCommitGetter get git commit
type GitCommitGetter interface {
	Interface
	GetGitCommit(ctx context.Context, option metav1alpha1.GitCommitOption) (metav1alpha1.GitCommit, error)
}

// GitCommitCreator create git commit
type GitCommitCreator interface {
	Interface
	CreateGitCommit(ctx context.Context, option coderepositoryv1alpha1.CreateGitCommitOption) (metav1alpha1.GitCommit, error)
}

// GitCommitLister List git commit
type GitCommitLister interface {
	Interface
	ListGitCommit(
		ctx context.Context,
		option metav1alpha1.GitCommitListOption,
		listOption metav1alpha1.ListOptions,
	) (metav1alpha1.GitCommitList, error)
}
