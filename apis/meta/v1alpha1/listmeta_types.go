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
)

// ListMeta extension of the metav1.ListMeta with paging related data
type ListMeta struct {
	metav1.ListMeta `json:",inline"`

	// TotalItems returned in the list
	TotalItems int `json:"totalItems"`
}

// ListOptions options for list
type ListOptions struct {
	// ItemsPerPage desired number of items to be returned in each page
	ItemsPerPage int `json:"itemsPerPage"`

	// Page desired to be returned
	Page int `json:"page"`

	// Custom search options
	// +optional
	Search map[string][]string `json:",inline"`
}

// RepositoryOptions list repositroy path params
type RepositoryOptions struct {
	// project name
	Project string         `json:"project"`
	SubType ProjectSubType `json:"subType"`
}

// ArtifactOptions path params
type ArtifactOptions struct {
	RepositoryOptions

	// repository name
	Repository string `json:"repository"`

	// artifact name
	Artifact string `json:"artifact"`
}
