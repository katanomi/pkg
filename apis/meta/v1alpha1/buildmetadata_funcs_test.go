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
		//
		// log level. It can be debug, info, warn, error, dpanic, panic, fatal.
		log, _ = logging.NewLogger("", "debug")
	)

	BeforeEach(func() {
		ctx = context.TODO()
		gitStatus = &BuildRunGitStatus{}
		Expect(ktesting.LoadYAML("testdata/gitstatus.golden.yaml", &gitStatus)).To(Succeed())
	})

	JustBeforeEach(func() {
		actual = gitStatus.GetValWithKey(ctx, field.NewPath("git"))
		log.Infow("BuildRunGitStatus.GetValWithKey", "gitStatus", gitStatus, "result", actual)
	})

	When("struct is empty", func() {
		BeforeEach(func() {
			gitStatus = &BuildRunGitStatus{}
		})

		It("should have values", func() {
			Expect(actual).To(Equal(map[string]string{
				"git.url": "",
				//
				"git.lastCommit.id":          "",
				"git.lastCommit.shortID":     "",
				"git.lastCommit.title":       "",
				"git.lastCommit.message":     "",
				"git.lastCommit.authorEmail": "",
				"git.lastCommit.pushedAt":    "",
				"git.lastCommit.webURL":      "",
				//
				"git.pullRequest.id":           "",
				"git.pullRequest.title":        "",
				"git.pullRequest.source":       "",
				"git.pullRequest.target":       "",
				"git.pullRequest.webURL":       "",
				"git.pullRequest.hasConflicts": "false",
				//
				"git.branch.name":      "",
				"git.branch.protected": "false",
				"git.branch.default":   "false",
				"git.branch.webURL":    "",
			}))
		})
	})

	When("struct is not empty", func() {
		It("should have values", func() {
			Expect(actual).To(Equal(map[string]string{
				"git.url": "https://github.com/katanomi/pkg",
				//
				"git.lastCommit.id":          "abe83942450308432a12e9679519795f938b2bed",
				"git.lastCommit.shortID":     "abe83942",
				"git.lastCommit.title":       "Initial commit 406",
				"git.lastCommit.message":     "Initial commit 406\n",
				"git.lastCommit.authorEmail": "alauda@github.com",
				"git.lastCommit.pushedAt":    "2020-01-01 01:02:03 +0000 UTC",
				"git.lastCommit.webURL":      "https://github.com",
				//
				"git.pullRequest.id":           "1",
				"git.pullRequest.title":        "test-build ==> master",
				"git.pullRequest.source":       "test-build",
				"git.pullRequest.target":       "master",
				"git.pullRequest.webURL":       "https://github.com/katanomi/pkg/merge_requests/1",
				"git.pullRequest.hasConflicts": "true",
				//
				"git.branch.name":      "test-build",
				"git.branch.protected": "true",
				"git.branch.default":   "true",
				"git.branch.webURL":    "https://github.com/katanomi/pkg/tree/test",
			}))
		})
	})

})
