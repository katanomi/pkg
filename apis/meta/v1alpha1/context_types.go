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

import (
	corev1 "k8s.io/api/core/v1"
)

// EnvironmentSpec describes a cluster/namespace environment
// for multi-cluster stage support
type EnvironmentSpec struct {
	// ClusteRef stores a Cluster object reference
	// currently only supports ClusterRegistry
	// +optional
	ClusterRef *corev1.ObjectReference `json:"clusterRef,omitempty"`

	// NamespaceRef defines the target namespace to run the Stage
	// and get StageRef and other reference objects from
	// If a ClusterRef is provided will use the cluster
	// and do cross-cluster access
	NamespaceRef *corev1.LocalObjectReference `json:"namespaceRef,omitempty"`
}
