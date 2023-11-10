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
	"sort"

	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
	"github.com/katanomi/pkg/plugin/route"
	"github.com/katanomi/pkg/pointer"
	. "github.com/katanomi/pkg/testing/framework"
	. "github.com/katanomi/pkg/testing/framework/base"
	. "github.com/katanomi/pkg/testing/testcases"
	. "github.com/katanomi/pkg/testing/testcases/gitplugin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
)

var ProjectCaseSet = func() {
	caseProjectList.Do()
	caseGetProject.Do()
	caseCreateProject.Do()
}

var caseProjectList = P0Case("test for getting project list").
	WithCondition(
		PluginImplementCondition{Interface: (*client.ProjectLister)(nil)},
	).WithFunc(func(testContext *TestContext) {
	var (
		ctx         context.Context
		instance    TestablePlugin
		projectList *v1alpha1.ProjectList
		queryOption *v1alpha1.ListOptions
		err         error
	)

	BeforeEach(func() {
		ctx = testContext.Context
		instance = GitPluginFromCtx(ctx)
		projectList = &v1alpha1.ProjectList{}
		queryOption = &v1alpha1.ListOptions{}
		err = nil
	})

	JustBeforeEach(func() {
		getter := instance.(client.ProjectLister)
		projectList, err = getter.ListProjects(ctx, *queryOption)
	})

	Context("list project successfully", func() {
		It("get return the specified projects", func() {
			Expect(err).Should(BeNil())
			Expect(projectList).ShouldNot(BeNil())
			Expect(projectList.Items).ShouldNot(BeEmpty())
		})
	})

	Context("page turning of project list", func() {
		var projectTwoPageList *v1alpha1.ProjectList
		JustBeforeEach(func() {
			getter := instance.(client.ProjectLister)
			projectTwoPageList, err = getter.ListProjects(ctx, v1alpha1.ListOptions{ItemsPerPage: 10, Page: 1})
			if projectTwoPageList.TotalItems < 10 {
				Skip("At least 10 projects must be included in the tool to run the current use case.")
			}
			projectSortRule(projectTwoPageList.Items, instance.GetTestUserProject())
		})

		When("page number 1", func() {
			BeforeEach(func() {
				projectList = &v1alpha1.ProjectList{}
				queryOption = &v1alpha1.ListOptions{ItemsPerPage: 5, Page: 1}
				err = nil
			})

			It("return to first page, and check order", func() {
				Expect(err).Should(BeNil())
				Expect(projectList).ShouldNot(BeNil())
				Expect(projectList.Items).ShouldNot(BeEmpty())
				Expect(projectList.Items).To(Equal(projectTwoPageList.Items[0:5]))
				index := rand.Intn(4)
				checkProjectRequiredFields(&projectList.Items[index], &projectTwoPageList.Items[index])
			})
		})

		When("page number 2", func() {
			BeforeEach(func() {
				projectList = &v1alpha1.ProjectList{}
				queryOption = &v1alpha1.ListOptions{ItemsPerPage: 5, Page: 2}
				err = nil
			})
			It("return to first page, and check order", func() {
				Expect(err).Should(BeNil())
				Expect(projectList).ShouldNot(BeNil())
				Expect(projectList.Items).ShouldNot(BeEmpty())
				Expect(projectList.Items).To(Equal(projectTwoPageList.Items[5:]))
				index := rand.Intn(4)
				checkProjectRequiredFields(&projectList.Items[index], &projectTwoPageList.Items[index+5])
			})
		})
	})

	Context("search project", func() {
		When("specified project not found", func() {
			BeforeEach(func() {
				queryOption.Search = map[string][]string{
					route.SearchQueryKey: {"not-exist"},
				}
			})

			It("should return empty list", func() {
				Expect(err).Should(BeNil())
				Expect(projectList).ShouldNot(BeNil())
				Expect(projectList.Items).Should(BeEmpty())
			})
		})

		When("specified project found", func() {
			var wantProject v1alpha1.Project
			BeforeEach(func() {
				projectName := "e2e-create-project-" + rand.String(4)
				wantProject = *createProject(ctx, instance, projectName)
				queryOption.Search = map[string][]string{
					route.SearchQueryKey: {wantProject.GetName()},
				}
				DeferCleanup(func() {
					cleanupProject(ctx, instance, projectName)
				})
			})
			It("should return the created project", func() {
				Expect(err).Should(BeNil())
				Expect(projectList).ShouldNot(BeNil())
				Expect(projectList.Items).ShouldNot(BeEmpty())
				gotProject := FindByName(ToPtrList(projectList.Items), wantProject.GetName())
				checkProjectRequiredFields(gotProject, &wantProject)
			})
		})
	})
})

