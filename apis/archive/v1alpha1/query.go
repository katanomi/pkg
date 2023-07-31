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

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
)

// ConditionOperator describe the operator for condition
type ConditionOperator string

const (
	// ConditionOperatorAnd is the operator for and
	ConditionOperatorAnd ConditionOperator = "and"
	// ConditionOperatorOr is the operator for or
	ConditionOperatorOr ConditionOperator = "or"
	// ConditionOperatorIn is the operator for in
	ConditionOperatorIn ConditionOperator = "in"

	// ConditionOperatorEqual is the operator for equal
	ConditionOperatorEqual ConditionOperator = "eq"
	// ConditionOperatorNotEqual is the operator for not equal
	ConditionOperatorNotEqual ConditionOperator = "ne"
	// ConditionOperatorGt is the operator for greater than
	ConditionOperatorGt ConditionOperator = "gt"
	// ConditionOperatorGte is the operator for greater than or equal
	ConditionOperatorGte ConditionOperator = "gte"
	// ConditionOperatorLt is the operator for less than
	ConditionOperatorLt ConditionOperator = "lt"
	// ConditionOperatorLte is the operator for less than or equal
	ConditionOperatorLte ConditionOperator = "lte"
	// ConditionOperatorLike is the operator for like
	ConditionOperatorLike ConditionOperator = "like"
)

// Query describe the query params
type Query struct {
	// Conditions is the query conditions
	Conditions []Condition `json:"conditions,omitempty"`
	// Fields is the fields which will be returned
	Fields []Field `json:"fields,omitempty"`
}

// Condition describe a query condition
type Condition struct {
	// Key is the key for condition
	Key string `json:"key,omitempty"`
	// Operator is the operator for condition
	Operator ConditionOperator `json:"operator,omitempty"`
	// Value is the value for condition
	Value interface{} `json:"value,omitempty"`
}

type conditionShadow struct {
	Key      string            `json:"key"`
	Operator ConditionOperator `json:"operator"`
	Value    json.RawMessage   `json:"value"`
}

// UnmarshalJSON override unmarshalJSON method to parse nested condition
func (p *Condition) UnmarshalJSON(data []byte) error {
	c := conditionShadow{}
	if err := json.Unmarshal(data, &c); err != nil {
		return err
	}
	p.Key = c.Key
	p.Operator = c.Operator
	switch c.Operator {
	case ConditionOperatorIn:
		var v []interface{}
		if err := json.Unmarshal(c.Value, &v); err != nil {
			return err
		}
		p.Value = v
	case ConditionOperatorOr, ConditionOperatorAnd:
		v := []Condition{}
		if err := json.Unmarshal(c.Value, &v); err != nil {
			return err
		}
		p.Value = v
	default:
		var v interface{}
		if err := json.Unmarshal(c.Value, &v); err != nil {
			return err
		}

		p.Value = v
	}

	return nil
}

// DeleteOption describe the delete option
type DeleteOption struct {
	// Direct delete directly, not soft delete
	Direct bool
}

// Order describe the collation of the results
type Order struct {
	Field string
	Asc   bool
}

// GetOptions options for querying a single record
type GetOptions struct {
}

// ListOptions options for querying a list of records
type ListOptions struct {
	metav1alpha1.Pager
	Orders []Order

	// WithDeletedData describe the deleted data should be returned
	WithDeletedData bool
}

// Field describe the field to be returned
type Field struct {
	// Name is the field name
	Name string `json:"name"`
	// Alias is the field alias
	Alias string `json:"alias"`
}
