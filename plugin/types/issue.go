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

//go:generate mockgen -package=types -destination=../../testing/mock/github.com/katanomi/pkg/plugin/types/issue.go github.com/katanomi/pkg/plugin/types IssueLister,IssueGetter,IssueBranchLister,IssueBranchCreator,IssueBranchDeleter,IssueAttributeGetter

// IssueLister issue lister
type IssueLister interface {
	Interface
	ListIssues(ctx context.Context, params metav1alpha1.IssueOptions, option metav1alpha1.ListOptions) (*metav1alpha1.IssueList, error)
}

type IssueGetter interface {
	Interface
	GetIssue(ctx context.Context, params metav1alpha1.IssueOptions, option metav1alpha1.ListOptions) (*metav1alpha1.Issue, error)
}

type IssueBranchLister interface {
	Interface
	ListIssueBranches(ctx context.Context, params metav1alpha1.IssueOptions, option metav1alpha1.ListOptions) (*metav1alpha1.BranchList, error)
}

type IssueBranchCreator interface {
	Interface
	CreateIssueBranch(ctx context.Context, params metav1alpha1.IssueOptions, payload metav1alpha1.Branch) (*metav1alpha1.Branch, error)
}

type IssueBranchDeleter interface {
	Interface
	DeleteIssueBranch(ctx context.Context, params metav1alpha1.IssueOptions, option metav1alpha1.ListOptions) error
}

type IssueAttributeGetter interface {
	Interface
	GetIssueAttribute(ctx context.Context, params metav1alpha1.IssueOptions, option metav1alpha1.ListOptions) (*metav1alpha1.Attribute, error)
}
