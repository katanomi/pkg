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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"knative.dev/pkg/apis"
)

const (
	// key for additional information type
	StatusInfoTypePipelineKey          = "pipeline"
	StatusInfoTypeVulnerabilityScanKey = "vulnerabilityScan"
	StatusInfoTypeCodeScanKey          = "codeScan"
	StatusInfoTypeArtifactKey          = "artifact"
	StatusInfoTypeUnittestKey          = "unittest"
	StatusInfoTypeVersionKey           = "version"

	// vulnerability scan severity level key
	VulnerabilityScanNoneLevelKey       = "None"
	VulnerabilityScanNegligibleLevelKey = "Negligible"
	VulnerabilityScanUnknownLevelKey    = "Unknown"
	VulnerabilityScanLowLevelKey        = "Low"
	VulnerabilityScanMediumLevelKey     = "Medium"
	VulnerabilityScanHighLevelKey       = "High"
	VulnerabilityScanCriticalLevelKey   = "Critical"
)

// BuildStateType represents a build state of the repository.
type BuildStateType string

// These constants represent all valid build states.
const (
	// Each tool needs to be converted to the corresponding value
	BuildStatePending  BuildStateType = "Pending"
	BuildStateCreated  BuildStateType = "Created"
	BuildStateRunning  BuildStateType = "Running"
	BuildStateSuccess  BuildStateType = "Success"
	BuildStateFailed   BuildStateType = "Failed"
	BuildStateCanceled BuildStateType = "Canceled"
	BuildStateSkipped  BuildStateType = "Skipped"
	BuildStateManual   BuildStateType = "Manual"
	BuildStateError    BuildStateType = "Error"
)

var (
	GitCommitStatusGVK     = GroupVersion.WithKind("GitCommitStatus")
	GitCommitStatusListGVK = GroupVersion.WithKind("GitCommitStatusList")
)

// GitCommitStatusInfo pipeline for code repository in commit statues
// https://github.com/katanomi/spec/blob/main/3.core.plugins.gitapi.md search 8.commit status
type GitCommitStatusInfo struct {
	// Type of additional information(pipeline, vulnerabilityScan, codeScan, artifact, unittest, version)
	Type string `json:"type"`
	// Status for additional information
	Status ConditionType `json:"status"`
	// URL result address
	URL        *apis.URL             `json:"url"`
	Properties *runtime.RawExtension `json:"properties,omitempty"`
}

type GitCommitStatusInfoList []GitCommitStatusInfo

type GitCommitStatusStatus struct {
	GitCommitStatusInfoList
	*apis.Condition
}

// GitCommitStatus object for plugin
type GitCommitStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GitCommitStatusSpec   `json:"spec"`
	Status GitCommitStatusStatus `json:"status"`
}

type GitCommitStatusSpec struct {
	// ID status id
	ID int `json:"id"`
	// TODO:
	// 1.SHA and CreatedAt will be removed, because they should be used in metadata.
	// 2.Status, Name, TargetURL will be removed that already have the same effect in status
	// For compatibility reasons, the above contents will be implemented in later versions.
	// SHA commit sha
	SHA string `json:"sha"`
	// Ref commit ref
	Ref string `json:"ref"`
	// Status
	Status string `json:"status"`
	// CreatedAt status create time
	CreatedAt metav1.Time `json:"createdAt"`
	// Name status name
	Name string `json:"name"`
	// Author status author
	Author GitUserBaseInfo `json:"author"`
	// Description status description
	Description string `json:"description"`
	// TargetURL
	TargetURL  string                `json:"targetUrl"`
	Properties *runtime.RawExtension `json:"properties,omitempty"`
}

// GitCommitStatusList list of commit status
type GitCommitStatusList struct {
	metav1.TypeMeta `json:",inline"`
	ListMeta        `json:"metadata,omitempty"`

	Items []GitCommitStatus `json:"items"`
}
