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

// ArtifactFilterRegexList string slice
// TODO: move to pkg/apis/artifacts/v1alpha1
type ArtifactFilterRegexList []string

// ArtifactTagFilter contains regular expressions used to match the artifact's tag.
// TODO: move to pkg/apis/artifacts/v1alpha1
type ArtifactTagFilter struct {
	// Regex regular expressions matches tag
	Regex ArtifactFilterRegexList `json:"regex,omitempty"`
}

// ArtifactEnvFilter contains name and regular expressions used to match the artifact's env.
// TODO: move to pkg/apis/artifacts/v1alpha1
type ArtifactEnvFilter struct {
	// Name represent env name
	Name string `json:"name,omitempty"`
	// Regex regular expressions matches env
	Regex ArtifactFilterRegexList `json:"regex,omitempty"`
}

// ArtifactLabelFilter contains name and regular expressions used to match the artifact's label.
type ArtifactLabelFilter struct {
	// Name represent label name
	Name string `json:"name,omitempty"`
	// Regex regular expressions matches label
	Regex ArtifactFilterRegexList `json:"regex,omitempty"`
}

// ArtifactFilter artifact filter.
// TODO: move to pkg/apis/artifacts/v1alpha1
type ArtifactFilter struct {
	// +optional
	Tags []ArtifactTagFilter `json:"tags,omitempty"`
	// +optional
	Envs []ArtifactEnvFilter `json:"envs,omitempty"`
	// +optional
	Labels []ArtifactLabelFilter `json:"labels,omitempty"`
}

// ArtifactFilterSet filters for ArtifactPromotionPolicy
// TODO: move to pkg/apis/artifacts/v1alpha1
type ArtifactFilterSet struct {
	//  all nested filter expressions MUST evaluate to true in order for the all filter expression to be true.
	// +optional
	All []ArtifactFilter `json:"all,omitempty"`

	//  at least one nested filter expressions MUST evaluate to true in order for any filter expression to be true
	// +optional
	Any []ArtifactFilter `json:"any,omitempty"`
}
