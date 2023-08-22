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

const (
	// BinaryArtifactTypeDirectory for directory type
	BinaryArtifactTypeDirectory = "directory"
	// BinaryArtifactTypeFile for file type
	BinaryArtifactTypeFile = "file"
)

// ArtifactTypeBinaryProperties properties for binary artifact type
type ArtifactTypeBinaryProperties struct {
	// Type of display type of artifact: directory, file
	Type string `json:"type"`
	// Name of artifact
	Name string `json:"name"`
	// Context of artifact
	Context string `json:"context"`
	// Size of artifact
	Size int64 `json:"size"`
}
