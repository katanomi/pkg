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
	"testing"

	// ktesting "github.com/katanomi/pkg/testing"
	. "github.com/onsi/gomega"
	// corev1 "k8s.io/api/core/v1"

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
