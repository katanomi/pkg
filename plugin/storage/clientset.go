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

package storage

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	client2 "github.com/katanomi/pkg/plugin/storage/client"
	archivev1alpha1 "github.com/katanomi/pkg/plugin/storage/client/versioned/archive/v1alpha1"
	corev1alpha1 "github.com/katanomi/pkg/plugin/storage/client/versioned/core/v1alpha1"
	"github.com/katanomi/pkg/plugin/storage/client/versioned/filestore/v1alpha1"
	v1 "knative.dev/pkg/apis/duck/v1"
)

//go:generate mockgen -source=clientset.go -destination=../../testing/mock/github.com/katanomi/pkg/plugin/storage/clientset.go -package=storage Interface

type Interface interface {
	CoreV1alpha1() corev1alpha1.CoreV1alpha1Interface
	FileStoreV1alpha1() v1alpha1.FileStoreV1alpha1Interface
	ArchiveV1alpha1() archivev1alpha1.ArchiveInterface
}

// Clientset contains the  core and capabilities plugin clients.
type Clientset struct {
	coreV1alpha1      corev1alpha1.CoreV1alpha1Interface
	fileStoreV1alpha1 v1alpha1.FileStoreV1alpha1Interface
	archiveV1alpha1   archivev1alpha1.ArchiveInterface
}

// CoreV1alpha1 return core v1alpha1 interface
func (c *Clientset) CoreV1alpha1() corev1alpha1.CoreV1alpha1Interface {
	return c.coreV1alpha1
}

// FileStoreV1alpha1 return file store v1alpha1 interface
func (c *Clientset) FileStoreV1alpha1() v1alpha1.FileStoreV1alpha1Interface {
	return c.fileStoreV1alpha1
}

// ArchiveV1alpha1 return archive v1alpha1 interface
func (c *Clientset) ArchiveV1alpha1() archivev1alpha1.ArchiveInterface {
	return c.archiveV1alpha1
}

// NewForClient return a new clientset instance
func NewForClient(classAddress *v1.Addressable, client *resty.Client) (*Clientset, error) {
	var (
		cs  Clientset
		err error
	)

	if classAddress == nil {
		err = fmt.Errorf("nil class address")
		return nil, err
	}

	pluginClient := client2.NewStoragePluginClient(classAddress, client2.WithRestClient(client))
	cs.coreV1alpha1 = corev1alpha1.NewForClient(pluginClient)
	cs.fileStoreV1alpha1 = v1alpha1.NewForClient(pluginClient)
	cs.archiveV1alpha1 = archivev1alpha1.NewForClient(pluginClient)

	return &cs, err
}
