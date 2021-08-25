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
	"testing"

	// ktesting "github.com/katanomi/pkg/testing"
	. "github.com/onsi/gomega"
	// corev1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

const (
	SomeCondition apis.ConditionType = "SomeCondition"
)

var stageTestCondSet = apis.NewLivingConditionSet(
	SomeCondition,
)

type StatusTest struct {
	duckv1.Status
}

func (s *StatusTest) GetCondition(t apis.ConditionType) *apis.Condition {
	return stageTestCondSet.Manage(s).GetCondition(t)
}

func TestIsConditionChanged(t *testing.T) {
	table := map[string]struct {
		current       apis.ConditionAccessor
		old           apis.ConditionAccessor
		conditionType apis.ConditionType
		expected      bool
	}{
		"Condition is nil": {
			current:       &StatusTest{},
			old:           &StatusTest{},
			conditionType: SomeCondition,
			expected:      false,
		},
		"Condition is the same": {
			current: &StatusTest{duckv1.Status{Conditions: duckv1.Conditions{
				apis.Condition{Type: SomeCondition, Status: corev1.ConditionTrue},
			}}},
			old: &StatusTest{duckv1.Status{Conditions: duckv1.Conditions{
				apis.Condition{Type: SomeCondition, Status: corev1.ConditionTrue},
			}}},
			conditionType: SomeCondition,
			expected:      false,
		},
		"Condition is nil on old": {
			current: &StatusTest{duckv1.Status{Conditions: duckv1.Conditions{
				apis.Condition{Type: SomeCondition, Status: corev1.ConditionTrue},
			}}},
			old:           &StatusTest{},
			conditionType: SomeCondition,
			expected:      true,
		},
		"Condition changed": {
			current: &StatusTest{duckv1.Status{Conditions: duckv1.Conditions{
				apis.Condition{Type: SomeCondition, Status: corev1.ConditionTrue},
			}}},
			old: &StatusTest{duckv1.Status{Conditions: duckv1.Conditions{
				apis.Condition{Type: SomeCondition, Status: corev1.ConditionFalse},
			}}},
			conditionType: SomeCondition,
			expected:      true,
		},
	}

	for name, test := range table {
		t.Run(name, func(t *testing.T) {
			result := IsConditionChanged(test.current, test.old, test.conditionType)
			g := NewGomegaWithT(t)
			g.Expect(result).To(Equal(test.expected))
		})
	}

}
