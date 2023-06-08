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

package conditioner

import (
	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
)

// RunConditionHandler contains RunCondition and interface to handle properties
type RunConditionHandler struct {
	v1alpha1.RunCondition
	// PropertiesHandler stores the interface that decide how to handle information in properties
	PropertiesHandler
}

// PropertiesHandler knows how to extract configuration and check the different between
// two properties
type PropertiesHandler interface {
	// Config return the configurations
	Config() string
	// Compare will compare two RawExtension and check whether they are different and return message users may interest
	Compare(old *runtime.RawExtension, new *runtime.RawExtension) (changed bool, message string, err error)
}
