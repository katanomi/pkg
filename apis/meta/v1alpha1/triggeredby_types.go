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

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DefinitionTriggeredType string

func (triggeredType DefinitionTriggeredType) String() string {
	return string(triggeredType)
}

type definitionTriggeredTypeValuesType struct {
	Manual    DefinitionTriggeredType
	Automated DefinitionTriggeredType
}

var DefinitionTriggeredTypeValues = definitionTriggeredTypeValuesType{
	// Indicates triggered manually
	Manual: "Manual",
	// Indicates triggered automatically
	Automated: "Automated",
}

// Stores a list of triggered information such as: Entity that triggered,
// reference of an object that could have triggered, and event that triggered.
type TriggeredBy struct {
	// Reference to the user that triggered the object. Any Kubernetes `Subject` is accepted.
	// +optional
	User *rbacv1.Subject `json:"user,omitempty"`

	// Cloud Event data for the event that triggered.
	// +optional
	CloudEvent *CloudEvent `json:"cloudEvent,omitempty"`

	// Reference to another object that might have triggered this object
	// +optional
	Ref *corev1.ObjectReference `json:"ref,omitempty"`

	// Date time of creation of triggered event. Will match a resource's metadata.creationTimestamp
	// it is added here for convinience only
	// +optional
	TriggeredTimestamp *metav1.Time `json:"triggeredTimestamp,omitempty"`

	// Indicates trigger type, such as Manual Automated.
	// +optional
	TriggeredType DefinitionTriggeredType `json:"triggeredType,omitempty"`
}

// IsZero basic function returns true when all attributes of the object are empty
func (by TriggeredBy) IsZero() bool {
	return by.User == nil &&
		by.CloudEvent == nil &&
		by.Ref == nil &&
		by.TriggeredTimestamp == nil &&
		by.TriggeredType.String() == ""
}

// FromAnnotation will set `by` from annotations
// it will find TriggeredByAnnotationKey and unmarshl content into struct type *TriggeredBy
// if not found TriggeredByAnnotationKey, error would be nil, and *TriggeredBy would be nil also.
// if some errors happened, error will not be nil and *TriggeredBy will be nil
func (by *TriggeredBy) FromAnnotation(annotations map[string]string) (*TriggeredBy, error) {
	jsonStr, ok := annotations[TriggeredByAnnotationKey]
	if !ok {
		return nil, nil
	}

	if by == nil {
		by = &TriggeredBy{}
	}

	err := json.Unmarshal([]byte(jsonStr), by)
	if err != nil {
		return nil, err
	}

	return by, nil
}

// SetIntoAnnotation will set TriggeredBy into annotations
// return annotations that with triggeredby.
func (by TriggeredBy) SetIntoAnnotation(annotations map[string]string) (map[string]string, error) {
	if by.CloudEvent != nil {
		// clean cloudevent data, it is so big limitted in annotations
		by.CloudEvent.Data = ""
	}

	jsonStr, err := json.Marshal(by)
	if err != nil {
		return annotations, err
	}
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[TriggeredByAnnotationKey] = string(jsonStr)
	return annotations, nil
}
