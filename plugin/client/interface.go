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

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// Interface base interface for plugins
type Interface interface {
	Path() string
}

// ProjectLister list project api
type ProjectLister interface {
	Interface
	ListProjects(ctx context.Context, option metav1alpha1.ListOptions) (*metav1alpha1.ProjectList, error)
}

// ProjectCreator create project api
type ProjectCreator interface {
	Interface
	CreateProject(ctx context.Context, project *metav1alpha1.Project) (*metav1alpha1.Project, error)
}

// ResourceLister list resource api
type ResourceLister interface {
	Interface
	ListResources(ctx context.Context, option metav1alpha1.ListOptions) (*metav1alpha1.ResourceList, error)
}

// Client inteface for PluginClient, client code shoud use the interface
// as dependency
type Client interface {
	Get(ctx context.Context, baseURL *duckv1.Addressable, uri string, options ...OptionFunc) error
	Post(ctx context.Context, baseURL *duckv1.Addressable, uri string, options ...OptionFunc) error
	Put(ctx context.Context, baseURL *duckv1.Addressable, uri string, options ...OptionFunc) error
	Delete(ctx context.Context, baseURL *duckv1.Addressable, uri string, options ...OptionFunc) error
}

type ClientProjectGetter interface {
	Project(meta Meta, secret corev1.Secret) ClientProject
}
