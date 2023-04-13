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
	"context"
	"encoding/json"
	"time"

	"github.com/katanomi/pkg/substitution"
	"k8s.io/apimachinery/pkg/util/validation/field"

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

// TriggeredBy stores a list of triggered information such as: Entity that triggered,
// reference of an object that could have triggered, and event that triggered.
type TriggeredBy struct {
	// Reference to the user that triggered the object. Any Kubernetes `Subject` is accepted.
	// +optional
	User *rbacv1.Subject `json:"user,omitempty" variable:"example=admin"`

	// Cloud Event data for the event that triggered.
	// +optional
	CloudEvent *CloudEvent `json:"cloudEvent,omitempty"`

	// Reference to another object that might have triggered this object
	// +optional
	Ref *corev1.ObjectReference `json:"ref,omitempty"`

	// Date time of creation of triggered event. Will match a resource's metadata.creationTimestamp
	// it is added here for convinience only
	// +optional
	TriggeredTimestamp *metav1.Time `json:"triggeredTimestamp,omitempty" variable:"example=2022-08-05T05:34:39Z"`

	// Indicates trigger type, such as Manual Automated.
	// +optional
	TriggeredType DefinitionTriggeredType `json:"triggeredType,omitempty" variable:"example=Automated"`
}

// IsZero basic function returns true when all attributes of the object are empty
func (by *TriggeredBy) IsZero() bool {
	return (by == nil) ||
		(by.User == nil &&
			by.CloudEvent == nil &&
			by.Ref == nil &&
			by.TriggeredTimestamp.IsZero() &&
			by.TriggeredType.String() == "")
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

	// this error is ignored because it will never happen
	jsonStr, _ := json.Marshal(by)
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[TriggeredByAnnotationKey] = string(jsonStr)
	return annotations, nil
}

// GetValWithKey returns the list of keys and values to support variable substitution
func (by *TriggeredBy) GetValWithKey(ctx context.Context, path *field.Path) (values map[string]string) {
	if by == nil {
		by = &TriggeredBy{}
	}

	values = map[string]string{
		path.String(): "",
	}

	// user
	values = substitution.MergeMap(values, RBACSubjectValGetter(by.User)(ctx, path.Child("user")))

	// cloud event
	values = substitution.MergeMap(values, by.CloudEvent.GetValWithKey(ctx, path.Child("cloudEvent")))

	// ref
	values = substitution.MergeMap(values, ObjectReferenceValGetter(by.Ref)(ctx, path.Child("ref")))

	// triggered by
	triggeredTimestamp := ""
	triggeredTimestampYyyyMMddmmss := ""
	if by.TriggeredTimestamp != nil {
		// by.TriggeredTimestamp.UTC().Format(layout string)
		triggeredTimestamp = by.TriggeredTimestamp.UTC().Format(time.RFC3339)
		triggeredTimestampYyyyMMddmmss = by.TriggeredTimestamp.UTC().Format("20060102150405")
	}
	values = substitution.MergeMap(values, map[string]string{
		path.Child("triggeredTimestamp").String():                       triggeredTimestamp,
		path.Child("triggeredType").String():                            string(by.TriggeredType),
		path.Child("triggeredTimestamp").Child("yyyyMMddmmss").String(): triggeredTimestampYyyyMMddmmss,
	})
	return
}
