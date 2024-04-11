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
	"fmt"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"

	// ktesting "github.com/katanomi/pkg/testing"
	. "github.com/onsi/gomega"
	// corev1 "k8s.io/api/core/v1"
	apismock "github.com/katanomi/pkg/testing/mock/knative.dev/pkg/apis"
	"go.uber.org/mock/gomock"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

//go:generate mockgen -package=apis -destination=../../../testing/mock/knative.dev/pkg/apis/condition_manager.go  knative.dev/pkg/apis ConditionManager

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

func TestGetCondition(t *testing.T) {
	table := map[string]struct {
		source        apis.Conditions
		t             apis.ConditionType
		expectedIndex int
	}{
		"length is 0":   {source: apis.Conditions{}, t: apis.ConditionType("WHAT"), expectedIndex: -1},
		"source is nil": {source: nil, t: apis.ConditionType("WHAT"), expectedIndex: -1},
		"contains in source": {
			source: apis.Conditions{
				{
					Type: "AType",
				},
				{
					Type: "BType",
				},
			},
			t:             "AType",
			expectedIndex: 0,
		},
		"not contains in source": {
			source: apis.Conditions{
				{
					Type: "AType",
				},
				{
					Type: "BType",
				},
			},
			t:             "CType",
			expectedIndex: -1,
		},
		"empty type": {
			source: apis.Conditions{
				{
					Type: "AType",
				},
				{
					Type: "BType",
				},
			},
			t:             "",
			expectedIndex: -1,
		},
	}

	for name, item := range table {
		t.Run(name, func(t *testing.T) {
			actual := GetCondition(item.source, item.t)
			g := NewGomegaWithT(t)
			if item.expectedIndex < 0 {
				g.Expect(actual).Should(BeNil())
			} else {
				g.Expect(fmt.Sprintf("%p", actual)).Should(BeEquivalentTo(fmt.Sprintf("%p", &item.source[item.expectedIndex])))
				actual.Message = "changed"
				g.Expect(item.source[item.expectedIndex].Message).Should(BeEquivalentTo(actual.Message))
			}
		})
	}
}

