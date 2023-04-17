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

import (
	authv1 "k8s.io/api/authorization/v1"
	"k8s.io/apimachinery/pkg/runtime"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// FileUploadParams Parameters required in ui upload file.
type FileUploadParams struct {
	// Address target repo's address
	Address *duckv1.Addressable `json:"address,omitempty"`

	// Type for support upload file type,ContainerImage, HelmChart, Binary
	// current only support ContainerImage.
	Type ArtifactType `json:"type,omitempty"`

	// Checksum generate unique path value
	Checksum string `json:"checksum"`

	// Properties Upload file expandable fields
	// if type is ContainerImage,The supported fields under properties is `tags`
	Properties *runtime.RawExtension `json:"properties,omitempty"`
}

// FileUploadResourceAttributes returns a ResourceAttribute object to be used in a filter
func FileUploadResourceAttributes(verb string) authv1.ResourceAttributes {
	return authv1.ResourceAttributes{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: "artifactuploads",
		Verb:     verb,
	}
}
