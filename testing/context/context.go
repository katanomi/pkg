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

package context

import (
	"context"

	"github.com/katanomi/pkg/apis/meta/v1alpha1"
)

type gitRepositoryCtxKey struct{}

// WithGitRepository returns a copy of parent in which the gitRepo value is set
func WithGitRepository(ctx context.Context, gitRepo *v1alpha1.GitRepository) context.Context {
	return context.WithValue(ctx, gitRepositoryCtxKey{}, gitRepo)
}

// GitRepositoryFromCtx returns the value of the gitRepo key on the ctx
func GitRepositoryFromCtx(ctx context.Context) *v1alpha1.GitRepository {
	value := ctx.Value(gitRepositoryCtxKey{})
	if value == nil {
		return nil
	}
	return value.(*v1alpha1.GitRepository)
}

type gitRepoCtxKey struct{}

// WithGitRepo returns a copy of parent in which the gitRepo value is set
func WithGitRepo(ctx context.Context, gitRepo *v1alpha1.GitRepo) context.Context {
	return context.WithValue(ctx, gitRepoCtxKey{}, gitRepo)
}

// GitRepoFromCtx returns the value of the gitRepo key on the ctx
func GitRepoFromCtx(ctx context.Context) *v1alpha1.GitRepo {
	value := ctx.Value(gitRepoCtxKey{})
	if value == nil {
		return nil
	}
	return value.(*v1alpha1.GitRepo)
}

type gitRepoLocalPathCtxKey struct{}

// WithLocalRepoPath returns a copy of parent in which the localRepoPath value is set
func WithLocalRepoPath(ctx context.Context, localRepoPath *string) context.Context {
	return context.WithValue(ctx, gitRepoLocalPathCtxKey{}, localRepoPath)
}

// LocalRepoPathFromCtx returns the value of the localRepoPath key on the ctx
func LocalRepoPathFromCtx(ctx context.Context) *string {
	value := ctx.Value(gitRepoLocalPathCtxKey{})
	if value == nil {
		return nil
	}
	return value.(*string)
}
