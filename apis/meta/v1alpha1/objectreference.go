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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

// IsTheSameObjRef compares two corev1.ObjectReference comparing:
// APIVersion, Kind, Name and Namespace. All other attributes are ignored
func IsTheSameObject(obj, compared corev1.ObjectReference) bool {
	return obj.APIVersion == compared.APIVersion &&
		obj.Kind == compared.Kind &&
		obj.Name == compared.Name &&
		obj.Namespace == compared.Namespace
}

// IsTheSameObjectReference uses pointers to make comparison of objects
func IsTheSameObjectReference(obj, compared *corev1.ObjectReference) bool {
	if (obj == nil && compared != nil) || (obj != nil && compared == nil) {
		return false
	}
	return (obj == nil && compared == nil) || (obj != nil && IsTheSameObject(*obj, *compared))
}

// GetObjectReferenceFromObject extracts an object reference from an object
func GetObjectReferenceFromObject(obj metav1.Object, opts ...ObjectRefOptionsFunc) (ref corev1.ObjectReference) {
	ref.Name = obj.GetName()
	for _, o := range opts {
		o(obj, &ref)
	}
	return
}

// GetNamespacedNameFromRef returns a types.NamespacedName from an object reference
func GetNamespacedNameFromRef(ref *corev1.ObjectReference) (named types.NamespacedName) {
	if ref != nil {
		named.Name = ref.Name
		named.Namespace = ref.Namespace
	}
	return
}

// +k8s:deepcopy-gen=false
type ObjectRefOptionsFunc func(obj metav1.Object, ref *corev1.ObjectReference)

func ObjectRefWithTypeMeta() ObjectRefOptionsFunc {
	return func(obj metav1.Object, ref *corev1.ObjectReference) {
		if runobj, ok := obj.(runtime.Object); ok {
			objkind := runobj.GetObjectKind()
			ref.APIVersion = objkind.GroupVersionKind().GroupVersion().String()
			ref.Kind = objkind.GroupVersionKind().Kind
		}
	}
}

func ObjectRefWithUID() ObjectRefOptionsFunc {
	return func(obj metav1.Object, ref *corev1.ObjectReference) {
		ref.UID = obj.GetUID()
	}
}

func ObjectRefWithNamespace() ObjectRefOptionsFunc {
	return func(obj metav1.Object, ref *corev1.ObjectReference) {
		ref.Namespace = obj.GetNamespace()
	}
}
