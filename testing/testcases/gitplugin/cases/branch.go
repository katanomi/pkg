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

	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
	"github.com/katanomi/pkg/pointer"
	. "github.com/katanomi/pkg/testing/framework"
	. "github.com/katanomi/pkg/testing/framework/base"
	. "github.com/katanomi/pkg/testing/testcases"
	. "github.com/katanomi/pkg/testing/testcases/gitplugin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/rand"
)

var BranchCaseSet = func() {
	_ = caseBranchGet.Do()
	_ = caseBranchList.Do()
}

var caseBranchGet = P0Case("test getting branch info").
	WithCondition(
		PluginImplementCondition{Interface: (*client.GitBranchGetter)(nil)},
		NewGitOrgRepoCondition("e2e-branch-"+rand.String(5)),
	).WithFunc(func(testContext *TestContext) {
	var (
		testBranch *string
		ctx        context.Context
		instance   TestablePlugin
		gitRepo    *v1alpha1.GitRepo
		branch     v1alpha1.GitBranch
		err        error
	)

	BeforeEach(func() {
		ctx = testContext.Context
		instance = GitPluginFromCtx(ctx)
		gitRepo = GitRepoFromCtx(ctx)
		testBranch = pointer.String("")
		branch = v1alpha1.GitBranch{}
		err = nil
	})

	JustBeforeEach(func() {
		getter := instance.(client.GitBranchGetter)
		branch, err = getter.GetGitBranch(ctx, *gitRepo, *testBranch)
	})

	Context("specified branch not exist", func() {
		BeforeEach(func() {
			*testBranch = "not-exist-branch"
		})
		It("should return 404", func() {
			Expect(err).Should(HaveOccurred())
			Expect(errors.IsNotFound(err)).Should(BeTrue())
		})
	})

	Context("specified branch is exist", func() {
		var expectedBranch v1alpha1.GitBranch
		BeforeEach(func() {
			*testBranch = "exist-branch"
			expectedBranch = createBranch(ctx, instance, v1alpha1.CreateBranchPayload{
				GitRepo: *gitRepo,
				CreateBranchParams: v1alpha1.CreateBranchParams{
					Branch: *testBranch,
				},
			})
		})
		It("should return the branch", func() {
			Expect(err).Should(Succeed())
			checkBranchRequiredFields(&branch, &expectedBranch)
			checkBranchOptionalFields(&branch, &expectedBranch)
		})
	})
})

