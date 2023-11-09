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
	"sort"

	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
	"github.com/katanomi/pkg/plugin/route"
	. "github.com/katanomi/pkg/testing/framework"
	. "github.com/katanomi/pkg/testing/framework/base"
	. "github.com/katanomi/pkg/testing/testcases"
	. "github.com/katanomi/pkg/testing/testcases/gitplugin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/rand"
)

var RepositoryCaseSet = func() {
	caseGitRepositoryList.Do()
	caseGetGitRepository.Do()
}

var caseGitRepositoryList = P0Case("test for getting repo list").
	WithCondition(
		PluginImplementCondition{Interface: (*client.GitRepositoryLister)(nil)},
	).WithFunc(func(testContext *TestContext) {
	var (
		ctx      context.Context
		instance TestablePlugin

		err     error
		repoErr error

		repoGitList   v1alpha1.GitRepositoryList
		wantGitRepo   v1alpha1.GitRepository
		gitRepo       v1alpha1.GitRepo
		gitRepoGetter client.GitRepositoryLister

		repoList   *v1alpha1.RepositoryList
		wantRepo   *v1alpha1.Repository
		repoParams v1alpha1.RepositoryOptions
		repoLister client.RepositoryLister
	)

	BeforeEach(func() {
		ctx = testContext.Context
		instance = GitPluginFromCtx(ctx)

		gitRepoGetter = instance.(client.GitRepositoryLister)
		repoLister = instance.(client.RepositoryLister)
		repoGitList = v1alpha1.GitRepositoryList{}
		repoList = nil
		err = nil
		repoErr = nil
		gitRepo = v1alpha1.GitRepo{}

	})

	Context("search repository", func() {
		var projectSubtype v1alpha1.ProjectSubType
		BeforeEach(func() {
			repoParams.SubType = ""
		})

		JustBeforeEach(func() {
			repoGitList, err = gitRepoGetter.ListGitRepository(ctx, gitRepo.Project, gitRepo.Repository, projectSubtype, v1alpha1.ListOptions{})
			repoList, repoErr = repoLister.ListRepositories(ctx, repoParams, v1alpha1.ListOptions{
				Search: map[string][]string{
					route.SearchQueryKey: {repoParams.Repository},
				},
			})
		})

		Context("subtype of project is User", func() {
			BeforeEach(func() {
				gitRepo.Project = instance.GetTestUserProject()
				gitRepo.Repository = "e2e-user-repo-" + rand.String(4)
				repoParams.Project = gitRepo.Project
				if gitRepo.Project == "" {
					Skip("plugin does not support user project")
				}
				repoParams.Repository = gitRepo.Repository
				projectSubtype = v1alpha1.GitUserProjectSubType
				repoParams.SubType = v1alpha1.GitUserProjectSubType

				wantGitRepo, _ = createRepository(ctx, instance, gitRepo, v1alpha1.GitRepositoryVisibilityPrivate)
				DeferCleanup(func() {
					cleanupRepository(ctx, instance, gitRepo)
				})

				repoGetter := instance.(client.RepositoryGetter)
				wantRepo, _ = repoGetter.GetRepository(ctx, repoParams)
			})

			It("get return the specified git repositories", func() {
				Expect(err).Should(Succeed())
				Expect(repoGitList.Items).ShouldNot(BeEmpty())
				gotRepo := FindByName(ToPtrList(repoGitList.Items), gitRepo.Repository)
				checkRequitedGitRepository(gotRepo, &wantGitRepo)
			})

			It("get return the specified repositories", func() {
				Expect(repoErr).Should(Succeed())
				Expect(repoList.Items).ShouldNot(BeEmpty())
				Expect(wantRepo).ShouldNot(BeNil())
				gotRepo := FindByName(ToPtrList(repoList.Items), gitRepo.Repository)
				checkRequitedRepository(gotRepo, wantRepo)
			})
		})

		Context("subtype of project is Group", func() {
			BeforeEach(func() {
				gitRepo.Project = instance.GetTestOrgProject()
				gitRepo.Repository = "e2e-org-repo-" + rand.String(4)
				repoParams.Project = gitRepo.Project
				repoParams.Repository = gitRepo.Repository

				projectSubtype = v1alpha1.GitGroupProjectSubType
				repoParams.SubType = v1alpha1.GitGroupProjectSubType
				if gitRepo.Project == "" {
					Skip("plugin does not support user project")
				}

				wantGitRepo, _ = createRepository(ctx, instance, gitRepo, v1alpha1.GitRepositoryVisibilityPrivate)
				DeferCleanup(func() {
					cleanupRepository(ctx, instance, gitRepo)
				})

				repoGetter := instance.(client.RepositoryGetter)
				wantRepo, _ = repoGetter.GetRepository(ctx, repoParams)
			})

			It("get return the specified git repositories", func() {
				Expect(err).Should(BeNil())
				Expect(repoGitList.Items).ShouldNot(BeEmpty())
				gotRepo := FindByName(ToPtrList(repoGitList.Items), gitRepo.Repository)
				checkRequitedGitRepository(gotRepo, &wantGitRepo)
			})

			It("get return the specified repositories", func() {
				Expect(repoErr).Should(Succeed())
				Expect(repoList.Items).ShouldNot(BeEmpty())
				Expect(wantRepo).ShouldNot(BeNil())
				gotRepo := FindByName(ToPtrList(repoList.Items), gitRepo.Repository)
				checkRequitedRepository(gotRepo, wantRepo)
			})
		})
	})

	Context("list git repository successfully", func() {
		var (
			listOption         v1alpha1.ListOptions
			projectSubtype     v1alpha1.ProjectSubType
			gitRepoTwoPageList v1alpha1.GitRepositoryList

			repoTwoPageList *v1alpha1.RepositoryList
		)

		BeforeAll(func() {
			gitRepo.Project = instance.GetTestOrgProject()
			projectSubtype = v1alpha1.GitGroupProjectSubType

			for i := 0; i < 5; i++ {
				_gitRepo := gitRepo
				_gitRepo.Repository = "e2e-repo-list-" + rand.String(4)
				createRepository(ctx, instance, _gitRepo, v1alpha1.GitRepositoryVisibilityPrivate)
				DeferCleanup(func() {
					cleanupRepository(ctx, instance, _gitRepo)
				})
			}
		})

		BeforeEach(func() {
			gitRepo.Project = instance.GetTestOrgProject()
			repoParams.Project = instance.GetTestOrgProject()
			repoParams.SubType = projectSubtype

			listOption = v1alpha1.ListOptions{}
			gitRepoTwoPageList = v1alpha1.GitRepositoryList{}
			repoTwoPageList = nil
		})

		Context("page turning of repository list", func() {
			JustBeforeEach(func() {
				listOption.ItemsPerPage = 5
				listOption.Page = 1
				gitRepoTwoPageList, err = gitRepoGetter.ListGitRepository(ctx, gitRepo.Project, "", projectSubtype, listOption)
				repoTwoPageList, repoErr = repoLister.ListRepositories(ctx, repoParams, listOption)
			})

			When("page number is 1", func() {
				BeforeEach(func() {
					listOption.Page = 1
					listOption.ItemsPerPage = 2
					repoGitList, err = gitRepoGetter.ListGitRepository(ctx, gitRepo.Project, "", projectSubtype, listOption)
					repoList, repoErr = repoLister.ListRepositories(ctx, repoParams, listOption)
				})

				It("get return first page git repositories", func() {
					Expect(err).Should(BeNil())
					Expect(repoGitList.Items).ShouldNot(BeEmpty())
					Expect(repoGitList.Items).To(Equal(gitRepoTwoPageList.Items[0:2]))
					index := rand.Intn(2)
					checkRequitedGitRepository(&repoGitList.Items[index], &gitRepoTwoPageList.Items[index])
				})

				It("get return first page repositories", func() {
					Expect(repoErr).Should(BeNil())
					Expect(repoList.Items).ShouldNot(BeEmpty())
					Expect(repoList.Items).To(Equal(repoTwoPageList.Items[0:2]))
					index := rand.Intn(2)
					checkRequitedRepository(&repoList.Items[index], &repoTwoPageList.Items[index])
				})
			})

			When("page number is 2", func() {
				BeforeEach(func() {
					listOption.Page = 2
					listOption.ItemsPerPage = 2
					repoGitList, err = gitRepoGetter.ListGitRepository(ctx, gitRepo.Project, "", projectSubtype, listOption)
					repoList, repoErr = repoLister.ListRepositories(ctx, repoParams, listOption)
				})

				It("get return second page git repositories", func() {
					Expect(err).Should(BeNil())
					Expect(repoGitList.Items).ShouldNot(BeEmpty())
					Expect(repoGitList.Items).To(Equal(gitRepoTwoPageList.Items[2:4]))
					index := rand.Intn(2)
					checkRequitedGitRepository(&repoGitList.Items[index], &gitRepoTwoPageList.Items[index+2])
				})

				It("get return second page repositories", func() {
					Expect(repoErr).Should(BeNil())
					Expect(repoList.Items).ShouldNot(BeEmpty())
					Expect(repoList.Items).To(Equal(repoTwoPageList.Items[2:4]))
					index := rand.Intn(2)
					checkRequitedRepository(&repoList.Items[index], &repoTwoPageList.Items[index+2])
				})
			})
		})

		Context("filter by subresources", func() {
			//check whether the warehouse meets the testing requirements,
			BeforeEach(func() {
				gitRepo.Project = instance.GetTestUserProject()
				projectSubtype = v1alpha1.GitUserProjectSubType
				repoParams.Project = gitRepo.Project
				repoParams.SubType = v1alpha1.GitUserProjectSubType

				gitRepoTwoPageList, err = gitRepoGetter.ListGitRepository(ctx, gitRepo.Project, "", projectSubtype, listOption)
				Expect(err).Should(BeNil())
				repoTwoPageList, repoErr = repoLister.ListRepositories(ctx, repoParams, listOption)
				Expect(repoErr).Should(BeNil())
			})

			When("subresources was set", func() {
				BeforeEach(func() {
					index := rand.Intn(1)
					wantGitRepo = gitRepoTwoPageList.Items[index]
					wantRepo = &repoTwoPageList.Items[index]

					listOption.SubResources = []string{wantGitRepo.GetName()}
					repoGitList, err = gitRepoGetter.ListGitRepository(ctx, gitRepo.Project, "", projectSubtype, listOption)
					repoList, repoErr = repoLister.ListRepositories(ctx, repoParams, listOption)
				})
				It("return with subresource git repositories", func() {
					Expect(err).Should(BeNil())
					Expect(repoGitList.Items).ShouldNot(BeEmpty())
					gotRepo := FindByName(ToPtrList(repoGitList.Items), wantGitRepo.GetName())
					checkRequitedGitRepository(gotRepo, &wantGitRepo)
				})

				It("return with subresource repositories", func() {
					Expect(repoErr).Should(BeNil())
					Expect(repoList.Items).ShouldNot(BeEmpty())
					gotRepo := FindByName(ToPtrList(repoList.Items), wantRepo.GetName())
					checkRequitedRepository(gotRepo, wantRepo)
				})
			})
		})

		Context("list git repository with sort", func() {
			JustBeforeEach(func() {
				repoParams.Project = instance.GetTestOrgProject()
				repoParams.SubType = projectSubtype

				repoGitList, err = gitRepoGetter.ListGitRepository(ctx, gitRepo.Project, "", projectSubtype, listOption)
				Expect(err).Should(BeNil())

				repoList, repoErr = repoLister.ListRepositories(ctx, repoParams, listOption)
				Expect(repoErr).Should(BeNil())
			})

			When("default sort is name,  and asc", func() {
				It("sort git repository by name", func() {
					Expect(err).Should(BeNil())
					Expect(repoGitList.Items).ShouldNot(BeEmpty())
					sortRpos := repoGitList.DeepCopy()
					sort.SliceStable(sortRpos.Items, func(i, j int) bool {
						return sortRpos.Items[i].GetName() < sortRpos.Items[j].GetName()
					})
					Expect(repoGitList.Items).To(Equal(sortRpos.Items))
				})

				It("sort by name", func() {
					Expect(repoErr).Should(BeNil())
					Expect(repoList.Items).ShouldNot(BeEmpty())
					sortRpos := repoList.DeepCopy()
					sort.SliceStable(sortRpos.Items, func(i, j int) bool {
						return sortRpos.Items[i].GetName() < sortRpos.Items[j].GetName()
					})
					Expect(repoList.Items).To(Equal(sortRpos.Items))
				})
			})
			When("sort by createTime, and asc", func() {
				BeforeEach(func() {
					listOption.Sort = []v1alpha1.SortOptions{{SortBy: v1alpha1.CreatedTimeSortKey, Order: v1alpha1.OrderAsc}}
				})
				It("sort git repository by createTime, and desc", func() {
					Expect(err).Should(BeNil())
					Expect(repoGitList.Items).ShouldNot(BeEmpty())
					sortRpos := repoGitList.DeepCopy()
					sort.SliceStable(sortRpos.Items, func(i, j int) bool {
						return sortRpos.Items[i].GetCreationTimestamp().Unix() < sortRpos.Items[j].GetCreationTimestamp().Unix()
					})
					Expect(repoGitList.Items).To(Equal(sortRpos.Items))
				})

				It("sort by createTime, and desc", func() {
					Expect(repoErr).Should(BeNil())
					Expect(repoList.Items).ShouldNot(BeEmpty())
					sortRpos := repoList.DeepCopy()
					sort.SliceStable(sortRpos.Items, func(i, j int) bool {
						return sortRpos.Items[i].GetCreationTimestamp().Unix() < sortRpos.Items[j].GetCreationTimestamp().Unix()
					})
					Expect(repoList.Items).To(Equal(sortRpos.Items))
				})
			})
			When("sort by updateTime, and desc", func() {
				BeforeEach(func() {
					listOption.Sort = []v1alpha1.SortOptions{{SortBy: v1alpha1.UpdatedTimeSortKey, Order: v1alpha1.OrderDesc}}
				})
				It("sort git repository by updateTime", func() {
					Expect(err).Should(BeNil())
					Expect(repoGitList.Items).ShouldNot(BeEmpty())
					sortRpos := repoGitList.DeepCopy()
					sort.SliceStable(sortRpos.Items, func(i, j int) bool {
						return sortRpos.Items[i].Spec.UpdatedAt.Unix() > sortRpos.Items[j].Spec.UpdatedAt.Unix()
					})
					Expect(repoGitList.Items).To(Equal(sortRpos.Items))
				})

				It("sort by updateTime", func() {
					Expect(repoErr).Should(BeNil())
					Expect(repoList.Items).ShouldNot(BeEmpty())
					sortRpos := repoList.DeepCopy()
					sort.SliceStable(sortRpos.Items, func(i, j int) bool {
						return sortRpos.Items[i].Spec.UpdatedTime.Unix() > sortRpos.Items[j].Spec.UpdatedTime.Unix()
					})
					Expect(repoList.Items).To(Equal(sortRpos.Items))
				})
			})
		})
	})
})

