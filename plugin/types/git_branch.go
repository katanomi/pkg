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

// GitBranchLister List git branch
type GitBranchLister interface {
	Interface
	ListGitBranch(ctx context.Context, branchOption metav1alpha1.GitBranchOption, option metav1alpha1.ListOptions) (metav1alpha1.GitBranchList, error)
}

// GitBranchGetter get git branch
type GitBranchGetter interface {
	Interface
	GetGitBranch(ctx context.Context, repoOption metav1alpha1.GitRepo, branch string) (metav1alpha1.GitBranch, error)
}

// GitBranchCreator create git branch,github, gogs don't support create branch
type GitBranchCreator interface {
	Interface
	CreateGitBranch(ctx context.Context, payload metav1alpha1.CreateBranchPayload) (metav1alpha1.GitBranch, error)
}
