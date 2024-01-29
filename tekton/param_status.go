/*
Copyright 2024 The Katanomi Authors.

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

package tekton

import (
	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	"knative.dev/pkg/apis"
)

// GitParameterStatuses list of all git parameter status named by paramName
// +listType=atomic
type GitParameterStatuses []GitParameterStatus

// GitParameterStatus for git parameter status
// +k8s:deepcopy-gen=true
type GitParameterStatus struct {
	// ParamName is the name of the parameter
	ParamName string `json:"paramName,omitempty"`
	// BaseGitStatus is the base git status
	v1alpha1.BaseGitStatus `json:",inline"`
	// Condition is the condition of the git parameter status
	Condition *apis.Condition `json:"condition,omitempty"`
}
