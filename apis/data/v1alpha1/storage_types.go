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
	"time"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	authv1 "k8s.io/api/authorization/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var StorageGVK = GroupVersion.WithKind("Storage")
var StorageListGVK = GroupVersion.WithKind("StorageList")

// Storage object for data service
type Storage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec StorageSpec `json:"spec"`
}

// StorageSpec spec for storage
type StorageSpec struct {
	// Resource describe storage base info
	Resource Resource `json:"resource"`

	// GC storage gc policy
	GC GC `json:"gc"`

	// Payloads save artifact data
	Payloads []*Payload `json:"payloads"`
}

// Resource describe storage base info
type Resource struct {
	// IntegrationClassName sets the name of IntegrationClass that this integration is implemented
	IntegrationClassName string `json:"integrationClassName"`

	// Uri stores the artifact address
	Uri string `json:"uri"`

	// ResourceType storage resource type
	ResourceType StorageResourceType `json:"resourceType"`

	// SecretRef stores a secret that is used to access the integrated service
	SecretRef *corev1.ObjectReference `json:"secretRef,omitempty"`
}

// GC storage gc policy
type GC struct {
	// Period duration for data to be purged from the system after creation. This is a hard deadline if GC fails to execute
	// +optional
	Period time.Duration `json:"period,omitempty"`
}

// Payload save artifact data
type Payload struct {
	// Type payload type, eg. build or deploy
	Type PayloadType `json:"type"`

	// payload uid
	// +optional
	Uid string `json:"uid,omitempty"`

	// CreatedTime describe payload upload time
	// +optional
	CreatedTime metav1.Time `json:"createdTime,omitempty"`

	// Properties extended properties for payload
	// +optional
	Properties *runtime.RawExtension `json:"properties,omitempty"`
}

// StorageList list of storages
type StorageList struct {
	metav1.TypeMeta       `json:",inline"`
	metav1alpha1.ListMeta `json:"metadata,omitempty"`

	Items []Storage `json:"items"`
}

// StorageResourceAttributes returns a ResourceAttribute object to be used in a filter
func StorageResourceAttributes(verb string) authv1.ResourceAttributes {
	return authv1.ResourceAttributes{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: "artifacts",
		Verb:     verb,
	}
}
