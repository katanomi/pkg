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

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// ClientAuth provides methods to verify authentication
type ClientAuth interface {
	Check(ctx context.Context, baseURL *duckv1.Addressable, options metav1alpha1.AuthCheckOptions, opts ...OptionFunc) (*metav1alpha1.AuthCheck, error)
	Token(ctx context.Context, baseURL *duckv1.Addressable, opts ...OptionFunc) (*metav1alpha1.AuthToken, error)
}

type authClient struct {
	client Client
}

func newAuthClient(client Client) ClientAuth {
	return &authClient{
		client: client,
	}
}

// Check checks auth data against an integrated service
func (auth *authClient) Check(ctx context.Context, baseURL *duckv1.Addressable, options metav1alpha1.AuthCheckOptions, opts ...OptionFunc) (authCheck *metav1alpha1.AuthCheck, err error) {
	uri := "auth/check"
	authCheck = &metav1alpha1.AuthCheck{}

	opts = append(opts, ResultOpts(authCheck), BodyOpts(options))
	err = auth.client.Post(ctx, baseURL, uri, opts...)
	return
}

// Token generates or refreshes access token
func (auth *authClient) Token(ctx context.Context, baseURL *duckv1.Addressable, opts ...OptionFunc) (authToken *metav1alpha1.AuthToken, err error) {
	uri := "auth/token"
	authToken = &metav1alpha1.AuthToken{}

	opts = append(opts, ResultOpts(authToken))
	err = auth.client.Post(ctx, baseURL, uri, opts...)
	return
}
