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
	pipev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

// Template defines the desired state of Template
type Template struct {
	// Name is the name of the Template
	// If not empty, must be unique in the list of templates
	// +optional
	Name string `json:"name,omitempty"`

	// ResolverRef can be used to refer to a Resource in a remote
	// location like a git repo.
	TemplateRef pipev1beta1.ResolverRef `json:"templateRef,omitempty"`

	// Params define a list of parameters for the Template
	// +optional
	Params []pipev1beta1.Param `json:"params,omitempty"`

	// Metadata contains the labels and annotations for template
	// This information will be automatically populated when obtaining it
	// without the need for the user to configure it actively.
	// +optional
	Metadata pipev1beta1.PipelineTaskMetadata `json:"metadata,omitempty"`
}
