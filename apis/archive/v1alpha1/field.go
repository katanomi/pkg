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

const (
	// TopClusterField the name of top cluster field
	TopClusterField = "topCluster"
	// ParentClusterField the name of parent cluster field
	ParentClusterField = "parentCluster"
	// ClusterField the name of cluster field
	ClusterField = "cluster"
	// TopUIDField the name of top uid field
	TopUIDField = "topUid"
	// ParentUIDField the name of parent uid field
	ParentUIDField = "parentUid"
	// UIDField the name of uid field
	UIDField = "uid"
	// TopNamespaceField the name of top namespace field
	TopNamespaceField = "topnamespace"
	// ParentNamespaceField the name of parent namespace field
	ParentNamespaceField = "parentNamespace"
	// NamespaceField the name of namespace field
	NamespaceField = "namespace"
	// TopNameField the name of top name field
	TopNameField = "topName"
	// ParentNameField the name of parent name field
	ParentNameField = "parentName"
	// NameField the name of name field
	NameField = "name"
	// GroupField the name of group field
	GroupField = "group"
	// VersionField the name of version field
	VersionField = "version"
	// KindField the name of kind field
	KindField = "kind"
	// MetadataFiled the name of metadata field
	MetadataFiled = "metadata"
	// DataField the name of data field
	DataField = "data"
	// CreationTimestampField the name of creation timestamp field
	CreationTimestampField = "creationTimestamp"
	// CleanupTimeField the name of deletion timestamp field
	CleanupTimeField = "cleanupAt"
)

// MetadataKey returns the metadata key
func MetadataKey(key string) string {
	return MetadataFiled + "." + key
}
