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

package v1alpha1

import (
	authv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

var (
	GitRepositoryTagGVK     = GroupVersion.WithKind("GitRepositoryTag")
	GitRepositoryTagListGVK = GroupVersion.WithKind("GitRepositoryTagList")
)

// GitRepositoryTag object for plugin
type GitRepositoryTag struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GitRepositoryTagSpec `json:"spec"`
}

// GitRepositoryTagSpec spec for commit
type GitRepositoryTagSpec struct {
	GitRepositoryTagInfo `json:",inline"`

	// Address for commit url for code repository web server
	// +optional
	Address *duckv1.Addressable `json:"address,omitempty"`

	// Message tag message
	// +optional
	Message *string `json:"message,omitempty"`

	// Properties extended properties for tag
	// +optional
	Properties *runtime.RawExtension `json:"properties,omitempty"`
}

// GitRepositoryTagInfo common tag info
type GitRepositoryTagInfo struct {
	// Name tag's name
	Name string `json:"name"`
	// SHA tags's sha
	SHA *string `json:"sha,omitempty"`
}

// GitRepositoryTagList list of commits
type GitRepositoryTagList struct {
	metav1.TypeMeta `json:",inline"`
	ListMeta        `json:"metadata,omitempty"`

	Items []GitRepositoryTag `json:"items"`
}

// GitRepositoryTagResourceAttributes returns a ResourceAttribute object to be used in a filter
func GitRepositoryTagResourceAttributes(verb string) authv1.ResourceAttributes {
	return authv1.ResourceAttributes{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: "gitrepositorytags",
		Verb:     verb,
	}
}
