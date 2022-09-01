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
			result := ArtifactVersion{}.GetBinaryObjectFromValues(values.ctx, values.values)

			if values.expected == nil {
				g.Expect(result).To(gomega.BeNil())
			} else {
				diff := cmp.Diff(*values.expected, result)
				g.Expect(diff).To(gomega.BeEmpty())
			}

		})
	}
}
