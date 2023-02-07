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
	filestoreintv1alpha1 "github.com/katanomi/pkg/plugin/storage/capabilities/filestore/v1alpha1"
	corev1alpha1 "github.com/katanomi/pkg/plugin/storage/client/versioned/core/v1alpha1"
	filestorev1alpha1 "github.com/katanomi/pkg/plugin/storage/client/versioned/filestore/v1alpha1"
)

// Clientset contains the  core and capabilities plugin clients.
type V1alpha1Clientset struct {
	client *StoragePluginClient
}

// CoreV1alpha1 return core v1alpha1 interface
func (c *V1alpha1Clientset) Core() corev1alpha1.CoreV1alpha1Interface {
	return corev1alpha1.New(c.client)
}

// FileStoreV1alpha1 return file store v1alpha1 interface
func (c *V1alpha1Clientset) FileStore(pluginName string) filestorev1alpha1.FileStoreV1alpha1Interface {
	return filestorev1alpha1.New(c.client.ForGroupVersion(&filestoreintv1alpha1.FileStoreV1alpha1GV), pluginName)
}
