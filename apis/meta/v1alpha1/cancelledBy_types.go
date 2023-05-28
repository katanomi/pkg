/*
Copyright 2023 The Katanomi Authors.

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

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CancelledBy contains the information for the cancelling
type CancelledBy struct {
	// Reference to the user that triggered the object. Any Kubernetes `Subject` is accepted.
	// +optional
	User *rbacv1.Subject `json:"user,omitempty"`

	// Reference to another object that might have cancelled this object
	// +optional
	Ref *corev1.ObjectReference `json:"ref,omitempty"`

	// CancelledTimestamp is time of cancelling the buildRun.
	// +optional
	CancelledTimestamp *metav1.Time `json:"cancelledTimestamp,omitempty"`
}

// FromAnnotationCancelledBy extract cancelledBy information from annotations
func (by *CancelledBy) FromAnnotationCancelledBy(annotations map[string]string) (*CancelledBy, error) {
	jsonStr, ok := annotations[CancelledByAnnotationKey]
	if !ok {
		return nil, nil
	}

	if by == nil {
		by = &CancelledBy{}
	}

	err := json.Unmarshal([]byte(jsonStr), by)
	if err != nil {
		return nil, err
	}

	return by, nil
}

// SetIntoAnnotationCancelledBy will set CancelledBy information into annotations
// return annotations that with CancelledBy.
func (by *CancelledBy) SetIntoAnnotationCancelledBy(annotations map[string]string) (map[string]string, error) {
	// this error is ignored because it will never happen
	jsonStr, err := json.Marshal(by)
	if err != nil {
		return map[string]string{}, err
	}
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[CancelledByAnnotationKey] = string(jsonStr)
	return annotations, nil
}
