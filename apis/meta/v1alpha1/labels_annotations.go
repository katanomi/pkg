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

package v1alpha1

// Common labels
const (
	// SourceLabelKey indicates the source of the resource
	SourceLabelKey = "katanomi.dev/source"
	// ManagerByLabelKey indicates the manager of the resource
	ManagerByLabelKey = "katanomi.dev/managedBy"

	// SecretLabelKey secret resource name
	SecretLabelKey = "core.kubernetes.io/secret" //nolint:gosec
	// NamespaceLabelKey namespace of a resource
	NamespaceLabelKey = "core.kubernetes.io/namespace"
	// IntegrationClassLabelKey for integrationclass resource
	IntegrationClassLabelKey = "core.katanomi.dev/integrationClass"
	// ProxyEnabledLabelKey
	ProxyEnabledLabelKey = "core.katanomi.dev/proxyEnabled"

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

	// IntegrationSecretResourcePathFmt annotation indicates resource path format for current secret
	IntegrationSecretResourcePathFmt = "integrations.katanomi.dev/secret.resourcePathFmt"
	// IntegrationSecretSubResourcePathFmt annotation indicates sub resource path format for current secret
	IntegrationSecretSubResourcePathFmt = "integrations.katanomi.dev/secret.subResourcePathFmt"

	// SecretSyncMutationLabelKey label key to select the suitable secret
	SecretSyncMutationLabelKey = "integrations.katanomi.dev/integration.mutation"

	// SettingsTypeLabelKey label key to select the settings secret
	SettingsTypeLabelKey = "settings.katanomi.dev/settingsType"

	// SecretSyncGeneratorLabelKey The name of the secret is the integration credentials that are automatically synced to the namespace
	SecretSyncGeneratorLabelKey = "integrations.katanomi.dev/integration.secretsync"

	// ClusterGitSourceLabelKey for cluster git source resources
	ClusterGitSourceLabelKey = "sources.katanomi.dev/clusterGitSource"

	// GitSourceLabelKey for git source resources
	GitSourceLabelKey = "sources.katanomi.dev/gitSource"
)

// Common Annotations
const (
	// DisplayNameAnnotationKey display name for objects
	DisplayNameAnnotationKey = "katanomi.dev/displayName"
	// CreatedTimeAnnotationKey creation time for objects
	CreatedTimeAnnotationKey = "katanomi.dev/creationTime"
	// UpdatedTimeAnnotationKey update time for objects
	UpdatedTimeAnnotationKey = "katanomi.dev/updateTime"
	// DeletedTimeAnnotationKey deletion time for objects
	DeletedTimeAnnotationKey = "katanomi.dev/deletionTime"
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
	// DeletedByAnnotationKey annotation key to store resource update username
	DeletedByAnnotationKey = "katanomi.dev/deletedBy"
	// CancelledByAnnotationKey annotation key to store a CancelledBy struct json with ref info
	CancelledByAnnotationKey = "katanomi.dev/cancelledBy"
	// SecretTypeAnnotationKey annotation key for an existed secret with a different type
	SecretTypeAnnotationKey = "katanomi.dev/secretType" //nolint:gosec
	// ClusterNameAnnotationKey annotation key to store resource cluster name
	ClusterNameAnnotationKey = "katanomi.dev/clusterName"
	// ClusterRefNamespaceAnnotationKey annotation key to store cluster reference namespace
	ClusterRefNamespaceAnnotationKey = "katanomi.dev/clusterRefNamespace"
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
	// ResourcePathFormat indicates project path format,
	// eg. maven project access url is /repository/maven
	// the value should be a json string like
	// {
	// 	"web-console": "/repository/%s",
	// 	"api": "/api/repository/%s",
	// }
	ResourcePathFormatAttributeKey = "resourcePathFormat"
	// SubResourcePathFormat indicates sub resource path format,
	// eg. bitbucket project access url is /scm/devops/demo
	// the value should be a json string like
	// {
	// 	"http-clone": "/scm/%s/%s",
	// 	"web-console": "/projects/%s/repo/%s",
	// }
	SubResourcePathFormatAttributeKey = "subResourcePathFormat"
	// GitPRRevisionPrefixes allows git related integrations
	// define custom PR revision prefixes
	GitPRRevisionPrefixes = "gitPRRevisionPrefixes"

	// GitPRRevisionSuffix allows git related integrations
	// define custom PR revision suffixes
	GitPRRevisionSuffixes = "gitPRRevisionSuffixes"

	//  PYPISubResourceExtendedAddressSuffixKey Add the corresponding suffix to the specified subResource which for pypi.
	// TODO: In the future, we need to adopt a better way to add
	// an extended dependency repository address to the plug-in
	PYPISubResourceExtendedAddressSuffixKey = "pypiSubResourceExtendedAddressSuffix"
)

// Attribute values for label source or manager
const (
	LabelSourceSystem = "system"
	LabelSourceUser   = "user"

	LabelKatanomi = "katanomi"
)

// Annotation keys for artifact parameter
const (
	// ArtifactAliasAnnotationKey indicates the alias of the artifact parameter in the delivery
	ArtifactAliasAnnotationKey = "alias"

	// ImageRegistryEndpoint for artifact endpoint
	ImageRegistryEndpoint = "imageRegistryEndpoint"
)

const (

	// UserOwnedAnnotationKey annotated the resource's owner is one user
	UserOwnedAnnotationKey = "katanomi.dev/owned.username" // NOSONAR // ignore: "Key" detected here, make sure this is not a hard-coded credential
)

const (
	// TrueValue represent string true
	TrueValue = "true"
	// FalseValue represent string false
	FalseValue = "false"
)
