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

package v1alpha1

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/strings/slices"
	"knative.dev/pkg/apis"
)

// ParamSpec defines arbitrary parameters needed beyond typed inputs (such as
// resources).
type ParamSpec struct {
	// Name declares the name by which a parameter is referenced.
	Name string `json:"name"`
	// Type is the user-specified type of the parameter. The possible types
	// are currently "string", "array" and "object", and "string" is the default.
	// +optional
	Type ParamType `json:"type,omitempty"`
	// Description is a user-facing description of the parameter that may be
	// used to populate a UI.
	// +optional
	Description string `json:"description,omitempty"`
	// Properties is the JSON Schema properties to support key-value pairs parameter.
	// +optional
	Properties map[string]PropertySpec `json:"properties,omitempty"`
	// Default is the value a parameter takes if no input value is supplied.
	// +optional
	Default *ParamValue `json:"default,omitempty"`
	// Enum declares a set of allowed param input values.
	// If Enum is not set, no input validation is performed for the param.
	// +optional
	Enum []string `json:"enum,omitempty"`
}

// ParamSpecs is a list of ParamSpec
type ParamSpecs []ParamSpec

// PropertySpec defines the struct for object keys
type PropertySpec struct {
	Type ParamType `json:"type,omitempty"`
}

// SetDefaults set the default type
func (pp *ParamSpec) SetDefaults(context.Context) {
	if pp == nil {
		return
	}

	// Propagate inferred type to the parent ParamSpec's type, and default type to the PropertySpec's type
	// The sequence to look at is type in ParamSpec -> properties -> type in default -> array/string/object value in default
	// If neither `properties` or `default` section is provided, ParamTypeString will be the default type.
	switch {
	case pp.Type != "":
		// If param type is provided by the author, do nothing but just set default type for PropertySpec in case `properties` section is provided.
		pp.setDefaultsForProperties()
	case pp.Properties != nil:
		pp.Type = ParamTypeObject
		// Also set default type for PropertySpec
		pp.setDefaultsForProperties()
	case pp.Default == nil:
		// ParamTypeString is the default value (when no type can be inferred from the default value)
		pp.Type = ParamTypeString
	case pp.Default.Type != "":
		pp.Type = pp.Default.Type
	case pp.Default.ArrayVal != nil:
		pp.Type = ParamTypeArray
	case pp.Default.ObjectVal != nil:
		pp.Type = ParamTypeObject
	default:
		pp.Type = ParamTypeString
	}
}

// setDefaultsForProperties sets default type for PropertySpec (string) if it's not specified
func (pp *ParamSpec) setDefaultsForProperties() {
	for key, propertySpec := range pp.Properties {
		if propertySpec.Type == "" {
			pp.Properties[key] = PropertySpec{Type: ParamTypeString}
		}
	}
}

// GetNames returns all the names of the declared parameters
func (ps ParamSpecs) GetNames() []string {
	var names = []string{}
	for _, p := range ps {
		names = append(names, p.Name)
	}
	return names
}

// SortByType splits the input params into string params, array params, and object params, in that order
func (ps ParamSpecs) SortByType() (ParamSpecs, ParamSpecs, ParamSpecs) {
	var stringParams, arrayParams, objectParams ParamSpecs
	for _, p := range ps {
		switch p.Type {
		case ParamTypeArray:
			arrayParams = append(arrayParams, p)
		case ParamTypeObject:
			objectParams = append(objectParams, p)
		case ParamTypeString:
			fallthrough
		default:
			stringParams = append(stringParams, p)
		}
	}
	return stringParams, arrayParams, objectParams
}

// ValidateNoDuplicateNames returns an error if any of the params have the same name
func (ps ParamSpecs) ValidateNoDuplicateNames() *apis.FieldError {
	var errs *apis.FieldError
	names := ps.GetNames()
	for dup := range findDups(names) {
		errs = errs.Also(apis.ErrGeneric("parameter appears more than once", "").ViaFieldKey("params", dup))
	}
	return errs
}

// validateParamEnum validates feature flag, duplication and allowed types for Param Enum
func (ps ParamSpecs) validateParamEnums(_ context.Context) *apis.FieldError {
	var errs *apis.FieldError
	for _, p := range ps {
		if len(p.Enum) == 0 {
			continue
		}
		if p.Type != ParamTypeString {
			errs = errs.Also(apis.ErrGeneric("enum can only be set with string type param", "").ViaKey(p.Name))
		}
		for dup := range findDups(p.Enum) {
			errs = errs.Also(apis.ErrGeneric(fmt.Sprintf("parameter enum value %v appears more than once", dup), "").ViaKey(p.Name))
		}
		if p.Default != nil && p.Default.StringVal != "" {
			if !slices.Contains(p.Enum, p.Default.StringVal) {
				errs = errs.Also(apis.ErrGeneric(fmt.Sprintf("param default value %v not in the enum list", p.Default.StringVal), "").ViaKey(p.Name))
			}
		}
	}
	return errs
}

// findDups returns the duplicate element in the given slice
func findDups(vals []string) sets.Set[string] {
	seen := sets.Set[string]{}
	dups := sets.Set[string]{}
	for _, val := range vals {
		if seen.Has(val) {
			dups.Insert(val)
		}
		seen.Insert(val)
	}
	return dups
}

