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

// ClientTestCase for test case
type ClientTestCase interface {
	List(ctx context.Context, baseURL *duckv1.Addressable, params metav1alpha1.TestProjectOptions, options ...OptionFunc) (*metav1alpha1.TestCaseList, error)
	Get(ctx context.Context, baseURL *duckv1.Addressable, params metav1alpha1.TestProjectOptions, options ...OptionFunc) (*metav1alpha1.TestCase, error)
}

type testCase struct {
	client Client
	meta   Meta
	secret corev1.Secret
}

func newTestCase(client Client, meta Meta, secret corev1.Secret) ClientTestCase {
	return &testCase{
		client: client,
		meta:   meta,
		secret: secret,
	}
}

// List get project using plugin
func (p *testCase) List(ctx context.Context, baseURL *duckv1.Addressable, params metav1alpha1.TestProjectOptions, options ...OptionFunc) (*metav1alpha1.TestCaseList, error) {
	list := &metav1alpha1.TestCaseList{}

	uri := fmt.Sprintf("projects/%s/testplans/%s/testcases", params.Project, params.TestPlanID)
	options = append(options, MetaOpts(p.meta), SecretOpts(p.secret), ResultOpts(list), QueryOpts(map[string]string{
		"buildID": params.BuildID,
	}))
	if err := p.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}

	return list, nil
}

func (p *testCase) Get(ctx context.Context, baseURL *duckv1.Addressable, params metav1alpha1.TestProjectOptions, options ...OptionFunc) (*metav1alpha1.TestCase, error) {
	tc := &metav1alpha1.TestCase{}

	uri := fmt.Sprintf("projects/%s/testplans/%s/testcases/%s", params.Project, params.TestPlanID, params.TestCaseID)
	options = append(options, MetaOpts(p.meta), SecretOpts(p.secret), ResultOpts(tc), QueryOpts(map[string]string{
		"buildID": params.BuildID,
	}))
	if err := p.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}

	return tc, nil
}
