/*
Copyright 2024 The Katanomi Authors.

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

package types

import (
	"context"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
)

//go:generate mockgen -package=types -destination=../../testing/mock/github.com/katanomi/pkg/plugin/types/plugin.go github.com/katanomi/pkg/plugin/types Interface,PluginRegister,PluginAddressable,DependentResourceGetter,AdditionalWebhookRegister,ResourcePathFormatter,PluginDisplayColumns,PluginAttributes,PluginVersionAttributes,LivenessChecker,Initializer,ToolMetadataGetter

// Interface base interface for plugins
type Interface interface {
	Path() string
	Setup(context.Context, *zap.SugaredLogger) error
}

// PluginRegister plugin registration methods to update IntegrationClass status
type PluginRegister interface {
	Interface
	PluginAddressable
	// GetIntegrationClassName returns integration class name
	GetIntegrationClassName() string
	// GetWebhookURL Returns a Webhook accessible URL for external tools
	// If not supported return nil, false
	GetWebhookURL() (*apis.URL, bool)
	// GetSupportedVersions Returns a list of supported versions by the plugin
	// For SaaS platform plugins use a "online" version.
	GetSupportedVersions() []string
	// GetSecretTypes Returns all secret types supported by the plugin
	GetSecretTypes() []string
	// GetReplicationPolicyTypes return replication policy types for ClusterIntegration
	GetReplicationPolicyTypes() []string
	// GetResourceTypes Returns a list of Resource types that can be used in ClusterIntegration and Integration
	GetResourceTypes() []string
	// GetAllowEmptySecret Returns if an empty secret is allowed with IntegrationClass
	GetAllowEmptySecret() []string
}

// PluginAddressable provides methods to get plugin address url
type PluginAddressable interface {
	// GetAddressURL Returns its own plugin access URL
	GetAddressURL() *apis.URL
}

// DependentResourceGetter checks and returns dependent resource references
type DependentResourceGetter interface {
	// GetDependentResources parses params and returns expected dependent resources
	GetDependentResources(ctx context.Context, params []metav1alpha1.Param) ([]corev1.ObjectReference, error)
}

type AdditionalWebhookRegister interface {
	// GetWebhookSupport get webhook support map
	GetWebhookSupport() map[metav1alpha1.WebhookEventSupportType][]string
}

// ResourcePathFormatter implements a formatter for resource path base on different scene
type ResourcePathFormatter interface {
	// GetResourcePathFmt resource path format
	GetResourcePathFmt() map[metav1alpha1.ResourcePathScene]string
	// GetSubResourcePathFmt resource path format
	GetSubResourcePathFmt() map[metav1alpha1.ResourcePathScene]string
}

// PluginDisplayColumns implements display columns manager.
//
// Used to record the format in which the front end should display data.
// example:
//
// projectColumns: ['{"name":"name","displayName":"_.integrations.project.columns.name","field":"metadata.name"}']
//
// projectColumns: the value agreed usage location.
// {"name":"name","displayName":"_.integrations.project.columns.name","field":"metadata.name"}: describe the details displayed.
// name: the name of the column.
// displayName: index used to find specific display data.
// field: the field of the data to be displayed.
type PluginDisplayColumns interface {
	SetDisplayColumns(k string, values ...metav1alpha1.DisplayColumn)
	GetDisplayColumns() map[string]metav1alpha1.DisplayColumns
}

type PluginAttributes interface {
	SetAttribute(k string, values ...string)
	GetAttribute(k string) []string
	Attributes() map[string][]string
}

// PluginVersionAttributes get diff configurations for different versions.
type PluginVersionAttributes interface {
	// GetVersionAttributes get the differential configuration of the specified version
	GetVersionAttributes(version string) map[string][]string

	// SetVersionAttributes set the differential configuration of the specified version
	SetVersionAttributes(version string, attributes map[string][]string)
}

// LivenessChecker check the tool service is alive
type LivenessChecker interface {
	// CheckAlive check the tool service is alive
	CheckAlive(ctx context.Context) error
}

// Initializer initialize the tool service
type Initializer interface {
	// Initialize  the tool service if desired
	Initialize(ctx context.Context) error
}

// ToolMetadataGetter get the version information corresponding to the address.
type ToolMetadataGetter interface {
	// GetToolMetadata get the version information corresponding to the address.
	GetToolMetadata(ctx context.Context) (*metav1alpha1.ToolMeta, error)
}
