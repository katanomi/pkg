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
	"context"
)

// ArtifactVersionCollection collection of artifacts versions
type ArtifactVersionCollection struct {
	// ArtifactVersions all the artifacts
	ArtifactVersions []ArtifactVersion `json:"artifactVersions"`
}

// ArtifactVersion artifacts
type ArtifactVersion struct {
	// Type of artifact
	Type ArtifactType `json:"type"`
	// URL of artifact
	URL string `json:"url"`
	// Digest means artifact digest
	// can be used to store a unique identifier
	// of the artifact version
	Digest string `json:"digest,omitempty"`

	// Versions of current artifact
	// +optional
	Versions []string `json:"versions,omitempty"`
}

func (ArtifactVersion) GetBinaryObjectFromValues(ctx context.Context, array []string) (versions []ArtifactVersion) {
	for _, item := range array {
		versions = append(versions, ArtifactVersion{
			Type: ArtifactTypeBinary,
			URL:  item,
		})
	}
	return
}
