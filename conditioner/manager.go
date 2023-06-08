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

package conditioner

import (
	"context"

	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ RunConditionerManager = &DefaultRunConditionerManager{}

// DefaultRunConditionerManager is default implement for RunConditionerManager
type DefaultRunConditionerManager struct {
}

// NewRunConditionManager initialize DefaultRunConditionManager
func NewRunConditionManager() DefaultRunConditionerManager {
	return DefaultRunConditionerManager{}
}

// CalculateRunConditons return all the condition for different runConditioner
func (manger DefaultRunConditionerManager) CalculateRunConditons(ctx context.Context, obj client.Object, list client.ObjectList, prepare func(client.ObjectList), conditioners []RunConditioner) ([]RunConditionHandler, error) {
	result := make([]RunConditionHandler, 0, len(conditioners))
	var errors []error
	prepare(list)

	for _, conditioner := range conditioners {
		cond, err := conditioner.RunCondition(ctx, obj, list)
		result = append(result, cond)
		if err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		return result, apierrors.NewAggregate(errors)
	}

	return result, nil
}

// GetAffected return all affected object base on different runConditioner
func (manager DefaultRunConditionerManager) GetAffected(ctx context.Context, obj client.Object, list client.ObjectList, prepare func(client.ObjectList), conditioners []RunConditioner) (map[v1alpha1.RunConditionerType]client.ObjectList, error) {

	result := make(map[v1alpha1.RunConditionerType]client.ObjectList)
	prepare(list)
	var errors []error

	for _, conditioner := range conditioners {
		affected, err := conditioner.GetAffected(ctx, obj, list)
		result[conditioner.Type()] = affected
		if err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		return result, apierrors.NewAggregate(errors)
	}
	return result, nil
}
