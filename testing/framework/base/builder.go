/*
Copyright 2021 The AlaudaDevops Authors.

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

package base

import (
	"fmt"

	"github.com/AlaudaDevops/pkg/testing"
)

// TestCaseLabel label for test case
type TestCaseLabel = string

// TestCasePriority priority for the testcase
type TestCasePriority uint16

const (
	// P0 critical priority test case
	P0 TestCasePriority = 0
	// P1 high priority test case
	P1 TestCasePriority = 1
	// P2 medium priority test case
	P2 TestCasePriority = 2
	// P3 low priority test case
	P3 TestCasePriority = 3
)

// TestCaseBuilder builder for TestCases
// helps provide methods to construct
type TestCaseBuilder struct {
	TestContextGetter

	// Name of the test case
	Name string

	// Priority of the test case
	Priority TestCasePriority

	// Scope defines what kind of permissions this test case needs
	// Labels used to filter test cases when executing testing
	Labels []string

	// Conditions condition list which will be checked before test case execution
	Conditions []Condition

	// FailedWhenConditionMismatch allow skip test case when test condition check failed
	// default to skip
	FailedWhenConditionMismatch bool

	// TestSpec the spec of the test case
	TestSpec TestSpecFunc
}

// CheckCondition check test case condition
func (b *TestCaseBuilder) CheckCondition(testCtx *TestContext) (skip bool, err error) {
	for _, condition := range b.Conditions {
		if condition == nil {
			continue
		}
		if err = condition.Condition(testCtx); err != nil {
			err = fmt.Errorf("condition %s check failed: %w", testing.ReflectName(condition), err)
			break
		}
	}
	if err != nil && !b.FailedWhenConditionMismatch {
		skip = true
	}

	return
}

// CaseName returns the formatted name of the test case
func (b *TestCaseBuilder) CaseName() string {
	return fmt.Sprintf("[P%d][%s]", b.Priority, b.Name)
}
