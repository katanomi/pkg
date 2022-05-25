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

// CreatedBy stores a list of created information.
type CreatedBy struct {
	// Reference to the user that created the object. Any Kubernetes `Subject` is accepted.
	// +optional
	User *rbacv1.Subject `json:"user,omitempty"`
}

// IsZero basic function returns true when all attributes of the object are empty
func (by *CreatedBy) IsZero() bool {
	return by == nil || by.User == nil
}

// FromAnnotation will set `by` from annotations
// it will find CreatedByAnnotationKey and unmarshl content into struct type *CreatedBy
// if not found CreatedByAnnotationKey, error would be nil, and *CreatedBy would be nil also.
// if some errors happened, error will not be nil and *CreatedBy will be nil
func (by *CreatedBy) FromAnnotation(annotations map[string]string) (*CreatedBy, error) {
	jsonStr, ok := annotations[CreatedByAnnotationKey]
	if !ok {
		return nil, nil
	}

	if by == nil {
		by = &CreatedBy{}
	}

	err := json.Unmarshal([]byte(jsonStr), by)
	if err != nil {
		return nil, err
	}

	return by, nil
}

// SetIntoAnnotation will set CreatedBy into annotations
// return annotations that with CreatedBy.
func (by *CreatedBy) SetIntoAnnotation(annotations map[string]string) (map[string]string, error) {
	// this error is ignored because it will never happen
	jsonStr, _ := json.Marshal(by)
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[CreatedByAnnotationKey] = string(jsonStr)
	return annotations, nil
}
