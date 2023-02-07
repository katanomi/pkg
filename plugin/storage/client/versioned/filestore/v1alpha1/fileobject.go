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
	"context"
	"fmt"

	"github.com/katanomi/pkg/apis/storage/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
	filestorev1alpha1 "github.com/katanomi/pkg/plugin/storage/capabilities/filestore/v1alpha1"
	sclient "github.com/katanomi/pkg/plugin/storage/client"
)

// FileObjectGetter returns FileObject getter object
type FileObjectGetter interface {
	FileObject(pluginName string) FileObjectInterface
}

// FileObjectInterface is interface for FileObject client
type FileObjectInterface interface {
	PUT(ctx context.Context, fileObj filestorev1alpha1.FileObject,
		options ...client.OptionFunc) (*v1alpha1.FileMeta, error)
	GET(ctx context.Context, key string) (*filestorev1alpha1.FileObject, error)
	DELETE(ctx context.Context, key string) error
}

type fileObjects struct {
	client     sclient.Interface
	pluginName string
}

func (f *fileObjects) PUT(ctx context.Context, fileObj filestorev1alpha1.FileObject,
	options ...client.OptionFunc) (*v1alpha1.FileMeta, error) {
	path := fmt.Sprintf("storageplugin/%s/fileobjects/%s", f.pluginName, fileObj.Spec.Key)
	fileMeta := v1alpha1.FileMeta{}
	err := f.client.Get(ctx, path, client.ResultOpts(&fileMeta))
	if err != nil {
		return nil, err
	}
	return &fileMeta, nil
}

func (f *fileObjects) GET(ctx context.Context, key string) (*filestorev1alpha1.FileObject, error) {
	path := fmt.Sprintf("storageplugin/%s/fileobjects/%s", f.pluginName, key)
	fileObject := filestorev1alpha1.FileObject{}
	err := f.client.Get(ctx, path, client.ResultOpts(&fileObject))
	if err != nil {
		return nil, err
	}
	return &fileObject, nil
}

func (f *fileObjects) DELETE(ctx context.Context, key string) error {
	path := fmt.Sprintf("storageplugin/%s/fileobjects/%s", f.pluginName, key)
	err := f.client.Get(ctx, path)
	if err != nil {
		return err
	}
	return nil
}

// newFileObjects returns a FileObjects
func newFileObjects(c *FileStoreV1alpha1Client, pluginName string) *fileObjects {
	return &fileObjects{
		client:     c.RESTClient(),
		pluginName: pluginName,
	}
}
