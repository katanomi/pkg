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

package v1alpha1

// Common labels
const (
	// SourceLabelKey indicates the source of the resource
	SourceLabelKey = "katanomi.dev/source"

	// SecretLabelKey secret resource name
	SecretLabelKey = "core.kubernetes.io/secret" //nolint:gosec
	// NamespaceLabelKey namespace of a resource
	NamespaceLabelKey = "core.kubernetes.io/namespace"
	// IntegrationClassLabelKey for integrationclass resource
	IntegrationClassLabelKey = "core.katanomi.dev/integrationClass"

	// ClusterIntegrationLabelKey for cluster integration resources
	ClusterIntegrationLabelKey = "integrations.katanomi.dev/clusterIntegration"
	// IntegrationLabelKey for integration resources
	IntegrationLabelKey = "integrations.katanomi.dev/integration"
	// ProjectLabelKey for integration resources
	ProjectLabelKey = "integrations.katanomi.dev/project"
	// RepositoryLabelKey for integration resources
	RepositoryLabelKey = "integrations.katanomi.dev/repository"

	// IntegrationAutoGenerateAnnotation annotation key to store generate flag.
	IntegrationAutoGenerateAnnotation = "integrations.katanomi.dev/resourceScope.autoGenerate"

	// IntegrationAddressAnnotation annotation key to store integration server address
	IntegrationAddressAnnotation = "integrations.katanomi.dev/integration.address"

	// IntegrationResourceScope annotation key to store integration resource scope
	IntegrationResourceScope = "integrations.katanomi.dev/integration.resourceScope"

	// IntegrationSecretApplyNamespaces annotation key to store apply namespace for current secret
	IntegrationSecretApplyNamespaces = "integrations.katanomi.dev/secret.applyNamespaces"

	// SecretSyncMutationLabelKey label key to select the suitable secret
	SecretSyncMutationLabelKey = "integrations.katanomi.dev/integration.mutation"

	// SettingsTypeLabelKey label key to select the settings secret
	SettingsTypeLabelKey = "settings.katanomi.dev/settingsType"
)

// Common Annotations
const (
	// CreatedTimeAnnotationKey creation time for objects
	CreatedTimeAnnotationKey = "katanomi.dev/creationTime"
	// CrossClusterAnnotationKey annotates a cross cluster resource/action
	CrossClusterAnnotationKey = "katanomi.dev/crossCluster"
	// ReconcileTriggeredAnnotationKey annotation key to trigger reconcile of objects
	ReconcileTriggeredAnnotationKey = "katanomi.dev/reconcileTriggeredOn"
	// NamespaceAnnotationKey namespace of objects
	NamespaceAnnotationKey = "katanomi.dev/namespace"
	// TriggeredByAnnotationKey annotation to store a TriggeredBy struct json
	TriggeredByAnnotationKey = "katanomi.dev/triggeredBy"
	// CreatedByAnnotationKey annotation key to store resource creation username
	CreatedByAnnotationKey = "katanomi.dev/createdBy"
	// UpdatedByAnnotationKey annotation key to store resource update username
	UpdatedByAnnotationKey = "katanomi.dev/updatedBy"
	// SecretTypeAnnotationKey annotation key for an existed secret with a different type
	SecretTypeAnnotationKey = "katanomi.dev/secretType" //nolint:gosec
	// ClusterNameAnnotationKey annotation key to store resource cluster name
	ClusterNameAnnotationKey = "integrations.katanomi.dev/clusterName"
	// TriggerNameAnnotationKey annotation key to store a friendly trigger name
	TriggerNameAnnotationKey = "katanomi.dev/triggerName"
	// SettingsConvertTypesKey annotation key to store the setting types need to be converted
	SettingsConvertTypesKey = "settings.katanomi.dev/convertTypes"
	// SettingsAutoGenerateKey annotation key to store whether the secret is automatically generated
	SettingsAutoGenerateKey = "settings.katanomi.dev/autoGenerate"
	// UIDescriptorsAnnotationKey annotation for storing ui descriptors in resources
	UIDescriptorsAnnotationKey = "ui.katanomi.dev/descriptors"
	// PodAnnotationKeyPrefix uses a prefix for pod annotations in katanomi
	PodAnnotationKeyPrefix = "pod.katanomi.dev/"
)

// Attribute keys for Integrations
// Keys used in IntegrationClass.status.attributes
const (
	AuthAttributeKey                   = "auth"
	ReplicationPolicyTypesAttributeKey = "replicationPolicyTypes"
	ResourceTypesAttributeKey          = "resourceTypes"
	SettingsTypesAttributeKey          = "settingsTypes"
	MethodsAttributeKey                = "methods"
	AllowEmptySecretAttributeKey       = "allowEmptySecret"
	DefaultProjectTypeAttributeKey     = "defaultProjectSubType"
	ResourcePathFormat                 = "resourcePathFormat"
	// GitPRRevisionPrefixes allows git related integrations
	// define custom PR revision prefixes
	GitPRRevisionPrefixes = "gitPRRevisionPrefixes"
)

// Attribute values for label source
const (
	LabelSourceSystem = "system"
	LabelSourceUser   = "user"
)

// Annotation keys for artifact parameter
const (
	// ArtifactAliasAnnotationKey indicates the alias of the artifact parameter in the delivery
	ArtifactAliasAnnotationKey = "alias"

	// ImageRegistryEndpoint for artifact endpoint
	ImageRegistryEndpoint = "imageRegistryEndpoint"
)
