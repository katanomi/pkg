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
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	authv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

var ArtifactGVK = GroupVersion.WithKind("Artifact")
var ArtifactListGVK = GroupVersion.WithKind("ArtifactList")

// Artifact object for plugins
type Artifact struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ArtifactSpec `json:"spec"`
}

// ArtifactSpec spec for repository
// TODO: add more necessary spec data
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

	// Size represent size of artifact
	Size int64 `json:"size,omitempty"`

	// Tags represent all tags of current artifact
	// +optional
	Tags []ArtifactTag `json:"tags,omitempty"`

	// imageconfig represent original container image config
	ImageConfig v1.ImageConfig `json:"config"`
}

// ArtifactProperties we should promote all shared field in properties into spec
// but just defined it before we do the promotion
type ArtifactProperties struct {
	// Size represent size of artifact
	Size int64 `json:"size,omitempty"`
	// Tags represent all tags of current artifact
	// +optional
	Tags       []ArtifactTag      `json:"tags,omitempty"`
	ExtraAttrs ArtifactExtraAttrs `json:"extra_attrs,omitempty"`
}

// ArtifactTag represent tag of artifact
type ArtifactTag struct {
	Name string `json:"tag,omitempty"`
}

type ArtifactExtraAttrs struct {
	// os ar
	OS           string `json:"os,omitempty"`
	Architecture string `json:"architecture,omitempty"`
	Variant      string `json:"Variant",omitempty`
	// Config represent original container image config
	Config v1.ImageConfig `json:"config"`
}

type ArtifactReference struct {
	ChildDigest string       `json:"child_digest"`
	Platform    *v1.Platform `json:"platform"`
}

// ArtifactList list of artifacts
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
