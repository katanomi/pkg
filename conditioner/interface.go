/*
Copyright 2023 The Katanomi Authors.

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

// Package conditioner contains conditioner related logic
package conditioner

import (
	"context"

	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate mockgen  -destination=../testing/mock/github.com/katanomi/pkg/conditioner/interface.go -package=mock github.com/katanomi/pkg/conditioner  RunConditionerManager

// RunConditionerManager uesd for answering two questions
// 1. If run meets all the condition that run can go on
// 2. If return the left runs that should go on if the run has done
// CalculateRunConditons answer the first question by returning the unsatisfied condition
// GetAffected answer the second question by return the runs base on different condition
type RunConditionerManager interface {
	// CalculateRunConditons return the unsatisfied condition
	CalculateRunConditons(ctx context.Context,
		obj client.Object,
		listrunList client.ObjectList,
		// prepare used to sort and filter the objects
		prepare func(client.ObjectList),
		conditions []RunConditioner) ([]RunConditionHandler, error)
	// GetAffected return all the affected objects that need handled after the current object completed
	GetAffected(ctx context.Context,
		obj client.Object,
		list client.ObjectList,
		// prepare used to sort and filter the objects
		prepare func(client.ObjectList),
		conditions []RunConditioner) (map[v1alpha1.RunConditionerType]client.ObjectList, error)
}

// RunConditioner indicate the state of the condition
type RunConditioner interface {
	// Type return the condtioner type
	Type() v1alpha1.RunConditionerType
	// RunCondition return the condition details that indicate the condition status
	RunCondition(context.Context, client.Object, client.ObjectList) (RunConditionHandler, error)
	// GetAffected return objects which be affected if the current object completed
	GetAffected(context.Context, client.Object, client.ObjectList) (client.ObjectList, error)
}
