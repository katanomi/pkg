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

package patchstatus

import (
	"context"
	"errors"

	kclient "github.com/katanomi/pkg/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("Test.PatchStatusAndRecordEvent", func() {
	var (
		ctx                 context.Context
		clt                 client.Client
		configmap           *corev1.ConfigMap
		recorder            record.EventRecorder
		fakeRecorder        *record.FakeRecorder
		err                 error
		isConditionChanged  func() bool
		getTopLevelConditon func() *apis.Condition
	)

	BeforeEach(func() {
		err = nil
		ctx = context.TODO()
		configmap = &corev1.ConfigMap{}
		clt = fake.NewClientBuilder().WithScheme(scheme).Build()
		ctx = kclient.WithClient(ctx, clt)
		ctx = logging.WithLogger(ctx, logger)
		recorder = record.NewFakeRecorder(100)
		fakeRecorder, _ = recorder.(*record.FakeRecorder)
		isConditionChanged = func() bool { return false }
		getTopLevelConditon = func() *apis.Condition { return nil }
	})

	JustBeforeEach(func() {
		PatchStatusAndRecordEvent(ctx, recorder, configmap, configmap, err, isConditionChanged, getTopLevelConditon)
	})

	When("err is not nil", func() {
		BeforeEach(func() {
			err = errors.New("test error")
		})
		It("should record event", func() {
			msg := <-fakeRecorder.Events
			Expect(msg).To(Equal("Warning Error error: test error"))
		})
	})

	Describe("condition changed", func() {
		BeforeEach(func() {
			isConditionChanged = func() bool { return true }
		})

		When("condition status is unknown", func() {
			BeforeEach(func() {
				isConditionChanged = func() bool { return true }
				getTopLevelConditon = func() *apis.Condition {
					return &apis.Condition{
						Type:    apis.ConditionSucceeded,
						Status:  corev1.ConditionUnknown,
						Message: "test message",
						Reason:  string(metav1.StatusReasonBadRequest),
					}
				}
			})
			It("should record event", func() {
				msg := <-fakeRecorder.Events
				Expect(msg).To(Equal("Normal BadRequest test message"))
			})
		})

		When("condition status is false", func() {
			BeforeEach(func() {
				isConditionChanged = func() bool { return true }
				getTopLevelConditon = func() *apis.Condition {
					return &apis.Condition{
						Type:    apis.ConditionSucceeded,
						Status:  corev1.ConditionFalse,
						Message: "test message",
						Reason:  string(metav1.StatusReasonBadRequest),
					}
				}
			})
			It("should record event", func() {
				msg := <-fakeRecorder.Events
				Expect(msg).To(Equal("Warning BadRequest test message"))
			})
		})
	})

})