var caseGetProject = P0Case("test for getting project info").
	WithCondition(
		PluginImplementCondition{Interface: (*client.ProjectGetter)(nil)},
	).WithFunc(func(testContext *TestContext) {
	var (
		ctx         context.Context
		instance    TestablePlugin
		projectName *string
		projectInfo *v1alpha1.Project
		err         error
	)

	BeforeEach(func() {
		ctx = testContext.Context
		instance = GitPluginFromCtx(ctx)
		projectName = pointer.String("")
		ctx = context.WithValue(ctx, v1alpha1.KeyForSubType, string(v1alpha1.GitUserProjectSubType))
		projectInfo = &v1alpha1.Project{}
		err = nil
	})

	JustBeforeEach(func() {
		getter := instance.(client.ProjectGetter)
		projectInfo, err = getter.GetProject(ctx, *projectName)
	})

	Context("project not exist", func() {
		BeforeEach(func() {
			*projectName = "not-exist-not-exist-not-exist"
		})
		It("get return 404", func() {
			Expect(errors.IsNotFound(err)).Should(BeTrue())
		})
	})

	Context("project exist", func() {
		BeforeEach(func() {
			*projectName = instance.GetTestOrgProject()
		})

		It("should return the project", func() {
			Expect(err).Should(BeNil())
			Expect(projectInfo).ShouldNot(BeNil())
			Expect(projectInfo.GetName()).Should(Equal(instance.GetTestOrgProject())) // TODO: more asserts
		})
	})
})

var caseCreateProject = P0Case("test for creating new project").
	WithCondition(
		PluginImplementCondition{Interface: (*client.ProjectCreator)(nil)},
	).WithFunc(func(testContext *TestContext) {
	var (
		ctx         context.Context
		instance    TestablePlugin
		project     *v1alpha1.Project
		err         error
		projectName = "create-project-" + rand.String(4)
		creator     client.ProjectCreator
	)

	BeforeEach(func() {
		ctx = testContext.Context
		instance = GitPluginFromCtx(ctx)
		ctx = context.WithValue(ctx, v1alpha1.KeyForSubType, string(v1alpha1.GitUserProjectSubType))
		project = &v1alpha1.Project{}
		err = nil
		creator = instance.(client.ProjectCreator)
	})

	JustBeforeEach(func() {
		project.Name = projectName
		project, err = creator.CreateProject(ctx, project)
	})

	Context("first create project", func() {
		It("should create successfully", func() {
			Expect(err).Should(BeNil())
			Expect(project).ShouldNot(BeNil())
			Expect(project.GetName()).Should(Equal(projectName))
			checkProjectRequired(project)
		})
	})

	Context("create an existed project", func() {
		It("should create failed", func() {
			Expect(err).ShouldNot(BeNil())
			Expect(project).Should(BeNil())
		})
	})

	AfterAll(func() {
		cleanupProject(ctx, instance, projectName)
	})
})

func checkProjectRequired(gotProject *v1alpha1.Project) {
	Expect(gotProject.GetName()).ShouldNot(BeEmpty())
	Expect(gotProject.Spec.Address).ShouldNot(BeNil())
	Expect(gotProject.Spec.Address.URL.String()).ShouldNot(BeEmpty())
}

func checkProjectRequiredFields(gotProject, wantProject *v1alpha1.Project) {
	checkProjectRequired(gotProject)
	Expect(gotProject.GetName()).To(Equal(wantProject.GetName()))
	Expect(gotProject.Spec.Public).To(Equal(wantProject.Spec.Public))
	Expect(gotProject.Spec.Address.URL.String()).To(Equal(gotProject.Spec.Address.URL.String()))

}

func checkProjectOptionalFields(gotProject, wantProject *v1alpha1.Project) {
	OptionalExpect(gotProject.CreationTimestamp.Time).To(Equal(wantProject.CreationTimestamp.Time))
	OptionalExpect(gotProject.Spec.SubType.String()).To(Equal(gotProject.Spec.SubType.String()))
	OptionalExpect(gotProject.Spec.Access).NotTo(BeEmpty())
	OptionalExpect(gotProject.Spec.Access.URL.String()).To(Equal(gotProject.Spec.Access.URL.String()))
}

func createProject(ctx context.Context, instance TestablePlugin, name string) *v1alpha1.Project {
	projectCreator, ok := instance.(client.ProjectCreator)
	if !ok {
		Skip("plugin not implement ProjectCreator interface")
	}
	wantProject, err := projectCreator.CreateProject(ctx, &v1alpha1.Project{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	})
	Expect(err).Should(Succeed())
	Expect(wantProject).ShouldNot(BeNil())
	return wantProject
}

func cleanupProject(testCtx context.Context, ins TestablePlugin, projectName string) {
	projectDeleter, ok := ins.(client.ProjectDeleter)
	if !ok {
		Skip("plugin does not implement ProjectDeleter interface")
		return
	}

	err := projectDeleter.DeleteProject(testCtx, &v1alpha1.Project{ObjectMeta: metav1.ObjectMeta{Name: projectName}})
	if err != nil {
		By(fmt.Sprintf("clean project failed: %s", err.Error()))
	}
}

// projectSortRule make the project list satisfy the order rules. The order rules of the project list are as follows:
// 1. GitGroup types are displayed first, and similar types are arranged in alphabetical order.
// 2. Then display the GitUser type, with the current user at the end, and the remaining types in alphabetical order.
func projectSortRule(projects []v1alpha1.Project, currentUser string) {
	sort.SliceStable(projects, func(i, j int) bool {
		if projects[i].Spec.SubType == projects[j].Spec.SubType {
			if projects[i].Spec.SubType == v1alpha1.GitUserProjectSubType {
				if currentUser == projects[i].GetName() {
					return false
				}

				if currentUser == projects[j].GetName() {
					return false
				}
			}
			return projects[i].GetName() < projects[j].GetName()
		}
		return projects[i].Spec.SubType == v1alpha1.GitGroupProjectSubType
	})
}
