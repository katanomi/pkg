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
	"testing"

	"github.com/google/go-cmp/cmp"
	. "github.com/katanomi/pkg/testing"
	"github.com/onsi/gomega"
)

func TestArtifactVersionGetBinaryObjectFromValues(t *testing.T) {
	table := map[string]struct {
		ctx    context.Context
		values []string

		expected *[]ArtifactVersion
	}{
		"full values with prefix": {
			context.Background(),
			[]string{"https://binary.katanomi.dev/abc/def.jpg", "https://binary.katanomi.dev/abc/anotherfolder"},
			MustLoadReturnObjectFromYAML("testdata/ArtifactVersion.GetBinaryObjectFromValues.golden.yaml", &[]ArtifactVersion{}).(*[]ArtifactVersion),
		},
		"nil values": {
			context.Background(),
			nil,
			nil,
		},
	}

	for test, values := range table {
		t.Run(test, func(t *testing.T) {
			g := gomega.NewGomegaWithT(t)
			result := GetBinaryObjectFromValues(values.ctx, values.values)

			if values.expected == nil {
				g.Expect(result).To(gomega.BeNil())
			} else {
				diff := cmp.Diff(*values.expected, result)
				g.Expect(diff).To(gomega.BeEmpty())
			}

		})
	}
}

func TestArtifactVersionGetHelmChartObjectFromURLValues(t *testing.T) {
	table := map[string]struct {
		ctx  context.Context
		url  string
		tags []string

		expected *[]ArtifactVersion
	}{
		"full values with multiple versions": {
			context.Background(),
			"registry.katanomi.dev/abc/def",
			[]string{"latest", "v1.0.1"},
			MustLoadReturnObjectFromYAML("testdata/ArtifactVersion.GetHelmChartObjectFromURLValues.golden.yaml", &[]ArtifactVersion{}).(*[]ArtifactVersion),
		},
		"empty values": {
			context.Background(),
			"   ",
			nil,
			nil,
		},
	}

	for test, values := range table {
		t.Run(test, func(t *testing.T) {
			g := gomega.NewGomegaWithT(t)
			result := GetHelmChartObjectFromURLValues(values.ctx, values.url, values.tags...)

			if values.expected == nil {
				g.Expect(result).To(gomega.BeNil())
			} else {
				diff := cmp.Diff(*values.expected, result)
				g.Expect(diff).To(gomega.BeEmpty())
			}

		})
	}
}

func TestArtifactVersionGetContainerImageObjectFromURLValues(t *testing.T) {
	table := map[string]struct {
		ctx    context.Context
		url    string
		digest string
		tags   []string

		expected *[]ArtifactVersion
	}{
		"full values with multiple versions": {
			context.Background(),
			"registry.katanomi.dev/abc/def",
			"sha256:b51bf41c3f3c64c818a8649156847d6b285d0c19f4bf0490378229b02ce5dafe",
			[]string{"latest", "v1.0.1"},
			MustLoadReturnObjectFromYAML("testdata/ArtifactVersion.GetContainerImageObjectFromURLValues.golden.yaml", &[]ArtifactVersion{}).(*[]ArtifactVersion),
		},
		"empty values": {
			context.Background(),
			"   ",
			"",
			nil,
			nil,
		},
	}

	for test, values := range table {
		t.Run(test, func(t *testing.T) {
			g := gomega.NewGomegaWithT(t)
			result := GetContainerImageObjectFromURLValues(values.ctx, values.url, values.digest, values.tags...)

			if values.expected == nil {
				g.Expect(result).To(gomega.BeNil())
			} else {
				diff := cmp.Diff(*values.expected, result)
				g.Expect(diff).To(gomega.BeEmpty())
			}

		})
	}
}

func TestArtifactVersionGetContainerImageFromValues(t *testing.T) {
	table := map[string]struct {
		ctx    context.Context
		values []string

		expected *[]ArtifactVersion
	}{
		"multiple values with tag, digest, digest and tag": {
			context.Background(),
			[]string{
				// the same artifact
				"registry.katanomi.dev/abc/def:latest",
				"registry.katanomi.dev/abc/def@sha256:b51bf41c3f3c64c818a8649156847d6b285d0c19f4bf0490378229b02ce5dafe",

				// another artifact with digest and two tags
				"registry.katanomi.dev/project/repository:v1@sha256:b51bf41c3f3c64c818a8649156847d6b285d0c19f4bf0490378229b02ce5dzzz",
				"registry.katanomi.dev/project/repository:v2",
				"registry.katanomi.dev/project/repository:v1.1@sha256:b51bf41c3f3c64c818a8649156847d6b285d0c19f4bf0490378229b02ce5dzzz",

				// ip address with port
				"192.168.1.1:32001/repo/test:tag",
			},
			MustLoadReturnObjectFromYAML("testdata/ArtifactVersion.GetContainerImageFromValues.golden.yaml", &[]ArtifactVersion{}).(*[]ArtifactVersion),
		},
		"empty values": {
			context.Background(),
			nil,
			nil,
		},
	}

	for test, values := range table {
		t.Run(test, func(t *testing.T) {
			g := gomega.NewGomegaWithT(t)
			result := GetContainerImageFromValues(values.ctx, values.values)

			if values.expected == nil {
				g.Expect(result).To(gomega.BeNil())
			} else {
				diff := cmp.Diff(*values.expected, result)
				g.Expect(diff).To(gomega.BeEmpty())
			}

		})
	}
}
