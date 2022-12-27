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
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
)

func TestParseURI(t *testing.T) {
	var data = []struct {
		desc     string
		uri      string
		t        ArtifactType
		u        URI
		hasError bool
	}{
		{
			desc: "with no tags no digest",
			uri:  "docker://docker.io/katanomi/pkg",
			u: URI{
				Host:     "docker.io",
				Protocol: string(ProtocolDocker),
				Path:     "katanomi/pkg",
				Raw:      "docker://docker.io/katanomi/pkg",
			},
		},
		{
			desc: "with no digest",
			uri:  "docker://docker.io/katanomi/pkg:v1",
			u: URI{
				Host:     "docker.io",
				Protocol: string(ProtocolDocker),
				Path:     "katanomi/pkg",
				Tag:      "v1",
				Raw:      "docker://docker.io/katanomi/pkg:v1",
			},
		},
		{
			desc: "with tag and digest",
			uri:  "docker://docker.io/katanomi/pkg:v1@sha256:asdf",
			u: URI{
				Host:      "docker.io",
				Protocol:  string(ProtocolDocker),
				Path:      "katanomi/pkg",
				Tag:       "v1",
				Digest:    "asdf",
				Algorithm: "sha256",
				Raw:       "docker://docker.io/katanomi/pkg:v1@sha256:asdf",
			},
		},
		{
			desc: "with digest",
			uri:  "docker://docker.io/katanomi/pkg@sha256:asdf",
			u: URI{
				Host:      "docker.io",
				Protocol:  string(ProtocolDocker),
				Path:      "katanomi/pkg",
				Digest:    "asdf",
				Algorithm: "sha256",
				Raw:       "docker://docker.io/katanomi/pkg@sha256:asdf",
			},
		},
		{
			desc: "with type",
			uri:  "docker.io/katanomi/pkg@sha256:asdf",
			t:    ArtifactTypeHelmChart,
			u: URI{
				Host:      "docker.io",
				Protocol:  string(ProtocolHelmChart),
				Path:      "katanomi/pkg",
				Digest:    "asdf",
				Algorithm: "sha256",
				Raw:       "docker.io/katanomi/pkg@sha256:asdf",
			},
		},
	}

	for _, item := range data {
		t.Run(item.desc, func(t *testing.T) {
			actualURI, err := ParseURI(item.uri, item.t)
			if item.hasError == false && err != nil {
				t.Errorf("expected error should be nil")
			}

			if !reflect.DeepEqual(actualURI, item.u) {
				t.Errorf("actual uri: %#v, expected uri: %#v", actualURI, item.u)
			}
		})
	}
}

func Test_StringWithDigestString(t *testing.T) {
	cases := []struct {
		name       string
		uri        string
		wantString string
		wantDigest string
	}{
		{
			name:       "with tag and digest",
			uri:        "docker.io/katanomi/pkg:v1@sha256:asdf",
			wantString: "docker.io/katanomi/pkg:v1",
			wantDigest: "docker.io/katanomi/pkg:v1@sha256:asdf",
		},
		{
			name:       "with digest",
			uri:        "docker.io/katanomi/pkg@sha256:asdf",
			wantString: "docker.io/katanomi/pkg@sha256:asdf",
			wantDigest: "docker.io/katanomi/pkg@sha256:asdf",
		},
		{
			name:       "with tag",
			uri:        "docker.io/katanomi/pkg:v1",
			wantString: "docker.io/katanomi/pkg:v1",
			wantDigest: "docker.io/katanomi/pkg:v1",
		},
		{
			name:       "only path",
			uri:        "docker.io/katanomi/pkg",
			wantString: "docker.io/katanomi/pkg",
			wantDigest: "docker.io/katanomi/pkg",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			uri, err := ParseURI(c.uri, ArtifactTypeHelmChart)
			g.Expect(err).To(BeNil())
			g.Expect(uri.String()).To(Equal(c.wantString))
			g.Expect(uri.WithDigestString()).To(Equal(c.wantDigest))
		})
	}
}

