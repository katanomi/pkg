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

package tekton

import (
	"context"

	"github.com/google/go-cmp/cmp"
	kclient "github.com/katanomi/pkg/client"
	"github.com/katanomi/pkg/errors"

	pkgtesting "github.com/katanomi/pkg/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	v1 "k8s.io/api/storage/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("Test.GetWorkspaceBindings", func() {
	var (
		ctx    context.Context
		pr     *v1beta1.PipelineRun
		trList *v1beta1.TaskRunList

		wsb []v1beta1.WorkspaceBinding
		err error
	)

	BeforeEach(func() {
		ctx = context.Background()
		pr = &v1beta1.PipelineRun{}
		trList = &v1beta1.TaskRunList{}

		err = nil

		pkgtesting.MustLoadYaml("testdata/PipelineRun.workspaces.yaml", pr)
	})

	JustBeforeEach(func() {
		clt := fake.NewClientBuilder().WithScheme(scheme).WithLists(trList).Build()
		ctx = kclient.WithClient(ctx, clt)
		wsb = make([]v1beta1.WorkspaceBinding, 0)
		wsb, err = GetWorkspaceBindings(ctx, pr)
	})

	When("task runs only contain template workspaces", func() {
		BeforeEach(func() {
			pkgtesting.MustLoadYaml("testdata/taskrunList.pipelineRunLabel.onlyTemplateWorkspaces.yaml", trList)
		})
		It("get all workspaces of pipelinerun", func() {
			Expect(err).NotTo(HaveOccurred())
			var expectedWorkspaces []v1beta1.WorkspaceBinding
			pkgtesting.MustLoadYaml("testdata/workspaces.full.golden.yaml", &expectedWorkspaces)
			Expect(cmp.Diff(expectedWorkspaces, wsb)).To(BeEmpty())
		})
	})

	When("task runs contain parts of template workspaces", func() {
		BeforeEach(func() {
			pkgtesting.MustLoadYaml("testdata/taskrunList.pipelineRunLabel.partialTemplateWorkspaces.yaml", trList)
		})
		It("get missing template workspaces without pvc name", func() {
			Expect(err).NotTo(HaveOccurred())
			var expectedWorkspaces []v1beta1.WorkspaceBinding
			pkgtesting.MustLoadYaml("testdata/workspaces.noCachePVC.golden.yaml", &expectedWorkspaces)
			Expect(cmp.Diff(expectedWorkspaces, wsb)).To(BeEmpty())
		})
	})

	When("task runs with full workspaces", func() {
		BeforeEach(func() {
			pkgtesting.MustLoadYaml("testdata/taskrunList.pipelineRunLabel.all.yaml", trList)
		})
		It("get all workspaces of pipelinerun", func() {
			Expect(err).NotTo(HaveOccurred())
			var expectedWorkspaces []v1beta1.WorkspaceBinding
			pkgtesting.MustLoadYaml("testdata/workspaces.full.golden.yaml", &expectedWorkspaces)
			Expect(cmp.Diff(expectedWorkspaces, wsb)).To(BeEmpty())
		})
	})
})

var _ = Describe("Test.CheckWorkspaceBindings", func() {
	var (
		ctx context.Context

		pr     *v1beta1.PipelineRun
		scList *v1.StorageClassList
		err    error
	)

	BeforeEach(func() {
		ctx = context.Background()
		scList = &v1.StorageClassList{}
		pr = nil
		err = nil

		v1.AddToScheme(scheme)
	})

	JustBeforeEach(func() {
		ctlClient := fake.NewClientBuilder().WithScheme(scheme).WithLists(scList).Build()
		ctx = kclient.WithClient(ctx, ctlClient)
		err = CheckWorkspaceBindings(ctx, pr)
	})

	Context("pipelineRun is nil", func() {
		It("got error", func() {
			Expect(err).Should(BeIdenticalTo(errors.ErrNilPointer))
		})
	})

	Context("pipelineRun has empty worksapces", func() {
		BeforeEach(func() {
			pr = &v1beta1.PipelineRun{}
		})
		It("has no error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("pipelineRun workspaces have VCT without SC name", func() {
		BeforeEach(func() {
			pkgtesting.MustLoadYaml("testdata/PipelineRun.workspaces.yaml", &pr)
		})

		Context("cluster has default SC", func() {
			BeforeEach(func() {
				pkgtesting.MustLoadYaml("testdata/scList.withDefault.yaml", &scList)
			})

			It("has no err", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("cluster has no default SC", func() {
			It("has DefaultStorageClassNotFound err", func() {
				Expect(err).Should(HaveOccurred())
				exists, reason := errors.Reason(err)
				Expect(exists).To(BeTrue())
				Expect(reason).To(Equal(errors.StatusReasonStorageClassNotFound))
			})
		})

	})

	Context("pipelineRun has VCT with SC name", func() {
		BeforeEach(func() {
			pkgtesting.MustLoadYaml("testdata/PipelineRun.workspaces.withPVCName.yaml", &pr)
		})

		It("has no error", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

})
