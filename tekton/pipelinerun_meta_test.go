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
	"encoding/base64"

	"github.com/google/go-cmp/cmp"
	pkgtesting "github.com/katanomi/pkg/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"sigs.k8s.io/yaml"
)

var _ = Describe("Test.Tekton.Pipeline.Funcs", func() {

	var (
		pipelinerun         v1beta1.PipelineRun
		injectionMeta       PipelineRunMeta
		expectedPipelineRun v1beta1.PipelineRun
		decodedMeta         PipelineRunMeta
	)

	Context("inject meta.json annotation", func() {

		BeforeEach(func() {
			pkgtesting.MustLoadYaml("testdata/PipelineRun.Normal.yaml", &pipelinerun)
		})

		JustAfterEach(func() {
			Expect(InjectPipelineRunMeta(&pipelinerun, injectionMeta)).To(Succeed())
			Expect(cmp.Diff(pipelinerun, expectedPipelineRun)).To(BeEmpty())

			decodeString, _ := base64.StdEncoding.DecodeString(pipelinerun.Annotations[PodAnnotationMetaJson])
			Expect(yaml.Unmarshal(decodeString, &decodedMeta)).To(Succeed())
			Expect(cmp.Diff(decodedMeta, injectionMeta)).To(BeEmpty())
		})

		Context("deliveryRun", func() {
			BeforeEach(func() {
				pkgtesting.MustLoadYaml("testdata/PipelineRunMeta.deliveryRun.yaml", &injectionMeta)
			})

			It("injects pipelinerun annotation", func() {
				pkgtesting.MustLoadYaml("testdata/PipelineRun.ownerDeliveryRun.golden.yaml", &expectedPipelineRun)
			})
		})

		Context("buildRun", func() {
			BeforeEach(func() {
				pkgtesting.MustLoadYaml("testdata/PipelineRunMeta.buildRun.yaml", &injectionMeta)
			})

			It("injects pipelinerun annotation", func() {
				pkgtesting.MustLoadYaml("testdata/PipelineRun.ownerBuildRun.golden.yaml", &expectedPipelineRun)
			})
		})

	})
})
