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
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

// CleanupObjectParamByProperties filters object parameter default values based on
// the parameter's Properties schema definition.
//
// This function handles two main scenarios:
//  1. When Properties is empty: Clears the Default.ObjectVal if it's non-empty,
//     preserving nil/empty values as-is
//  2. When Properties is defined: Only retains properties that exist in the Properties
//     schema, filtering out any undefined properties
//
// This filtering ensures that object parameters maintain schema consistency during
// template rendering and prevents undefined properties from being propagated.
//
// Parameters:
//   - param: Pointer to ParamSpec that should be of type Object with Default.ObjectVal
//
// The function returns early if param is nil, not an object type, or has no default value.
func CleanupObjectParamByProperties(param *v1beta1.ParamSpec) {
	if param == nil || param.Type != v1beta1.ParamTypeObject || param.Default == nil {
		return
	}

	defaultMap := param.Default.ObjectVal

	// Case 1: No properties schema defined - clear non-empty ObjectVal
	if len(param.Properties) == 0 {
		if len(defaultMap) > 0 {
			param.Default.ObjectVal = map[string]string{}
		}
		return
	}

	// Case 2: Properties schema exists - filter to only include defined properties
	// Pre-allocate map with exact capacity needed for better performance
	filteredDefaults := make(map[string]string, len(param.Properties))
	for propertyName := range param.Properties {
		if value, exists := defaultMap[propertyName]; exists {
			filteredDefaults[propertyName] = value
		}
	}
	param.Default.ObjectVal = filteredDefaults
}
