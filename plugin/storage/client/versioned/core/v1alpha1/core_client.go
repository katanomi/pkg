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

package v1alpha1

import (
	"github.com/katanomi/pkg/plugin/storage/client"
	"github.com/katanomi/pkg/plugin/storage/core/v1alpha1"
)

type CoreV1alpha1Interface interface {
	RESTClient() client.Interface
	AuthGetter
}

// New creates a new CoreV1alpha1Client for the given RESTClient.
func New(c client.Interface) *CoreV1alpha1Client {
	return &CoreV1alpha1Client{restClient: c}
}

// CoreV1alpha1Client is client for core v1alpha1
type CoreV1alpha1Client struct {
	restClient client.Interface
}

func (c *CoreV1alpha1Client) RESTClient() client.Interface {
	return c.restClient
}

func (c *CoreV1alpha1Client) Auth() AuthInterface {
	return newAuth(c)
}

func NewForClient(pClient *client.StoragePluginClient) *CoreV1alpha1Client {
	return &CoreV1alpha1Client{restClient: pClient.ForGroupVersion(&v1alpha1.CoreV1alpha1GV)}
}
