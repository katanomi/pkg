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
	"strconv"

	v1alpha12 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/apis/storage/v1alpha1"
	pclient "github.com/katanomi/pkg/plugin/client"
	"github.com/katanomi/pkg/plugin/storage/client"
	"knative.dev/pkg/logging"
)

//go:generate mockgen -source=filemeta.go -destination=../../../../../../testing/mock/github.com/katanomi/pkg/plugin/storage/client/versioned/filestore/v1alpha1/filemeta.go -package=v1alpha1 FileMetaInterface

// FileMetaInterface for file meta restful resource methods
type FileMetaInterface interface {
	GET(ctx context.Context, key string) (*v1alpha1.FileMeta, error)
	List(ctx context.Context, opts v1alpha1.FileMetaListOptions) ([]v1alpha1.FileMeta, error)
}

type fileMetas struct {
	client     client.Interface
	pluginName string
}

// newFileMetas returns a FileMetas
func newFileMetas(c *FileStoreV1alpha1Client, pluginName string) *fileMetas {
	return &fileMetas{
		client:     c.RESTClient(),
		pluginName: pluginName,
	}
}

func (f *fileMetas) GET(ctx context.Context, key string) (*v1alpha1.FileMeta, error) {
	path := fmt.Sprintf("storageplugins/%s/filemetas/%s", f.pluginName, key)
	fileMeta := v1alpha1.FileMeta{}
	err := f.client.Get(ctx, path, pclient.ResultOpts(&fileMeta))
	if err != nil {
		return nil, err
	}
	return &fileMeta, nil
}

func (f *fileMetas) List(ctx context.Context, opts v1alpha1.FileMetaListOptions) ([]v1alpha1.FileMeta, error) {
	logger := logging.FromContext(ctx)
	path := fmt.Sprintf("storageplugins/%s/filemetas", f.pluginName)
	listOpt := v1alpha12.ListOptions{}
	// default recursive is true
	if opts.Recursive == false {
		listOpt.SearchSet("recursive", "false")
	}
	if opts.Prefix != "" {
		listOpt.SearchSet("prefix", opts.Prefix)
	}
	if opts.StartAfter != "" {
		listOpt.SearchSet("startAfter", opts.StartAfter)
	}
	if opts.Limit > 0 {
		listOpt.SearchSet("limit", strconv.Itoa(opts.Limit))
	}

	logger.Debugw("set list meta options for client", "opts", listOpt)
	fileMetas := make([]v1alpha1.FileMeta, 0)

	err := f.client.Get(ctx, path, pclient.ResultOpts(&fileMetas), pclient.ListOpts(listOpt))
	if err != nil {
		return nil, err
	}
	return fileMetas, nil
}
