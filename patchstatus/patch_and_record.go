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
	"fmt"

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
// If the input err is not nil, or patch failed, it will return an error.
func PatchStatusAndRecordEvent(
	ctx context.Context,
	eventRecorder record.EventRecorder,
	obj, old runtime.Object, err error,
	isConditionChanged func() bool,
	getTopLevelConditon func() *apis.Condition) error {

	log := logging.FromContext(ctx)
	clt := kclient.Client(ctx)

	// makes this casting to patch data in the methods below
	// there were a few method changes inside controller-runtime
	// TODO(danielfbm): consider changing the parameter type to avoid this
	oldObj, oldIsObject := old.(client.Object)
	if !oldIsObject {
		log.Warnw("old object could not be patched because it is not a client.Object", "old", old)
		// just return the original error, not creating a new one.
		return err
	}
	clientObj, _ := obj.(client.Object)
	patch := client.MergeFrom(oldObj)
	patchData, _ := patch.Data(clientObj)

	// Empty patch is {}, hence we check for that.
	// Execute only when patch is required
	if len(patchData) > 2 {
		patchErr := clt.Status().Patch(ctx, clientObj, patch)
		if patchErr != nil {
			log.Warnw("object patch failed", "err", patchErr, "patchData", string(patchData))
			return fmt.Errorf("failed to patch object: %w", patchErr)
		} else {
			log.Debugw("object patch success", "patchData", string(patchData))
		}
	}

	if err != nil {
		reason := metav1alpha1.ReasonForError(err)
		if reason == "" {
			reason = metav1alpha1.ErrorReason
		}
		eventRecorder.Eventf(obj, corev1.EventTypeWarning, reason, "error: %s", err)
		return err
	}

	if !isConditionChanged() {
		return err
	}

	top := getTopLevelConditon()
	if top == nil {
		return err
	}
	switch top.Status {
	case corev1.ConditionTrue, corev1.ConditionUnknown:
		eventRecorder.Eventf(obj, corev1.EventTypeNormal, top.Reason, top.Message)
	case corev1.ConditionFalse:
		eventRecorder.Eventf(obj, corev1.EventTypeWarning, top.Reason, top.Message)
	default:
		logging.FromContext(ctx).Warnw("unknown top condition status", "status", top.Status)
	}
	return err
}
