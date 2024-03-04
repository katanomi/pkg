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

package v1alpha1

import (
	"context"

	"github.com/google/go-cmp/cmp"
	. "github.com/katanomi/pkg/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var _ = Describe("GitRepositoryStatus.GetValWithkey", func() {
	var (
		ctx       context.Context
		gitStatus *GitRepositoryStatus
		path      *field.Path
		// results
		stringVars, expectedStringVars map[string]string
		arrayVars, expectedArrayVars   map[string][]string
		objVars, expectedObjVars       map[string]map[string]string
	)

	BeforeEach(func() {
		ctx = context.TODO()
		gitStatus = &GitRepositoryStatus{}
		path = field.NewPath("git")
		expectedStringVars = map[string]string{}
		expectedArrayVars = map[string][]string{}
		expectedObjVars = map[string]map[string]string{}

	})

	JustBeforeEach(func() {
		stringVars, arrayVars, objVars = gitStatus.GetValWithKey(ctx, path)
	})

	When("Is a pull request revision with complete data", func() {
		BeforeEach(func() {
			MustLoadYaml("testdata/gitrepositorystatus/getvalwithkey.pullrequest.data.yaml", &gitStatus)
			MustLoadYaml("testdata/gitrepositorystatus/getvalwithkey.pullrequest.string.yaml", &expectedStringVars)
			MustLoadYaml("testdata/gitrepositorystatus/getvalwithkey.pullrequest.array.yaml", &expectedArrayVars)
			MustLoadYaml("testdata/gitrepositorystatus/getvalwithkey.pullrequest.object.yaml", &expectedObjVars)
		})
		It("should print pullRequest, lastCommit, branch, target and revision variables with name", func() {
			Expect(cmp.Diff(expectedStringVars, stringVars)).To(BeEmpty(), "should have same string variables")
			Expect(cmp.Diff(expectedArrayVars, arrayVars)).To(BeEmpty(), "should have same array variables")
			Expect(cmp.Diff(expectedObjVars, objVars)).To(BeEmpty(), "should have same object variables")
		})
	})
	When("Is a branch revision with related data", func() {
		BeforeEach(func() {
			MustLoadYaml("testdata/gitrepositorystatus/getvalwithkey.branch.data.yaml", &gitStatus)
			MustLoadYaml("testdata/gitrepositorystatus/getvalwithkey.branch.string.yaml", &expectedStringVars)
			MustLoadYaml("testdata/gitrepositorystatus/getvalwithkey.branch.array.yaml", &expectedArrayVars)
			MustLoadYaml("testdata/gitrepositorystatus/getvalwithkey.branch.object.yaml", &expectedObjVars)
		})
		It("should return branch, lastCommit and revision variables with name", func() {
			Expect(cmp.Diff(expectedStringVars, stringVars)).To(BeEmpty(), "should have same string variables")
			Expect(cmp.Diff(expectedArrayVars, arrayVars)).To(BeEmpty(), "should have same array variables")
			Expect(cmp.Diff(expectedObjVars, objVars)).To(BeEmpty(), "should have same object variables")
		})
	})
	When("Is an nil object", func() {
		BeforeEach(func() {
			gitStatus = nil
		})
		It("should return non-nil empty responses", func() {
			Expect(stringVars).To(BeEmpty())
			Expect(stringVars).ToNot(BeNil())
			Expect(arrayVars).To(BeEmpty())
			Expect(arrayVars).ToNot(BeNil())
			Expect(objVars).To(BeEmpty())
			Expect(objVars).ToNot(BeNil())
		})
	})
})
