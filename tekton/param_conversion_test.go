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

package tekton

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/katanomi/pkg/testing"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	pipelinev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

var _ = Describe("Test.Params.Conversion", func() {
	var (
		paramsBeta1 = []pipelinev1beta1.Param{}
		paramsV1    = []pipelinev1.Param{}
	)

	BeforeEach(func() {
		Expect(LoadYAML("testdata/params.conversion.v1beta1.yaml", &paramsBeta1)).To(Succeed())
		Expect(LoadYAML("testdata/params.conversion.v1.yaml", &paramsV1)).To(Succeed())
	})

	Describe("ConvertParamsBeta1ToV1", func() {
		It("should convert params from v1beta1 to v1", func() {
			actual := ConvertParamsBeta1ToV1(paramsBeta1)
			Expect(actual).To(Equal(paramsV1))
		})
	})

	Describe("ConvertParamsV1ToBeta1", func() {
		It("should convert params from v1 to v1beta1", func() {
			actual := ConvertParamsV1ToBeta1(paramsV1)
			Expect(actual).To(Equal(paramsBeta1))
		})
	})
})
