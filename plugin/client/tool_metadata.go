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

package client

import (
	"context"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

type ClientToolMetadata interface {
	GetToolMetadata(ctx context.Context, baseURL *duckv1.Addressable) (*metav1alpha1.ToolMeta, error)
}

type toolMetadata struct {
	client Client
}

func newToolMetadata(client Client) ClientToolMetadata {
	return &toolMetadata{
		client: client,
	}
}

// GetVersion get the version metadata corresponding to the address.
func (p *toolMetadata) GetToolMetadata(ctx context.Context, baseURL *duckv1.Addressable) (*metav1alpha1.ToolMeta, error) {
	toolMate := &metav1alpha1.ToolMeta{}
	if err := p.client.Get(ctx, baseURL, "tools/metadata", ResultOpts(toolMate)); err != nil {
		return nil, err
	}
	return toolMate, nil
}
