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

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client/base"
)

// ListArtifacts list artifacts
func (p *PluginClient) ListArtifacts(ctx context.Context, params metav1alpha1.ArtifactOptions, option metav1alpha1.ListOptions) (*metav1alpha1.ArtifactList, error) {
	list := &metav1alpha1.ArtifactList{}

	uri := fmt.Sprintf("projects/%s/repositories/%s/artifacts", params.Project, params.Repository)
	options := []base.OptionFunc{base.ResultOpts(list), base.ListOpts(option)}
	if err := p.Get(ctx, p.ClassAddress, uri, options...); err != nil {
		return nil, err
	}

	return list, nil
}

// GetArtifact get artifact detail
func (p *PluginClient) GetArtifact(ctx context.Context, params metav1alpha1.ArtifactOptions) (*metav1alpha1.Artifact, error) {
	artifact := &metav1alpha1.Artifact{}

	uri := fmt.Sprintf("projects/%s/repositories/%s/artifacts/%s", params.Project, params.Repository, params.Artifact)
	options := []base.OptionFunc{base.ResultOpts(artifact)}
	if err := p.Get(ctx, p.ClassAddress, uri, options...); err != nil {
		return nil, err
	}

	return artifact, nil
}

// DeleteArtifact delete a artifact
func (p *PluginClient) DeleteArtifact(ctx context.Context, params metav1alpha1.ArtifactOptions) error {
	uri := fmt.Sprintf("projects/%s/repositories/%s/artifacts/%s", params.Project, params.Repository, params.Artifact)
	if err := p.Delete(ctx, p.ClassAddress, uri); err != nil {
		return err
	}

	return nil
}

// DeleteArtifactTag delete a specific tag of the artifact.
func (p *PluginClient) DeleteArtifactTag(ctx context.Context, params metav1alpha1.ArtifactTagOptions) error {
	uri := fmt.Sprintf("projects/%s/repositories/%s/artifacts/%s/tags/%s", params.Project, params.Repository, params.Artifact, params.Tag)
	if err := p.Delete(ctx, p.ClassAddress, uri); err != nil {
		return err
	}
	return nil
}
