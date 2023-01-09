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

	duckv1 "knative.dev/pkg/apis/duck/v1"
)

type ClientToolService interface {
	CheckAlive(ctx context.Context, options ...OptionFunc) error
	Initialize(ctx context.Context, options ...OptionFunc) error
}

type toolService struct {
	client  Client
	baseURL *duckv1.Addressable
}

func newToolService(client Client, baseURL *duckv1.Addressable) ClientToolService {
	return &toolService{
		client:  client,
		baseURL: baseURL,
	}
}

// CheckAlive to check if the tool service is alive
func (p *toolService) CheckAlive(ctx context.Context, options ...OptionFunc) error {
	return p.client.Get(ctx, p.baseURL, "tools/liveness", options...)
}

// Initialize to initialize the tool service
func (p *toolService) Initialize(ctx context.Context, options ...OptionFunc) error {
	return p.client.Get(ctx, p.baseURL, "tools/initialize", options...)
}
