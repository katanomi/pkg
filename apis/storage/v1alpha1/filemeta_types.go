/*
Copyright 2023 The Katanomi Authors.

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
)

// FileMetaGVK for GVK of FileMeta
var FileMetaGVK = GroupVersion.WithKind("FileMeta")

// FileMetaListGVK for GVK of FileMetaList
var FileMetaListGVK = GroupVersion.WithKind("FileMetaList")

// FileMeta object for sources
type FileMeta struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec FileMetaSpec `json:"spec"`
}

// FileMetaSpec spec for FileMeta
type FileMetaSpec struct {
	// Entry shows the index file of a directory if the contentType is a site
	// +optional
	Entry string `json:"entry,omitempty"`

	// Key for file key of file object, following format:  {storage-plugin-name}/-/{file-object-name}
	// +optional
	Key string `json:"key,omitempty"`

	// ContentType for file content type
	ContentType string `json:"contentType"`

	// ContentLength for file content size
	ContentLength int64 `json:"contentLength"`

	// FileType for file type
	// +optional
	FileType string `json:"fileType,omitempty"`
}

// FileMetaListOptions for list option of FileMeta
type FileMetaListOptions struct {
	Prefix     string `json:"prefix,omitempty"`
	Recursive  bool   `json:"recursive,omitempty"`
	StartAfter string `json:"startAfter,omitempty"`
	Limit      int    `json:"limit,omitempty"`
}

// FileMetaList list of FileMetas
type FileMetaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []FileMeta `json:"items"`
}

// FileMetaResourceAttributes returns a ResourceAttribute object to be used in a filter
func FileMetaResourceAttributes(verb string) authv1.ResourceAttributes {
	return authv1.ResourceAttributes{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: "filemetas",
		Verb:     verb,
	}
}
