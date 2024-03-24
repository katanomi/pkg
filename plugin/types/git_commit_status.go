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

//go:generate mockgen -package=types -destination=../../testing/mock/github.com/katanomi/pkg/plugin/types/git_commit_status.go github.com/katanomi/pkg/plugin/types GitCommitStatusLister,GitCommitStatusCreator,GitCommitCommentLister,GitCommitCommentCreator

// GitCommitStatusLister list git commit status
type GitCommitStatusLister interface {
	Interface
	ListGitCommitStatus(
		ctx context.Context,
		option metav1alpha1.GitCommitOption,
		listOption metav1alpha1.ListOptions,
	) (metav1alpha1.GitCommitStatusList, error)
}

// GitCommitStatusCreator create git commit status
type GitCommitStatusCreator interface {
	Interface
	CreateGitCommitStatus(ctx context.Context, payload metav1alpha1.CreateCommitStatusPayload) (metav1alpha1.GitCommitStatus, error)
}

// GitCommitCommentLister list git commit comment
type GitCommitCommentLister interface {
	Interface
	ListGitCommitComment(
		ctx context.Context,
		option metav1alpha1.GitCommitOption,
		listOption metav1alpha1.ListOptions,
	) (metav1alpha1.GitCommitCommentList, error)
}

// GitCommitCommentCreator create git commit comment
type GitCommitCommentCreator interface {
	Interface
	CreateGitCommitComment(ctx context.Context, payload metav1alpha1.CreateCommitCommentPayload) (metav1alpha1.GitCommitComment, error)
}
