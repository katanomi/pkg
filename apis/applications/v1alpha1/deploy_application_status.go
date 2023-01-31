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

import (
	corev1 "k8s.io/api/core/v1"
)

// NamedDeployApplicationResults list of NamedDeployApplicationResult
// with helpful methods to manage data
type NamedDeployApplicationResults []NamedDeployApplicationResult

// NamedDeployApplicationResult adds a name for DeployApplicationResults
// useful for store deployment data in a list
type NamedDeployApplicationResult struct {
	// Name for the specific deployment application result
	Name string `json:"name"`

	// DeployApplicationResults result of the deployment
	DeployApplicationResults `json:",inline"`
}

// IsSameResult implements equal method for generic comparable usage
func (n NamedDeployApplicationResult) IsSameResult(y NamedDeployApplicationResult) bool {
	return n.Name == y.Name
}

// DeployApplicationResults describes a result of an application deployment
// stating the specific modified object and the modified components and versions
type DeployApplicationResults struct {
	// ApplicationRef stores the reference of the deployed application
	// as a kubernetes resource
	ApplicationRef *corev1.ObjectReference `json:"applicationRef"`

	// Before store status for the application or its components before the
	// deployment and may store previous versions and status
	Before []DeployApplicationStatus `json:"before,omitempty"`

	// After store status for the application or its components after the
	// deployment and may store the update versions and status
	After []DeployApplicationStatus `json:"after,omitempty"`
}

// DeployApplicationStatus represents the status of an Application or component
// deployed. Can be used to store status before update or after the update
type DeployApplicationStatus struct {
	// Name of the application or application subcomponent
	// if empty simbolizes the entire application update status
	Name string `json:"name"`
	// Status of the component/application. The contents are strictly tied
	// to the application/component type
	Status string `json:"status,omitempty"`
	// Version of the deployed application/component as a unique identifier
	// i.e container image url, artifact url, etc.
	Version string `json:"version,omitempty"`
}
