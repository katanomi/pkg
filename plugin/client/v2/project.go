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

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client/base"
)

// ListProjects list projects
func (p *PluginClient) ListProjects(ctx context.Context, option metav1alpha1.ListOptions) (*metav1alpha1.ProjectList, error) {
	list := &metav1alpha1.ProjectList{}

	options := []base.OptionFunc{base.ResultOpts(list), base.ListOpts(option)}
	if err := p.Get(ctx, p.ClassAddress, "projects", options...); err != nil {
		return nil, err
	}

	return list, nil
}

// CreateProject create project
func (p *PluginClient) CreateProject(ctx context.Context, project *metav1alpha1.Project) (*metav1alpha1.Project, error) {
	createdProject := &metav1alpha1.Project{}
	options := []base.OptionFunc{base.BodyOpts(project), base.ResultOpts(createdProject)}
	if err := p.Post(ctx, p.ClassAddress, "projects", options...); err != nil {
		return nil, err
	}

	return createdProject, nil
}

// GetProject get project
func (p *PluginClient) GetProject(ctx context.Context, id string) (*metav1alpha1.Project, error) {
	resp := &metav1alpha1.Project{}
	options := []base.OptionFunc{base.ResultOpts(resp)}
	if err := p.Get(ctx, p.ClassAddress, "projects/"+id, options...); err != nil {
		return nil, err
	}

	return resp, nil
}