// Param declares an ParamValues to use for the parameter called name.
type Param struct {
	Name  string     `json:"name"`
	Value ParamValue `json:"value"`
}

// ExtractNames returns a set of unique names
func (ps Params) ExtractNames() sets.Set[string] {
	names := sets.Set[string]{}
	for _, p := range ps {
		names.Insert(p.Name)
	}
	return names
}

func (ps Params) extractValues() []string {
	pvs := []string{}
	for i := range ps {
		pvs = append(pvs, ps[i].Value.StringVal)
		pvs = append(pvs, ps[i].Value.ArrayVal...)
		for _, v := range ps[i].Value.ObjectVal {
			pvs = append(pvs, v)
		}
	}
	return pvs
}

// extractParamMapArrVals creates a param map with the key: param.Name and
// val: param.Value.ArrayVal
func (ps Params) extractParamMapArrVals() map[string][]string {
	paramsMap := make(map[string][]string)
	for _, p := range ps {
		paramsMap[p.Name] = p.Value.ArrayVal
	}
	return paramsMap
}

// Params is a list of Param
type Params []Param

// validateDuplicateParameters checks if a parameter with the same name is defined more than once
func (ps Params) validateDuplicateParameters() (errs *apis.FieldError) {
	paramNames := sets.NewString()
	for i, param := range ps {
		if paramNames.Has(param.Name) {
			errs = errs.Also(apis.ErrGeneric(fmt.Sprintf("parameter names must be unique,"+
				" the parameter \"%s\" is also defined at", param.Name), fmt.Sprintf("[%d].name", i)))
		}
		paramNames.Insert(param.Name)
	}
	return errs
}

// ParamType indicates the type of an input parameter;
// Used to distinguish between a single string and an array of strings.
type ParamType string

// Valid ParamTypes:
const (
	ParamTypeString ParamType = "string"
	ParamTypeArray  ParamType = "array"
	ParamTypeObject ParamType = "object"
)

// AllParamTypes can be used for ParamType validation.
var AllParamTypes = []ParamType{ParamTypeString, ParamTypeArray, ParamTypeObject}

// ParamValues is modeled after IntOrString in kubernetes/apimachinery:

// ParamValue is a type that can hold a single string, string array, or string map.
// Used in JSON unmarshalling so that a single JSON field can accept
// either an individual string or an array of strings.
type ParamValue struct {
	Type      ParamType // Represents the stored type of ParamValues.
	StringVal string
	// +listType=atomic
	ArrayVal  []string
	ObjectVal map[string]string
}

// UnmarshalJSON implements the json.Unmarshaller interface.
func (paramValues *ParamValue) UnmarshalJSON(value []byte) error {
	// ParamValues is used for Results Value as well, the results can be any kind of
	// data so we need to check if it is empty.
	if len(value) == 0 {
		paramValues.Type = ParamTypeString
		return nil
	}
	if value[0] == '[' {
		// We're trying to Unmarshal to []string, but for cases like []int or other types
		// of nested array which we don't support yet, we should continue and Unmarshal
		// it to String. If the Type being set doesn't match what it actually should be,
		// it will be captured by validation in reconciler.
		// if failed to unmarshal to array, we will convert the value to string and marshal it to string
		var a []string
		if err := json.Unmarshal(value, &a); err == nil {
			paramValues.Type = ParamTypeArray
			paramValues.ArrayVal = a
			return nil
		}
	}
	if value[0] == '{' {
		// if failed to unmarshal to map, we will convert the value to string and marshal it to string
		var m map[string]string
		if err := json.Unmarshal(value, &m); err == nil {
			paramValues.Type = ParamTypeObject
			paramValues.ObjectVal = m
			return nil
		}
	}

	// By default we unmarshal to string
	paramValues.Type = ParamTypeString
	if err := json.Unmarshal(value, &paramValues.StringVal); err == nil {
		return nil
	}
	paramValues.StringVal = string(value)

	return nil
}

// MarshalJSON implements the json.Marshaller interface.
func (paramValues ParamValue) MarshalJSON() ([]byte, error) {
	switch paramValues.Type {
	case ParamTypeString:
		return json.Marshal(paramValues.StringVal)
	case ParamTypeArray:
		return json.Marshal(paramValues.ArrayVal)
	case ParamTypeObject:
		return json.Marshal(paramValues.ObjectVal)
	default:
		return []byte{}, fmt.Errorf("impossible ParamValues.Type: %q", paramValues.Type)
	}
}

// NewStructuredValues creates an ParamValues of type ParamTypeString or ParamTypeArray, based on
// how many inputs are given (>1 input will create an array, not string).
func NewStructuredValues(value string, values ...string) *ParamValue {
	if len(values) > 0 {
		return &ParamValue{
			Type:     ParamTypeArray,
			ArrayVal: append([]string{value}, values...),
		}
	}
	return &ParamValue{
		Type:      ParamTypeString,
		StringVal: value,
	}
}

// NewObject creates an ParamValues of type ParamTypeObject using the provided key-value pairs
func NewObject(pairs map[string]string) *ParamValue {
	return &ParamValue{
		Type:      ParamTypeObject,
		ObjectVal: pairs,
	}
}
