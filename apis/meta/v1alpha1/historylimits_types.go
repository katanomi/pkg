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

import "k8s.io/apimachinery/pkg/util/validation/field"

// HistoryLimits limits the number of executed items are preserved
// It only calculates already completed items
type HistoryLimits struct {
	// Sets a hard count for all finished items
	// to be cleared from storage
	Count *int `json:"count,omitempty"`
}

// Validate make sure the data is legitimate..
func (h *HistoryLimits) Validate(path *field.Path) (errs field.ErrorList) {
	errs = field.ErrorList{}
	if h == nil {
		return
	}
	if h.Count != nil && *h.Count < 0 {
		errs = append(errs, field.Invalid(path.Child("count"), *h.Count, "should be greater than zero"))
	}
	return
}
