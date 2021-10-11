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

// storage resource type
type StorageResourceType string

const (
	// OCI resource type
	ResourceTypeOCI StorageResourceType = "OCI"
)

func (s StorageResourceType) String() string {
	return string(s)
}

// Payload type
type PayloadType string

const (
	// build type
	PayloadTypeBuild PayloadType = "build"

	// deploy type
	PayloadTypeDeploy PayloadType = "deploy"
)

func (p PayloadType) String() string {
	return string(p)
}

// backend storage type
type BackendType string

const (
	// pv type
	BackendTypePV BackendType = "pv"

	// memory type
	BackendTypeMemory BackendType = "memory"
)
