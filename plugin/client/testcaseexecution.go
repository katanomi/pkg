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

package client

import (
	"context"
	"fmt"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// ClientTestCaseExecution client for test case execution
type ClientTestCaseExecution interface {
	List(ctx context.Context,
		baseURL *duckv1.Addressable,
		params metav1alpha1.TestProjectOptions,
		options ...OptionFunc) (*metav1alpha1.TestCaseExecutionList, error)
	Create(ctx context.Context,
		baseURL *duckv1.Addressable,
		params metav1alpha1.TestProjectOptions,
		payload metav1alpha1.TestCaseExecution,
		options ...OptionFunc,
	) (*metav1alpha1.TestCaseExecution, error)
}

type testCaseExecution struct {
	client Client
	meta   Meta
	secret corev1.Secret
}

func newTestCaseExecution(client Client, meta Meta, secret corev1.Secret) ClientTestCaseExecution {
	return &testCaseExecution{
		client: client,
		meta:   meta,
		secret: secret,
	}
}

// List get project using plugin
func (p *testCaseExecution) List(ctx context.Context,
	baseURL *duckv1.Addressable,
	params metav1alpha1.TestProjectOptions,
	options ...OptionFunc) (*metav1alpha1.TestCaseExecutionList, error) {
	list := &metav1alpha1.TestCaseExecutionList{}

	uri := fmt.Sprintf(
		"projects/%s/testplans/%s/testcases/%s/executions",
		params.Project,
		params.TestPlanID,
		params.TestCaseID,
	)
	options = append(options, MetaOpts(p.meta), SecretOpts(p.secret), ResultOpts(list), QueryOpts(map[string]string{
		"buildID": params.BuildID,
	}))
	if err := p.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}

	return list, nil
}

func (p *testCaseExecution) Create(ctx context.Context,
	baseURL *duckv1.Addressable,
	params metav1alpha1.TestProjectOptions,
	payload metav1alpha1.TestCaseExecution,
	options ...OptionFunc) (*metav1alpha1.TestCaseExecution, error) {
	tc := &metav1alpha1.TestCaseExecution{}

	uri := fmt.Sprintf("projects/%s/testplans/%s/testcases/%s/executions", params.Project, params.TestPlanID,
		params.TestCaseID)
	options = append(options, MetaOpts(p.meta), SecretOpts(p.secret), BodyOpts(payload), ResultOpts(tc))
	if err := p.client.Post(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}

	return tc, nil
}
