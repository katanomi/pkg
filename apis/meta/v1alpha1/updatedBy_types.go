/*
Copyright 2021 The Katanomi Authors.

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
	"encoding/json"

	rbacv1 "k8s.io/api/rbac/v1"
)

// updatedBy stores a list of created information.
type UpdatedBy struct {
	// Reference to the user that created the object. Any Kubernetes `Subject` is accepted.
	// +optional
	User *rbacv1.Subject `json:"user,omitempty"`
}

// IsZero basic function returns true when all attributes of the object are empty
func (by *UpdatedBy) IsZero() bool {
	return by == nil || by.User == nil
}

// FromAnnotation will set `by` from annotations
// it will find UpdatedByAnnotationKey and unmarshl content into struct type *UpdatedBy
// if not found UpdatedByAnnotationKey, error would be nil, and *UpdatedBy would be nil also.
// if some errors happened, error will not be nil and *UpdatedBy will be nil
func (by *UpdatedBy) FromAnnotation(annotations map[string]string) (*UpdatedBy, error) {
	jsonStr, ok := annotations[UpdatedByAnnotationKey]
	if !ok {
		return nil, nil
	}

	if by == nil {
		by = &UpdatedBy{}
	}

	err := json.Unmarshal([]byte(jsonStr), by)
	if err != nil {
		return nil, err
	}

	return by, nil
}

// SetIntoAnnotation will set UpdatedBy into annotations
// return annotations that with UpdatedBy.
func (by *UpdatedBy) SetIntoAnnotation(annotations map[string]string) (map[string]string, error) {
	// this error is ignored because it will never happen
	jsonStr, _ := json.Marshal(by)
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[UpdatedByAnnotationKey] = string(jsonStr)
	return annotations, nil
}
