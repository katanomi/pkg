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
	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StorageAuthCheck consists of result for an auth check request
// +k8s:deepcopy-gen=false
type StorageAuthCheck struct {
	metav1.TypeMeta `json:",inline"`
	Status          v1alpha1.AuthCheckStatus `json:"status"`
}

// StorageAuthCheckRequest is used for request entity of auth check
// use this struct rather than corev1.StoragePlugin to avoid import cycle
type StorageAuthCheckRequest struct {
	StoragePluginName string           `json:"storagePluginName"`
	Params            []v1alpha1.Param `json:"params"`
}
