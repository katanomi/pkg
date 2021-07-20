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
	"k8s.io/apimachinery/pkg/api/errors"
	"knative.dev/pkg/apis"
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
