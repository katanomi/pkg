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

//go:generate mockgen -package=types -destination=../../testing/mock/github.com/katanomi/pkg/plugin/types/test_plan.go github.com/katanomi/pkg/plugin/types TestPlanLister,TestPlanGetter

// TestPlanLister list test plans
type TestPlanLister interface {
	Interface
	ListTestPlans(ctx context.Context, params metav1alpha1.TestProjectOptions, option metav1alpha1.ListOptions) (*metav1alpha1.TestPlanList, error)
}

// TestPlanGetter get a test plan
type TestPlanGetter interface {
	Interface
	GetTestPlan(ctx context.Context, params metav1alpha1.TestProjectOptions) (*metav1alpha1.TestPlan, error)
}
