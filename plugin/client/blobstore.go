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
	corev1 "k8s.io/api/core/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

type ClientBlobStore interface {
	List(ctx context.Context, baseURL *duckv1.Addressable, options ...OptionFunc) (*metav1alpha1.BlobStoreList, error)
}

type blobStore struct {
	client Client
	meta   Meta
	secret corev1.Secret
}

func newBlobStore(client Client, meta Meta, secret corev1.Secret) ClientBlobStore {
	return &blobStore{
		client: client,
		meta:   meta,
		secret: secret,
	}
}

// List get blob stores using plugin
func (p *blobStore) List(ctx context.Context, baseURL *duckv1.Addressable, options ...OptionFunc) (*metav1alpha1.BlobStoreList, error) {
	list := &metav1alpha1.BlobStoreList{}

	options = append(options, MetaOpts(p.meta), SecretOpts(p.secret), ResultOpts(list))
	if err := p.client.Get(ctx, baseURL, "blobStores", options...); err != nil {
		return nil, err
	}

	return list, nil
}
