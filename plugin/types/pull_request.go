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

// GitPullRequestCommentCreator create pull request comment functions
type GitPullRequestCommentCreator interface {
	Interface
	CreatePullRequestComment(ctx context.Context, option metav1alpha1.CreatePullRequestCommentPayload) (metav1alpha1.GitPullRequestNote, error)
}

// GitPullRequestCommentUpdater updates pull request comment
type GitPullRequestCommentUpdater interface {
	Interface
	UpdatePullRequestComment(ctx context.Context, option metav1alpha1.UpdatePullRequestCommentPayload) (metav1alpha1.GitPullRequestNote, error)
}

// GitPullRequestCommentLister list pull request comment functions
type GitPullRequestCommentLister interface {
	Interface
	ListPullRequestComment(
		ctx context.Context,
		option metav1alpha1.GitPullRequestOption,
		listOption metav1alpha1.ListOptions,
	) (metav1alpha1.GitPullRequestNoteList, error)
}

// GitPullRequestHandler list, get and create pr function
type GitPullRequestHandler interface {
	Interface
	GitPullRequestLister
	GitPullRequestGetter
	GitPullRequestCreator
}

// GitPullRequestLister list pull requests
type GitPullRequestLister interface {
	Interface
	ListGitPullRequest(
		ctx context.Context,
		option metav1alpha1.GitPullRequestListOption,
		listOption metav1alpha1.ListOptions,
	) (metav1alpha1.GitPullRequestList, error)
}

// GitPullRequestGetter get a pull request
type GitPullRequestGetter interface {
	Interface
	GetGitPullRequest(ctx context.Context, option metav1alpha1.GitPullRequestOption) (metav1alpha1.GitPullRequest, error)
}

// GitPullRequestCreator create a new pull request
type GitPullRequestCreator interface {
	Interface
	CreatePullRequest(ctx context.Context, payload metav1alpha1.CreatePullRequestPayload) (metav1alpha1.GitPullRequest, error)
}