var caseBranchList = P0Case("test getting branch list").
	WithCondition(
		PluginImplementCondition{Interface: (*client.GitBranchLister)(nil)},
		NewGitUserRepoCondition("e2e-branch-list-"+rand.String(5)),
	).WithFunc(func(testContext *TestContext) {
	var (
		ctx      context.Context
		instance TestablePlugin
		gitRepo  *v1alpha1.GitRepo

		branchList        v1alpha1.GitBranchList
		branchTwoPageList v1alpha1.GitBranchList
		branchOption      *v1alpha1.GitBranchOption
		listOption        *v1alpha1.ListOptions
		err               error

		getter client.GitBranchLister
	)

	BeforeAll(func() {
		branchList := []string{
			"test-a", "test-b", "page-a", "page-b", "page-c", "page-d",
		}

		ctx = testContext.Context
		gitRepo = GitRepoFromCtx(ctx)
		for _, item := range branchList {
			createBranch(ctx, instance, v1alpha1.CreateBranchPayload{
				GitRepo: *gitRepo,
				CreateBranchParams: v1alpha1.CreateBranchParams{
					Branch: item,
				},
			})
		}
	})

	BeforeEach(func() {
		ctx = testContext.Context
		instance = GitPluginFromCtx(ctx)
		gitRepo = GitRepoFromCtx(ctx)
		branchList = v1alpha1.GitBranchList{}
		branchTwoPageList = v1alpha1.GitBranchList{}
		branchOption = &v1alpha1.GitBranchOption{
			GitRepo: *gitRepo,
		}
		listOption = &v1alpha1.ListOptions{}
		err = nil
		getter = instance.(client.GitBranchLister)
	})

	JustBeforeEach(func() {
		branchList, err = getter.ListGitBranch(ctx, *branchOption, *listOption)
	})

	Context("specified branch not found", func() {
		notExistBranchName := "not-found"
		BeforeEach(func() {
			listOption.Search = map[string][]string{
				v1alpha1.SearchValueKey: {notExistBranchName},
			}
		})

		It("should return empty list", func() {
			Expect(err).Should(BeNil())
			Expect(branchList).ShouldNot(BeNil())
			gotBranch := FindByName(ToPtrList(branchList.Items), notExistBranchName)
			Expect(gotBranch).Should(BeNil())
		})
	})

	Context("get branch list successfully", func() {
		BeforeEach(func() {
			listOption.Search = map[string][]string{
				v1alpha1.SearchValueKey: {"test"},
			}
		})

		It("get return the branch list which contains the test branch", func() {
			Expect(err).Should(BeNil())

			By("the number of return branches should be greater than 2")
			Expect(branchList).ShouldNot(BeNil())
			Expect(len(branchList.Items) > 0).Should(BeTrue())

			By("the branch test-a should come before test-b")
			branchAIndex := FindIndexByName(ToPtrList(branchList.Items), "test-a")
			branchBIndex := FindIndexByName(ToPtrList(branchList.Items), "test-b")
			Expect(branchAIndex < branchBIndex).Should(BeTrue())

			By("information of branch a should as expected")
			branchA := branchList.Items[branchAIndex]
			Expect(branchA.Spec.Project).To(Equal(gitRepo.Project))
			Expect(branchA.Spec.Repository).To(Equal(gitRepo.Repository))
			Expect(branchA.Spec.Name).To(Equal("test-a"))
			Expect(*branchA.Spec.Protected).To(BeFalse())
			Expect(*branchA.Spec.Default).To(BeFalse())
		})
	})

	Context("page turning of branch list", func() {
		BeforeEach(func() {
			branchTwoPageList, err = getter.ListGitBranch(ctx, *branchOption, *listOption)
			Expect(err).Should(BeNil())
		})

		When("set page number 1", func() {
			BeforeEach(func() {
				listOption.Page = 1
				listOption.ItemsPerPage = 2
			})
			It("return first page items", func() {
				Expect(err).Should(BeNil())
				Expect(branchList).ShouldNot(BeNil())
				Expect(branchList.Items).Should(HaveLen(2))
				Expect(branchList.Items).Should(Equal(branchTwoPageList.Items[:2]))
			})
		})

		When("set page number 2", func() {
			BeforeEach(func() {
				listOption.Page = 2
				listOption.ItemsPerPage = 2
			})

			It("return second page items", func() {
				Expect(err).Should(BeNil())
				Expect(branchList).ShouldNot(BeNil())
				Expect(branchList.Items).To(HaveLen(2))
				Expect(branchList.Items).Should(Equal(branchTwoPageList.Items[2:4]))
			})
		})

	})
})

func checkBranchOptionalFields(wantBranch, gotBranch *v1alpha1.GitBranch) {
	// these fields are optional
	OptionalExpect(gotBranch.Spec.Commit.CreatedAt).To(Equal(wantBranch.Spec.Commit.CreatedAt))
	OptionalExpect(gotBranch.Spec.Default).ToNot(BeNil(), "field of `default branch` is not set")
	OptionalExpect(gotBranch.Spec.Protected).ToNot(BeNil(), "field of `protected` is not set")
	OptionalExpect(gotBranch.Spec.DevelopersCanPush).ToNot(BeNil(), "field of `developers can push` is not set")
	OptionalExpect(gotBranch.Spec.DevelopersCanMerge).ToNot(BeNil(), "field of `developers can merge` is not set")
}

func checkBranchRequiredFields(wantBranch, gotBranch *v1alpha1.GitBranch) {
	Expect(gotBranch.GetName()).To(Equal(wantBranch.GetName()))
	Expect(gotBranch.Spec.Project).To(Equal(wantBranch.Spec.Project))
	Expect(gotBranch.Spec.Repository).To(Equal(wantBranch.Spec.Repository))
	Expect(gotBranch.Spec.Name).To(Equal(wantBranch.Spec.Name))
	Expect(gotBranch.Spec.Commit.SHA).To(Equal(wantBranch.Spec.Commit.SHA))
}

func getBranch(textCtx context.Context, ins TestablePlugin, gitRepo v1alpha1.GitRepo, branchName string) v1alpha1.GitBranch {
	branchGetter, ok := ins.(client.GitBranchGetter)
	if !ok {
		Skip("plugin does not implement GitBranchGetter interface")
	}

	branch, err := branchGetter.GetGitBranch(textCtx, gitRepo, branchName)
	Expect(err).Should(Succeed())
	return branch
}

func createBranch(testCtx context.Context, ins TestablePlugin, payload v1alpha1.CreateBranchPayload) v1alpha1.GitBranch {
	err := CreateNewBranch(testCtx, payload.Branch)
	Expect(err).Should(Succeed())
	return getBranch(testCtx, ins, payload.GitRepo, payload.Branch)
}
