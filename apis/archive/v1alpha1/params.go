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

package v1alpha1

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

const (
	// LikeValuePrefixForEqualOperator Like operator value prefix
	LikeValuePrefixForEqualOperator = "like--"
)

// AggregateParams params for aggregating archive data
type AggregateParams struct {
	Query   AggregateQuery `json:"query,omitempty"`
	Options *ListOptions   `json:"options,omitempty"`
}

// ListParams params for fetching list archive data
type ListParams struct {
	Query   Query        `json:"query,omitempty"`
	Options *ListOptions `json:"options,omitempty"`
}

// DeleteParams params for deleting archive data
type DeleteParams struct {
	Conditions []Condition   `json:"conditions,omitempty"`
	Options    *DeleteOption `json:"options,omitempty"`
}

// ConvertMedataSelectorToConditions convert a string representing a metadata selector and returns Condition array
// The input will cause an error if it does not follow this form:
// 'in' operator: key in (value1,value2,value3)
// 'equal' operator: key=value
// 'exist' operator: key
// You can join multi operator, example: "key in (value1,value2,value3),key=value,key2"
func ConvertMedataSelectorToConditions(metadataSelectorStr string) ([]Condition, error) {
	requirements, err := convertMedataSelectorToRequirements(metadataSelectorStr)
	if err != nil {
		return nil, err
	}

	var conditions []Condition
	for _, req := range requirements {
		switch req.Operator() {
		case selection.Equals:
			if req.Values().Len() > 0 {
				val := req.Values().List()[0]
				if strings.HasPrefix(val, LikeValuePrefixForEqualOperator) {
					conditions = append(conditions, Like(MetadataKey(req.Key()), strings.TrimPrefix(val, "like--")))
				} else {
					conditions = append(conditions, Equal(MetadataKey(req.Key()), val))
				}
			}
		case selection.In:
			var inConditions []Condition
			for _, val := range req.Values().List() {
				inConditions = append(inConditions, Equal(MetadataKey(req.Key()), val))
			}
			conditions = append(conditions, Or(inConditions...))
		case selection.Exists:
			conditions = append(conditions, Exist(MetadataKey(req.Key())))
		default:
			return nil, fmt.Errorf("%s is not a valid label selector operator", req.Operator())
		}
	}
	return conditions, nil
}

func convertMedataSelectorToRequirements(selectorStr string) (labels.Requirements, error) {
	selector, err := labels.Parse(selectorStr)
	if err != nil {
		return nil, err
	}

	requirements, selectable := selector.Requirements()
	if !selectable {
		return nil, fmt.Errorf("%s doesn't want to select anything", selectorStr)
	}

	return requirements, nil
}
