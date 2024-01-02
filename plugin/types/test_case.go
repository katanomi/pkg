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

package types

import (
	"context"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
)

// TestCaseLister list test cases
type TestCaseLister interface {
	Interface
	ListTestCases(ctx context.Context, params metav1alpha1.TestProjectOptions, options metav1alpha1.ListOptions) (*metav1alpha1.TestCaseList, error)
}

// TestCaseGetter get a test case
type TestCaseGetter interface {
	Interface
	GetTestCase(ctx context.Context, params metav1alpha1.TestProjectOptions) (*metav1alpha1.TestCase, error)
}

// TestCaseExecutionLister list test case executions
type TestCaseExecutionLister interface {
	Interface
	ListTestCaseExecutions(ctx context.Context, params metav1alpha1.TestProjectOptions, options metav1alpha1.ListOptions) (*metav1alpha1.TestCaseExecutionList, error)
}

// TestCaseExecutionCreator create a new test case execution
type TestCaseExecutionCreator interface {
	Interface
	CreateTestCaseExecution(ctx context.Context, params metav1alpha1.TestProjectOptions, payload metav1alpha1.TestCaseExecution) (*metav1alpha1.TestCaseExecution, error)
}
