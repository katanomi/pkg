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
	// "regexp"
	"strings"
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

func GetBinaryObjectFromValues(ctx context.Context, array []string) (versions []ArtifactVersion) {
	for _, item := range array {
		versions = append(versions, ArtifactVersion{
			Type: ArtifactTypeBinary,
			URL:  item,
		})
	}
	return
}

// GetHelmChartObjectFromURLValues return a helm chart artifact using url and a list of tags
func GetHelmChartObjectFromURLValues(ctx context.Context, url string, tags ...string) (versions []ArtifactVersion) {
	if strings.TrimSpace(url) != "" {
		versions = append(versions, ArtifactVersion{
			Type:     ArtifactTypeHelmChart,
			URL:      url,
			Versions: tags,
		})
	}
	return
}

// GetContainerImageObjectFromURLValues return a container image artifact using url, digest and a list of tags
func GetContainerImageObjectFromURLValues(ctx context.Context, url, digest string, tags ...string) (versions []ArtifactVersion) {
	if strings.TrimSpace(url) != "" {
		versions = append(versions, ArtifactVersion{
			Type:     ArtifactTypeContainerImage,
			URL:      url,
			Digest:   digest,
			Versions: tags,
		})
	}
	return
}

// GetContainerImageFromValues return a list of container image artifacts using url, digest and tags
func GetContainerImageFromValues(ctx context.Context, array []string) (versions []ArtifactVersion) {
	// will use the digest as an index
	// to attach tags to the same artifact
	// must provide the same digest otherwise will consider to be
	// different artifacts
	digestIndex := map[string]int{}
	for _, value := range array {
		url, digest, tag := ExtractRepositoryDigestTag(value)

		artifact := ArtifactVersion{
			Type:   ArtifactTypeContainerImage,
			URL:    url,
			Digest: digest,
		}
		idx, hasDigest := digestIndex[digest]
		if digest != "" && hasDigest {
			artifact = versions[idx]
		}
		if tag != "" {
			artifact.Versions = append(artifact.Versions, tag)
		}
		if hasDigest {
			versions[idx] = artifact
		} else {
			if digest != "" {
				digestIndex[digest] = len(versions)
			}
			versions = append(versions, artifact)
		}
	}
	return
}

// ExtractRepositoryDigestTag takes a oci artifact url and extracts
// url, digest and tag
func ExtractRepositoryDigestTag(value string) (url, digest, tag string) {
	digestIndex := strings.Index(value, "@sha256")
	if digestIndex > 0 {
		digest = value[digestIndex+1:]
		value = value[:digestIndex]
	}
	tagIndex := strings.LastIndex(value, ":")
	if tagIndex > 0 && tagIndex > strings.LastIndex(value, "/") {
		tag = value[tagIndex+1:]
		value = value[:tagIndex]
	}
	url = value
	return
}
