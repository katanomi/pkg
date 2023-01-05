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
	"testing"

	. "github.com/onsi/gomega"
	knativeApis "knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

func TestBranchSpec_Equal(t *testing.T) {
	t.Run("nil object call equal", func(t *testing.T) {
		g := NewGomegaWithT(t)
		var branchSpec *BranchSpec
		item := BranchSpec{
			CodeInfo: CodeInfo{
				Project: "project1",
			},
		}
		g.Expect(branchSpec.Equal(item)).To(BeFalse())
	})

	t.Run("code info not equal", func(t *testing.T) {
		g := NewGomegaWithT(t)
		branchSpec := BranchSpec{
			CodeInfo: CodeInfo{
				Project: "project",
			},
		}
		item := BranchSpec{
			CodeInfo: CodeInfo{
				Project: "project1",
			},
		}
		g.Expect(branchSpec.Equal(item)).To(BeFalse())
	})

	t.Run("author not equal", func(t *testing.T) {
		g := NewGomegaWithT(t)
		branchSpec := BranchSpec{
			Author: UserSpec{Name: "name"},
		}
		item := BranchSpec{
			Author: UserSpec{Email: "email"},
		}
		g.Expect(branchSpec.Equal(item)).To(BeFalse())
	})

	t.Run("issue not equal", func(t *testing.T) {
		g := NewGomegaWithT(t)
		branchSpec := BranchSpec{Issue: IssueInfo{Type: "task"}}
		item := BranchSpec{Issue: IssueInfo{Type: "bug"}}
		g.Expect(branchSpec.Equal(item)).To(BeFalse())
	})

	t.Run("compare equal", func(t *testing.T) {
		g := NewGomegaWithT(t)
		branchSpec := BranchSpec{
			CodeInfo: CodeInfo{
				Project: "project",
			},
			Author: UserSpec{Name: "name"},
			Issue:  IssueInfo{Type: "task"},
		}
		item := BranchSpec{
			CodeInfo: CodeInfo{
				Project: "project",
			},
			Author: UserSpec{Name: "name"},
			Issue:  IssueInfo{Type: "task"},
		}
		g.Expect(branchSpec.Equal(item)).To(BeTrue())
	})
}

func TestCodeInfo_Equal(t *testing.T) {
	t.Run("when integration not equal, compare return true", func(t *testing.T) {
		g := NewGomegaWithT(t)
		codeAddress, _ := knativeApis.ParseURL("https://code.address.com/project/code")
		codeInfo := CodeInfo{
			Project:         "project",
			IntegrationName: "GitLab",
			Repository:      "code",
			Branch:          "test",
			BaseBranch:      "master",
			Address:         &duckv1.Addressable{URL: codeAddress},
		}

		item := CodeInfo{
			Project:         "project",
			IntegrationName: "Gitlab",
			Repository:      "code",
			Branch:          "test",
			BaseBranch:      "master",
			Address:         &duckv1.Addressable{URL: codeAddress},
		}
		g.Expect(codeInfo.Equal(item)).To(BeTrue())
	})

	t.Run("nil object call equal", func(t *testing.T) {
		g := NewGomegaWithT(t)
		var codeInfo *CodeInfo
		item := CodeInfo{}
		g.Expect(codeInfo.Equal(item)).To(BeFalse())
	})

	t.Run("address is nil equal", func(t *testing.T) {
		g := NewGomegaWithT(t)
		codeAddress, _ := knativeApis.ParseURL("https://code.address.com/project/code")
		codeInfo := CodeInfo{
			Project:         "project",
			IntegrationName: "Gitlab",
			Repository:      "code",
			Branch:          "test",
			BaseBranch:      "master",
			Address:         &duckv1.Addressable{URL: codeAddress},
		}

		item := CodeInfo{
			Project:         "project",
			IntegrationName: "Gitlab",
			Repository:      "code",
			Branch:          "test",
			BaseBranch:      "master",
			Address:         nil,
		}
		g.Expect(codeInfo.Equal(item)).To(BeFalse())
	})

	t.Run("compare host not equal", func(t *testing.T) {
		g := NewGomegaWithT(t)
		codeAddress, _ := knativeApis.ParseURL("https://code.address.com/project/code")
		codeAddress1, _ := knativeApis.ParseURL("https://code1.address.com/project/code")
		codeInfo := CodeInfo{
			Project:         "project",
			IntegrationName: "Gitlab",
			Repository:      "code",
			Branch:          "test",
			BaseBranch:      "master",
			Address:         &duckv1.Addressable{URL: codeAddress},
		}

		item := CodeInfo{
			Project:         "project",
			IntegrationName: "Gitlab",
			Repository:      "code",
			Branch:          "test",
			BaseBranch:      "master",
			Address:         &duckv1.Addressable{URL: codeAddress1},
		}
		g.Expect(codeInfo.Equal(item)).To(BeFalse())
	})
}
