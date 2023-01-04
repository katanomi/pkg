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
	"strings"

	authv1 "k8s.io/api/authorization/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

const (
	// KeyForSubType subtype key in context
	KeyForSubType ContextKey = "subtype" // NOSONAR
)

// ProjectSubType stores a specific project subtype
type ProjectSubType string

func (r ProjectSubType) String() string {
	return string(r)
}

// Validate if is known types
func (r ProjectSubType) Validate(fld *field.Path) field.ErrorList {
	errs := field.ErrorList{}

	supportedTypes := map[ProjectSubType]struct{}{
		DefaultProjectSubType:              {},
		ImageRegistryProjectSubType:        {},
		GitUserProjectSubType:              {},
		GitGroupProjectSubType:             {},
		MavenRepositoryProjectSubType:      {},
		RawRepositoryProjectSubType:        {},
		ProjectManagementSubtype:           {},
		TestProjectSubType:                 {},
		MavenProxyRepositoryProjectSubType: {},
		MavenGroupRepositoryProjectSubType: {},
		RawProxyRepositoryProjectSubType:   {},
		RawGroupRepositoryProjectSubType:   {},
		NPMRepositoryProjectSubType:        {},
		NPMProxyRepositoryProjectSubType:   {},
		NPMGroupRepositoryProjectSubType:   {},
		PYPIRepositoryProjectSubType:       {},
		PYPIProxyRepositoryProjectSubType:  {},
		PYPIGroupRepositoryProjectSubType:  {},
		GoProxyRepositoryProjectSubType:    {},
		GoGroupRepositoryProjectSubType:    {},
	}

	types := strings.Split(r.String(), ",")

	for _, t := range types {
		if _, exist := supportedTypes[ProjectSubType(t)]; !exist {
			errs = append(errs, field.Invalid(fld, r, "resource subtype is invalid"))
		}
	}

	return errs
}

const (
	// DefaultProjectSubType default project subtype
	DefaultProjectSubType ProjectSubType = "Project"

	// ImageRegistryProjectSubType image registry project subtype
	ImageRegistryProjectSubType ProjectSubType = "ImageRegistry"

	// GitUserProjectSubType git user project subtype
	GitUserProjectSubType ProjectSubType = "GitUser"

	// GitGroupProjectSubType git group project subtype
	GitGroupProjectSubType ProjectSubType = "GitGroup"

	// RawRepositoryProjectSubType raw repository project subtype
	RawRepositoryProjectSubType ProjectSubType = "RawRepository"

	// MavenRepositoryProjectSubType maven repository project subtype
	MavenRepositoryProjectSubType ProjectSubType = "MavenRepository"

	// ProjectManagementSubtype project management subtype
	ProjectManagementSubtype ProjectSubType = "ProjectManagement"

	// TestProjectSubType test project subtype
	TestProjectSubType ProjectSubType = "TestProject"

	// MavenProxyRepositoryProjectSubType maven-proxy repository project subtype
	MavenProxyRepositoryProjectSubType ProjectSubType = "MavenProxyRepository"

	// MavenGroupRepositoryProjectSubType maven-group repository project subtype
	MavenGroupRepositoryProjectSubType ProjectSubType = "MavenGroupRepository"

	// RawProxyRepositoryProjectSubType raw-proxy repository project subtype
	RawProxyRepositoryProjectSubType ProjectSubType = "RawProxyRepository"

	// RawGroupRepositoryProjectSubType raw-group repository project subtype
	RawGroupRepositoryProjectSubType ProjectSubType = "RawGroupRepository"

	// NPMRepositoryProjectSubType npm repository project subtype
	NPMRepositoryProjectSubType ProjectSubType = "NPMRepository"

	// NPMProxyRepositoryProjectSubType npm-proxy repository project subtype
	NPMProxyRepositoryProjectSubType ProjectSubType = "NPMProxyRepository"

	// NPMGroupRepositoryProjectSubType npm-group repository project subtype
	NPMGroupRepositoryProjectSubType ProjectSubType = "NPMGroupRepository"

	// PYPIRepositoryProjectSubType python repository project subtype
	PYPIRepositoryProjectSubType ProjectSubType = "PYPIRepository"

	// PYPIProxyRepositoryProjectSubType python-proxy repository project subtype
	PYPIProxyRepositoryProjectSubType ProjectSubType = "PYPIProxyRepository"

	// PYPIGroupRepositoryProjectSubType python-group repository project subtype
	PYPIGroupRepositoryProjectSubType ProjectSubType = "PYPIGroupRepository"

	// GoProxyRepositoryProjectSubType go-proxy repository project subtype
	GoProxyRepositoryProjectSubType ProjectSubType = "GoProxyRepository"

	// GoGroupRepositoryProjectSubType go-group repository project subtype
	GoGroupRepositoryProjectSubType ProjectSubType = "GoGroupRepository"

	// TODO: add more subtypes
)

var ProjectGVK = GroupVersion.WithKind("Project")
var ProjectListGVK = GroupVersion.WithKind("ProjectList")

// Project object for plugins
type Project struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ProjectSpec `json:"spec"`
}

// ProjectSpec spec for project
// TODO: add more necessary spec data
type ProjectSpec struct {
	// Public defines if a project is public or not
	Public bool `json:"public"`

	// Address API related access URL
	// +optional
	Address *duckv1.Addressable `json:"address,omitempty"`

	// Access stores the webconsole address if any
	// +optional
	Access *duckv1.Addressable `json:"access,omitempty"`

	// project subtype
	// +kubebuilder:default="Project"
	SubType ProjectSubType `json:"subType"`

	// NamespaceRefs for which this project is already bound to
	// +optional
	NamespaceRefs []corev1.ObjectReference `json:"namespaceRefs,omitempty"`

	// Properties extended properties for Project
	Properties *runtime.RawExtension `json:"properties,omitempty"`
}

// ProjectList list of projects
type ProjectList struct {
	metav1.TypeMeta `json:",inline"`
	ListMeta        `json:"metadata,omitempty"`

	Items []Project `json:"items"`
}

// ProjectResourceAttributes returns a ResourceAttribute object to be used in a filter
func ProjectResourceAttributes(verb string) authv1.ResourceAttributes {
	return authv1.ResourceAttributes{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: "projects",
		Verb:     verb,
	}
}
