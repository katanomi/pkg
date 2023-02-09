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
	GitRepoFileGVK     = GroupVersion.WithKind("GitRepositoryFile")
	GitRepoFileListGVK = GroupVersion.WithKind("GitRepositoryFileList")
)

// GitRepoFile object for plugins
type GitRepoFile struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GitRepoFileSpec `json:"spec"`
}

// GitRepoFileSpec spec for repository's file
type GitRepoFileSpec struct {
	GitCommitBasicInfo
	// FileName file name
	FileName string `json:"fileName" yaml:"fileName"`
	// FilePath file path
	FilePath string `json:"filePath" yaml:"filePath"`
	// Size file size
	Size int64 `json:"size"`
	// Encoding maybe text or base64 when path is a file
	Encoding *string `json:"encoding"`
	// Content file content
	Content []byte `json:"content"`
	// NodeSHA same as sha for file tree node.
	NodeSHA    string                `json:"NodeSHA"`
	Properties *runtime.RawExtension `json:"properties,omitempty"`
}
