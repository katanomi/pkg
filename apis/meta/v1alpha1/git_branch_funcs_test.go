/*
Copyright 2021 The Katanomi Authors.

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
	. "github.com/katanomi/pkg/testing"

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
)

var _ = Describe("GitBranch.GetBranchStatus", func() {

	var (
		branch   *GitBranch
		result   *BuildGitBranchStatus
		expected *BuildGitBranchStatus
	)

	BeforeEach(func() {
		branch = &GitBranch{}
		expected = &BuildGitBranchStatus{}
	})
	JustBeforeEach(func() {
		result = branch.GetBranchStatus()
	})
	When("struct has all data", func() {
		BeforeEach(func() {
			LoadYAML("testdata/git_branch_funcs_getbranchstatus.full.yaml", branch)
			LoadYAML("testdata/git_branch_funcs_getbranchstatus.full.golden.yaml", expected)
		})
		It("should return full BuildRunGitStatus", func() {
			Expect(result).ToNot(BeNil())

			diff := cmp.Diff(result, expected)
			Expect(diff).To(BeEmpty())
		})
	})
	When("struct weburl is in properties", func() {
		BeforeEach(func() {
			LoadYAML("testdata/git_branch_funcs_getbranchstatus.weburl.yaml", branch)
			LoadYAML("testdata/git_branch_funcs_getbranchstatus.full.golden.yaml", expected)
		})
		It("should return full BuildRunGitStatus", func() {
			Expect(result).ToNot(BeNil())

			diff := cmp.Diff(result, expected)
			Expect(diff).To(BeEmpty())
		})
	})
})
