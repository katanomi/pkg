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
	"context"

	"github.com/google/go-cmp/cmp"
	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
	. "github.com/katanomi/pkg/testing/framework"
	. "github.com/katanomi/pkg/testing/framework/base"
	. "github.com/katanomi/pkg/testing/testcases/gitplugin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/rand"
)

var FileTreeCaseSet = func() {
	caseFileTreeList.Do()
}

var caseFileTreeList = P0Case("test for getting file tree").
	WithCondition(
		PluginImplementCondition{Interface: (*client.GitRepositoryFileTreeGetter)(nil)},
		NewGitUserRepoCondition("e2e-filetree-"+rand.String(5)),
		NewGitBranchCondition("main"),
	).WithFunc(func(testContext *TestContext) {
	var (
		ctx        context.Context
		instance   TestablePlugin
		repoOption *v1alpha1.GitRepoFileTreeOption
		listOption *v1alpha1.ListOptions
		fileTree   v1alpha1.GitRepositoryFileTree
		getter     client.GitRepositoryFileTreeGetter
		err        error
	)

	BeforeEach(func() {
		ctx = testContext.Context
		instance = GitPluginFromCtx(ctx)
		gitrepo := GitRepoFromCtx(ctx)
		repoOption = &v1alpha1.GitRepoFileTreeOption{
			GitRepo: *gitrepo,
		}
		listOption = &v1alpha1.ListOptions{}
		fileTree = v1alpha1.GitRepositoryFileTree{}
		getter = instance.(client.GitRepositoryFileTreeGetter)
	})

	JustBeforeEach(func() {
		fileTree, err = getter.GetGitRepositoryFileTree(ctx, *repoOption, *listOption)
	})

	Context("get fileTree list with ref", func() {
		BeforeEach(func() {
			repoOption.TreeSha = "e2e-filetree-branch-" + rand.String(5)
			gitrepo := GitRepoFromCtx(ctx)
			createBranch(ctx, instance, v1alpha1.CreateBranchPayload{
				GitRepo: *gitrepo,
				CreateBranchParams: v1alpha1.CreateBranchParams{
					Branch: repoOption.TreeSha,
				},
			})

			createCommit(ctx, repoOption.TreeSha, "commit message")
		})

		It("return specify branch filetree", func() {
			Expect(err).Should(BeNil())
			Expect(fileTree.Spec.Tree).Should(HaveLen(3))
		})
	})

	Context("get fileTree list with recursive", func() {
		BeforeEach(func() {
			repoOption.Recursive = true
		})

		It("return all node", func() {
			Expect(err).Should(BeNil())
			Expect(fileTree.Spec.Tree).Should(HaveLen(5))
		})
	})

	Context("get fileTree list with path", func() {
		When("list recursive is true", func() {
			BeforeEach(func() {
				repoOption.Path = "a/b"
				repoOption.Recursive = true
			})

			It("return specify path node", func() {
				Expect(err).Should(BeNil())
				Expect(fileTree.Spec.Tree).Should(HaveLen(2))
			})
		})

		When("list recursive is false", func() {
			BeforeEach(func() {
				repoOption.Path = "a/b"
			})

			It("return specify path node", func() {
				Expect(err).Should(BeNil())
				Expect(fileTree.Spec.Tree).Should(HaveLen(1))
			})
		})

		When("list recursive with /", func() {
			BeforeEach(func() {
				repoOption.Path = "/"
			})

			It("return specify path node", func() {
				Expect(err).Should(BeNil())
				Expect(fileTree.Spec.Tree).Should(HaveLen(2))
			})
		})

		When("list recursive with ./", func() {
			BeforeEach(func() {
				repoOption.Path = "./"
			})

			It("return specify path node", func() {
				Expect(err).Should(BeNil())
				Expect(fileTree.Spec.Tree).Should(HaveLen(2))
			})
		})

		When("list recursive with .", func() {
			BeforeEach(func() {
				repoOption.Path = "."
			})

			It("return specify path node", func() {
				Expect(err).Should(BeNil())
				Expect(fileTree.Spec.Tree).Should(HaveLen(2))
			})
		})
	})
})

func checkRequitedFileTree(got, want *v1alpha1.GitRepositoryFileTree) {
	Expect(got.GetName()).ShouldNot(BeEmpty())

	Expect(got.GetName()).Should(Equal(want.GetName()))
}

func checkRequitedTreeNode(got, want *v1alpha1.GitRepositoryFileTreeNode) {
	Expect(got.Sha).ShouldNot(BeEmpty())
	Expect(got.Name).ShouldNot(BeEmpty())
	Expect(got.Path).ShouldNot(BeEmpty())
	Expect(got.Type).ShouldNot(BeEmpty())

	diff := cmp.Diff(got, want)
	Expect(diff).Should(BeEmpty())
}
