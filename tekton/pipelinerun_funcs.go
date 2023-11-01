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

import "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"

// GetPipelineSpec return pipelineSpec from pipelinerun
func GetPipelineSpec(pr *v1beta1.PipelineRun) *v1beta1.PipelineSpec {
	if pr == nil {
		return nil
	}

	if pr.Spec.PipelineSpec != nil {
		return pr.Spec.PipelineSpec
	}
	return pr.Status.PipelineSpec
}
