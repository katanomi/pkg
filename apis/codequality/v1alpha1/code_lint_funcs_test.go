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

	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestCodeLintResultGetObjectWithValues(t *testing.T) {
	table := map[string]struct {
		ctx    context.Context
		path   *field.Path
		values map[string]string

		expected *CodeLintResult
	}{
		"full values with prefix": {
			context.Background(),
			field.NewPath("value"),
			map[string]string{
				"value.result":       "Failed",
				"value.issues.count": "100",
			},
			MustLoadReturnObjectFromYAML("testdata/CodeLintResult.GetObjectWthValues.full.golden.yaml", &CodeLintResult{}).(*CodeLintResult),
		},
		"full values without prefix": {
			context.Background(),
			nil,
			map[string]string{
				"result":       "Failed",
				"issues.count": "100",
			},
			MustLoadReturnObjectFromYAML("testdata/CodeLintResult.GetObjectWthValues.full.golden.yaml", &CodeLintResult{}).(*CodeLintResult),
		},
		"nil values": {
			context.Background(),
			field.NewPath("value"),
			nil,
			nil,
		},
	}

	for test, values := range table {
		t.Run(test, func(t *testing.T) {
			g := gomega.NewGomegaWithT(t)
			result := CodeLintResult{}.GetObjectWithValues(values.ctx, values.path, values.values)

			diff := cmp.Diff(values.expected, result)
			g.Expect(diff).To(gomega.BeEmpty())
		})
	}
}

func TestCodeLintResultIsEmpty(t *testing.T) {
	t.Run("is empty struct", func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)
		object := CodeLintResult{}

		g.Expect(object.IsEmpty()).To(gomega.BeTrue())
	})

	t.Run("has issues but empty/zero", func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)
		object := CodeLintResult{Issues: &CodeLintIssues{}}

		g.Expect(object.IsEmpty()).To(gomega.BeFalse())
	})

	t.Run("has any attribute with nil issues", func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)
		object := CodeLintResult{Result: "a"}

		g.Expect(object.IsEmpty()).To(gomega.BeFalse())
	})
}
