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

package tekton

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

var _ = Describe("GetPipelineSpec", func() {
	var (
		pr           *v1beta1.PipelineRun
		pipelineSpec *v1beta1.PipelineSpec
	)

	BeforeEach(func() {
		pr = &v1beta1.PipelineRun{}
	})

	JustBeforeEach(func() {
		pipelineSpec = GetPipelineSpec(pr)
	})

	Context("nil PipelineRun", func() {
		BeforeEach(func() {
			pr = nil
		})

		It("should return nil", func() {
			Expect(pipelineSpec).To(BeNil())
		})
	})

	Context("PipelineSpec in Spec", func() {
		BeforeEach(func() {
			pr = &v1beta1.PipelineRun{
				Spec: v1beta1.PipelineRunSpec{
					PipelineSpec: &v1beta1.PipelineSpec{},
				},
			}
		})

		It("should return PipelineSpec in Spec", func() {
			Expect(pipelineSpec).To(Equal(pr.Spec.PipelineSpec))
		})
	})

	Context("PipelineSpec in Status", func() {
		BeforeEach(func() {
			pr = &v1beta1.PipelineRun{
				Status: v1beta1.PipelineRunStatus{
					PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{
						PipelineSpec: &v1beta1.PipelineSpec{},
					},
				},
			}
		})

		It("should return PipelineSpec in Status", func() {
			Expect(pipelineSpec).To(Equal(pr.Status.PipelineSpec))
		})
	})
})
