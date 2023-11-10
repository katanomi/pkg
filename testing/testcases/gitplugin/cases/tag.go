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
	. "github.com/katanomi/pkg/testing/framework"
	. "github.com/katanomi/pkg/testing/framework/base"
	. "github.com/katanomi/pkg/testing/testcases"
	. "github.com/katanomi/pkg/testing/testcases/gitplugin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/rand"
)

var TagCaseSet = func() {
	caseTagList.Do()
}

var caseTagList = P0Case("test for getting tag list").
	WithCondition(
		PluginImplementCondition{Interface: (*client.GitRepositoryTagLister)(nil)},
		NewGitOrgRepoCondition("e2e-tag-"+rand.String(5)),
	).WithFunc(func(testContext *TestContext) {
	var (
		ctx              context.Context
		instance         TestablePlugin
		gitRepo          *v1alpha1.GitRepo
		listOption       v1alpha1.ListOptions
		tagList          v1alpha1.GitRepositoryTagList
		err              error
		tag1, tag2, tag3 v1alpha1.GitRepositoryTag
	)

	BeforeEach(func() {
		ctx = testContext.Context
		instance = GitPluginFromCtx(ctx)
		gitRepo = GitRepoFromCtx(ctx)
		listOption = v1alpha1.ListOptions{}
		tagList = v1alpha1.GitRepositoryTagList{}
		err = nil
	})

	JustBeforeEach(func() {
		getter := instance.(client.GitRepositoryTagLister)
		option := v1alpha1.GitRepositoryTagListOption{GitRepo: *gitRepo}
		tagList, err = getter.ListGitRepositoryTag(ctx, option, listOption)
	})

	Context("no tags in repository", func() {
		It("should return empty list", func() {
			Expect(err).Should(Succeed())
			Expect(len(tagList.Items)).Should(Equal(0))
		})
	})

	Context("page turning of tag list", func() {
		BeforeAll(func() {
			tag1 = createTag(ctx, instance, *gitRepo, "v1.0.0", "test")
			tag2 = createTag(ctx, instance, *gitRepo, "v2.0.0", "test2")
			tag3 = createTag(ctx, instance, *gitRepo, "v3.0.0", "test2")
		})

		When("page number is 1", func() {
			BeforeEach(func() {
				listOption.Page = 1
				listOption.ItemsPerPage = 2
			})

			It("should return first page tags", func() {
				Expect(err).Should(Succeed())
				Expect(len(tagList.Items)).Should(Equal(2))
				Expect(tagList.Items[0].GetName()).Should(Equal("v3.0.0"))
				Expect(tagList.Items[1].GetName()).Should(Equal("v2.0.0"))
				Expect(tagList.Continue).Should(Equal("true"))
				checkTagRequiredFields(&tagList.Items[0], &tag3)
				checkTagOptionalFields(&tagList.Items[0], &tag3)

				checkTagRequiredFields(&tagList.Items[1], &tag2)
				checkTagOptionalFields(&tagList.Items[1], &tag2)
			})
		})

		When("page number is 2", func() {
			BeforeEach(func() {
				listOption.Page = 2
				listOption.ItemsPerPage = 2
			})

			It("should return second page tags", func() {
				Expect(err).Should(Succeed())
				Expect(len(tagList.Items)).Should(Equal(1))
				Expect(tagList.Items[0].GetName()).Should(Equal("v1.0.0"))
				Expect(tagList.Continue).Should(Equal("false"))
				checkTagRequiredFields(&tagList.Items[0], &tag1)
				checkTagOptionalFields(&tagList.Items[0], &tag1)
			})
		})

	})
})

func createTag(ctx context.Context, instance TestablePlugin, gitRepo v1alpha1.GitRepo, tagName string, branchName string) v1alpha1.GitRepositoryTag {
	branch := createBranch(ctx, instance, v1alpha1.CreateBranchPayload{
		GitRepo: gitRepo,
		CreateBranchParams: v1alpha1.CreateBranchParams{
			Branch: branchName,
		},
	})

	err := CreateNewTag(ctx, branch.GetName(), "e2e tag", tagName)
	Expect(err).Should(Succeed())
	return getTag(ctx, instance, gitRepo, tagName)
}

func getTag(textCtx context.Context, ins TestablePlugin, gitRepo v1alpha1.GitRepo, tagName string) v1alpha1.GitRepositoryTag {
	tegGetter, ok := ins.(client.GitRepositoryTagGetter)
	if !ok {
		Skip("plugin does not implement GitBranchGetter interface")
	}

	tag, err := tegGetter.GetGitRepositoryTag(textCtx, v1alpha1.GitRepositoryTagOption{
		GitRepo: gitRepo,
		Tag:     tagName,
	})
	Expect(err).Should(Succeed())
	return tag
}

func checkTagRequiredFields(gotProject, wantProject *v1alpha1.GitRepositoryTag) {
	Expect(gotProject.GetName()).ShouldNot(BeEmpty())
	Expect(gotProject.Spec.Name).ShouldNot(BeEmpty())
	Expect(gotProject.Spec.Name).Should(Equal(gotProject.GetName()))
	Expect(gotProject.Spec.SHA).ShouldNot(BeNil())

	Expect(gotProject.GetName()).To(Equal(wantProject.GetName()))
	Expect(gotProject.Spec.Name).To(Equal(wantProject.Spec.Name))
}

func checkTagOptionalFields(gotProject, wantProject *v1alpha1.GitRepositoryTag) {
	OptionalExpect(gotProject.CreationTimestamp.Time).To(Equal(wantProject.CreationTimestamp.Time))
	OptionalExpect(gotProject.Spec.Address).NotTo(BeEmpty())
	OptionalExpect(gotProject.Spec.Address.URL.String()).To(Equal(gotProject.Spec.Address.URL.String()))
}
