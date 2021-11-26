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
	GitBranchGVK     = GroupVersion.WithKind("GitBranch")
	GitBranchListGVK = GroupVersion.WithKind("GitBranchList")
)

// GitBranch object for plugin
type GitBranch struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GitBranchSpec `json:"spec"`
}

// GitBranchBaseInfo branch base info
type GitBranchBaseInfo struct {
	GitRepo
	// Name branch name
	Name string `json:"name"`
}

// GitBranchSpec spec of branch
type GitBranchSpec struct {
	GitBranchBaseInfo
	// Protected the branch is protected
	Protected *bool `json:"protected,omitempty" yaml:"protected,omitempty"`
	// Default the branch is default branch (repo only have one default branch like main)
	Default *bool `json:"default,omitempty" yaml:"default,omitempty"`
	// DevelopersCanPush developer can push to this branch
	DevelopersCanPush *bool `json:"developersCanPush,omitempty" yaml:"developersCanPush,omitempty"`
	// DevelopersCanMerge developer can merge to this branch
	DevelopersCanMerge *bool `json:"developersCanMerge,omitempty" yaml:"developersCanMerge,omitempty"`
	// Commit latest commit's sha in this branch
	Commit      GitCommitInfo         `json:"commit"`
	WebURL      string                `json:"webURL"`
	DownloadURL DownloadURL           `json:"downloadURL"`
	Properties  *runtime.RawExtension `json:"properties,omitempty"`
}

// GitBranchList list of branch
type GitBranchList struct {
	metav1.TypeMeta `json:",inline"`
	ListMeta        `json:"metadata,omitempty"`

	Items []GitBranch `json:"items"`
}
