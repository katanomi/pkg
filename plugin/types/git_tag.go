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

//go:generate mockgen -package=types -destination=../../testing/mock/github.com/katanomi/pkg/plugin/types/git_tag.go github.com/katanomi/pkg/plugin/types GitRepositoryTagCreator,GitRepositoryTagGetter,GitRepositoryTagLister

// GitRepositoryTagCreator create git repository tag
type GitRepositoryTagCreator interface {
	Interface
	CreateGitRepositoryTag(ctx context.Context, option metav1alpha1.CreateGitTagPayload) (metav1alpha1.GitRepositoryTag, error)
}

// GitRepositoryTagGetter get git repository Tag
type GitRepositoryTagGetter interface {
	Interface
	GetGitRepositoryTag(
		ctx context.Context,
		option metav1alpha1.GitRepositoryTagOption,
	) (metav1alpha1.GitRepositoryTag, error)
}

// GitRepositoryTagLister list git repository Tag
type GitRepositoryTagLister interface {
	Interface
	ListGitRepositoryTag(
		ctx context.Context,
		option metav1alpha1.GitRepositoryTagListOption,
		listOption metav1alpha1.ListOptions,
	) (metav1alpha1.GitRepositoryTagList, error)
}
