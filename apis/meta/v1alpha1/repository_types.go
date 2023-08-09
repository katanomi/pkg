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

	authv1 "k8s.io/api/authorization/v1"
)

// RepositorySubType stores a specific repository subtype
type RepositorySubType string

func (r RepositorySubType) String() string {
	return string(r)
}

const (
	// DefaultRepositorySubType default repository subtype
	DefaultRepositorySubType RepositorySubType = "Repository"

	// ImageRepositorySubType OCI artifact repository subtype
	ImageRepositorySubType RepositorySubType = "ImageRepository"

	// CodeRepositorySubType Code repository subtype
	CodeRepositorySubType RepositorySubType = "CodeRepository"

	// FileDirectorySubType Raw repository subtype
	FileDirectorySubType RepositorySubType = "FileDirectory"
)

var RepositoryGVK = GroupVersion.WithKind("Repository")
var RepositoryListGVK = GroupVersion.WithKind("RepositoryList")

// Repository object for plugins
type Repository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec RepositorySpec `json:"spec"`
}

// RepositorySpec spec for repository
// TODO: add more necessary spec data
type RepositorySpec struct {
	// Address API related access URL
	// +optional
	Address *duckv1.Addressable `json:"address,omitempty"`

	// Access stores the webconsole address if any
	// +optional
	Access *duckv1.Addressable `json:"access,omitempty"`

	// Type of repository content
	Type RepositorySubType `json:"type"`

	// NamespaceRefs for which this project is already bound to
	// +optional
	NamespaceRefs []corev1.ObjectReference `json:"namespaceRefs,omitempty"`

	// UpdatedTime updated time for repository
	// +optional
	UpdatedTime metav1.Time `json:"updatedTime"`

	// Properties extended properties for Repository
	// +optional
	Properties *runtime.RawExtension `json:"properties,omitempty"`
}

// RepositoryList list of repositories
type RepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	ListMeta        `json:"metadata,omitempty"`

	Items []Repository `json:"items"`
}

// RepositoryResourceAttributes returns a ResourceAttribute object to be used in a filter
func RepositoryResourceAttributes(verb string) authv1.ResourceAttributes {
	return authv1.ResourceAttributes{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: "repositories",
		Verb:     verb,
	}
}
