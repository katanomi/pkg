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
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
)

// ObjectConditionAccessor sets and gets ObjectConditions
// +k8s:deepcopy-gen=false
type ObjectConditionAccessor interface {
	GetObjectConditions() ObjectConditions
	SetObjectConditions(ObjectConditions)
}

// ObjectConditionChanger sets and removes objects
// +k8s:deepcopy-gen=false
type ObjectConditionChanger interface {
	SetObjectCondition(objc ObjectCondition)
	GetObjectConditionByObjRef(objref corev1.ObjectReference) *ObjectCondition
	RemoveObjectConditionByObjRef(objref corev1.ObjectReference)
}

// ObjectConditionManager manages
// +k8s:deepcopy-gen=false
type ObjectConditionManager interface {
	ObjectConditionAccessor
	ObjectConditionChanger

	// MarkTrue sets the status of t to true, and then marks the happy condition to
	// true if all dependents are true.
	MarkTrue(t corev1.ObjectReference)

	// MarkTrueWithReason sets the status of t to true with the reason, and then marks the happy
	// condition to true if all dependents are true.
	MarkTrueWithReason(t corev1.ObjectReference, reason, messageFormat string, messageA ...interface{})

	// MarkUnknown sets the status of t to Unknown and also sets the happy condition
	// to Unknown if no other dependent condition is in an error state.
	MarkUnknown(t corev1.ObjectReference, reason, messageFormat string, messageA ...interface{})

	// MarkFalse sets the status of t and the happy condition to False.
	MarkFalse(t corev1.ObjectReference, reason, messageFormat string, messageA ...interface{})

	// SetConditionType sets a condition type for object condition
	SetConditionType(t corev1.ObjectReference, conditionType apis.ConditionType)

	// SetSeverity sets a severity for object condition
	SetSeverity(t corev1.ObjectReference, severity apis.ConditionSeverity)
}

// ManageObjectCondition returns a ObjectConditionManager
func ManageObjectCondition(accessor ObjectConditionAccessor) ObjectConditionManager {
	return &ObjectConditionSet{accessor: accessor}
}

// ObjectConditionSet set of object conditions managed in a controller and added to specific object types
// it helps iterate and update its contents
// +k8s:deepcopy-gen=false
type ObjectConditionSet struct {
	accessor ObjectConditionAccessor
}

// GetObjectConditions get conditions
func (o *ObjectConditionSet) GetObjectConditions() ObjectConditions {
	return o.accessor.GetObjectConditions()
}

// SetObjectConditions set conditions
func (o *ObjectConditionSet) SetObjectConditions(objcs ObjectConditions) {
	o.accessor.SetObjectConditions(objcs)
}

// SetObjectCondition sets object into condition slice in a upsert method
func (o *ObjectConditionSet) SetObjectCondition(objc ObjectCondition) {
	o.accessor.SetObjectConditions(o.accessor.GetObjectConditions().SetObjectCondition(objc))

}

// RemoveObjectConditionByObjRef removes item by corev1.ObjectReference
func (o *ObjectConditionSet) RemoveObjectConditionByObjRef(objref corev1.ObjectReference) {
	o.accessor.SetObjectConditions(o.accessor.GetObjectConditions().RemoveObjectConditionByObjRef(objref))
}

// GetObjectConditionByObjRef returns object condition by object reference, returns nil if not found
func (o *ObjectConditionSet) GetObjectConditionByObjRef(objref corev1.ObjectReference) *ObjectCondition {
	return o.accessor.GetObjectConditions().GetObjectConditionByObjRef(objref)
}

// MarkTrue sets the status of t to true, and then marks the happy condition to
// true if all dependents are true.
func (o *ObjectConditionSet) MarkTrue(objref corev1.ObjectReference) {
	o.markStatus(objref, corev1.ConditionTrue, "", "")
}

// MarkTrueWithReason sets the status of t to true with the reason
func (o *ObjectConditionSet) MarkTrueWithReason(objref corev1.ObjectReference, reason, messageFormat string, messageA ...interface{}) {
	o.markStatus(objref, corev1.ConditionTrue, reason, messageFormat, messageA...)
}

