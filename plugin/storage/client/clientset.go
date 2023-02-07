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
	"fmt"

	"github.com/go-resty/resty/v2"
	v1 "knative.dev/pkg/apis/duck/v1"
)

type VersionedClient interface {
	V1alpha1() *V1alpha1Clientset
}

type versionedClient struct {
	client *StoragePluginClient
}

func (v versionedClient) V1alpha1() *V1alpha1Clientset {
	return &V1alpha1Clientset{client: v.client}
}

// NewForClient return a new clientset instance
func NewForClient(classAddress *v1.Addressable, client *resty.Client) (VersionedClient, error) {
	if classAddress == nil {
		err := fmt.Errorf("nil class address")
		return nil, err
	}

	cs := versionedClient{
		client: NewStoragePluginClient(classAddress, WithRestClient(client)),
	}
	return &cs, nil
}
