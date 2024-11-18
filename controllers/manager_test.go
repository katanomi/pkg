/*
Copyright 2023 The AlaudaDevops Authors.

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

package controllers

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cliengorecord "k8s.io/client-go/tools/record"

	"github.com/AlaudaDevops/pkg/testing/mock/sigs.k8s.io/controller-runtime/pkg/manager"
	"github.com/golang/mock/gomock"
)

func TestGetEventRecorderFor(t *testing.T) {
	mockctl := gomock.NewController(t)
	defer func() {
		mockctl.Finish()
	}()

	mockManaer := manager.NewMockManager(mockctl)
	mockManaer.EXPECT().GetEventRecorderFor(gomock.Any()).Return(cliengorecord.NewFakeRecorder(5)).Times(1)
	var cm = ControllerManager{
		Manager: mockManaer,
	}

	t.Run("it should return local EventRecorder", func(t *testing.T) {
		recorder := cm.GetEventRecorderFor("Demo")
		if _, ok := recorder.(*EventRecorder); !ok {
			t.Errorf("should be *EventRecorder, but got %T", recorder)
		}
	})
}

func TestEventRecorder_Event(t *testing.T) {
	fakeRecorder := cliengorecord.NewFakeRecorder(5)
	rec := &EventRecorder{
		EventRecorder:     fakeRecorder,
		MessageShortenLen: 10,
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "demo", Namespace: "default"},
	}
	var data = []struct {
		desc string
		msg  string

		expectedEvt string
	}{
		{
			desc: "message is long than max length, it should be cut",
			msg:  "this-is-long-message",

			expectedEvt: "Warning PullImageFailed this-is-lo",
		},
		{
			desc: "message is less than max length, it should be keep original string",
			msg:  "message",

			expectedEvt: "Warning PullImageFailed message",
		},
	}

	for _, item := range data {
		t.Run(item.desc, func(t *testing.T) {
			rec.Event(pod, corev1.EventTypeWarning, "PullImageFailed", item.msg)

			select {
			case evt := <-fakeRecorder.Events:
				{
					diff := cmp.Diff(item.expectedEvt, evt)
					if diff != "" {
						t.Errorf("expected event is different with actual: %s", diff)
					}
				}
			}
		})
	}
}

func TestEventRecorder_Eventf(t *testing.T) {
	fakeRecorder := cliengorecord.NewFakeRecorder(5)
	rec := &EventRecorder{
		EventRecorder:     fakeRecorder,
		MessageShortenLen: 10,
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "demo", Namespace: "default"},
	}

	var data = []struct {
		desc   string
		msgFmt string
		args   []interface{}

		expectedEvt string
	}{
		{
			desc:   "arg is less than max length, it should keep original",
			msgFmt: "arg1: %s, arg2: %s",
			args: []interface{}{
				"message1",
				"message2",
			},

			expectedEvt: "Warning PullImageFailed arg1: message1, arg2: message2",
		},

		{
			desc:   "message is long than max length and with other types, it should be cut",
			msgFmt: "arg-long: %s, arg-short: %s, arg-error1: %s, arg-error2: %s, arg-int: %d, arg-bool: %t, arg-q: %q, arg-w: %w",
			args: []interface{}{
				"this-is-long-message",
				"message",
				fmt.Errorf("1-errors-created-by-fmt"),
				fmt.Errorf("2-errors-created-by-fmt"),
				1,
				true,
				"this-is-long-message-q",
				fmt.Errorf("3-errors-created-by-fmt"),
			},

			expectedEvt: "Warning PullImageFailed arg-long: this-is-lo, arg-short: message, arg-error1: 1-errors-c, arg-error2: 2-errors-c, arg-int: 1, arg-bool: true, arg-q: \"this-is-lo\", arg-w: %!w(string=3-errors-c)",
		},
	}

	for _, item := range data {
		t.Run(item.desc, func(t *testing.T) {
			rec.Eventf(pod, corev1.EventTypeWarning, "PullImageFailed", item.msgFmt, item.args...)

			select {
			case evt := <-fakeRecorder.Events:
				{
					diff := cmp.Diff(item.expectedEvt, evt)
					if diff != "" {
						t.Errorf("expected event is different with actual: %s", diff)
					}
				}
			}

			rec.AnnotatedEventf(pod, map[string]string{}, corev1.EventTypeWarning, "PullImageFailed", item.msgFmt, item.args...)
			select {
			case evt := <-fakeRecorder.Events:
				{
					diff := cmp.Diff(item.expectedEvt, evt)
					if diff != "" {
						t.Errorf("expected event is different with actual: %s", diff)
					}
				}
			}
		})
	}
}
