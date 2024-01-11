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

	"github.com/go-resty/resty/v2"
	"github.com/katanomi/pkg/apis/storage/v1alpha1"
	"github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/plugin/client"
	filestorev1alpha1 "github.com/katanomi/pkg/plugin/storage/capabilities/filestore/v1alpha1"
	sclient "github.com/katanomi/pkg/plugin/storage/client"
)

//go:generate mockgen -source=fileobject.go -destination=../../../../../../testing/mock/github.com/katanomi/pkg/plugin/storage/client/versioned/filestore/v1alpha1/interface.go -package=v1alpha1 FileObjectInterface

// FileObjectInterface is interface for FileObject client
type FileObjectInterface interface {
	PUT(ctx context.Context, fileObj filestorev1alpha1.FileObject,
		options ...client.OptionFunc) (*v1alpha1.FileMeta, error)
	GET(ctx context.Context, fileObjectName string) (*filestorev1alpha1.FileObject, error)
	DELETE(ctx context.Context, fileObjectName string) error
}

// fileObjectClientKey for key of fileObjectClient in context
type fileObjectClientKey struct{}

type fileObjects struct {
	client     sclient.Interface
	pluginName string
}

func (f *fileObjects) PUT(ctx context.Context, fileObj filestorev1alpha1.FileObject,
	options ...client.OptionFunc) (*v1alpha1.FileMeta, error) {
	path := fmt.Sprintf("storageplugins/%s/fileobjects/%s",
		f.pluginName, fileObj.Name)

	retMeta := v1alpha1.FileMeta{}
	err := f.client.Put(ctx, path, client.ResultOpts(&retMeta),
		client.HeaderOpts(v1alpha1.HeaderFileMeta, fileObj.FileMeta.Encode()),
		client.HeaderOpts("Content-Type", "application/octet-stream"),
		// this must be set or 406 status code will be returned
		client.HeaderOpts("Accept", "application/json"),
		client.BodyOpts(fileObj.FileReadCloser),
	)
	if err != nil {
		return nil, err
	}
	return &retMeta, nil
}

func (f *fileObjects) GET(ctx context.Context, key string) (*filestorev1alpha1.FileObject, error) {
	path := fmt.Sprintf("storageplugins/%s/fileobjects/%s", f.pluginName, key)
	fileObject := filestorev1alpha1.FileObject{}
	resp, err := f.client.GetResponse(ctx, path, func(request *resty.Request) {
		request.SetDoNotParseResponse(true)
	})
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, errors.AsStatusError(resp)
	}

	fileMetaEncoded := resp.Header().Get(v1alpha1.HeaderFileMeta)
	fileMeta, err := v1alpha1.DecodeAsFileMeta(fileMetaEncoded)
	if err == nil {
		fileObject.FileMeta = *fileMeta
	}
	fileObject.FileReadCloser = resp.RawBody()
	return &fileObject, nil
}

func (f *fileObjects) DELETE(ctx context.Context, key string) error {
	path := fmt.Sprintf("storageplugins/%s/fileobjects/%s", f.pluginName, key)
	err := f.client.Delete(ctx, path)
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

// WithFileObjectClient inject fileObjectClient to context
func WithFileObjectClient(ctx context.Context, fileObjectClient FileObjectInterface) context.Context {
	return context.WithValue(ctx, fileObjectClientKey{}, fileObjectClient)
}

// FileObjectClientFrom get fileObjectClient from context
func FileObjectClientFrom(ctx context.Context) FileObjectInterface {
	fileObjectClientVal, _ := ctx.Value(fileObjectClientKey{}).(FileObjectInterface)
	return fileObjectClientVal
}
