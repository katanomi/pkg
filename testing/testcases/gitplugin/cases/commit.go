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
	"fmt"

	codev1alpha1 "github.com/katanomi/pkg/apis/coderepository/v1alpha1"
	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
	. "github.com/katanomi/pkg/testing/framework"
	. "github.com/katanomi/pkg/testing/framework/base"
	. "github.com/katanomi/pkg/testing/testcases/gitplugin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/rand"
)

var CommitCaseSet = func() {
	caseCommitList.Do()
}

var caseCommitList = P0Case("test for getting commit list").
	WithCondition(
		PluginImplementCondition{Interface: (*client.GitCommitLister)(nil)},
		NewGitOrgRepoCondition("e2e-commit-"+rand.String(5)),
		//NewGitBranchCondition("commit-branch"),
	).WithFunc(func(testContext *TestContext) {
	var (
		ctx          context.Context
		instance     TestablePlugin
		option       *v1alpha1.GitCommitListOption
		commitList   v1alpha1.GitCommitList
		listOption   *v1alpha1.ListOptions
		getter       client.GitCommitLister
		commitBranch = "commit-branch"
		err          error
	)

	BeforeAll(func() {
		ctx = testContext.Context
		instance = GitPluginFromCtx(ctx)
		payload := codev1alpha1.CreateGitCommitOption{
			GitCreateCommit: codev1alpha1.GitCreateCommit{
				Spec: codev1alpha1.GitCreateCommitSpec{
					Branch: commitBranch,
				},
			},
		}

		for i := 0; i <= 10; i++ {
			payload.Spec.Message = fmt.Sprintf("commit message %d", i)
			createCommit(ctx, instance, payload)
		}
	})

	BeforeEach(func() {
		gitrepo := GitRepoFromCtx(ctx)
		option = &v1alpha1.GitCommitListOption{
			GitRepo: *gitrepo,
			Ref:     commitBranch,
		}
		listOption = &v1alpha1.ListOptions{}
		commitList = v1alpha1.GitCommitList{}
		err = nil
		getter = instance.(client.GitCommitLister)
	})

	JustBeforeEach(func() {
		commitList, err = getter.ListGitCommit(ctx, *option, *listOption)
		Expect(err).Should(BeNil())
	})

	Context("search commit with message", func() {
		When("commit message not found", func() {
			BeforeEach(func() {
				listOption.Search = map[string][]string{
					v1alpha1.SearchValueKey: {"not-found"},
				}
			})

			It("not commit return", func() {
				Expect(err).Should(BeNil())
				Expect(commitList.Items).Should(BeEmpty())
			})
		})

		When("commit message exits", func() {
			BeforeEach(func() {
				listOption.Search = map[string][]string{
					v1alpha1.SearchValueKey: {"1"},
				}
				listOption.Page = 1
				listOption.ItemsPerPage = 5
			})

			It("return commit messgae include search value", func() {
				Expect(err).Should(BeNil())
				Expect(commitList.Items).Should(HaveLen(1))
				Expect(commitList.Continue).Should(Equal("true"))
				Expect(*commitList.Items[0].Spec.Message).Should(ContainSubstring("1"))
			})
		})
	})

	Context("page turning of commit list", func() {
		When("page number is 1", func() {
			BeforeEach(func() {
				listOption.Page = 1
				listOption.ItemsPerPage = 2
			})

			It("return first commit", func() {
				Expect(err).Should(BeNil())
				Expect(commitList.Items).Should(HaveLen(2))
				Expect(*commitList.Items[0].Spec.Message).Should(Equal("commit message 10"))
				Expect(*commitList.Items[1].Spec.Message).Should(Equal("commit message 9"))
				Expect(commitList.Continue).Should(Equal("true"))
			})
		})

		When("page number is 2", func() {
			BeforeEach(func() {
				listOption.Page = 2
				listOption.ItemsPerPage = 2
			})

			It("return second commit", func() {
				Expect(err).Should(BeNil())
				Expect(commitList.Items).Should(HaveLen(2))
				Expect(*commitList.Items[0].Spec.Message).Should(Equal("commit message 8"))
				Expect(*commitList.Items[1].Spec.Message).Should(Equal("commit message 7"))
				Expect(commitList.Continue).Should(Equal("true"))
			})
		})
	})

	Context("commit list query since/util", func() {
		var wantCommit v1alpha1.GitCommit
		BeforeEach(func() {
			allCommitList, err := getter.ListGitCommit(ctx, *option, v1alpha1.ListOptions{})
			Expect(err).Should(BeNil())

			wantCommit = allCommitList.Items[2]
			option.Since = wantCommit.CreationTimestamp.DeepCopy()
			option.Until = wantCommit.CreationTimestamp.DeepCopy()
		})

		It("return special commit", func() {
			Expect(err).Should(BeNil())
			Expect(commitList.Items).Should(HaveLen(1))
			Expect(commitList.Continue).Should(Equal("false"))
			checkCommitRequiredFields(&commitList.Items[0], &wantCommit)
		})
	})
})

func createCommit(testCtx context.Context, _ TestablePlugin, payload codev1alpha1.CreateGitCommitOption) {
	err := CreatNewCommit(testCtx, payload.Spec.Branch, payload.Spec.Message)
	Expect(err).Should(Succeed())
	return
}

func checkCommitRequiredFields(gotCommit, wantCommit *v1alpha1.GitCommit) {
	Expect(gotCommit.GetName()).ShouldNot(BeEmpty())
	Expect(gotCommit.Spec.SHA).ShouldNot(BeNil())
	Expect(gotCommit.Spec.CreatedAt).ShouldNot(BeZero())
	Expect(gotCommit.Spec.Address).ShouldNot(BeNil())
	Expect(gotCommit.Spec.Address.URL.String()).ShouldNot(BeEmpty())
	Expect(gotCommit.Spec.Author).ShouldNot(BeNil())
	Expect(gotCommit.Spec.Committer).ShouldNot(BeNil())
	Expect(gotCommit.Spec.Message).ShouldNot(BeNil())

	Expect(gotCommit.GetName()).Should(Equal(wantCommit.GetName()))
	Expect(*gotCommit.Spec.SHA).Should(Equal(*wantCommit.Spec.SHA))
	Expect(gotCommit.Spec.CreatedAt).Should(Equal(wantCommit.Spec.CreatedAt))
	Expect(gotCommit.Spec.Address.URL.String()).Should(Equal(wantCommit.Spec.Address.URL.String()))
	Expect(*gotCommit.Spec.Message).Should(Equal(*wantCommit.Spec.Message))
}
