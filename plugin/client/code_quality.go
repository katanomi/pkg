/*
Copyright 2021 The Katanomi Authors.

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

package client

import (
	"context"
	"fmt"
	"time"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

type ClientCodeQuality interface {
	Get(ctx context.Context, baseURL *duckv1.Addressable, projectKey string, options ...OptionFunc) (*metav1alpha1.CodeQuality, error)
	GetOverview(ctx context.Context, baseURL *duckv1.Addressable, options ...OptionFunc) (*metav1alpha1.CodeQualityProjectOverview, error)
	GetByBranch(ctx context.Context, baseURL *duckv1.Addressable, opt metav1alpha1.CodeQualityBaseOption, options ...OptionFunc) (*metav1alpha1.CodeQuality, error)
	GetLineCharts(ctx context.Context, baseURL *duckv1.Addressable, opt metav1alpha1.CodeQualityLineChartOption, options ...OptionFunc) (*metav1alpha1.CodeQualityLineChart, error)
	GetMetricsByTaskID(ctx context.Context, baseURL *duckv1.Addressable, opt metav1alpha1.CodeQualityTaskOption, options ...OptionFunc) (*metav1alpha1.CodeQualityTaskMetrics, error)
}

type codeQuality struct {
	client Client
}

func newCodeQuality(client Client) ClientCodeQuality {
	return &codeQuality{
		client: client,
	}
}

func (c *codeQuality) Get(ctx context.Context, baseURL *duckv1.Addressable, projectKey string, options ...OptionFunc) (*metav1alpha1.CodeQuality, error) {
	codeQualityResult := &metav1alpha1.CodeQuality{}
	options = append(options, ResultOpts(codeQualityResult))
	uri := fmt.Sprintf("/codeQuality/%s", projectKey)
	if err := c.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}
	return codeQualityResult, nil
}

func (c *codeQuality) GetByBranch(ctx context.Context, baseURL *duckv1.Addressable, opt metav1alpha1.CodeQualityBaseOption, options ...OptionFunc) (*metav1alpha1.CodeQuality, error) {
	codeQualityResult := &metav1alpha1.CodeQuality{}
	options = append(options, ResultOpts(codeQualityResult))
	uri := fmt.Sprintf("/codeQuality/%s/branches/%s", opt.ProjectKey, opt.BranchKey)
	if err := c.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}
	return codeQualityResult, nil
}

func (c *codeQuality) GetLineCharts(ctx context.Context, baseURL *duckv1.Addressable, opt metav1alpha1.CodeQualityLineChartOption, options ...OptionFunc) (*metav1alpha1.CodeQualityLineChart, error) {
	lineChartResult := &metav1alpha1.CodeQualityLineChart{}
	query := map[string]string{
		"metricKeys": opt.Metrics,
	}
	if opt.StartTime != nil {
		query["startTime"] = opt.StartTime.Format(time.RFC3339)
	}
	if opt.CompletionTime != nil {
		query["completionTime"] = opt.CompletionTime.Format(time.RFC3339)
	}
	options = append(options, QueryOpts(query), ResultOpts(lineChartResult))
	uri := fmt.Sprintf("/codeQuality/%s/branches/%s/lineCharts", opt.ProjectKey, opt.BranchKey)
	if err := c.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}
	return lineChartResult, nil
}

func (c *codeQuality) GetOverview(ctx context.Context, baseURL *duckv1.Addressable, options ...OptionFunc) (*metav1alpha1.CodeQualityProjectOverview, error) {
	overview := &metav1alpha1.CodeQualityProjectOverview{}
	options = append(options, ResultOpts(overview))
	if err := c.client.Get(ctx, baseURL, "/codeQuality", options...); err != nil {
		return nil, err
	}
	return overview, nil
}

func (c *codeQuality) GetMetricsByTaskID(ctx context.Context, baseURL *duckv1.Addressable, opt metav1alpha1.CodeQualityTaskOption, options ...OptionFunc) (*metav1alpha1.CodeQualityTaskMetrics, error) {
	taskMetrics := &metav1alpha1.CodeQualityTaskMetrics{}
	query := map[string]string{
		"project-key": opt.ProjectKey,
		"branch":      opt.Branch,
		"pullRequest": opt.PullRequest,
	}
	options = append(options, ResultOpts(taskMetrics), QueryOpts(query))
	if err := c.client.Get(ctx, baseURL, fmt.Sprintf("/codeQuality/task/%s/summary", opt.TaskID), options...); err != nil {
		return nil, err
	}
	return taskMetrics, nil
}
