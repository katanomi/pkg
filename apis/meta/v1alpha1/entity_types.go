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
	pipev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	authv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var EntityGVK = GroupVersion.WithKind("Entity")
var EntityListGVK = GroupVersion.WithKind("EntityList")

// Entity object for sources
type Entity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec EntitySpec `json:"spec"`
}

// EntitySpec spec for Entity
type EntitySpec struct {
	// Resolver is the type of resolver to which the entity belongs
	Resolver string `json:"resolver,omitempty"`

	// Params are the parameters required to use the entity.
	Params []pipev1beta1.Param `json:"params,omitempty"`

	// Raw is original content of the entity
	Raw []byte `json:"raw,omitempty"`

	// metadata of the original resource in the raw
	// Easy to filter resources at the front end
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

// EntityResourceAttributes returns a ResourceAttribute object to be used in a filter
func EntityResourceAttributes(verb string) authv1.ResourceAttributes {
	return authv1.ResourceAttributes{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: "entities",
		Verb:     verb,
	}
}
