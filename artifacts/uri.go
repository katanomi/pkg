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

package artifacts

import (
	"fmt"
	"strings"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
)

// Protocol represent artifact transport protocol
type Protocol string

var (
	// ProtocolDocker docker protocol
	ProtocolDocker Protocol = "docker"
	// ProtocolHelmChart helm chart protocol
	ProtocolHelmChart Protocol = "chart"
	// ProtocolOCI oci protocol
	ProtocolOCI Protocol = "oci"
)

// DigestAlgorithm artifact digest algorithm
type DigestAlgorithm string

const (
	// SHA256 artifact digest algorithm is sha256
	SHA256 DigestAlgorithm = "sha256"
)

// URI represents artifact uri , like docker://gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/nop:v0.37.1@sha256:04411f239bc7144c3248b53af6741c2726eaddbe9b9cf62a24cf812689cc3223
type URI struct {
	// Protocol represent artifact transport protocol, it is optional, default value is docker
	Protocol string
	// Host represent artifact host
	Host string
	// Path represent artifact path
	Path string
	// Tag represent artifact tag
	Tag string
	// Digest represent artifact digest
	Digest string
	// Algorithm digest algorithm
	Algorithm DigestAlgorithm

	// Raw represent original uri
	Raw string
}

// DigestString will return format of Algorithm:Digest
func (u URI) DigestString() string {
	if u.Digest == "" {
		return ""
	}
	return fmt.Sprintf("%s:%s", u.Algorithm, u.Digest)
}

// Version will return sha256:xx or tag .
func (u URI) Version() string {
	if u.Digest != "" {
		return u.DigestString()
	} else {
		return u.Tag
	}
}

// ParseURI parse uri to URI struct
func ParseURI(uri string, t metav1alpha1.ArtifactType) (URI, error) {
	var u = URI{Raw: uri}

	str := uri

	protocolIndex := strings.Index(uri, "://")
	if protocolIndex <= 0 {
		u.Protocol = string(ProtocolDocker)
	} else {
		u.Protocol = uri[0:protocolIndex]
		str = str[protocolIndex+len("://"):]
	}

	if t == metav1alpha1.OCIContainerImageArtifactParameterType {
		u.Protocol = string(ProtocolDocker)
	}
	if t == metav1alpha1.OCIHelmChartArtifactParameterType {
		u.Protocol = string(ProtocolHelmChart)
	}

	hostIndex := strings.Index(str, "/")
	if hostIndex <= 0 {
		return u, fmt.Errorf("error format of uri: %s", uri)
	}
	u.Host = str[0:hostIndex]

	str = str[hostIndex+1:]

	digestStartIndex := strings.Index(str, "@")
	if digestStartIndex > 0 {
		digestStr := str[digestStartIndex+1:]
		algorithmIndex := strings.Index(digestStr, ":")
		if algorithmIndex <= 0 {
			return u, fmt.Errorf("algorithm is invalid")
		}
		u.Algorithm = DigestAlgorithm(digestStr[0:algorithmIndex])
		u.Digest = digestStr[algorithmIndex+1:]

		str = str[0:digestStartIndex]
	}

	tagIndex := strings.Index(str, ":")
	if tagIndex > 0 {
		u.Tag = str[tagIndex+1:]
		u.Path = str[0:tagIndex]
	} else {
		u.Path = str
	}

	return u, nil
}
