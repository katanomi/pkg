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

func TestAnalysisResultGetObjectWithValues(t *testing.T) {
	table := map[string]struct {
		ctx    context.Context
		path   *field.Path
		values map[string]string

		expected *AnalysisResult
	}{
		"full values with prefix": {
			context.Background(),
			field.NewPath("value"),
			map[string]string{
				"value.reportURL":                               "http://sonar.katanomi.dev",
				"value.projectID":                               "github-com-katanomi-spec",
				"value.taskID":                                  "abc123",
				"value.result":                                  "Succeeded",
				"value.metrics.branch.coverage.new":             "56",
				"value.metrics.branch.coverage.total":           "30",
				"value.metrics.branch.duplications.new":         "0.34",
				"value.metrics.branch.duplications.total":       "43",
				"value.metrics.target.coverage.total":           "30",
				"value.metrics.target.coverage.new":             "76.2",
				"value.metrics.target.duplications.total":       "23",
				"value.metrics.target.duplications.new":         "1.43",
				"value.metrics.ratings.reliability.rate":        "A",
				"value.metrics.ratings.reliability.issues":      "123",
				"value.metrics.ratings.vulnerability.rate":      "D",
				"value.metrics.ratings.vulnerability.issues":    "321",
				"value.metrics.ratings.maintainability.rate":    "E",
				"value.metrics.ratings.maintainability.issues":  "512",
				"value.metrics.ratings.securityHotspots.rate":   "B",
				"value.metrics.ratings.securityHotspots.issues": "2",
				"value.metrics.languages":                       "java,go",
				"value.metrics.codeSize.linesOfCode":            "1234",
			},
			MustLoadReturnObjectFromYAML("testdata/AnalysisResult.GetObjectWithValues.full.golden.yaml", &AnalysisResult{}).(*AnalysisResult),
		},
		"full values without prefix": {
			context.Background(),
			nil,
			map[string]string{
				"reportURL":                               "http://sonar.katanomi.dev",
				"projectID":                               "github-com-katanomi-spec",
				"taskID":                                  "abc123",
				"result":                                  "Succeeded",
				"metrics.branch.coverage.new":             "56",
				"metrics.branch.coverage.total":           "30",
				"metrics.branch.duplications.new":         "0.34",
				"metrics.branch.duplications.total":       "43",
				"metrics.target.coverage.total":           "30",
				"metrics.target.coverage.new":             "76.2",
				"metrics.target.duplications.total":       "23",
				"metrics.target.duplications.new":         "1.43",
				"metrics.ratings.reliability.rate":        "A",
				"metrics.ratings.reliability.issues":      "123",
				"metrics.ratings.vulnerability.rate":      "D",
				"metrics.ratings.vulnerability.issues":    "321",
				"metrics.ratings.maintainability.rate":    "E",
				"metrics.ratings.maintainability.issues":  "512",
				"metrics.ratings.securityHotspots.rate":   "B",
				"metrics.ratings.securityHotspots.issues": "2",
				"metrics.languages":                       "java,go",
				"metrics.codeSize.linesOfCode":            "1234",
			},
			MustLoadReturnObjectFromYAML("testdata/AnalysisResult.GetObjectWithValues.full.golden.yaml", &AnalysisResult{}).(*AnalysisResult),
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
			result := AnalysisResult{}.GetObjectWithValues(values.ctx, values.path, values.values)

			diff := cmp.Diff(values.expected, result)
			g.Expect(diff).To(gomega.BeEmpty())
		})
	}
}

func TestAnalisysMetricsGetObjectWithValues(t *testing.T) {
	table := map[string]struct {
		ctx    context.Context
		path   *field.Path
		values map[string]string

		expected *AnalisysMetrics
	}{
		"full values without prefix": {
			context.Background(),
			nil,
			map[string]string{
				"branch.coverage.new":             "56",
				"branch.coverage.total":           "30",
				"branch.duplications.new":         "0.34",
				"branch.duplications.total":       "43",
				"target.coverage.total":           "30",
				"target.coverage.new":             "76.2",
				"target.duplications.total":       "23",
				"target.duplications.new":         "1.43",
				"ratings.reliability.rate":        "A",
				"ratings.reliability.issues":      "123",
				"ratings.vulnerability.rate":      "D",
				"ratings.vulnerability.issues":    "321",
				"ratings.maintainability.rate":    "E",
				"ratings.maintainability.issues":  "512",
				"ratings.securityHotspots.rate":   "B",
				"ratings.securityHotspots.issues": "2",
				"languages":                       "java,go",
				"codeSize.linesOfCode":            "1234",
			},
			MustLoadReturnObjectFromYAML("testdata/AnalisysMetrics.GetObjectWithValues.full.golden.yaml", &AnalisysMetrics{}).(*AnalisysMetrics),
		},
		"full values with prefix": {
			context.Background(),
			field.NewPath("abc", "def"),
			map[string]string{
				"abc.def.branch.coverage.new":             "56",
				"abc.def.branch.coverage.total":           "30",
				"abc.def.branch.duplications.new":         "0.34",
				"abc.def.branch.duplications.total":       "43",
				"abc.def.target.coverage.total":           "30",
				"abc.def.target.coverage.new":             "76.2",
				"abc.def.target.duplications.total":       "23",
				"abc.def.target.duplications.new":         "1.43",
				"abc.def.ratings.reliability.rate":        "A",
				"abc.def.ratings.reliability.issues":      "123",
				"abc.def.ratings.vulnerability.rate":      "D",
				"abc.def.ratings.vulnerability.issues":    "321",
				"abc.def.ratings.maintainability.rate":    "E",
				"abc.def.ratings.maintainability.issues":  "512",
				"abc.def.ratings.securityHotspots.rate":   "B",
				"abc.def.ratings.securityHotspots.issues": "2",
				"abc.def.languages":                       "java,go",
				"abc.def.codeSize.linesOfCode":            "1234",
			},
			MustLoadReturnObjectFromYAML("testdata/AnalisysMetrics.GetObjectWithValues.full.golden.yaml", &AnalisysMetrics{}).(*AnalisysMetrics),
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
			result := AnalisysMetrics{}.GetObjectWithValues(values.ctx, values.path, values.values)

			diff := cmp.Diff(values.expected, result)
			g.Expect(diff).To(gomega.BeEmpty())
		})
	}
}
