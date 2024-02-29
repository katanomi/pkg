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

package testing

import (
	"fmt"
	"regexp"

	. "github.com/onsi/ginkgo/v2"
)

var caseNameReg = regexp.MustCompile(`{case:(\w+)}`)
var unitTestCaseNameReg = regexp.MustCompile(`^Test(\w+)(/.*)*$`)

// Case the unique case name label
func Case(caseName string) Labels {
	if caseName == "" {
		return Labels{}
	}
	return Labels{fmt.Sprintf("{case:%s}", caseName)}
}

// GetCaseNames Resolve the case identifier from the testcase name in the junit report
// For go test junit report, the case name may be started with `Test`, e.g: TestGetProject
// For ginkgo test junit report, the case name may contain the {case:%s} string, e.g: [It] when xxxx [{case:GetProject}]
func GetCaseNames(name string) []string {
	parts := caseNameReg.FindAllStringSubmatch(name, -1)

	if len(parts) == 0 {
		parts = unitTestCaseNameReg.FindAllStringSubmatch(name, -1)
	}

	matchedNames := make([]string, 0, len(parts))
	for _, item := range parts {
		matchedNames = append(matchedNames, item[1])
	}

	return matchedNames
}