func TestIsConditionChanged(t *testing.T) {
	now := metav1.Now()
	oneSecondAgo := metav1.NewTime(now.Add(-time.Second))
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
		"Condition status is the same, last transition changed": {
			current: &StatusTest{duckv1.Status{Conditions: duckv1.Conditions{
				apis.Condition{Type: SomeCondition, Status: corev1.ConditionUnknown, LastTransitionTime: apis.VolatileTime{Inner: now}},
			}}},
			old: &StatusTest{duckv1.Status{Conditions: duckv1.Conditions{
				apis.Condition{Type: SomeCondition, Status: corev1.ConditionUnknown, LastTransitionTime: apis.VolatileTime{Inner: oneSecondAgo}},
			}}},
			conditionType: SomeCondition,
			expected:      false,
		},
		"Condition status is the same, last transition same, reason changed": {
			current: &StatusTest{duckv1.Status{Conditions: duckv1.Conditions{
				apis.Condition{Type: SomeCondition, Status: corev1.ConditionUnknown, LastTransitionTime: apis.VolatileTime{Inner: now}, Reason: "abc"},
			}}},
			old: &StatusTest{duckv1.Status{Conditions: duckv1.Conditions{
				apis.Condition{Type: SomeCondition, Status: corev1.ConditionUnknown, LastTransitionTime: apis.VolatileTime{Inner: now}, Reason: "def"},
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

func TestSetConditonByErrorReason(t *testing.T) {
	t.Run("Error with empty reason", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conditionManager := apismock.NewMockConditionManager(ctrl)

		err := errors.NewBadRequest("some reason")

		conditionManager.EXPECT().GetCondition(apis.ConditionReady).Return(nil)

		conditionManager.EXPECT().
			MarkFalse(apis.ConditionReady, string(metav1.StatusReasonBadRequest), err.Error()).
			Times(1)

		SetConditionByErrorReason(conditionManager, apis.ConditionReady, err, "")
	})

	t.Run("Error with specific reason", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conditionManager := apismock.NewMockConditionManager(ctrl)

		err := errors.NewBadRequest("some reason")
		reason := "another reason"

		conditionManager.EXPECT().GetCondition(apis.ConditionReady).Return(nil)

		conditionManager.EXPECT().
			MarkFalse(apis.ConditionReady, reason, err.Error()).
			Times(1)

		SetConditionByErrorReason(conditionManager, apis.ConditionReady, err, reason)
	})
}

func TestSetConditonByError(t *testing.T) {

	t.Run("Explicit reason error case", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conditionManager := apismock.NewMockConditionManager(ctrl)

		err := errors.NewBadRequest("some reason")

		conditionManager.EXPECT().GetCondition(apis.ConditionReady).Return(nil)

		conditionManager.EXPECT().
			MarkFalse(apis.ConditionReady, string(metav1.StatusReasonBadRequest), err.Error()).
			Times(1)

		SetConditionByError(conditionManager, apis.ConditionReady, err)
	})

	t.Run("Error condition reason changed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conditionManager := apismock.NewMockConditionManager(ctrl)

		err := errors.NewBadRequest("some reason")

		conditionManager.EXPECT().GetCondition(apis.ConditionReady).Return(&apis.Condition{
			Type:    apis.ConditionReady,
			Status:  corev1.ConditionFalse,
			Reason:  "another reason",
			Message: err.Error(),
		})

		conditionManager.EXPECT().
			MarkFalse(apis.ConditionReady, string(metav1.StatusReasonBadRequest), err.Error()).
			Times(1)

		SetConditionByError(conditionManager, apis.ConditionReady, err)
	})

	t.Run("Error condition message changed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conditionManager := apismock.NewMockConditionManager(ctrl)

		err := errors.NewBadRequest("some reason")

		conditionManager.EXPECT().GetCondition(apis.ConditionReady).Return(&apis.Condition{
			Type:    apis.ConditionReady,
			Status:  corev1.ConditionFalse,
			Reason:  string(metav1.StatusReasonBadRequest),
			Message: "some message",
		})

		conditionManager.EXPECT().
			MarkFalse(apis.ConditionReady, string(metav1.StatusReasonBadRequest), err.Error()).
			Times(1)

		SetConditionByError(conditionManager, apis.ConditionReady, err)
	})

	t.Run("Error condition not changed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conditionManager := apismock.NewMockConditionManager(ctrl)

		err := errors.NewBadRequest("some reason")

		conditionManager.EXPECT().GetCondition(apis.ConditionReady).Return(&apis.Condition{
			Type:    apis.ConditionReady,
			Status:  corev1.ConditionFalse,
			Reason:  string(metav1.StatusReasonBadRequest),
			Message: err.Error(),
		})

		conditionManager.EXPECT().
			MarkFalse(apis.ConditionReady, string(metav1.StatusReasonBadRequest), err.Error()).
			Times(0)

		SetConditionByError(conditionManager, apis.ConditionReady, err)
	})

	t.Run("Success condition change to error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conditionManager := apismock.NewMockConditionManager(ctrl)

		err := errors.NewBadRequest("some reason")

		conditionManager.EXPECT().GetCondition(apis.ConditionReady).Return(&apis.Condition{
			Type:   apis.ConditionReady,
			Status: corev1.ConditionTrue,
		})

		conditionManager.EXPECT().
			MarkFalse(apis.ConditionReady, string(metav1.StatusReasonBadRequest), err.Error()).
			Times(1)

		SetConditionByError(conditionManager, apis.ConditionReady, err)
	})

	t.Run("random error case", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conditionManager := apismock.NewMockConditionManager(ctrl)

		err := fmt.Errorf("some random error")

		conditionManager.EXPECT().GetCondition(apis.ConditionSucceeded).Return(nil)

		conditionManager.EXPECT().
			MarkFalse(apis.ConditionSucceeded, "", err.Error()).
			Times(1)

		SetConditionByError(conditionManager, apis.ConditionSucceeded, err)
	})

	t.Run("successful case", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conditionManager := apismock.NewMockConditionManager(ctrl)

		conditionManager.EXPECT().GetCondition(apis.ConditionSucceeded).Return(nil)

		conditionManager.EXPECT().
			MarkTrue(apis.ConditionSucceeded).
			Times(1)

		SetConditionByError(conditionManager, apis.ConditionSucceeded, nil)
	})

	t.Run("err is nil and from unknown status", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conditionManager := apismock.NewMockConditionManager(ctrl)

		conditionManager.EXPECT().GetCondition(apis.ConditionSucceeded).Return(&apis.Condition{
			Type:   apis.ConditionSucceeded,
			Status: corev1.ConditionUnknown,
		})

		conditionManager.EXPECT().
			MarkTrue(apis.ConditionSucceeded).
			Times(1)

		SetConditionByError(conditionManager, apis.ConditionSucceeded, nil)
	})

	t.Run("err is not nil and from unknown status", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conditionManager := apismock.NewMockConditionManager(ctrl)

		err := errors.NewBadRequest("some reason")

		conditionManager.EXPECT().GetCondition(apis.ConditionSucceeded).Return(&apis.Condition{
			Type:    apis.ConditionSucceeded,
			Status:  corev1.ConditionUnknown,
			Message: err.Error(),
			Reason:  string(metav1.StatusReasonBadRequest),
		})

		conditionManager.EXPECT().
			MarkFalse(apis.ConditionSucceeded, string(metav1.StatusReasonBadRequest), err.Error()).
			Times(1)

		SetConditionByError(conditionManager, apis.ConditionSucceeded, err)
	})

	t.Run("Success condition not changed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conditionManager := apismock.NewMockConditionManager(ctrl)

		conditionManager.EXPECT().GetCondition(apis.ConditionSucceeded).Return(&apis.Condition{
			Type:   apis.ConditionSucceeded,
			Status: corev1.ConditionTrue,
		})

		conditionManager.EXPECT().
			MarkTrue(apis.ConditionSucceeded).
			Times(0)

		SetConditionByError(conditionManager, apis.ConditionSucceeded, nil)
	})

	t.Run("Error condition change to success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conditionManager := apismock.NewMockConditionManager(ctrl)

		conditionManager.EXPECT().GetCondition(apis.ConditionSucceeded).Return(&apis.Condition{
			Type:   apis.ConditionSucceeded,
			Status: corev1.ConditionFalse,
		})

		conditionManager.EXPECT().
			MarkTrue(apis.ConditionSucceeded).
			Times(1)

		SetConditionByError(conditionManager, apis.ConditionSucceeded, nil)
	})

}

