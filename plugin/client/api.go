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
	"context"
)

// ListOption list option for list option
type ListOption struct {
	Keyword      string
	ItemsPerPage int
	Page         int
}

// PluginClient plugin api sub path
type PluginClient interface {
	Path() string
}

// ProjectLister list project api
type ProjectLister interface {
	PluginClient
	ListProjects(ctx context.Context, option ListOption) (ProjectList, error)
}

// ProjectCreator create project api
type ProjectCreator interface {
	PluginClient
	CreateProject(ctx context.Context, project *Project) (*Project, error)
}

// ResourceLister list resource api
type ResourceLister interface {
	PluginClient
	ListResources(ctx context.Context, option ListOption) (ResourceList, error)
}
