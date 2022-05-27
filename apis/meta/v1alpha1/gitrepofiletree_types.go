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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	GitRepositoryFileTreeGVK = GroupVersion.WithKind("GitRepositoryFileTree")
)

// GitRepositoryFileTree object for plugins
type GitRepositoryFileTree struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GitRepositoryFileTreeSpec `json:"spec"`
}

// GitRepositoryFileTreeSpec spec for repository's file
type GitRepositoryFileTreeSpec struct {
	Tree []GitRepositoryFileTreeNode `json:"tree"`
}

type GitRepositoryFileTreeNodeType string

const (
	// TreeNodeBlobType represents a file
	TreeNodeBlobType GitRepositoryFileTreeNodeType = "blob"
	// TreeNodeTreeType represents a folder
	TreeNodeTreeType GitRepositoryFileTreeNodeType = "tree"
)

// GitRepositoryFileTreeNode represents a node in the file system
type GitRepositoryFileTreeNode struct {
	// Sha is the ID of the node
	Sha string `json:"sha"`
	// Name is the name of the node
	Name string `json:"name"`
	// Path is the path of the node
	Path string `json:"path"`
	// Type is the type of the node
	Type GitRepositoryFileTreeNodeType `json:"type"`
	// Mode indicates the permission level of the file
	Mode string `json:"mode"`
}

// GitRepoFileTreeOption requesting parameters for the File Tree API
type GitRepoFileTreeOption struct {
	GitRepo
	Path      string `json:"path"`
	TreeSha   string `json:"tree_sha"`
	Recursive bool   `json:"recursive"`
}
