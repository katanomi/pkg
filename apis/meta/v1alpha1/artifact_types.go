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

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	authv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

var ArtifactGVK = GroupVersion.WithKind("Artifact")
var ArtifactListGVK = GroupVersion.WithKind("ArtifactList")

// Artifact object for plugins
// TODO: move to pkg/apis/artifacts/v1alpha1
type Artifact struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ArtifactSpec `json:"spec"`
}

// ArtifactSpec spec for repository
// TODO: add more necessary spec data
// TODO: move to pkg/apis/artifacts/v1alpha1
type ArtifactSpec struct {
	// Address API related access URL
	// +optional
	Address *duckv1.Addressable `json:"address,omitempty"`

	// Access stores the webconsole address if any
	// +optional
	Access *duckv1.Addressable `json:"access,omitempty"`

	// Type of repository content
	Type string `json:"type"`

	// Version of specified artifact
	Version string `json:"version"`

	// UpdatedTime updated time for repository
	// +optional
	UpdatedTime metav1.Time `json:"updatedTime"`

	// PullTime latest pull time for repository
	// +optional
	PullTime *metav1.Time `json:"pullTime,omitempty"`

	// Properties extended properties for Artifact
	// +optional
	Properties *runtime.RawExtension `json:"properties,omitempty"`
}

// ArtifactProperties we should promote all shared field in properties into spec
// but just defined it before we do the promotion
// TODO: move to pkg/apis/artifacts/v1alpha1
// +k8s:deepcopy-gen=false
type ArtifactProperties struct {
	// Size represent size of artifact
	Size int64 `json:"size,omitempty"`
	// Tags represent all tags of current artifact
	// +optional
	Tags []ArtifactTag `json:"tags,omitempty"`

	//ExtraAttrs artifact extra attributes
	// +optional
	ExtraAttrs ArtifactExtraAttrs `json:"extra_attrs,omitempty"`
	//References artifact references
	References []ArtifactReference `json:"references,omitempty"`
}

// ArtifactTag represent tag of artifact
// TODO: move to pkg/apis/artifacts/v1alpha1
// +k8s:deepcopy-gen=false
type ArtifactTag struct {
	// Name artifact tag name
	Name string `json:"name,omitempty"`
}

// ArtifactExtraAttrs represents artifact extra attributes like os
// TODO: move to pkg/apis/artifacts/v1alpha1
// +k8s:deepcopy-gen=false
type ArtifactExtraAttrs struct {
	// OS system os type
	OS string `json:"os,omitempty"`
	// Architecture artifact arch
	Architecture string `json:"architecture,omitempty"`
	// Variant system variant
	Variant string `json:"Variant,omitempty"`
	// Config represent original container image config
	Config v1.ImageConfig `json:"config,omitempty"`
}

// ArtifactReference represents reference of artifact
// TODO: move to pkg/apis/artifacts/v1alpha1
// +k8s:deepcopy-gen=false
type ArtifactReference struct {
	// ChildDigest child artifact digests
	ChildDigest string `json:"child_digest"`
	// Platform artifact platform
	Platform *v1.Platform `json:"platform"`
}

// PlatformString will return format os/arch
func (ref ArtifactReference) PlatformString() string {
	if ref.Platform == nil {
		return ""
	}

	return fmt.Sprintf("%s/%s", ref.Platform.OS, ref.Platform.Architecture)
}

// ArtifactList list of artifacts
// TODO: move to pkg/apis/artifacts/v1alpha1
type ArtifactList struct {
	metav1.TypeMeta `json:",inline"`
	ListMeta        `json:"metadata,omitempty"`

	Items []Artifact `json:"items"`
}

// ArtifactResourceAttributes returns a ResourceAttribute object to be used in a filter
func ArtifactResourceAttributes(verb string) authv1.ResourceAttributes {
	return authv1.ResourceAttributes{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: "artifacts",
		Verb:     verb,
	}
}
