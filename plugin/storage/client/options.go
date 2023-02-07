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
	"github.com/go-resty/resty/v2"
	"k8s.io/apimachinery/pkg/runtime/schema"
	v1 "knative.dev/pkg/apis/duck/v1"
)

// BuildOptions Options to build the plugin client
type BuildOptions func(client *StoragePluginClient)

// WithRestClient adds a custom client build options for plugin client
func WithRestClient(clt *resty.Client) BuildOptions {
	return func(client *StoragePluginClient) {
		client.client = clt
	}
}

// WithGroupVersion adds a custom client build options for plugin client
func WithGroupVersion(gv *schema.GroupVersion) BuildOptions {
	return func(client *StoragePluginClient) {
		client.groupVersion = gv
	}
}

// WithClassAddress sets client based address url for plugin client
func WithClassAddress(address *v1.Addressable) BuildOptions {
	return func(client *StoragePluginClient) {
		client.classAddress = address
	}
}