func Test_Validate(t *testing.T) {

	cases := []struct {
		name     string
		uris     []string
		wantErrs []error
	}{
		{
			name: "success validate",
			uris: []string{
				"docker.io/katanomi/pkg",
				"docker.io/katanomi/pkg:v1",
				"docker.io/katanomi/pkg@sha256:744c8b3d4c8f5b30a1a78c5e3893c4d3f793919d1e14bcaee61028931e9f9929",
				"docker.io/katanomi/pkg:v1@sha256:744c8b3d4c8f5b30a1a78c5e3893c4d3f793919d1e14bcaee61028931e9f9929",
				"127.0.0.1/katanomi/pkg",
				"127.0.0.1:8080/katanomi/pkg",
			},
		},
		{
			name: "failed validate host",
			uris: []string{
				"doc KHal.io/katanomi/pkg",
				"doc中文Hal.io/katanomi/pkg",
				"doc$Hal.io/katanomi/pkg",
			},
			wantErrs: []error{
				fmt.Errorf("%s: invalid registry parse \"sample://doc KHal.io\": invalid character \" \" in host name", ErrInvalidReference),
			},
		},
		{
			name: "failed validate path",
			uris: []string{
				"docker.io/katano#mi/pkg",
				"docker.io/katano mi/pkg",
				"docker.io/katano中文mi/pkg",
			},
			wantErrs: []error{
				fmt.Errorf("%s: invalid repository katano#mi/pkg", ErrInvalidReference),
				fmt.Errorf("%s: invalid repository katano mi/pkg", ErrInvalidReference),
				fmt.Errorf("%s: invalid repository katano中文mi/pkg", ErrInvalidReference),
			},
		},
		{
			name: "failed validate tag",
			uris: []string{
				"docker.io/katanomi/pkg:v2&",
				"docker.io/katanomi/pkg:v2中文",
				"docker.io/katanomi/pkg: v2",
				"docker.io/katanomi/pkg:v2:v1",
			},
			wantErrs: []error{
				fmt.Errorf("%s: invalid tag v2&", ErrInvalidReference),
				fmt.Errorf("%s: invalid tag v2中文", ErrInvalidReference),
				fmt.Errorf("%s: invalid tag  v2", ErrInvalidReference),
				fmt.Errorf("%s: invalid tag v2:v1", ErrInvalidReference),
			},
		},
		{
			name: "failed validate digest",
			uris: []string{
				"docker.io/katanomi/pkg@sha256:fdsa",
			},
			wantErrs: []error{
				fmt.Errorf("%s: invalid digest sha256:fdsa; invalid checksum digest length", ErrInvalidReference),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			errs := []error{}
			for _, ref := range c.uris {
				uri, err := ParseURI(ref, ArtifactTypeHelmChart)
				g.Expect(err).To(BeNil())
				err = uri.Validate()
				if err != nil {
					errs = append(errs, err)
				}
			}

			if c.wantErrs == nil {
				c.wantErrs = []error{}
			}

			g.Expect(errs).To(Equal(c.wantErrs))
		})
	}
}

func TestParseFirstDigest(t *testing.T) {
	table := map[string]struct {
		Input  string
		Algo   DigestAlgorithm
		Digest string
		Ok     bool
	}{
		"invalid output": {
			Input:  "this is a fake output sha256:f27cbdf8b3e7325658 not full digest",
			Digest: "",
			Ok:     false,
		},
		"oras cli output": {
			Input:  "Digest: sha256:f27cbdf8b3e732565880a7c372cb4015a2afc98dd6d708887db99cf5036027bb",
			Algo:   SHA256,
			Digest: "f27cbdf8b3e732565880a7c372cb4015a2afc98dd6d708887db99cf5036027bb",
			Ok:     true,
		},
		"other complex output": {
			Input: ` some other output
				Pushed successfully with returning digest sha256:f27cbdf8b3e732565880a7c372cb4015a2afc98dd6d708887db99cf5036027zz
			`,
			Algo:   SHA256,
			Digest: "f27cbdf8b3e732565880a7c372cb4015a2afc98dd6d708887db99cf5036027zz",
			Ok:     true,
		},

		"multiple digests only returns the first": {
			Input: `first digest sha256:f27cbdf8b3e732565880a7c372cb4015a2afc98dd6d708887db99cf5036027zz
			second digest Digest: sha256:f27cbdf8b3e732565880a7c372cb4015a2afc98dd6d708887db99cf5036027bb`,
			Algo:   SHA256,
			Digest: "f27cbdf8b3e732565880a7c372cb4015a2afc98dd6d708887db99cf5036027zz",
			Ok:     true,
		},
	}
	for name, c := range table {
		t.Run(name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			algo, result, okResult := ParseFirstDigest(c.Input)

			g.Expect(result).To(Equal(c.Digest))
			g.Expect(algo).To(Equal(c.Algo))
			g.Expect(okResult).To(Equal(c.Ok))
		})
	}
}

func TestAsDigestStringArray(t *testing.T) {
	table := map[string]struct {
		Input  []URI
		Result []string
	}{
		"empty input": {
			Input:  []URI{},
			Result: []string{},
		},
		"multiple uris with and without digest": {
			Input: []URI{
				{Host: "docker.io", Path: "/katanomi/repo", Tag: "latest", Algorithm: SHA256, Digest: "0123456789012345678901234567890123456789012345678901234567890123"},
				{Host: "index.docker.com", Path: "/someproject/somerepo", Tag: "v1.1.1"},
			},
			Result: []string{
				"docker.io/katanomi/repo:latest@sha256:0123456789012345678901234567890123456789012345678901234567890123",
				"index.docker.com/someproject/somerepo:v1.1.1",
			},
		},
	}
	for name, c := range table {
		t.Run(name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			result := AsDigestStringArray(c.Input...)

			g.Expect(result).To(Equal(c.Result))
		})
	}
}
