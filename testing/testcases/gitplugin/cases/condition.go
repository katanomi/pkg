//go:build e2e
// +build e2e

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

package cases

import (
	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/pointer"
	. "github.com/katanomi/pkg/testing/framework/base"
	. "github.com/katanomi/pkg/testing/testcases/gitplugin"
	. "github.com/onsi/ginkgo/v2"
)

func NewGitBranchCondition(branch string) *GitBranchCondition {
	return &GitBranchCondition{branch: branch}
}

type GitBranchCondition struct {
	branch string
}

func (g GitBranchCondition) Condition(testCtx *TestContext) error {
	instance := GitPluginFromCtx(testCtx.Context)
	gitRepo := GitRepoFromCtx(testCtx.Context)
	createBranch(testCtx.Context, instance, v1alpha1.CreateBranchPayload{
		GitRepo: *gitRepo,
		CreateBranchParams: v1alpha1.CreateBranchParams{
			Branch: g.branch,
		},
	})
	return nil
}

func NewGitOrgRepoCondition(repository string) *GitOrgRepoCondition {
	return &GitOrgRepoCondition{repository: repository}
}

type GitOrgRepoCondition struct {
	repository string
}

func (g GitOrgRepoCondition) Condition(testCtx *TestContext) error {
	instance := GitPluginFromCtx(testCtx.Context)
	project := instance.GetTestOrgProject()
	return NewGitRepoCondition(project, g.repository).Condition(testCtx)
}

func NewGitUserRepoCondition(repository string) *GitUserRepoCondition {
	return &GitUserRepoCondition{repository: repository}
}

type GitUserRepoCondition struct {
	repository string
}

func (g GitUserRepoCondition) Condition(testCtx *TestContext) error {
	instance := GitPluginFromCtx(testCtx.Context)
	project := instance.GetTestUserProject()
	return NewGitRepoCondition(project, g.repository).Condition(testCtx)
}

func NewGitRepoCondition(project, repository string) *GitRepoCondition {
	return &GitRepoCondition{
		Project:    project,
		Repository: repository,
	}
}

type GitRepoCondition struct {
	Project    string
	Repository string
}

func (g GitRepoCondition) Condition(testCtx *TestContext) error {
	var (
		instance = GitPluginFromCtx(testCtx.Context)
	)

	var (
		gitRepo = v1alpha1.GitRepo{
			Project:    g.Project,
			Repository: g.Repository,
		}
		localRepoPath = pointer.String("")
	)

	_, *localRepoPath = createRepository(testCtx.Context, instance, gitRepo, v1alpha1.GitRepositoryVisibilityPrivate)
	DeferCleanup(func() {
		cleanupRepository(testCtx.Context, instance, gitRepo)
	})

	testCtx.Context = WithGitRepo(testCtx.Context, &gitRepo)
	testCtx.Context = WithLocalRepoPath(testCtx.Context, localRepoPath)
	return nil
}
