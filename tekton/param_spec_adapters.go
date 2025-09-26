/*
Copyright 2025 The Katanomi Authors.

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

import pipev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"

// Ensure pipev1beta1.ParamSpec implements NamedItem interface
var _ NamedItem = ParamSpec{}

// ParamSpec wraps pipev1beta1.ParamSpec to implement NamedItem
type ParamSpec pipev1beta1.ParamSpec

// GetName returns the parameter name
func (p ParamSpec) GetName() string {
	return p.Name
}

// TektonParamSpecs converts slice of pipev1beta1.ParamSpec to slice of NamedItem
func TektonParamSpecs(params []pipev1beta1.ParamSpec) []NamedItem {
	result := make([]NamedItem, len(params))
	for i, param := range params {
		result[i] = ParamSpec(param)
	}
	return result
}
