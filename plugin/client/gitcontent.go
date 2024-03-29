/*
Copyright 2021 The Katanomi Authors.

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
	"bytes"
	"context"
	"encoding/base64"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/path"
	"k8s.io/apimachinery/pkg/api/errors"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

//go:generate mockgen -source=gitcontent.go -destination=../../testing/mock/github.com/katanomi/pkg/plugin/client/gitcontent.go -package=client ClientGitContent
type ClientGitContent interface {
	Get(ctx context.Context, baseURL *duckv1.Addressable, option metav1alpha1.GitRepoFileOption, options ...OptionFunc) (*metav1alpha1.GitRepoFile, error)
	Create(ctx context.Context, baseURL *duckv1.Addressable, payload metav1alpha1.CreateRepoFilePayload, options ...OptionFunc) (*metav1alpha1.GitCommit, error)
}

type gitContent struct {
	client Client
}

func newGitContent(client Client) ClientGitContent {
	return &gitContent{
		client: client,
	}
}

func (g *gitContent) Get(ctx context.Context, baseURL *duckv1.Addressable, option metav1alpha1.GitRepoFileOption, options ...OptionFunc) (*metav1alpha1.GitRepoFile, error) {
	fileInfo := &metav1alpha1.GitRepoFile{}
	options = append(options, QueryOpts(map[string]string{"ref": option.Ref}), ResultOpts(fileInfo))
	if option.Repository == "" {
		return nil, errors.NewBadRequest("repo is empty string")
	}
	if option.Path == "" {
		return nil, errors.NewBadRequest("file path is empty string")
	}
	uri := path.Format("projects/%s/coderepositories/%s/content/%s", option.Project, option.Repository, option.Path)
	if err := g.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}
	return fileInfo, nil
}

func (g *gitContent) Create(ctx context.Context, baseURL *duckv1.Addressable, payload metav1alpha1.CreateRepoFilePayload, options ...OptionFunc) (*metav1alpha1.GitCommit, error) {
	commitInfo := &metav1alpha1.GitCommit{}
	var b bytes.Buffer
	w := base64.NewEncoder(base64.StdEncoding, &b)
	_, err := w.Write(payload.Content)
	if err != nil {
		return nil, err
	}
	w.Close()
	payload.Content = b.Bytes()

	options = append(options, BodyOpts(payload.CreateRepoFileParams), ResultOpts(commitInfo))
	if payload.Repository == "" {
		return nil, errors.NewBadRequest("repo is empty string")
	}

	uri := path.Format("projects/%s/coderepositories/%s/content/%s", payload.Project, payload.Repository, payload.FilePath)

	if err := g.client.Post(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}

	return commitInfo, nil
}
