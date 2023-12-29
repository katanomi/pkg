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
	"k8s.io/apimachinery/pkg/runtime"
)

// DisplayColumn tells the outside how the backend fields should be displayed.
type DisplayColumn struct {
	// Name column name, should be the only value.
	Name string `json:"name,omitempty"`

	// Field the path to read data, usually jsonPath.
	Field string `json:"field,omitempty"`

	// DisplayName index of impression data used for matching.
	DisplayName string `json:"displayName,omitempty"`

	// TranslationPrefix index of translation prefixes.
	// +optional
	TranslationPrefix string `json:"translationPrefix,omitempty"`

	// Properties extended properties for DisplayColum
	// +optional
	Properties *runtime.RawExtension `json:"properties,omitempty"`
}

// DisplayColumns tells the outside how the backend fields should be displayed.
type DisplayColumns []DisplayColumn
