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

package options

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/katanomi/pkg/report"
	pkgTesting "github.com/katanomi/pkg/testing"
	. "github.com/onsi/gomega"
)

func TestReportOption(t *testing.T) {
	g := NewGomegaWithT(t)
	ctx := context.Background()
	obj := struct {
		ReportPathsByTypesOption
	}{}

	args := []string{
		"--report-configs", "invalid-type=junit.xml",
	}
	err := RegisterSetup(&obj, ctx, nil, args)
	g.Expect(err).Should(HaveOccurred())

	args = []string{
		"--report-configs", "junit-xml=junit.xml",
	}
	err = RegisterSetup(&obj, ctx, nil, args)
	g.Expect(err).Should(Succeed())
	g.Expect(obj.ReportPathByTypes).To(Equal(map[report.ReportType]string{
		"junit-xml": "junit.xml",
	}))

}

func TestSummariesByType(t *testing.T) {
	g := NewGomegaWithT(t)
	reportPath := "./testdata/reports"

	obj := struct {
		ReportPathsByTypesOption
	}{}

	// empty types
	obj.ReportPathByTypes = map[report.ReportType]string{}
	g.Expect(obj.TestSummariesByType(reportPath)).To(BeEmpty())

	// valid type
	obj.ReportPathByTypes = map[report.ReportType]string{
		report.TypeJunitXml: "junit.xml",
	}
	var expectedResult report.SummariesByType
	pkgTesting.MustLoadJSON("./testdata/summary.byType.golden.json", &expectedResult)

	result, err := obj.TestSummariesByType(reportPath)
	g.Expect(err).ShouldNot(HaveOccurred())
	diff := cmp.Diff(result, expectedResult)
	g.Expect(diff).To(BeEmpty())
}
