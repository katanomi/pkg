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
	"net/url"
	"regexp"
	"strings"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/opencontainers/go-digest"
	"oras.land/oras-go/v2/errdef"
)

// Protocol represent artifact transport protocol
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
type Protocol string

var (
	// ProtocolDocker docker protocol
	//
	// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
	ProtocolDocker Protocol = "docker"
	// ProtocolHelmChart helm chart protocol
	//
	// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
	ProtocolHelmChart Protocol = "chart"
	// ProtocolOCI oci protocol
	//
	// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
	ProtocolOCI Protocol = "oci"

	// repositoryRegexp is adapted from the distribution implementation. The
	// repository name set under OCI distribution spec is a subset of the docker
	// spec. For maximum compatibility, the docker spec is verified client-side.
	// Further checks are left to the server-side.
	// References:
	// - https://github.com/distribution/distribution/blob/v2.7.1/reference/regexp.go#L53
	// - https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pulling-manifests
	//
	// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
	repositoryRegexp = regexp.MustCompile(`^[a-z0-9]+(?:(?:[._]|__|[-]*)[a-z0-9]+)*(?:/[a-z0-9]+(?:(?:[._]|__|[-]*)[a-z0-9]+)*)*$`)

	// tagRegexp checks the tag name.
	// The docker and OCI spec have the same regular expression.
	// Reference: https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pulling-manifests
	//
	// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
	tagRegexp = regexp.MustCompile(`^[\w][\w.-]{0,127}$`)

	//
	// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
	ErrInvalidReference = "invalid reference"
)

// DigestAlgorithm artifact digest algorithm
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
type DigestAlgorithm string

const (
	// SHA256 artifact digest algorithm is sha256
	//
	// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
	SHA256 DigestAlgorithm = "sha256"
)

// URI represents artifact uri , like docker://gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/nop:v0.37.1@sha256:04411f239bc7144c3248b53af6741c2726eaddbe9b9cf62a24cf812689cc3223
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
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
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func (u URI) DigestString() string {
	if u.Digest == "" {
		return ""
	}
	return fmt.Sprintf("%s:%s", u.Algorithm, u.Digest)
}

// Version will return sha256:xx or tag .
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func (u URI) Version() string {
	if u.Digest != "" {
		return u.DigestString()
	} else {
		return u.Tag
	}
}

// String returns the uri string, preferably using tag as a reference.
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func (u URI) String() string {
	ref := u.Host
	if u.Path != "" {
		ref = fmt.Sprintf("%s/%s", strings.Trim(ref, "/"), strings.Trim(u.Path, "/"))
	}

	if u.Tag != "" {
		return fmt.Sprintf("%s:%s", ref, u.Tag)
	}

	if u.Digest != "" {
		return fmt.Sprintf("%s@%s", ref, u.DigestString())
	}

	return ref
}

// WithDigestString return uri string, including tag and digest.
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func (u URI) WithDigestString() string {
	ref := u.Host
	if u.Path != "" {
		ref = fmt.Sprintf("%s/%s", strings.Trim(ref, "/"), strings.Trim(u.Path, "/"))
	}

	if u.Tag != "" {
		ref = fmt.Sprintf("%s:%s", ref, u.Tag)
	}

	if u.Digest != "" {
		ref = fmt.Sprintf("%s@%s", ref, u.DigestString())
	}

	return ref
}

// ParseURI parse uri to URI struct
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
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

// ValidateHost validates the host.
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func (u URI) ValidateHost() error {
	if _, err := url.ParseRequestURI("sample://" + u.Host); err != nil {
		return fmt.Errorf("%s: invalid registry %s", errdef.ErrInvalidReference, err.Error())
	}
	return nil
}

// ValidatePath validates the path.
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func (u URI) ValidatePath() error {
	if !repositoryRegexp.MatchString(u.Path) {
		return fmt.Errorf("%s: invalid repository %s", ErrInvalidReference, u.Path)
	}
	return nil
}

// ValidateTag validates the tag.
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func (u URI) ValidateTag() error {
	if !tagRegexp.MatchString(u.Tag) {
		return fmt.Errorf("%s: invalid tag %s", ErrInvalidReference, u.Tag)
	}
	return nil
}

// ValidateDigest validates the a digest.
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func (u URI) ValidateDigest() error {
	if _, err := digest.Parse(u.DigestString()); err != nil {
		return fmt.Errorf("%s: invalid digest %s; %v", ErrInvalidReference, u.DigestString(), err)
	}
	return nil
}

// Validate returns an error if the uri is invalid, otherwise returns empty.
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func (u URI) Validate() error {
	if err := u.ValidateHost(); err != nil {
		return err
	}

	if err := u.ValidatePath(); err != nil {
		return err
	}

	if err := u.ValidateTag(); err != nil && u.Tag != "" {
		return err
	}

	if err := u.ValidateDigest(); err != nil && u.DigestString() != "" {
		return err
	}

	return nil
}
