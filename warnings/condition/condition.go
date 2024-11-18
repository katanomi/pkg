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

// Package condition provides a set of utilities for managing warning conditions and emitting warning events
package condition

import (
	"context"

	krecord "github.com/AlaudaDevops/pkg/record"
	"github.com/AlaudaDevops/pkg/warnings"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/conversion"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate mockgen -package=warnings -destination=../../testing/mock/github.com/AlaudaDevops/pkg/warnings/condition/condition_marker.go github.com/AlaudaDevops/pkg/warnings/condition WarningConditionManager

// WarningConditionManager used to manager warning condition and emit warning events
type WarningConditionManager interface {
	// GetObject returns the object, used to emit warning events
	GetObject() client.Object

	// MarkUniqueWarning marks a unique warning condition
	MarkUniqueWarning(*warnings.WarningRecord)

	// GetWarningCondition returns the current warning condition
	GetWarningCondition() *apis.Condition
}

// MarkAndRecordWarning marks warning condition and emit warning events if the condition has changed
func MarkAndRecordWarning(ctx context.Context, manager WarningConditionManager, warning *warnings.WarningRecord) {

	if warning == nil {
		return
	}

	before := manager.GetWarningCondition()
	manager.MarkUniqueWarning(warning)
	after := manager.GetWarningCondition()

	// emit the warning event
	recorder := krecord.FromContext(ctx)
	if recorder != nil && hasConditionChanged(before, after) {
		recorder.Eventf(manager.GetObject(), corev1.EventTypeWarning, warning.Reason, warning.Message)
		logger := logging.FromContext(ctx)
		logger.Debugw("Recorded warning", "reason", warning.Reason, "message", warning.Message)
	}

	return
}

// hasConditionChanged returns true if the condition has changed, ignore the LastTransitionTime field
func hasConditionChanged(before, after *apis.Condition) bool {
	// ignore the LastTransitionTime field
	equalities := conversion.EqualitiesOrDie(
		func(a, b apis.VolatileTime) bool {
			return true
		},
	)
	return !equalities.DeepEqual(before, after)
}
