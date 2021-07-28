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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

type ResourceSubType string

const (
	// OCI artifact registry project
	ResourceSubTypeImageRegistry ResourceSubType = "ImageRegistry"
	// Code repository project
	ResourceSubTypeCodeRepository ResourceSubType = "CodeRepository"
)

var ResourceGVK = GroupVersion.WithKind("Resource")
var ResourceListGVK = GroupVersion.WithKind("ResourceList")

// Resource object for plugins
type Resource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ResourceSpec `json:"spec"`
}

// ResourceSpec spec for a generic resource response
type ResourceSpec struct {
	// Address API related access URL
	// +optional
	Address *duckv1.Addressable `json:"address,omitempty"`

	// Access stores the webconsole address if any
	// +optional
	Access *duckv1.Addressable `json:"access,omitempty"`

	// Type of resource
	Type string `json:"type"`

	// SubType of resource
	SubType string `json:"subType"`

	// Version of specified resource
	// +optional
	Version string `json:"version,omitempty"`

	// Properties extended properties for Resource
	// +optional
	Properties *runtime.RawExtension `json:"properties,omitempty"`
}

// ResourceList list of resources
type ResourceList struct {
	metav1.TypeMeta `json:",inline"`
	ListMeta        `json:"metadata,omitempty"`

	Items []Resource `json:"items"`
}