var caseGetGitRepository = P0Case("test for getting repo list").
	WithCondition(
		PluginImplementCondition{Interface: (*client.GitRepositoryGetter)(nil)},
	).WithFunc(func(testContext *TestContext) {
	var (
		ctx      context.Context
		instance TestablePlugin
		gitRepo  v1alpha1.GitRepo
		repo     v1alpha1.GitRepository
		err      error
	)

	BeforeEach(func() {
		ctx = testContext.Context
		instance = GitPluginFromCtx(ctx)
		repo = v1alpha1.GitRepository{}
		gitRepo.Project = instance.GetTestUserProject()
		if gitRepo.Project == "" {
			Skip("plugin does not support user gitrepository")
		}
		err = nil
	})

	JustBeforeEach(func() {
		gitRepoGetter := instance.(client.GitRepositoryGetter)
		repo, err = gitRepoGetter.GetGitRepository(ctx, gitRepo)
	})

	Context("gitrepository not exist", func() {
		BeforeEach(func() {
			gitRepo.Repository = "not-exist-not-exist-not-exist"
		})
		It("get return 404", func() {
			Expect(errors.IsNotFound(err)).Should(BeTrue())
		})
	})

	Context("gitrepository exist", func() {
		var wantGitRepo v1alpha1.GitRepository
		BeforeEach(func() {
			wantGitRepo, _ = createRepository(ctx, instance, gitRepo, v1alpha1.GitRepositoryVisibilityPrivate)
			DeferCleanup(func() {
				cleanupRepository(ctx, instance, gitRepo)
			})
			gitRepo.Repository = wantGitRepo.GetName()
		})

		It("should return the project", func() {
			Expect(err).Should(BeNil())
			Expect(repo.GetName()).Should(Equal(wantGitRepo.GetName())) // TODO: more asserts
			checkRequitedGitRepository(&repo, &wantGitRepo)
			checkOptionalGitRepository(&repo, &wantGitRepo)
		})
	})
})

