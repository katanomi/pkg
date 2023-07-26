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
	"fmt"

	authv1 "k8s.io/api/authorization/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// RecordList is a list of Record
type RecordList struct {
	metav1.TypeMeta `json:",inline"`

	Items []Record `json:"items"`
}

// Record describe the archive record of a resource
type Record struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec the spec of the resource
	Spec RecordSpec `json:"spec"`

	// RelatedRecords the related records of the resource
	RelatedRecords []Record `json:"relatedRecords,omitempty"`
}

// NamespaceName returns the namespace/name of the record
func (p *Record) NamespaceName() string {
	return fmt.Sprintf("%s/%s", p.Spec.Namespace, p.Spec.Name)
}

// RecordSpec describe the spec of an archive record
type RecordSpec struct {
	// ID the unique id of the resource in the archive storage
	ID uint `json:"id,omitempty"`
	// TopCluster the cluster name of the top level resource
	TopCluster string `json:"topCluster,omitempty"`
	// TopNamespace the namespace name of the top level resource
	TopNamespace string `json:"topNamespace,omitempty"`
	// TopName the name of the top level resource
	TopName string `json:"topName,omitempty"`
	// TopUID the uid of the top level resource
	TopUID string `json:"topUid,omitempty"`
	// ParentCluster the cluster name of the parent resource
	ParentCluster string `json:"parentCluster,omitempty"`
	// ParentNamespace the namespace name of the parent resource
	ParentNamespace string `json:"parentNamespace,omitempty"`
	// ParentName the name of the parent resource
	ParentName string `json:"parentName,omitempty"`
	// ParentUID the uid of the parent resource
	ParentUID string `json:"parentUid,omitempty"`
	// TopOwnerCluster the cluster name of the top owner resource
	TopOwnerCluster string `json:"topOwnerCluster,omitempty"`
	// TopOwnerNamespace the namespace name of the top owner resource
	TopOwnerNamespace string `json:"topOwnerNamespace,omitempty"`
	// TopOwnerName the name of the top owner resource
	TopOwnerName string `json:"topOwnerName,omitempty"`
	// TopOwnerUID the uid of the top owner resource
	TopOwnerUID string `json:"topOwnerUid,omitempty"`
	// Cluster the cluster name of the resource
	Cluster string `json:"cluster,omitempty"`
	// Namespace the namespace name of the resource
	Namespace string `json:"namespace"`
	// Name the name of the resource
	Name string `json:"name"`
	// UID the uid of the resource
	UID string `json:"uid"`
	// Group the group of the resource
	Group string `json:"group"`
	// Version the version of the resource
	Version string `json:"version"`
	// Kind the kind of the resource
	Kind string `json:"kind"`
	// CreationTimestamp the creation timestamp of the resource
	CreationTimestamp int64 `json:"creationTimestamp,omitempty"`
	// CleanupTime the cleanup time of the resource
	CleanupTime int64 `json:"cleanupTime,omitempty"`

	// Metadata the metadata of the resource
	Metadata map[string]string `json:"metadata"`
	// Data original data of the resource
	Data map[string]interface{} `json:"data,omitempty"`
}

func (p *RecordSpec) DeepCopy() *RecordSpec {
	spec := *p
	spec.Data = runtime.DeepCopyJSON(p.Data)
	spec.Metadata = make(map[string]string, len(p.Metadata))
	for k, v := range p.Metadata {
		spec.Metadata[k] = v
	}
	return &spec
}

// RecordResourceAttributes returns a ResourceAttribute object to be used in a filter
func RecordResourceAttributes(verb string) authv1.ResourceAttributes {
	return authv1.ResourceAttributes{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: "records",
		Verb:     verb,
	}
}
