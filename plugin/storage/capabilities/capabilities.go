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

package capabilities

import (
	"reflect"

	archivev1alpha1 "github.com/katanomi/pkg/plugin/storage/capabilities/archive/v1alpha1"
	filestorev1alpha1 "github.com/katanomi/pkg/plugin/storage/capabilities/filestore/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// reflectedCapElmMap for reflection in init func
var reflectedCapElmMap map[string]reflect.Type

func init() {
	reflectedCapElmMap = make(map[string]reflect.Type)
	for capability, intf := range RegisteredCapabilities {
		iElem := reflect.TypeOf(intf).Elem()
		reflectedCapElmMap[capability.String()] = iElem
	}
}

// RegisteredCapabilities declares which interface a versioned capability should implement.
var RegisteredCapabilities = map[schema.GroupVersion]interface{}{
	filestorev1alpha1.FileStoreV1alpha1GV: (*filestorev1alpha1.FileStoreCapable)(nil),
	archivev1alpha1.ArchiveV1alpha1:       (*archivev1alpha1.ArchiveCapable)(nil),
}

// GetImplementedCapabilities returns string list of capabilities an object implemented
func GetImplementedCapabilities(obj interface{}) []string {
	var capabilities []string
	if obj == nil {
		return nil
	}
	typeOfObj := reflect.TypeOf(obj)
	for capability, iElem := range reflectedCapElmMap {
		// TODO: to cover non-pointer receiver
		if typeOfObj.Implements(iElem) {
			capabilities = append(capabilities, capability)
		}
	}

	return capabilities
}
