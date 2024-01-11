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
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	pipelinev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

// ConvertParamBeta1ToV1 convert param from v1beta1 to v1
func ConvertParamBeta1ToV1(param pipelinev1beta1.Param) pipelinev1.Param {
	return pipelinev1.Param{
		Name: param.Name,
		Value: pipelinev1.ParamValue{
			Type:      pipelinev1.ParamType(param.Value.Type),
			StringVal: param.Value.StringVal,
			ArrayVal:  param.Value.ArrayVal,
			ObjectVal: param.Value.ObjectVal,
		},
	}
}

// ConvertParamV1ToBeta1 convert param from v1 to v1beta1
func ConvertParamV1ToBeta1(param pipelinev1.Param) pipelinev1beta1.Param {
	return pipelinev1beta1.Param{
		Name: param.Name,
		Value: pipelinev1beta1.ParamValue{
			Type:      pipelinev1beta1.ParamType(param.Value.Type),
			StringVal: param.Value.StringVal,
			ArrayVal:  param.Value.ArrayVal,
			ObjectVal: param.Value.ObjectVal,
		},
	}
}

// ConvertParamsBeta1ToV1 convert params from v1beta1 to v1
func ConvertParamsBeta1ToV1(params []pipelinev1beta1.Param) []pipelinev1.Param {
	result := make([]pipelinev1.Param, 0, len(params))
	for _, p := range params {
		result = append(result, ConvertParamBeta1ToV1(p))
	}
	return result
}

// ConvertParamsV1ToBeta1 convert params from v1 to v1beta1
func ConvertParamsV1ToBeta1(params []pipelinev1.Param) []pipelinev1beta1.Param {
	result := make([]pipelinev1beta1.Param, 0, len(params))
	for _, p := range params {
		result = append(result, ConvertParamV1ToBeta1(p))
	}
	return result
}
