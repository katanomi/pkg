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
	GitCommitCommentGVK     = GroupVersion.WithKind("GitCommitComment")
	GitCommitCommentListGVK = GroupVersion.WithKind("GitCommitCommentList")
)

// GitCommitComment object for plugin
type GitCommitComment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GitCommitCommentSpec `json:"spec"`
}

// GitCommitCommentSpec spec for commit comment
type GitCommitCommentSpec struct {
	// Note content
	Note string `json:"note"`
	// Path file path
	Path string `json:"path"`
	// Line comment line number
	Line int `json:"line"`
	// LineType
	LineType *string `json:"lineType"`
	// Author comment author
	Author     GitUserBaseInfo       `json:"author"`
	Properties *runtime.RawExtension `json:"properties,omitempty"`
}

// GitCommitCommentList list of commit comment
type GitCommitCommentList struct {
	metav1.TypeMeta `json:",inline"`
	ListMeta        `json:"metadata,omitempty"`

	Items []GitCommitComment `json:"items"`
}
