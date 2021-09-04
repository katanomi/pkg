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
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Stores a list of triggered information such as: Entity that triggered,
// reference of an object that could have triggered, and event that triggered.
type TriggeredBy struct {
	// Reference to the user that triggered the object. Any Kubernetes `Subject` is accepted.
	// +optional
	User *rbacv1.Subject `json:"user,omitempty"`

	// Cloud Event data for the event that triggered.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:nullable
	// +optional
	CloudEvent *runtime.RawExtension `json:"cloudEvent,omitempty"`

	// Reference to another object that might have triggered this object
	// +optional
	Ref *corev1.ObjectReference `json:"ref,omitempty"`

	// Date time of creation of triggered event. Will match a resource's metadata.creationTimestamp
	// it is added here for convinience only
	// +optional
	TriggeredTimestamp *metav1.Time `json:"triggeredTimestamp,omitempty"`
}
