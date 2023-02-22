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
	"github.com/google/go-cmp/cmp"
	pkgtesting "github.com/katanomi/pkg/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

var _ = Describe("Test.PipelineRunEnvironments.Injection", func() {

	var (
		pipelineRun         v1beta1.PipelineRun
		variables           = map[string]string{}
		expectedPipelineRun v1beta1.PipelineRun
	)

	BeforeEach(func() {
		pkgtesting.MustLoadYaml("testdata/PipelineRun.Normal.yaml", &pipelineRun)
		pkgtesting.MustLoadYaml("testdata/PipelineRun.environemntVariables.yaml", &variables)
	})

	Context("InjectPipelineRunENVs", func() {
		It("annotates pipeline run with environment variables", func() {
			pkgtesting.MustLoadYaml("testdata/PipelineRun.envInjection.golden.yaml", &expectedPipelineRun)
			InjectPipelineRunENVs(&pipelineRun, variables)
			Expect(cmp.Diff(pipelineRun, expectedPipelineRun)).To(BeEmpty())
		})
	})
})
