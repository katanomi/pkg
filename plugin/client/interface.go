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

	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"

	cloudevent "github.com/cloudevents/sdk-go/v2"
)

// Interface base interface for plugins
type Interface interface {
	Path() string
}

// PluginRegister plugin registration methods to update IntegationClass status
type PluginRegister interface {
	Interface
	GetIntegrationClassName() string
	// Returns its own plugin access URL
	GetAddressURL() *apis.URL
	// Returns a Webhook accessible URL for external tools
	// If not supported return nil, false
	GetWebhookURL() (*apis.URL, bool)
	// Returns a list of supported versions by the plugin
	// For SaaS platform plugins use a "online" version.
	GetSupportedVersions() []string
	// Returns all secret types supported by the plugin
	GetSecretTypes() []string
	// GetReplicationPolicyTypes return replication policy types for ClusterIntegration
	GetReplicationPolicyTypes() []string
	// Returns a list of Resource types that can be used in ClusterIntegration and Integration
	GetResourceTypes() []string
}

// ProjectLister list project api
type ProjectLister interface {
	Interface
	ListProjects(ctx context.Context, option metav1alpha1.ListOptions) (*metav1alpha1.ProjectList, error)
}

// ProjectGetter list project api
type ProjectGetter interface {
	Interface
	GetProject(ctx context.Context, id string) (*metav1alpha1.Project, error)
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

// RepositoryLister list repository
type RepositoryLister interface {
	Interface
	ListRepositories(ctx context.Context, params metav1alpha1.RepositoryOptions, option metav1alpha1.ListOptions) (*metav1alpha1.RepositoryList, error)
}

// ArtifactLister list artifact
type ArtifactLister interface {
	Interface
	ListArtifacts(ctx context.Context, params metav1alpha1.ArtifactOptions, option metav1alpha1.ListOptions) (*metav1alpha1.ArtifactList, error)
}

// ArtifactGetter get artifact detail
type ArtifactGetter interface {
	Interface
	GetArtifact(ctx context.Context, params metav1alpha1.ArtifactOptions) (*metav1alpha1.Artifact, error)
}

// ArtifactDeleter delete artifact
type ArtifactDeleter interface {
	Interface
	DeleteArtifact(ctx context.Context, params metav1alpha1.ArtifactOptions) error
}

// ScanImage scan image
type ScanImage interface {
	Interface
	ScanImage(ctx context.Context, params metav1alpha1.ArtifactOptions) error
}

// WebhookRegister used to register and manage webhooks
type WebhookRegister interface {
	// Use the methods below to manage webhooks in the target platform
	CreateWebhook(ctx context.Context, spec metav1alpha1.WebhookRegisterSpec, secret corev1.Secret) (metav1alpha1.WebhookRegisterStatus, error)
	UpdateWebhook(ctx context.Context, spec metav1alpha1.WebhookRegisterSpec, secret corev1.Secret) (metav1alpha1.WebhookRegisterStatus, error)
	DeleteWebhook(ctx context.Context, spec metav1alpha1.WebhookRegisterSpec, secret corev1.Secret) error
}

// WebhookResourceDiffer used to compare different webhook resources in order to provide
// a way to merge webhook registration requests. If not provided, the resource's URI will be directly compared
type WebhookResourceDiffer interface {
	// IsSameResource will provide two ResourceURI
	// the plugin should discern if they are the same.
	// If this method is not implemented a standard comparisons will be used
	IsSameResource(ctx context.Context, i, j metav1alpha1.ResourceURI) bool
}

// WebhookReceiver receives a webhook request with validation and transform it into a cloud event
type WebhookReceiver interface {
	Interface
	ReceiveWebhook(ctx context.Context, req *restful.Request, secret string) (cloudevent.Event, error)
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
