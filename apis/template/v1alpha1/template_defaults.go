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

package v1alpha1

import (
	"context"

	"github.com/katanomi/pkg/tekton"
)

// SetDefaults sets the default values for the Template resource.
// This method is typically called by Kubernetes admission webhooks
// to ensure that Template parameters have proper default values.
//
// The method performs the following operations:
// - Validates the Template is not nil
// - Sets default values for all parameters in the Template.Params slice
// - Sets default values for all parameters in the TemplateRef.Params slice
// - Ensures parameter types are properly initialized
func (t *Template) SetDefaults(ctx context.Context) {
	if t == nil {
		return
	}

	// Set defaults for template parameters if any exist
	if len(t.Params) > 0 {
		tekton.ParamSlice(t.Params).SetDefaults(ctx)
	}

	// Set defaults for template reference parameters if any exist
	if len(t.TemplateRef.Params) > 0 {
		tekton.ParamSlice(t.TemplateRef.Params).SetDefaults(ctx)
	}
}
