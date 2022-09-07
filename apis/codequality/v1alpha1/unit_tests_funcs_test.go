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

func TestUnitTestsResultGetObjectWithValues(t *testing.T) {
	table := map[string]struct {
		ctx    context.Context
		path   *field.Path
		values map[string]string

		expected *UnitTestsResult
	}{
		"full values with prefix": {
			context.Background(),
			field.NewPath("value"),
			map[string]string{
				"value.coverage.lines":             "23%",
				"value.coverage.branches":          "45%",
				"value.testResults.passed":         "1",
				"value.testResults.failed":         "2",
				"value.testResults.skipped":        "3",
				"value.testResults.passedTestRate": "100%",
			},
			MustLoadReturnObjectFromYAML("testdata/UnitTestsResult.GetObjectWithValues.full.golden.yaml", &UnitTestsResult{}).(*UnitTestsResult),
		},
		"full values without prefix": {
			context.Background(),
			nil,
			map[string]string{
				"coverage.lines":             "23%",
				"coverage.branches":          "45%",
				"testResults.passed":         "1",
				"testResults.failed":         "2",
				"testResults.skipped":        "3",
				"testResults.passedTestRate": "100%",
			},
			MustLoadReturnObjectFromYAML("testdata/UnitTestsResult.GetObjectWithValues.full.golden.yaml", &UnitTestsResult{}).(*UnitTestsResult),
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
			result := UnitTestsResult{}.GetObjectWithValues(values.ctx, values.path, values.values)

			diff := cmp.Diff(values.expected, result)
			g.Expect(diff).To(gomega.BeEmpty())
		})
	}
}

func TestUnitTestsResultIsEmpty(t *testing.T) {
	t.Run("is empty struct", func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)
		object := UnitTestsResult{}

		g.Expect(object.IsEmpty()).To(gomega.BeTrue())
	})

	t.Run("has non nil coverage but empty", func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)
		object := UnitTestsResult{Coverage: &TestCoverage{}}

		g.Expect(object.IsEmpty()).To(gomega.BeTrue())
	})

	t.Run("has non nil results with empty values nil coverage", func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)
		object := UnitTestsResult{TestResult: &TestResult{}}

		g.Expect(object.IsEmpty()).To(gomega.BeTrue())
	})

	t.Run("has results with nil coverage", func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)
		object := UnitTestsResult{TestResult: &TestResult{Passed: 1}}

		g.Expect(object.IsEmpty()).To(gomega.BeFalse())
	})

	t.Run("nil results with coverage", func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)
		object := UnitTestsResult{Coverage: &TestCoverage{Lines: "1"}}

		g.Expect(object.IsEmpty()).To(gomega.BeFalse())
	})
}
