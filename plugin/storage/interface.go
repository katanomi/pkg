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

package storage

import (
	"github.com/katanomi/pkg/plugin/client"
	"github.com/katanomi/pkg/plugin/route"
	corev1alpha1 "github.com/katanomi/pkg/plugin/storage/core/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// VersionedRouter defines route with group version method
type VersionedRouter interface {
	route.ContextRoute
	GroupVersion() schema.GroupVersion
}

// PluginRegister provides storage plugin registration methods to update StoragePluginClass and StoragePlugin status
type PluginRegister interface {
	client.Interface
	client.LivenessChecker
	client.Initializer
	client.DependentResourceGetter
	client.PluginAddressable

	// core interface must be implemented, will be served
	corev1alpha1.CoreInterface

	// GetStoragePluginClassName returns storage plugin class name
	GetStoragePluginClassName() string
}
