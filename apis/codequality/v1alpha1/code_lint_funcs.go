/*
Copyright 2022 The Katanomi Authors.

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

	"k8s.io/apimachinery/pkg/util/validation/field"
)

// IsEmpty returns true if the struct is empty
func (a CodeLintResult) IsEmpty() bool {
	return a.Result == "" && a.Issues == nil
}

// GetObjectWithValues inits an object based on a json.path values map
// returns nil if values is nil
func (CodeLintResult) GetObjectWithValues(ctx context.Context, path *field.Path, values map[string]string) (result *CodeLintResult) {
	if values != nil {
		result = &CodeLintResult{
			Result: values[path.Child("result").String()],
			Issues: CodeLintIssues{}.GetObjectWithValues(ctx, path.Child("issues"), values),
		}
	}
	return
}

// GetObjectWithValues inits an object based on a json.path values map
// returns nil if values is nil
func (CodeLintIssues) GetObjectWithValues(ctx context.Context, path *field.Path, values map[string]string) (result *CodeLintIssues) {
	if values != nil {
		result = &CodeLintIssues{
			Count: strconvAtoi(values[path.Child("count").String()]),
		}
	}
	return
}
