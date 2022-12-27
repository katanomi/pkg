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

type ArtifactType string

const (
	// ArtifactTypeContainerImage artifact type as container image
	ArtifactTypeContainerImage              ArtifactType = "ContainerImage"
	DeprecatedArtifactTypeOCIContainerImage              = "OCIContainerImage"

	// ArtifactTypeHelmChart artifact type helm chart
	ArtifactTypeHelmChart              ArtifactType = "HelmChart"
	DeprecatedArtifactTypeOCIHelmChart              = "OCIHelmChart"

	// ArtifactTypeBinary binary artifact
	ArtifactTypeBinary ArtifactType = "Binary"
	// ArtifactTypeMaven maven artifact
	ArtifactTypeMaven ArtifactType = "Maven"

	// OCIHelmMediaType media type used for OCI helm chart artifact
	OCIHelmMediaType = "application/vnd.cncf.helm.config.v1+json"
	// OCIHelmChartContentType content type for a chart tar file when using OCI
	// registry as storage
	OCIHelmChartContentType = "application/vnd.cncf.helm.chart.content.v1.tar+gzip"

	// HelmChartDigestAnnotationKey annotation key used to store a digest
	// generated by pkg/hash.HashFolder method and used to compare content equality
	HelmChartDigestAnnotationKey = "digest.katanomi.dev/chart"
)
