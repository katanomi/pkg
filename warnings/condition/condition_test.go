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

package condition

import (
	"context"
	"time"

	krecord "github.com/AlaudaDevops/pkg/record"
	conditionmock "github.com/AlaudaDevops/pkg/testing/mock/github.com/AlaudaDevops/pkg/warnings/condition"
	"github.com/AlaudaDevops/pkg/warnings"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"knative.dev/pkg/apis"
)

var _ = Describe("Test.MarkAndRecordWarning", func() {
	var (
		ctx          context.Context
		mockCtrl     *gomock.Controller
		managerMock  *conditionmock.MockWarningConditionManager
		warning      *warnings.WarningRecord
		recorder     record.EventRecorder
		fakeRecorder *record.FakeRecorder
		c            *apis.Condition
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockCtrl = gomock.NewController(GinkgoT())
		managerMock = conditionmock.NewMockWarningConditionManager(mockCtrl)

		recorder = record.NewFakeRecorder(100)
		fakeRecorder, _ = recorder.(*record.FakeRecorder)
		ctx = krecord.WithRecorder(ctx, fakeRecorder)

		warning = &warnings.WarningRecord{
			Reason:  "reason",
			Message: "message",
		}
		c = &apis.Condition{
			Type:    "type",
			Reason:  warning.Reason,
			Message: warning.Message,
		}

		managerMock.EXPECT().GetObject().Return(&corev1.Pod{}).AnyTimes()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	JustBeforeEach(func() {
		MarkAndRecordWarning(ctx, managerMock, warning)
	})

	When("warning is nil", func() {
		BeforeEach(func() {
			warning = nil
			managerMock.EXPECT().GetWarningCondition().Times(0)
		})
		It("should do nothing", func() {
			Expect(fakeRecorder.Events).To(HaveLen(0))
		})
	})

	When("warning is not null and does not exist", func() {
		BeforeEach(func() {
			managerMock.EXPECT().GetWarningCondition().Return(nil).Times(1)
			managerMock.EXPECT().GetWarningCondition().Return(c).Times(1)
			managerMock.EXPECT().MarkUniqueWarning(gomock.Any()).Times(1)
		})
		When("context has recorder", func() {
			It("should record warning", func() {
				Expect(fakeRecorder.Events).To(HaveLen(1))
				msg := <-fakeRecorder.Events
				expectedMsg := "Warning reason message"
				Expect(msg).To(Equal(expectedMsg), msg)
			})
		})
		When("context does not have recorder", func() {
			BeforeEach(func() {
				ctx = context.Background()
			})
			It("should not record warning", func() {
				Expect(fakeRecorder.Events).To(HaveLen(0))
			})
		})
	})

	When("warning is not null but exists", func() {
		BeforeEach(func() {
			managerMock.EXPECT().GetWarningCondition().Return(c).Times(2)
			managerMock.EXPECT().MarkUniqueWarning(gomock.Any()).Times(1)
		})
		It("should not record warning", func() {
			Expect(fakeRecorder.Events).To(HaveLen(0))
		})
	})

	When("warning already exists but with different type", func() {
		BeforeEach(func() {
			c2 := &apis.Condition{
				Type:    "type-2",
				Reason:  warning.Reason,
				Message: warning.Message,
			}
			managerMock.EXPECT().GetWarningCondition().Return(c2).Times(1)
			managerMock.EXPECT().GetWarningCondition().Return(c).Times(1)
			managerMock.EXPECT().MarkUniqueWarning(gomock.Any()).Times(1)
		})
		When("context has recorder", func() {
			It("should record warning", func() {
				Expect(fakeRecorder.Events).To(HaveLen(1))
				msg := <-fakeRecorder.Events
				expectedMsg := "Warning reason message"
				Expect(msg).To(Equal(expectedMsg), msg)
			})
		})
	})

})

var _ = Describe("Test.HasConditionChanged", func() {
	var (
		before *apis.Condition
		after  *apis.Condition
	)

	BeforeEach(func() {
		before = &apis.Condition{
			Type:    "type",
			Status:  "status",
			Reason:  "reason",
			Message: "message",
			LastTransitionTime: apis.VolatileTime{
				Inner: metav1.Time{
					Time: time.Now(),
				},
			},
		}
		after = &apis.Condition{
			Type:               "type",
			Status:             "status",
			Reason:             "reason",
			Message:            "message",
			LastTransitionTime: apis.VolatileTime{},
		}
	})

	It("should return false if the condition has not changed, ignore the LastTransitionTime field", func() {
		result := hasConditionChanged(before, after)
		Expect(result).To(BeFalse())
	})

	It("should return true if the condition has changed", func() {
		after.Status = "status2"
		result := hasConditionChanged(before, after)
		Expect(result).To(BeTrue())
	})

})
