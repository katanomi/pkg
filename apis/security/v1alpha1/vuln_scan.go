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
	"strconv"
	"strings"
)

// VulnScanTargetType defines the type of the target to scan
type VulnScanTargetType string

const (
	// VulnScanTargetTypeImage defines the target type as image
	VulnScanTargetTypeImage VulnScanTargetType = "ContainerImage"

	// VulnScanTargetTypeFileSystem defines the target type as fs
	VulnScanTargetTypeFileSystem VulnScanTargetType = "FileSystem"

	// VulnScanTargetTypeRepository defines the target type as repository
	VulnScanTargetTypeRepository VulnScanTargetType = "GitRepository"
)

// VulnSeverity defines the severity of the vulnerability
type VulnSeverity string

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
func (v VulnScanResult) IsEmpty() bool {
	return v.Result == "" && len(v.Targets) == 0
}

// ToVulnScanResultShadow convert VulnScanResult to VulnScanResultShadow
func (v VulnScanResult) ToVulnScanResultShadow() VulnScanResultShadow {
	targets := make([]VulnScanTargetShadow, 0, len(v.Targets))
	for _, item := range v.Targets {
		targets = append(targets, item.ToVulnScanTargetShadow())
	}
	return VulnScanResultShadow{
		Result:  v.Result,
		Targets: targets,
	}
}

// VulnScanResultShadow stores code linting results
type VulnScanResultShadow struct {
	// Result for the linting process
	//  - Succeeded: successful code linting with passing quality gates
	//  - Failed: failed code linting
	//  - Canceled: canceled code linting due to canceled task
	Result string `json:"result"`

	Targets []VulnScanTargetShadow `json:"targets,omitempty"`
}

// ToVulnScanResult convert VulnScanResultShadow to VulnScanResult
func (v *VulnScanResultShadow) ToVulnScanResult() VulnScanResult {
	targets := make([]VulnScanTarget, 0, len(v.Targets))
	for _, item := range v.Targets {
		targets = append(targets, item.ToVulnScanTarget())
	}
	return VulnScanResult{
		Result:  v.Result,
		Targets: targets,
	}
}

// VulnScanTargetShadow Describe the target for vulnerability scan
type VulnScanTargetShadow struct {
	// Uri identify of the target
	Uri string `json:"uri"`
	// Type the type of the target
	Type VulnScanTargetType `json:"type"`
	Cvss CVSS               `json:"cvss"`

	// Compress multiple metrics into a single field.
	// because the tekton result has a limit on length
	Statistic string `json:"statistic"`
}

// ToVulnScanTarget convert VulnScanTargetShadow to VulnScanTarget
func (v *VulnScanTargetShadow) ToVulnScanTarget() VulnScanTarget {
	list := strings.Split(v.Statistic, ",")
	counts := make([]int, len(list))
	for index, item := range list {
		counts[index], _ = strconv.Atoi(item)
	}
	target := VulnScanTarget{
		Uri:  v.Uri,
		Type: v.Type,
		Cvss: v.Cvss,
	}
	if len(counts) > 0 {
		target.CriticalCount = counts[0]
	}
	if len(counts) > 1 {
		target.HighCount = counts[1]
	}
	if len(counts) > 2 {
		target.MediumCount = counts[2]
	}
	if len(counts) > 3 {
		target.LowCount = counts[3]
	}
	if len(counts) > 4 {
		target.UnknownCount = counts[4]
	}
	return target
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

// ToVulnScanTargetShadow convert VulnScanTarget to VulnScanTargetShadow
func (v *VulnScanTarget) ToVulnScanTargetShadow() VulnScanTargetShadow {
	statistic := fmt.Sprintf("%d,%d,%d,%d,%d",
		v.CriticalCount, v.HighCount, v.MediumCount, v.LowCount, v.UnknownCount,
	)
	return VulnScanTargetShadow{
		Uri:       v.Uri,
		Type:      v.Type,
		Cvss:      v.Cvss,
		Statistic: statistic,
	}
}

// CVSS Describe the vulnerability with the highest severity
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
	// CriticalCount Count of critical severity vulnerabilities
	CriticalCount int `json:"criticalCount"`

	// HighCount Count of high severity vulnerabilities
	HighCount int `json:"highCount"`

	// MediumCount Count of medium severity vulnerabilities
	MediumCount int `json:"mediumCount"`

	// LowCount Count of low severity vulnerabilities
	LowCount int `json:"lowCount"`

	// UnknownCount Count of unknown severity vulnerabilities
	UnknownCount int `json:"unknownCount"`
}