func checkRequitedGitRepository(got, want *v1alpha1.GitRepository) {
	Expect(got.GetName()).ShouldNot(BeEmpty())
	Expect(got.Spec.DefaultBranch).ShouldNot(BeEmpty())
	Expect(got.Spec.HttpCloneURL).ShouldNot(BeEmpty())

	Expect(got.GetName()).Should(Equal(want.GetName()))
	Expect(got.Spec.Name).Should(Equal(want.Spec.Name))
	Expect(got.Spec.DefaultBranch).Should(Equal(want.Spec.DefaultBranch))
}

func checkOptionalGitRepository(got, want *v1alpha1.GitRepository) {
	OptionalExpect(got.Spec.Owner.Name).Should(Equal(got.Spec.Owner.Name))
	OptionalExpect(got.Spec.HtmlURL).Should(Equal(want.Spec.HtmlURL))
	OptionalExpect(got.Spec.HttpCloneURL).Should(Equal(want.Spec.HttpCloneURL))
	OptionalExpect(got.Spec.SshCloneURL).Should(Equal(want.Spec.SshCloneURL))
}

func checkRequitedRepository(got, want *v1alpha1.Repository) {
	Expect(got.GetName()).ShouldNot(BeEmpty())
	Expect(got.GetCreationTimestamp()).ShouldNot(BeZero())
	Expect(got.Spec.Access).ShouldNot(BeNil())
	Expect(got.Spec.Access.URL.String()).ShouldNot(BeEmpty())
	Expect(got.Spec.Type).ShouldNot(BeEmpty())

	Expect(got.GetName()).Should(Equal(want.GetName()))
	Expect(got.Spec.Access.URL.String()).Should(Equal(want.Spec.Access.URL.String()))
	Expect(got.Spec.Type).Should(Equal(want.Spec.Type))
}

