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
	"reflect"
	"testing"

	"github.com/katanomi/pkg/apis/meta/v1alpha1"
)

func TestParseURI(t *testing.T) {
	var data = []struct {
		desc     string
		uri      string
		t        v1alpha1.ArtifactType
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
			t:    v1alpha1.OCIHelmChartArtifactParameterType,
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
