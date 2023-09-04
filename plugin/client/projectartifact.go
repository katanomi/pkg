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
	"context"
	"fmt"
	"io"
	"net/url"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

type ClientProjectArtifact interface {
	List(ctx context.Context, baseURL *duckv1.Addressable, project string, options ...OptionFunc) (*metav1alpha1.ArtifactList, error)
	Get(ctx context.Context, baseURL *duckv1.Addressable, project string, artifact string, options ...OptionFunc) (*metav1alpha1.Artifact, error)
	Put(ctx context.Context, baseURL *duckv1.Addressable, project string, artifact string, options ...OptionFunc) error
	Download(ctx context.Context, baseURL *duckv1.Addressable, project string, artifact string, options ...OptionFunc) (io.ReadCloser, error)
	Delete(ctx context.Context, baseURL *duckv1.Addressable, project string, artifact string, options ...OptionFunc) error
}

type projectArtifact struct {
	client Client
}

// Download implements ClientProjectArtifact.
func (p *projectArtifact) Download(ctx context.Context, baseURL *duckv1.Addressable, project string, artifact string, options ...OptionFunc) (io.ReadCloser, error) {
	uri := fmt.Sprintf("projects/%s/artifacts/%s/file", project, url.PathEscape(artifact))
	options = append(options, ResultOpts(artifact), DoNotParseResponseOpts())

	resp, err := p.client.GetResponse(ctx, baseURL, uri, options...)
	if err != nil {
		return nil, err
	}
	return resp.RawBody(), nil
}

// Put implements ClientProjectArtifact.
func (p *projectArtifact) Put(ctx context.Context, baseURL *duckv1.Addressable, project string, artifact string, options ...OptionFunc) error {
	uri := fmt.Sprintf("projects/%s/artifacts/%s", project, url.PathEscape(artifact))
	return p.client.Put(ctx, baseURL, uri, options...)
}

func newProjectArtifact(client Client) ClientProjectArtifact {
	return &projectArtifact{
		client: client,
	}
}

// List get project artifacts using plugin
func (p *projectArtifact) List(ctx context.Context,
	baseURL *duckv1.Addressable,
	project string,
	options ...OptionFunc) (*metav1alpha1.ArtifactList, error) {

	list := &metav1alpha1.ArtifactList{}

	uri := fmt.Sprintf("projects/%s/artifacts", project)
	options = append(options, ResultOpts(list))
	if err := p.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}

	return list, nil
}

// Get gets artifact using plugin
func (p *projectArtifact) Get(ctx context.Context,
	baseURL *duckv1.Addressable,
	project, artifactName string,
	options ...OptionFunc) (*metav1alpha1.Artifact, error) {
	artifact := &metav1alpha1.Artifact{}
	uri := fmt.Sprintf("projects/%s/artifacts/%s", project, url.PathEscape(artifactName))
	options = append(options, ResultOpts(artifact))
	if err := p.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}

	return artifact, nil
}

// Delete artifact using plugin
func (p *projectArtifact) Delete(ctx context.Context, baseURL *duckv1.Addressable, project string, artifact string, options ...OptionFunc) error {
	uri := fmt.Sprintf("projects/%s/artifacts/%s", project, url.PathEscape(artifact))
	if err := p.client.Delete(ctx, baseURL, uri, options...); err != nil {
		return err
	}

	return nil
}
