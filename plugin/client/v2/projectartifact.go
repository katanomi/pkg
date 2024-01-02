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

package v2

import (
	"context"
	"fmt"
	"io"
	"net/url"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client/base"
)

// GetProjectArtifactFile download project artifact file
// Note: you must close the io.ReadCloser after use or it will cause a memory leak
func (p *PluginClient) GetProjectArtifactFile(ctx context.Context, params metav1alpha1.ProjectArtifactOptions) (io.ReadCloser, error) {
	options := []base.OptionFunc{base.DoNotParseResponseOpts()}

	uri := fmt.Sprintf("projects/%s/artifacts/%s/file", params.Project, url.PathEscape(params.Artifact))
	resp, err := p.GetResponse(ctx, p.ClassAddress, uri, options...)
	if err != nil {
		return nil, err
	}
	return resp.RawBody(), nil
}

// ListProjectArtifacts list project artifacts
func (p *PluginClient) ListProjectArtifacts(ctx context.Context, params metav1alpha1.ArtifactOptions, option metav1alpha1.ListOptions) (*metav1alpha1.ArtifactList, error) {
	list := &metav1alpha1.ArtifactList{}
	options := []base.OptionFunc{base.ResultOpts(list), base.ListOpts(option)}

	uri := fmt.Sprintf("projects/%s/artifacts", params.Project)
	if err := p.Get(ctx, p.ClassAddress, uri, options...); err != nil {
		return nil, err
	}

	return list, nil
}

// UploadArtifact upload artifact file
func (p *PluginClient) UploadArtifact(ctx context.Context, params metav1alpha1.ProjectArtifactOptions, r io.Reader) error {
	options := []base.OptionFunc{base.BodyOpts(r)}
	uri := fmt.Sprintf("projects/%s/artifacts/%s", params.Project, url.PathEscape(params.Artifact))
	return p.Put(ctx, p.ClassAddress, uri, options...)
}

// GetProjectArtifact get project artifact
func (p *PluginClient) GetProjectArtifact(ctx context.Context, params metav1alpha1.ProjectArtifactOptions) (*metav1alpha1.Artifact, error) {
	artifact := &metav1alpha1.Artifact{}
	uri := fmt.Sprintf("projects/%s/artifacts/%s", params.Project, url.PathEscape(params.Artifact))
	options := []base.OptionFunc{base.ResultOpts(artifact)}
	if err := p.Get(ctx, p.ClassAddress, uri, options...); err != nil {
		return nil, err
	}

	return artifact, nil
}

// DeleteProjectArtifact delete project artifact
func (p *PluginClient) DeleteProjectArtifact(ctx context.Context, params metav1alpha1.ProjectArtifactOptions) error {
	uri := fmt.Sprintf("projects/%s/artifacts/%s", params.Project, url.PathEscape(params.Artifact))
	return p.Delete(ctx, p.ClassAddress, uri)
}
