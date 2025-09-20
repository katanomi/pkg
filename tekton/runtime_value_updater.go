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

import (
	"maps"

	pipelinev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

// MergeDefaultWithRuntimeValue merges the original ParamSpec default value with runtime value.
// For object type parameters, only properties defined in original.Properties are merged from runtimeValue.
// For other types, the default value is completely replaced with a deep copy of runtimeValue.
// The update modifies original directly while keeping runtimeValue unchanged through deep copying.
// If the types don't match, no modification is made.
func MergeDefaultWithRuntimeValue(original *pipelinev1beta1.ParamSpec, runtimeValue *pipelinev1beta1.ParamValue) {
	if original == nil || runtimeValue == nil {
		return
	}

	if original.Type != runtimeValue.Type {
		return
	}

	switch original.Type {
	case pipelinev1beta1.ParamTypeObject:
		original.Default = MergeObjectDefaults(original.Default, runtimeValue)
		// Ensure only properties defined in original.Properties are retained
		CleanupObjectParamByProperties(original)
	default:
		original.Default = runtimeValue.DeepCopy()
	}
}

// MergeObjectDefaults merges multiple ParamValue objects of type 'object' into a single result.
// It combines the ObjectVal maps from all provided ParamValues, with later parameters having higher priority.
// Only ParamValues of type 'object' are processed; others are skipped.
//
// Merge behavior:
// - Properties from later parameters override those from earlier parameters
// - If no valid object parameters are found, returns a deep copy of the first parameter
// - If only one parameter is provided, returns a deep copy of it
// - Returns nil if no parameters are provided
//
// Example:
//
//	param1 := &ParamValue{Type: "object", ObjectVal: {"key1": "value1", "key2": "old"}}
//	param2 := &ParamValue{Type: "object", ObjectVal: {"key2": "new", "key3": "value3"}}
//	result := MergeObjectDefaults(param1, param2)
//	// result.ObjectVal = {"key1": "value1", "key2": "new", "key3": "value3"}
func MergeObjectDefaults(defaults ...*pipelinev1beta1.ParamValue) *pipelinev1beta1.ParamValue {
	if len(defaults) == 0 {
		return nil
	}

	// If only one parameter, return a deep copy
	if len(defaults) == 1 {
		return defaults[0].DeepCopy()
	}

	resultObjectVal := map[string]string{}

	// Merge parameters in order (later parameters have higher priority)
	for i := range len(defaults) {
		param := defaults[i]
		// Skip nil parameters, non-object types, or empty ObjectVal
		if param == nil || param.Type != pipelinev1beta1.ParamTypeObject || len(param.ObjectVal) == 0 {
			continue
		}

		// Merge properties with later parameters overriding earlier ones
		maps.Copy(resultObjectVal, param.ObjectVal)
	}

	// If no valid object properties were found, return a deep copy of the first parameter
	if len(resultObjectVal) == 0 {
		return defaults[0].DeepCopy()
	}

	// Return the merged result
	return &pipelinev1beta1.ParamValue{
		Type:      pipelinev1beta1.ParamTypeObject,
		ObjectVal: resultObjectVal,
	}
}
