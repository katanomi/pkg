/*
Copyright 2022 The Katanomi Authors.

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
	"testing"

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"knative.dev/pkg/logging"

	ktesting "github.com/katanomi/pkg/testing"
)

func TestAssignByGitBranch(t *testing.T) {
	g := NewGomegaWithT(t)

	gitBranch := GitBranch{}
	gitBranchStatus := BuildGitBranchStatus{}
	g.Expect(ktesting.LoadYAML("testdata/gitbranch.golden.yaml", &gitBranch)).To(Succeed())
	g.Expect(ktesting.LoadYAML("testdata/gitbranch.status.yaml", &gitBranchStatus)).To(Succeed())

	expectStatus := gitBranchStatus
	actualGitBranchStatus := gitBranchStatus.AssignByGitBranch(&gitBranch)
	g.Expect(*actualGitBranchStatus).To(Equal(expectStatus))

	// empty struct should not panic
	gitBranch = GitBranch{}
	gitBranchStatus = BuildGitBranchStatus{}
	actualGitBranchStatus = gitBranchStatus.AssignByGitBranch(&gitBranch)
	g.Expect(actualGitBranchStatus.Default).To(BeFalse())
}

func TestAssignByGitCommit(t *testing.T) {
	g := NewGomegaWithT(t)

	gitCommit := GitCommit{}
	gitCommitStatus := BuildGitCommitStatus{}
	g.Expect(ktesting.LoadYAML("testdata/gitcommit.golden.yaml", &gitCommit)).To(Succeed())
	g.Expect(ktesting.LoadYAML("testdata/gitcommit.status.yaml", &gitCommitStatus)).To(Succeed())

	expectStatus := gitCommitStatus
	actualGitCommitStatus := gitCommitStatus.AssignByGitCommit(&gitCommit)
	g.Expect(*actualGitCommitStatus).To(Equal(expectStatus))
}

func TestAssignByGitPullRequest(t *testing.T) {
	g := NewGomegaWithT(t)

	gitPullRequest := GitPullRequest{}
	gitPullRequestStatus := BuildGitPullRequestStatus{}
	g.Expect(ktesting.LoadYAML("testdata/gitpullrequest.golden.yaml", &gitPullRequest)).To(Succeed())
	g.Expect(ktesting.LoadYAML("testdata/gitpullrequest.status.yaml", &gitPullRequestStatus)).To(Succeed())

	expectStatus := gitPullRequestStatus
	actualGitPullRequestStatus := gitPullRequestStatus.AssignByGitPullRequest(&gitPullRequest)
	g.Expect(*actualGitPullRequestStatus).To(Equal(expectStatus))
}

var _ = Describe("Test.BuildRunGitStatus.GetValWithKey", func() {
	var (
		ctx       context.Context
		gitStatus *BuildRunGitStatus
		actual    map[string]string
		expected  map[string]string
		//
		// log level. It can be debug, info, warn, error, dpanic, panic, fatal.
		log, _ = logging.NewLogger("", "debug")
	)

	BeforeEach(func() {
		ctx = context.TODO()
		gitStatus = &BuildRunGitStatus{}
		expected = map[string]string{}
		Expect(ktesting.LoadYAML("testdata/gitstatus.golden.yaml", &gitStatus)).To(Succeed())
	})

	JustBeforeEach(func() {
		actual = gitStatus.GetValWithKey(ctx, field.NewPath("git"))
		log.Infow("BuildRunGitStatus.GetValWithKey", "gitStatus", gitStatus, "result", actual)
	})

	When("struct is empty", func() {
		BeforeEach(func() {
			gitStatus = &BuildRunGitStatus{}
			ktesting.MustLoadYaml("testdata/gitstatus.emptymap.golden.yaml", &expected)
		})

		It("should have values", func() {
			diff := cmp.Diff(actual, expected)
			Expect(diff).To(BeEmpty())
		})
	})

	When("struct is not empty", func() {
		BeforeEach(func() {
			ktesting.MustLoadYaml("testdata/gitstatus.map.golden.yaml", &expected)
		})
		It("should have values", func() {
			diff := cmp.Diff(actual, expected)
			Expect(diff).To(BeEmpty())
		})
	})

})

var _ = Describe("Test.BaseGitStatus", func() {
	Context("IsPR", func() {
		It("should return true if GitStatus is a pull request", func() {
			baseGitStatus := &BaseGitStatus{
				PullRequest: &BuildGitPullRequestStatus{
					ID: "123",
				},
			}
			Expect(baseGitStatus.IsPR()).To(BeTrue())
		})

		It("should return false if GitStatus is not a pull request", func() {
			baseGitStatus := &BaseGitStatus{
				PullRequest: nil,
			}
			Expect(baseGitStatus.IsPR()).To(BeFalse())
		})

		It("should return false if GitStatus is a pull request but ID is empty", func() {
			baseGitStatus := &BaseGitStatus{
				PullRequest: &BuildGitPullRequestStatus{
					ID: "",
				},
			}
			Expect(baseGitStatus.IsPR()).To(BeFalse())
		})
	})
})
