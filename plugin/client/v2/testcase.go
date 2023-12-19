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

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client/base"
)

// ListTestCases list test cases
func (p *PluginClientV2) ListTestCases(ctx context.Context, params metav1alpha1.TestProjectOptions, option metav1alpha1.ListOptions) (*metav1alpha1.TestCaseList, error) {
	list := &metav1alpha1.TestCaseList{}

	options := []base.OptionFunc{
		base.ResultOpts(list),
		base.ListOpts(option),
		base.QueryOpts(map[string]string{
			"buildID": params.BuildID,
		}),
	}
	uri := fmt.Sprintf("projects/%s/testplans/%s/testcases", params.Project, params.TestPlanID)
	if err := p.Get(ctx, p.ClassAddress, uri, options...); err != nil {
		return nil, err
	}

	return list, nil
}

// GetTestCase get test case
func (p *PluginClientV2) GetTestCase(ctx context.Context, params metav1alpha1.TestProjectOptions) (*metav1alpha1.TestCase, error) {
	tc := &metav1alpha1.TestCase{}

	options := []base.OptionFunc{
		base.ResultOpts(tc),
		base.QueryOpts(map[string]string{
			"buildID": params.BuildID,
		}),
	}
	uri := fmt.Sprintf("projects/%s/testplans/%s/testcases/%s", params.Project, params.TestPlanID, params.TestCaseID)
	if err := p.Get(ctx, p.ClassAddress, uri, options...); err != nil {
		return nil, err
	}

	return tc, nil
}
