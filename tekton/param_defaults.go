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

// Package tekton provides utilities for working with Tekton pipeline resources.
// This file implements default value handling for Tekton parameters.
package tekton

import (
	"context"

	pipev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

// ParamSlice wraps a slice of Tekton parameters to provide default value operations.
// This type follows the same pattern as other wrapper types in this package
// like ParamSpec for consistent API design.
type ParamSlice []pipev1beta1.Param

// SetDefaults ensures all parameters in the slice have proper default values.
// This method is primarily used in admission webhooks and template processing
// to normalize parameter specifications.
//
// The method performs the following normalizations:
// - Infers and sets ParamType based on the parameter value content
// - Sets ParamType to String for string values
// - Sets ParamType to Array for array values
// - Sets ParamType to Object for object values
// - Handles edge cases where parameter values might be malformed
//
// Parameters with explicit types are left unchanged to preserve user intent.
func (ps ParamSlice) SetDefaults(_ context.Context) {
	for i := range ps {
		param := &ps[i]

		// Ensure parameter value has a valid type
		// Infer type based on the actual value content if not specified
		if param.Value.Type != "" {
			continue
		}

		// Infer type based on which value field is populated
		if param.Value.ArrayVal != nil {
			param.Value.Type = pipev1beta1.ParamTypeArray
		} else if param.Value.ObjectVal != nil {
			param.Value.Type = pipev1beta1.ParamTypeObject
		} else {
			// Default to string type for simple values or empty values
			param.Value.Type = pipev1beta1.ParamTypeString
		}
	}
}