// MarkUnknown sets the status of t to Unknown and also sets the happy condition
// to Unknown if no other dependent condition is in an error state.
func (o *ObjectConditionSet) MarkUnknown(objref corev1.ObjectReference, reason, messageFormat string, messageA ...interface{}) {
	o.markStatus(objref, corev1.ConditionUnknown, reason, messageFormat, messageA...)
}

// MarkFalse sets the status of t and the happy condition to False.
func (o *ObjectConditionSet) MarkFalse(objref corev1.ObjectReference, reason, messageFormat string, messageA ...interface{}) {
	o.markStatus(objref, corev1.ConditionFalse, reason, messageFormat, messageA...)
}

// SetConditionType sets a condition type for object condition
func (o *ObjectConditionSet) SetConditionType(objref corev1.ObjectReference, conditionType apis.ConditionType) {
	if objCondition := o.GetObjectConditionByObjRef(objref); objCondition != nil {
		objCondition.Type = conditionType
		o.SetObjectCondition(*objCondition)
	}
}

// SetSeverity sets a severity for object condition
func (o *ObjectConditionSet) SetSeverity(objref corev1.ObjectReference, severity apis.ConditionSeverity) {
	if objCondition := o.GetObjectConditionByObjRef(objref); objCondition != nil {
		objCondition.Severity = severity
		o.SetObjectCondition(*objCondition)
	}
}

// MarkStatus set status
func (o *ObjectConditionSet) markStatus(objref corev1.ObjectReference, cond corev1.ConditionStatus, reason, messageFormat string, messageA ...interface{}) {
	var objCondition *ObjectCondition
	if objCondition = o.GetObjectConditionByObjRef(objref); objCondition == nil {
		objCondition = &ObjectCondition{ObjectReference: objref}
	}
	objCondition.Status = cond
	objCondition.Reason = reason
	objCondition.Message = fmt.Sprintf(messageFormat, messageA...)
	if objCondition.Status != cond || objCondition.Reason != reason || objCondition.LastTransitionTime.Inner.IsZero() {
		objCondition.LastTransitionTime = apis.VolatileTime{Inner: metav1.NewTime(time.Now())}
	}
	o.SetObjectCondition(*objCondition)
}

// ReplaceObjectConditions will replace all conditions in source used by replaced, and remove all conditions that not exists in replaced
func ReplaceObjectConditions(source []ObjectCondition, replaced []ObjectCondition) (res []ObjectCondition) {
	shouldRemoved := []ObjectCondition{}
	if len(source) == 0 {
		return replaced
	}

	for _, item := range source {
		contains := false
		for _, cond := range replaced {
			if IsTheSameObject(item.ObjectReference, cond.ObjectReference) {
				contains = true
			}
		}

		if !contains {
			shouldRemoved = append(shouldRemoved, item)
		}
	}

	for _, cond := range replaced {
		source = ObjectConditions(source).SetObjectCondition(cond)
	}

	for _, cond := range shouldRemoved {
		source = ObjectConditions(source).RemoveObjectConditionByObjRef(cond.ObjectReference)
	}
	return source
}

// AggregateObjectCondition aggregated object conditions to apis.Conditioion
func AggregateObjectCondition(conds []ObjectCondition, condType apis.ConditionType) *apis.Condition {
	cond := apis.Condition{
		Type:     condType,
		Status:   corev1.ConditionTrue,
		Severity: apis.ConditionSeverityInfo,
	}

	if len(conds) == 0 {
		cond.Message = "No targets need be synced"
		return &cond
	}

	falseContains := false
	unknownContains := false
	messagesIndex := map[string]struct{}{}
	messages := []string{}

	for _, item := range conds {
		key := item.ObjectReference.Namespace + "/" + item.ObjectReference.Name
		if item.Status == corev1.ConditionFalse {
			falseContains = true
		}
		if item.Status == corev1.ConditionUnknown {
			unknownContains = true
		}
		if item.Status != corev1.ConditionTrue {
			if _, ok := messagesIndex[key]; !ok {
				messagesIndex[key] = struct{}{}
				messages = append(messages, fmt.Sprintf("%s: %s", key, item.Message))
			}
		}
	}

	if falseContains {
		cond.Status = corev1.ConditionFalse
	}
	if unknownContains {
		cond.Status = corev1.ConditionUnknown
	}

	cond.Message = strings.Join(messages, ";  ")

	return &cond
}
