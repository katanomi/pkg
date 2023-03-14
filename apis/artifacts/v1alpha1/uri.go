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
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/opencontainers/go-digest"
	"oras.land/oras-go/v2/errdef"
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

	// repositoryRegexp is adapted from the distribution implementation. The
	// repository name set under OCI distribution spec is a subset of the docker
	// spec. For maximum compatibility, the docker spec is verified client-side.
	// Further checks are left to the server-side.
	// References:
	// - https://github.com/distribution/distribution/blob/v2.7.1/reference/regexp.go#L53
	// - https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pulling-manifests
	repositoryRegexp = regexp.MustCompile(`^[a-z0-9]+(?:(?:[._]|__|[-]*)[a-z0-9]+)*(?:/[a-z0-9]+(?:(?:[._]|__|[-]*)[a-z0-9]+)*)*$`)

	// tagRegexp checks the tag name.
	// The docker and OCI spec have the same regular expression.
	// Reference: https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pulling-manifests
	tagRegexp = regexp.MustCompile(`^[\w][\w.-]{0,127}$`)

	digestRegexp = regexp.MustCompile(`(?P<algo>sha256):(?P<hash>[a-z0-9]{64})`)

	ErrInvalidReference = "invalid reference"
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

// String returns the uri string, preferably using tag as a reference.
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

func (u URI) Repository() string {
	ref := u.Host
	if u.Path != "" {
		ref = fmt.Sprintf("%s/%s", strings.Trim(ref, "/"), strings.Trim(u.Path, "/"))
	}

	return ref
}

// ParseURI parse uri to URI struct
func ParseURI(uri string, t ArtifactType) (URI, error) {
	var u = URI{Raw: uri}

	str := uri

	protocolIndex := strings.Index(uri, "://")
	if protocolIndex <= 0 {
		u.Protocol = string(ProtocolDocker)
	} else {
		u.Protocol = uri[0:protocolIndex]
		str = str[protocolIndex+len("://"):]
	}

	switch t {
	case ArtifactTypeHelmChart:
		u.Protocol = string(ProtocolHelmChart)
	case ArtifactTypeContainerImage:
		u.Protocol = string(ProtocolDocker)
	default:
		// no-op
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
func (u URI) ValidateHost() error {
	if _, err := url.ParseRequestURI("sample://" + u.Host); err != nil {
		return fmt.Errorf("%s: invalid registry %s", errdef.ErrInvalidReference, err.Error())
	}
	return nil
}

// ValidatePath validates the path.
func (u URI) ValidatePath() error {
	if !repositoryRegexp.MatchString(u.Path) {
		return fmt.Errorf("%s: invalid repository %s", ErrInvalidReference, u.Path)
	}
	return nil
}

// ValidateTag validates the tag.
func (u URI) ValidateTag() error {
	if !tagRegexp.MatchString(u.Tag) {
		return fmt.Errorf("%s: invalid tag %s", ErrInvalidReference, u.Tag)
	}
	return nil
}

// ValidateDigest validates the a digest.
func (u URI) ValidateDigest() error {
	if _, err := digest.Parse(u.DigestString()); err != nil {
		return fmt.Errorf("%s: invalid digest %s; %v", ErrInvalidReference, u.DigestString(), err)
	}
	return nil
}

// Validate returns an error if the uri is invalid, otherwise returns empty.
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

// ParseFirstDigest parse first digest from text string
func ParseFirstDigest(line string) (algorithm DigestAlgorithm, digest string, ok bool) {
	groupNames := digestRegexp.SubexpNames()
	hash, algo := "", ""

	for _, match := range digestRegexp.FindAllStringSubmatch(line, -1) {
		for groupIdx, group := range match {
			name := groupNames[groupIdx]
			switch name {
			case "hash":
				hash = group
			case "algo":
				algo = group
			default:
				// no-op
			}
		}
		if hash != "" && algo != "" {
			digest = hash
			algorithm = DigestAlgorithm(algo)
			ok = true
			break
		}
	}
	return
}

// AsDigestStringArray returns list of uri as a string array
func AsDigestStringArray(uris ...URI) (array []string) {
	array = make([]string, 0, len(uris))
	for _, uri := range uris {
		array = append(array, uri.WithDigestString())
	}
	return
}
