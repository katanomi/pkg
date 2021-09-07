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
)

var (
	GitCommitStatusGVK     = GroupVersion.WithKind("GitCommitStatus")
	GitCommitStatusListGVK = GroupVersion.WithKind("GitCommitStatusList")
)

// GitCommitStatus object for plugin
type GitCommitStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GitCommitStatusSpec `json:"spec"`
}

type GitCommitStatusSpec struct {
	// ID status id
	ID int `json:"id"`
	// SHA commit sha
	SHA string `json:"sha"`
	// Ref commit ref
	Ref string `json:"ref"`
	// Status
	Status string `json:"status"`
	// CreatedAt status create time
	CreatedAt metav1.Time `json:"createdAt"`
	// Name status name
	Name string `json:"name"`
	// Author status author
	Author GitUserBaseInfo `json:"author"`
	// Description status description
	Description string `json:"description"`
	// TargetURL
	TargetURL  string                `json:"targetUrl"`
	Properties *runtime.RawExtension `json:"properties,omitempty"`
}

// GitCommitStatusList list of commit status
type GitCommitStatusList struct {
	metav1.TypeMeta `json:",inline"`
	ListMeta        `json:"metadata,omitempty"`

	Items []GitCommitStatus `json:"items"`
}
