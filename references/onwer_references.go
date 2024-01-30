/*
Copyright 2022 The Katanomi Authors.

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

// Package references contains methods to manage owner references
package references

import (
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// FindOwnerReference find the build owner in the owner references
func FindOwnerReference(owners []metav1.OwnerReference, gvk schema.GroupVersionKind) *metav1.OwnerReference {
	for _, owner := range owners {
		if owner.Kind == gvk.Kind &&
			owner.APIVersion == gvk.GroupVersion().String() {
			return &owner
		}
	}
	return nil
}

// FindOwnerReferenceWithGroupKind find the owner in the owner references with group kind
func FindOwnerReferenceWithGroupKind(owners []metav1.OwnerReference, gk schema.GroupKind) *metav1.OwnerReference {
	for _, owner := range owners {
		if owner.Kind != gk.Kind {
			continue
		}
		// For the legacy v1 resource which apiVersion is v1 without group
		// https://github.com/kubernetes/kubernetes/blob/5310e4f30e212a3d58b37fd07633c3b249627b53/pkg/apis/core/register.go#L26-L29
		// https://github.com/kubernetes/kubernetes/blob/5310e4f30e212a3d58b37fd07633c3b249627b53/pkg/apis/core/register.go#L55-L99
		// https://kubernetes.io/zh-cn/docs/reference/using-api/#api-versioning
		if owner.APIVersion == "v1" && gk.Group == "" {
			return &owner
		}

		if strings.HasPrefix(owner.APIVersion, gk.Group) {
			return &owner
		}
	}
	return nil
}
