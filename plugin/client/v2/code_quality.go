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
	"fmt"
	"time"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client/base"
)

// GetCodeQuality get code quality data of a specific project
func (p *PluginClientV2) GetCodeQuality(ctx context.Context, projectKey string) (*metav1alpha1.CodeQuality, error) {
	codeQualityResult := &metav1alpha1.CodeQuality{}

	options := []base.OptionFunc{base.ResultOpts(codeQualityResult)}
	uri := fmt.Sprintf("/codeQuality/%s", projectKey)
	if err := p.Get(ctx, p.ClassAddress, uri, options...); err != nil {
		return nil, err
	}
	return codeQualityResult, nil
}

// GetCodeQualityOverviewByBranch get code quality data of a specific branch
func (p *PluginClientV2) GetCodeQualityOverviewByBranch(ctx context.Context, opt metav1alpha1.CodeQualityBaseOption) (*metav1alpha1.CodeQuality, error) {
	codeQualityResult := &metav1alpha1.CodeQuality{}
	options := []base.OptionFunc{base.ResultOpts(codeQualityResult)}
	uri := fmt.Sprintf("/codeQuality/%s/branches/%s", opt.ProjectKey, opt.BranchKey)
	if err := p.Get(ctx, p.ClassAddress, uri, options...); err != nil {
		return nil, err
	}
	return codeQualityResult, nil
}

// GetCodeQualityLineCharts get code quality data used for line charts
func (p *PluginClientV2) GetCodeQualityLineCharts(ctx context.Context, opt metav1alpha1.CodeQualityLineChartOption) (*metav1alpha1.CodeQualityLineChart, error) {
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
	options := []base.OptionFunc{base.QueryOpts(query), base.ResultOpts(lineChartResult)}
	uri := fmt.Sprintf("/codeQuality/%s/branches/%s/lineCharts", opt.ProjectKey, opt.BranchKey)
	if err := p.Get(ctx, p.ClassAddress, uri, options...); err != nil {
		return nil, err
	}
	return lineChartResult, nil
}

// GetOverview get code quality overview data
func (p *PluginClientV2) GetOverview(ctx context.Context) (*metav1alpha1.CodeQualityProjectOverview, error) {
	overview := &metav1alpha1.CodeQualityProjectOverview{}
	options := []base.OptionFunc{base.ResultOpts(overview)}
	if err := p.Get(ctx, p.ClassAddress, "/codeQuality", options...); err != nil {
		return nil, err
	}
	return overview, nil
}

// GetSummaryByTaskID get code quality summary data by task id
func (p *PluginClientV2) GetSummaryByTaskID(ctx context.Context, opt metav1alpha1.CodeQualityTaskOption) (*metav1alpha1.CodeQualityTaskMetrics, error) {
	taskMetrics := &metav1alpha1.CodeQualityTaskMetrics{}
	query := map[string]string{
		"project-key": opt.ProjectKey,
		"branch":      opt.Branch,
		"pullRequest": opt.PullRequest,
	}
	options := []base.OptionFunc{base.ResultOpts(taskMetrics), base.QueryOpts(query)}
	if err := p.Get(ctx, p.ClassAddress, fmt.Sprintf("/codeQuality/task/%s/summary", opt.TaskID), options...); err != nil {
		return nil, err
	}
	return taskMetrics, nil
}
