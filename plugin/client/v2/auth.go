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

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client/base"
)

// AuthCheck check authorization
func (p *PluginClientV2) AuthCheck(ctx context.Context, option metav1alpha1.AuthCheckOptions) (*metav1alpha1.AuthCheck, error) {
	authCheck := &metav1alpha1.AuthCheck{}

	options := []base.OptionFunc{base.ResultOpts(authCheck), base.BodyOpts(option)}
	err := p.Post(ctx, p.ClassAddress, "auth/check", options...)
	return authCheck, err
}

// AuthToken generate token or refresh token
func (p *PluginClientV2) AuthToken(ctx context.Context) (*metav1alpha1.AuthToken, error) {
	authToken := &metav1alpha1.AuthToken{}

	options := []base.OptionFunc{base.ResultOpts(authToken)}
	err := p.Post(ctx, p.ClassAddress, "auth/token", options...)
	return authToken, err
}
