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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"knative.dev/pkg/apis"
)

const (
	ConditionReasonNotSet = "NotSet"
)

// ReasonForError returns a string for Reason from an error. If the error is a apimachinery.StatusError
// will return its reason, otherwise will return unknown
func ReasonForError(err error) string {
	return string(errors.ReasonForError(err))
}

// SetConditionByError sets a condition using an error to determine if it is True or False.
// If the given error is not nil will examine its reason and mark the condition as False
// otherwise will set condition to True
func SetConditionByError(conditionManager apis.ConditionManager, condition apis.ConditionType, err error) {
	if err == nil {
		conditionManager.MarkTrue(condition)
	} else {
		conditionManager.MarkFalse(condition, ReasonForError(err), err.Error())
	}
}

// PropagateCondition propagates a condition based on an external conditions
// used mainly when a resource depends on another and its conditions needs to be synced
func PropagateCondition(conditionManager apis.ConditionManager, conditionType apis.ConditionType, condition *apis.Condition) {
	if condition == nil {
		conditionManager.MarkUnknown(conditionType, ConditionReasonNotSet, "condition is empty")
		return
	}
	switch condition.Status {
	case corev1.ConditionTrue:
		conditionManager.MarkTrueWithReason(conditionType, condition.Reason, condition.Message)
	case corev1.ConditionFalse:
		conditionManager.MarkFalse(conditionType, condition.Reason, condition.Message)
	case corev1.ConditionUnknown:
		fallthrough
	default:
		conditionManager.MarkUnknown(conditionType, condition.Reason, condition.Message)
	}
}

// IsConditionChanged given two condition accessors and a condition type will check if conditions changed
func IsConditionChanged(current, old apis.ConditionAccessor, conditionType apis.ConditionType) bool {
	currentCondition := current.GetCondition(conditionType)
	oldCondition := old.GetCondition(conditionType)
	if (currentCondition == nil && oldCondition != nil) ||
		(currentCondition != nil && oldCondition == nil) {
		return true
	}
	if currentCondition == nil {
		return false
	}
	return (currentCondition.Status != oldCondition.Status) ||
		!currentCondition.LastTransitionTime.Inner.Equal(&oldCondition.LastTransitionTime.Inner) ||
		currentCondition.Reason != oldCondition.Reason
}