func checkOptionalRepository(got, want *v1alpha1.Repository) {
	OptionalExpect(got.Spec.UpdatedTime).Should(Equal(got.Spec.UpdatedTime))
	OptionalExpect(got.Spec.Properties).Should(Equal(want.Spec.Properties))
}

func getRepository(testCtx context.Context, ins TestablePlugin, gitRepo v1alpha1.GitRepo) v1alpha1.GitRepository {
	repoGetter, ok := ins.(client.GitRepositoryGetter)
	if !ok {
		Skip("plugin does not implement GitRepositoryCreator interface")
		return v1alpha1.GitRepository{}
	}

	repo, err := repoGetter.GetGitRepository(testCtx, gitRepo)
	Expect(err).Should(Succeed())
	Expect(repo).ShouldNot(BeNil())
	return repo
}

func cleanupRepository(testCtx context.Context, ins TestablePlugin, gitRepo v1alpha1.GitRepo) {
	repoDeleter, ok := ins.(client.GitRepositoryDeleter)
	if !ok {
		Skip("plugin does not implement GitRepositoryCreator interface")
		return
	}

	err := repoDeleter.DeleteGitRepository(testCtx, gitRepo)
	Expect(err).Should(Succeed())
}

func createRepository(testCtx context.Context, ins TestablePlugin, gitRepo v1alpha1.GitRepo,
	visibility v1alpha1.GitRepositoryVisibility) (v1alpha1.GitRepository, string) {
	repoCreator, ok := ins.(client.GitRepositoryCreator)
	if !ok {
		Skip("plugin does not implement GitRepositoryCreator interface")
	}

	repo, err := repoCreator.CreateGitRepository(testCtx, v1alpha1.CreateGitRepositoryPayload{
		GitRepo:     gitRepo,
		DisplayName: "Test-" + gitRepo.Repository,
		Visibility:  visibility,
		AutoInit:    false,
	})

	Expect(err).Should(Succeed())
	Expect(repo).ShouldNot(BeNil())
	Expect(repo.GetName()).Should(Equal(gitRepo.Repository))

	username, password := GetUsernamePasswordFromCtx(testCtx)
	localRepoPath, err := InitRepo(testCtx, repo.Spec.HttpCloneURL, username, password)
	Expect(err).Should(Succeed())

	return getRepository(testCtx, ins, gitRepo), localRepoPath
}
