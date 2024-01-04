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

// ListBlobStores list blob stores
func (p *PluginClient) ListBlobStores(ctx context.Context, listOption metav1alpha1.ListOptions) (*metav1alpha1.BlobStoreList, error) {
	list := &metav1alpha1.BlobStoreList{}

	options := []base.OptionFunc{base.ResultOpts(list), base.ListOpts(listOption)}
	if err := p.Get(ctx, p.ClassAddress, "blobStores", options...); err != nil {
		return nil, err
	}

	return list, nil
}
