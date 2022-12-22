/*
Copyright 2022 The Katanomi Authors.

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

// ArtifactType type of artifacts parameter enum
type ArtifactType string

const (
	// OCIHelmChartArtifactParameterType helm charts as OCI artifact
	// Deprecated: use pkg/apis/artifacts/v1alpha1.ArtifactTypeHelmChart instead
	OCIHelmChartArtifactParameterType ArtifactType = "OCIHelmChart"
	// OCIContainerImageArtifactParameterType runnable container image used to deploy workloads'
	// Deprecated: use pkg/apis/artifacts/v1alpha1.ArtifactTypeContainerImage instead
	OCIContainerImageArtifactParameterType ArtifactType = "OCIContainerImage"
)

// ArtifactParameterSpec specs for an strong typed parameter as an artifact
// TODO: move to pkg/apis/artifacts/v1alpha1
type ArtifactParameterSpec struct {
	// URI for artifact, must be a complete identifier, i.e docker.io/katanomi/repository
	// +optional
	URI string `json:"uri,omitempty"`

	// Type of artifact to be expected in this parameter
	// +optional
	Type ArtifactType `json:"type,omitempty"`

	// Annotations for the artifact.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// IntegrationClassName is instance name of IntegrationClass.
	// +optional
	IntegrationClassName string `json:"integrationClassName,omitempty"`
}

// NamedValue can use the NamedValue structure to set some special parameters in the artifact.
// i.e artifact promotion use NamedValue to record artifact detail info.
// TODO: move to pkg/apis/artifacts/v1alpha1
type NamedValue struct {
	// Name parameter name.
	// +optional
	Name string `json:"name,omitempty"`

	// Value The specific value of name, you can get the corresponding value according to the name.
	// +optional
	Value string `json:"value,omitempty"`
}