func TestPropagateCondition(t *testing.T) {

	condType := apis.ConditionReady

	t.Run("Explicit reason error case", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conditionManager := apismock.NewMockConditionManager(ctrl)

		condition := &apis.Condition{
			Type:    condType,
			Reason:  "SomeReason",
			Status:  corev1.ConditionFalse,
			Message: "some message",
		}

		conditionManager.EXPECT().
			MarkFalse(condType, "SomeReason", "some message").
			Times(1)

		PropagateCondition(conditionManager, condType, condition)
	})

	t.Run("nil condition", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conditionManager := apismock.NewMockConditionManager(ctrl)

		conditionManager.EXPECT().
			MarkUnknown(condType, ConditionReasonNotSet, "condition is empty").
			Times(1)

		PropagateCondition(conditionManager, condType, nil)
	})

	t.Run("condition is unknown", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conditionManager := apismock.NewMockConditionManager(ctrl)

		condition := &apis.Condition{
			Type:    condType,
			Reason:  "abc",
			Status:  corev1.ConditionUnknown,
			Message: "msg",
		}

		conditionManager.EXPECT().
			MarkUnknown(condType, "abc", "msg").
			Times(1)

		PropagateCondition(conditionManager, condType, condition)
	})

	t.Run("condition is true", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conditionManager := apismock.NewMockConditionManager(ctrl)

		condition := &apis.Condition{
			Type:    condType,
			Reason:  "abc",
			Status:  corev1.ConditionTrue,
			Message: "msg",
		}

		conditionManager.EXPECT().
			MarkTrueWithReason(condType, "abc", "msg").
			Times(1)

		PropagateCondition(conditionManager, condType, condition)
	})
}
