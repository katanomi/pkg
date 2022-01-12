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

	// IntegrationAddressAnnotation annotation key to store integration server address
	IntegrationAddressAnnotation = "integrations.katanomi.dev/integration.address"

	// IntegrationResourceScope annotation key to store integration resource scope
	IntegrationResourceScope = "integrations.katanomi.dev/integration.resourceScope"
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
	// UpdatedByAnnotationKey annotation key annotation key
	UpdatedByAnnotationKey = "katanomi.dev/updatedBy"
	// SecretTypeAnnotationKey annotation key for an existed secret with a different type
	SecretTypeAnnotationKey = "katanomi.dev/secretType" //nolint:gosec
	//ClusterNameAnnotationKey annotation key to store resource cluster name
	ClusterNameAnnotationKey = "integrations.katanomi.dev/clusterName"
)

// Attribute keys for Integrations
const (
	// Keys used in IntegrationClass.status.attributes
	AuthAttributeKey                   = "auth"
	ReplicationPolicyTypesAttributeKey = "replicationPolicyTypes"
	ResourceTypesAttributeKey          = "resourceTypes"
	MethodsAttributeKey                = "methods"
	AllowEmptySecretAttributeKey       = "allowEmptySecret"
	DefaultProjectTypeAttributeKey     = "defaultProjectSubType"
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
)
