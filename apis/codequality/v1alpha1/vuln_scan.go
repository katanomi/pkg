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

// VulnScanTargetType defines the type of the target to scan
type VulnScanTargetType string

const (
	// VulnScanTargetTypeImage defines the target type as image
	VulnScanTargetTypeImage VulnScanTargetType = "Image"

	// VulnScanTargetTypeFileSystem defines the target type as fs
	VulnScanTargetTypeFileSystem VulnScanTargetType = "FileSystem"

	// VulnScanTargetTypeRepository defines the target type as repository
	VulnScanTargetTypeRepository VulnScanTargetType = "Repository"
)

// VulnSeverity defines the severity of the vulnerability
type VulnSeverity = string

const (
	// VulnSeverityCritical defines the critical severity
	VulnSeverityCritical VulnSeverity = "Critical"

	// VulnSeverityHigh defines the high severity
	VulnSeverityHigh VulnSeverity = "High"

	// VulnSeverityMedium defines the medium severity
	VulnSeverityMedium VulnSeverity = "Medium"

	// VulnSeverityLow defines the low severity
	VulnSeverityLow VulnSeverity = "Low"

	// VulnSeverityUnknown defines the unknown severity
	VulnSeverityUnknown VulnSeverity = "Unknown"
)

// AvailableVulnSeverities returns the available severities
var AvailableVulnSeverities = []VulnSeverity{
	VulnSeverityCritical,
	VulnSeverityHigh,
	VulnSeverityMedium,
	VulnSeverityLow,
	VulnSeverityUnknown,
}

// NamedVulnScanResult adds name over integrated VulnScanResult
type NamedVulnScanResult struct {
	// Name of a specific lint result
	Name string `json:"name,omitempty"`

	VulnScanResult `json:",inline"`
}

// VulnScanResult stores code linting results
type VulnScanResult struct {
	// Result for the linting process
	//  - Succeeded: successful code linting with passing quality gates
	//  - Failed: failed code linting
	//  - Canceled: canceled code linting due to canceled task
	Result string `json:"result"`

	Targets []VulnScanTarget `json:"targets,omitempty"`
}

// IsEmpty returns true if the struct is empty
func (a VulnScanResult) IsEmpty() bool {
	return a.Result == "" && len(a.Targets) == 0
}

// VulnScanTarget Describe the target for vulnerability scan
type VulnScanTarget struct {
	// Uri identify of the target
	Uri string `json:"uri"`
	// Type the type of the target
	Type VulnScanTargetType `json:"type"`
	Cvss CVSS               `json:"cvss"`

	VulnStatistic `json:",inline" path:",squash"`
}

// CVSS Describes the Highest vulnerability severity
type CVSS struct {
	// Source the source of cvss score for the highest vulnerability
	Source string `json:"source"`

	// Severity the severity of the highest vulnerability
	Severity string `json:"severity"`

	// Score the score of the highest vulnerability
	Score float64 `json:"score"`
}

// VulnStatistic Describes the vulnerability statistic
type VulnStatistic struct {
	// CriticalCount Number of critical vulnerabilities
	CriticalCount int `json:"criticalCount"`

	// HighCount Number of high vulnerabilities
	HighCount int `json:"highCount"`

	// MediumCount Number of medium vulnerabilities
	MediumCount int `json:"mediumCount"`

	// LowCount Number of low vulnerabilities
	LowCount int `json:"lowCount"`

	// UnknownCount Number of unknown vulnerabilities
	UnknownCount int `json:"unknownCount"`
}
