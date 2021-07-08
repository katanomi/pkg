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
	"knative.dev/pkg/apis"
)

// ObjectCondition is integration of an ObjectReference and a apis.Condition to give a condition for a specific object
type ObjectCondition struct {
	apis.Condition         `json:",inline"`
	corev1.ObjectReference `json:",inline"`
}

// IsTheSame compares both object conditions and returns true if its the same object
// it uses the ObjectReference as comparisson method and ignore all condition related attributes
// making this useful when managing ObjectConditions
func (o ObjectCondition) IsTheSame(obj ObjectCondition) bool {
	return IsTheSameObject(o.ObjectReference, obj.ObjectReference)
}

// ObjectConditions collection of object conditions
// useful to store and iterate over different object conditions in the status of an object
type ObjectConditions []ObjectCondition

// SetObjectCondition updates a condition using the object reference, and if not found will append to the end and return a new slice
func (o ObjectConditions) SetObjectCondition(objc ObjectCondition) ObjectConditions {
	found := false
	for i, each := range o {
		if each.IsTheSame(objc) {
			o[i] = objc
			found = true
			break
		}
	}
	if !found {
		o = append(o, objc)
	}
	return o
}

// GetObjectConditionByObjRef get object conditon by object reference, returns nil if not found
func (o ObjectConditions) GetObjectConditionByObjRef(objref corev1.ObjectReference) *ObjectCondition {
	for _, each := range o {
		if IsTheSameObject(each.ObjectReference, objref) {
			return &each
		}
	}
	return nil
}

// RemoveObjectConditionByObjRef removes a conditions based on a object reference and returns a new slice
func (o ObjectConditions) RemoveObjectConditionByObjRef(objref corev1.ObjectReference) ObjectConditions {
	for i, each := range o {
		if IsTheSameObject(each.ObjectReference, objref) {
			return append(o[:i], o[i+1:]...)
		}
	}
	return o
}

// Manage returns a ObjectConditionSet to manage its contents
func (o ObjectConditions) Manage(accessor ObjectConditionAccessor) ObjectConditionManager {
	accessor.SetObjectConditions(o)
	return ManageObjectCondition(accessor)
}
