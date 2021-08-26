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
	// SecretLabelKey secret resource name
	SecretLabelKey = "core.kubernetes.io/secret" //nolint:gosec
	// NamespaceLabelKey namespace of a resource
	NamespaceLabelKey = "core.kubernetes.io/namespace"

	// ClusterIntegrationLabelKey for cluster integration resources
	ClusterIntegrationLabelKey = "integrations.katanomi.dev/clusterIntegration"
	// IntegrationLabelKey for integration resources
	IntegrationLabelKey = "integrations.katanomi.dev/integration"
	// ProjectLabelKey for integration resources
	ProjectLabelKey = "integrations.katanomi.dev/project"
	// RepositoryLabelKey for integration resources
	RepositoryLabelKey = "integrations.katanomi.dev/repository"
)

// Common Annotations
const (
	// CreatedTimeAnnotationKey creation time for objects
	CreatedTimeAnnotationKey = "katanomi.dev/creationTime"
	// CrossClusterAnnotationKey annotates a cross cluster resource/action
	CrossClusterAnnotationKey = "katanomi.dev/crossCluster"
	// ReconcileTriggeredAnnotationKey annotation key to trigger reconcile of objects
	ReconcileTriggeredAnnotationKey = "katanomi.dev/reconcileTriggeredOn"
)

// Attribute keys for Integrations
const (
	// Keys used in IntegrationClass.status.attributes
	AuthAttributeKey                   = "auth"
	ReplicationPolicyTypesAttributeKey = "replicationPolicyTypes"
	ResourceTypesAttributeKey          = "resourceTypes"
	MethodsAttributeKey                = "methods"
)
