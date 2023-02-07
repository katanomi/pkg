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
	filestorev1alpha1 "github.com/katanomi/pkg/plugin/storage/capabilities/filestore/v1alpha1"
	"github.com/katanomi/pkg/plugin/storage/client"
)

type FileStoreV1alpha1Interface interface {
	RESTClient() client.Interface
	FileObjectGetter
	FileMetaGetter
}

// FileStoreV1alpha1Client is client for core v1alpha1
type FileStoreV1alpha1Client struct {
	restClient client.Interface
}

// New creates a new FileStoreV1alpha1Client for the given RESTClient.
func New(c client.Interface) *FileStoreV1alpha1Client {
	return &FileStoreV1alpha1Client{restClient: c}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FileStoreV1alpha1Client) RESTClient() client.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}

func (c *FileStoreV1alpha1Client) FileObject() FileObjectInterface {
	return newFileObjects(c)
}

func (c *FileStoreV1alpha1Client) FileMeta() FileMetaInterface {
	return newFileMetas(c)
}

func NewForClient(pClient *client.StoragePluginClient) *FileStoreV1alpha1Client {
	return &FileStoreV1alpha1Client{restClient: pClient.ForGroupVersion(&filestorev1alpha1.FileStoreV1alpha1GV)}
}
