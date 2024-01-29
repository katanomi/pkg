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

// Package v1alpha1 defines versioned interfaces for file-store capability
package v1alpha1

import (
	"context"
	"io"

	"github.com/katanomi/pkg/apis/storage/v1alpha1"
)

// FileStoreCapable defines methods of file-store capability
//
//go:generate mockgen -source=interface.go -destination=../../../../../testing/mock/github.com/katanomi/pkg/plugin/storage/capabilities/filestore/v1alpha1/interface.go -package=v1alpha1 FileStoreCapable
type FileStoreCapable interface {
	FileObjectInterface
	FileMetaInterface
}

// FileObjectInterface for file object
type FileObjectInterface interface {
	// GetFileObject for get file object
	GetFileObject(ctx context.Context, objectName string) (*FileObject, error)
	// PutFileObject for put file object, use FileObject.Name to store objectName
	PutFileObject(ctx context.Context, obj *FileObject) (*v1alpha1.FileMeta,
		error)
	// DeleteFileObject for delete file object
	DeleteFileObject(ctx context.Context, objectName string) error
}

// FileMetaInterface for file meta
type FileMetaInterface interface {
	// ListFileMetas for list file metas
	ListFileMetas(ctx context.Context, opt v1alpha1.FileMetaListOptions) ([]v1alpha1.FileMeta, error)
	// GetFileMeta for get file meta
	GetFileMeta(ctx context.Context, objectName string) (*v1alpha1.FileMeta, error)
}

// FileObject wraps FileMeta with file reader for implementing file download
type FileObject struct {
	v1alpha1.FileMeta
	FileReadCloser io.ReadCloser
}
