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
	SetConditionByErrorReason(conditionManager, condition, err, "")
}

// SetConditionByErrorReason sets a condition with error and reason
// If the given error is not nil will examine its reason and mark the condition as False
// otherwise will set condition to True
// If the given reason is empty, it will parse a string for reason from the error
func SetConditionByErrorReason(conditionManager apis.ConditionManager, condition apis.ConditionType, err error, reason string) {
	old := conditionManager.GetCondition(condition)

	if err != nil {
		if reason == "" {
			reason = ReasonForError(err)
		}
		message := err.Error()
		if old == nil || !old.IsFalse() || old.GetMessage() != message || old.GetReason() != reason {
			conditionManager.MarkFalse(condition, reason, message)
		}
	} else {
		if old == nil || !old.IsTrue() {
			conditionManager.MarkTrue(condition)
		}
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

// GetCondition will return the first condition pointer filter by type in conditions
func GetCondition(conditions apis.Conditions, t apis.ConditionType) *apis.Condition {
	if len(conditions) == 0 {
		return nil
	}

	for i := range conditions {
		if conditions[i].Type == t {
			return &conditions[i]
		}
	}

	return nil
}

// ConditionType is a camel-cased condition type.
type ConditionType string

const (
	// ConditionReady specifies that the resource is ready.
	// For long-running resources.
	ConditionReady ConditionType = "Ready"
	// ConditionSucceeded specifies that the resource has finished.
	// For resource which run to completion.
	ConditionSucceeded ConditionType = "Succeeded"
	// ConditionPending specifies that the resource is pending.
	// For resource which run is waiting to be executed.
	ConditionPending ConditionType = "Pending"
	// ConditionRunning specifies that the resource is running.
	// For resource which run is running.
	ConditionRunning ConditionType = "Running"
	// ConditionFailed specifies that the resource is failed.
	// For resource which run to failed.
	ConditionFailed ConditionType = "Failed"
	// ConditionDisabled specifies that the resource is disabled.
	// For resource which can't be run.
	ConditionDisabled ConditionType = "Disabled"
	// ConditionCanceled specifies that the resource is canceled.
	// For resource which run to canceled.
	ConditionCanceled ConditionType = "Canceled"
)
