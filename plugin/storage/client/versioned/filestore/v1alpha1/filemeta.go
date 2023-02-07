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
	client2 "github.com/katanomi/pkg/plugin/client"
	"github.com/katanomi/pkg/plugin/storage/client"
)

type FileMetaGetter interface {
	FileMeta() FileMetaInterface
}

type FileMetaInterface interface {
	GET(ctx context.Context, pluginName string, key string) (*v1alpha1.FileMeta, error)
	// TODO: Add List methods
}

type fileMetas struct {
	client client.Interface
}

// newFileMetas returns a FileMetas
func newFileMetas(c *FileStoreV1alpha1Client) *fileMetas {
	return &fileMetas{
		client: c.RESTClient(),
	}
}

func (f *fileMetas) GET(ctx context.Context, pluginName, key string) (*v1alpha1.FileMeta, error) {
	path := fmt.Sprintf("storageplugin/%s/filemetas/%s", pluginName, key)
	fileMeta := v1alpha1.FileMeta{}
	err := f.client.Get(ctx, path, client2.ResultOpts(&fileMeta))
	if err != nil {
		return nil, err
	}
	return &fileMeta, nil
}
