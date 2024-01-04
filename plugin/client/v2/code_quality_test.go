/*
Copyright 2023 The Katanomi Authors.

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

package v2

import (
	"context"

	"github.com/jarcoal/httpmock"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test GetCodeQuality", func() {
	It("should generate the correct url and got expected response", func() {
		expected := fakeStruct[metav1alpha1.CodeQuality]()
		httpmock.RegisterResponder(
			"GET",
			"https://example.com/codeQuality/projectKey",
			httpmock.NewJsonResponderOrPanic(200, expected),
		)

		got, err := pluginClient.GetCodeQuality(context.Background(), "projectKey")
		Expect(err).To(Succeed())
		Expect(diff(got, expected)).To(BeEmpty())
	})
})

var _ = Describe("Test GetCodeQualityOverviewByBranch", func() {
	It("should generate the correct url and got expected response", func() {
		expected := fakeStruct[metav1alpha1.CodeQuality]()
		httpmock.RegisterResponder(
			"GET",
			"https://example.com/codeQuality/projectKey/branches/branchKey",
			httpmock.NewJsonResponderOrPanic(200, expected),
		)

		got, err := pluginClient.GetCodeQualityOverviewByBranch(context.Background(), metav1alpha1.CodeQualityBaseOption{ProjectKey: "projectKey", BranchKey: "branchKey"})
		Expect(err).To(Succeed())
		Expect(diff(got, expected)).To(BeEmpty())
	})
})

var _ = Describe("Test GetCodeQualityLineCharts", func() {
	It("should generate the correct url and got expected response", func() {
		expected := fakeStruct[metav1alpha1.CodeQualityLineChart]()
		httpmock.RegisterResponder(
			"GET",
			"https://example.com/codeQuality/projectKey/branches/branchKey/lineCharts",
			httpmock.NewJsonResponderOrPanic(200, expected),
		)
		option := metav1alpha1.CodeQualityLineChartOption{Metrics: "metrics"}
		option.ProjectKey = "projectKey"
		option.BranchKey = "branchKey"
		option.Metrics = "metrics"
		got, err := pluginClient.GetCodeQualityLineCharts(context.Background(), option)
		Expect(err).To(Succeed())
		Expect(diff(got, expected)).To(BeEmpty())
	})
})

var _ = Describe("Test GetOverview", func() {
	It("should generate the correct url and got expected response", func() {
		expected := fakeStruct[metav1alpha1.CodeQualityProjectOverview]()
		httpmock.RegisterResponder(
			"GET",
			"https://example.com/codeQuality",
			httpmock.NewJsonResponderOrPanic(200, expected),
		)

		got, err := pluginClient.GetOverview(context.Background())
		Expect(err).To(Succeed())
		Expect(diff(got, expected)).To(BeEmpty())
	})
})

var _ = Describe("Test GetSummaryByTaskID", func() {
	It("should generate the correct url and got expected response", func() {
		expected := fakeStruct[metav1alpha1.CodeQualityTaskMetrics]()
		httpmock.RegisterResponder(
			"GET",
			"https://example.com/codeQuality/task/taskID/summary",
			httpmock.NewJsonResponderOrPanic(200, expected),
		)

		got, err := pluginClient.GetSummaryByTaskID(context.Background(), metav1alpha1.CodeQualityTaskOption{TaskID: "taskID", ProjectKey: "projectKey", Branch: "branch", PullRequest: "pullRequest"})
		Expect(err).To(Succeed())
		Expect(diff(got, expected)).To(BeEmpty())
	})
})
