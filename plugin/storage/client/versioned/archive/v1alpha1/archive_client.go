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
	archivev1alpha1 "github.com/katanomi/pkg/plugin/storage/capabilities/archive/v1alpha1"
	"github.com/katanomi/pkg/plugin/storage/client"
)

//go:generate mockgen -source=archive.go -destination=../../../../../../testing/mock/github.com/katanomi/pkg/plugin/storage/client/versioned/archive/v1alpha1/interface.go -package=v1alpha1 ArchiveInterface
type ArchiveInterface interface {
	RESTClient() client.Interface
	RecordGetter
}

// ArchiveClient is client for core v1alpha1
type ArchiveClient struct {
	restClient client.Interface
	pluginName string
}

// New creates a new ArchiveClient for the given RESTClient.
func New(c client.Interface) *ArchiveClient {
	return &ArchiveClient{restClient: c}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *ArchiveClient) RESTClient() client.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}

func (c *ArchiveClient) Record(pluginName string) RecordInterface {
	return newRecord(c, pluginName)
}

func NewForClient(pClient *client.StoragePluginClient) *ArchiveClient {
	return &ArchiveClient{restClient: pClient.ForGroupVersion(&archivev1alpha1.ArchiveV1alpha1GV)}
}
