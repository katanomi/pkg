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
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

var ProjectGVK = GroupVersion.WithKind("Project")
var ProjectListGVK = GroupVersion.WithKind("ProjectList")

// Project object for plugins
type Project struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ProjectSpec `json:"spec"`
}

// ProjectSpec spec for project
// TODO: add more necessary spec data
type ProjectSpec struct {
	// Public defines if a project is public or not
	Public bool `json:"public"`

	// Address API related access URL
	// +optional
	Address *duckv1.Addressable `json:"address,omitempty"`

	// Access stores the webconsole address if any
	// +optional
	Access *duckv1.Addressable `json:"access,omitempty"`

	// NamespaceRefs for which this project is already bound to
	// +optional
	NamespaceRefs []*corev1.ObjectReference `json:"namespaceRefs,omitempty"`

	// Properties extended properties for Project
	Properties *runtime.RawExtension `json:"properties,omitempty"`
}

// ProjectList list of projects
type ProjectList struct {
	metav1.TypeMeta `json:",inline"`
	ListMeta        `json:"metadata,omitempty"`

	Items []Project `json:"items"`
}
