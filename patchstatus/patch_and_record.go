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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	kclient "github.com/katanomi/pkg/client"
)

// PatchStatusAndRecordEvent patch status and record event
func PatchStatusAndRecordEvent(
	ctx context.Context,
	eventRecorder record.EventRecorder,
	obj, old runtime.Object, err error,
	isConditionChanged func() bool,
	getTopLevelConditon func() *apis.Condition) {

	log := logging.FromContext(ctx)
	clt := kclient.Client(ctx)

	patch := client.MergeFrom(old)
	patchData, _ := patch.Data(obj)
	clientObj, _ := obj.(client.Object)
	patchErr := clt.Status().Patch(ctx, clientObj, patch)
	log.Debugw("object patch result", "err", patchErr, "patchData", string(patchData))

	if err != nil {
		reason := metav1alpha1.ReasonForError(err)
		if reason == "" {
			reason = metav1alpha1.ErrorReason
		}
		eventRecorder.Eventf(obj, corev1.EventTypeWarning, reason, "error: %s", err)
		return
	}

	if !isConditionChanged() {
		return
	}

	top := getTopLevelConditon()
	if top == nil {
		return
	}
	switch top.Status {
	case corev1.ConditionTrue, corev1.ConditionUnknown:
		eventRecorder.Eventf(obj, corev1.EventTypeNormal, top.Reason, top.Message)
	case corev1.ConditionFalse:
		eventRecorder.Eventf(obj, corev1.EventTypeWarning, top.Reason, top.Message)
	default:
		logging.FromContext(ctx).Warnw("unknown top condition status", "status", top.Status)
	}
}
