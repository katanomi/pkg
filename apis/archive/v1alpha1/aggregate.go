/*
Copyright 2022 The Katanomi Authors.

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
	"encoding/json"
	"fmt"
	"reflect"
)

// AggregateOperator describe the operator type for aggregate
type AggregateOperator string

const (
	// AggregateOperatorMax is the operator for max
	AggregateOperatorMax AggregateOperator = "max"
	// AggregateOperatorMin is the operator for min
	AggregateOperatorMin AggregateOperator = "min"
	// AggregateOperatorSum is the operator for sum
	AggregateOperatorSum AggregateOperator = "sum"
	// AggregateOperatorCount is the operator for count
	AggregateOperatorCount AggregateOperator = "count"
)

// AggregateField is a field for aggregate
type AggregateField struct {
	Field `json:",inline"`
	// Operator is the operator for aggregate
	Operator AggregateOperator `json:"operator"`
}

// AggregateQuery describe the query params for aggregate
type AggregateQuery struct {
	// Conditions is the conditions for aggregate query
	Conditions []Condition `json:"conditions,omitempty"`
	// GroupFields is the fields for group by
	GroupFields []Field `json:"group_fields"`
	// AggregateFields is the fields for aggregate
	AggregateFields []AggregateField `json:"aggregate_fields"`
}

// Max generate aggregate field with max operator
func Max(name, alias string) AggregateField {
	return AggregateField{
		Field:    Field{Name: name, Alias: alias},
		Operator: AggregateOperatorMax,
	}
}

// Min generate aggregate field with min operator
func Min(name, alias string) AggregateField {
	return AggregateField{
		Field:    Field{Name: name, Alias: alias},
		Operator: AggregateOperatorMin,
	}
}

// Sum generate aggregate field with sum operator
func Sum(name, alias string) AggregateField {
	return AggregateField{
		Field:    Field{Name: name, Alias: alias},
		Operator: AggregateOperatorSum,
	}
}

// Count generate aggregate field with count operator
func Count(alias string) AggregateField {
	return AggregateField{
		Field:    Field{Alias: alias},
		Operator: AggregateOperatorCount,
	}
}

// AggregateResult describe the result of aggregate query
type AggregateResult []map[string]interface{}

// Unmarshal convert aggregate result to list
func (a AggregateResult) Unmarshal(list interface{}) error {
	v := reflect.ValueOf(list)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("list param must be a pointer to slice")
	}
	data, err := json.Marshal(a)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, list)
}
