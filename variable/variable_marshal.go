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

package variable

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	"k8s.io/utils/field"
)

// MarshalFuncManager defines an interface for getting and returning conversion functions.
type MarshalFuncManager interface {
	NameFuncs() map[string]VariableMarshalFunc
	KindFuncs() map[reflect.Kind]VariableMarshalFunc
}

// VariableMarshalFunc define variable transformation function
type VariableMarshalFunc func(reflect.Type, *field.Path, MarshalFuncManager) ([]v1alpha1.Variable, error)

// DefaultKindMarshalFuncs provides a default Kind conversion function.
var DefaultKindMarshalFuncs = map[reflect.Kind]VariableMarshalFunc{
	reflect.Struct:  MarshalStructToVariable,
	reflect.Pointer: MarshalPointerToVariable,
}

// Marshal returns a list of variables based on v.
// use the default parameters of VariableMarshaller to call Marshal for processing.
func Marshal(v interface{}) (v1alpha1.VariableList, error) {
	marshaller := VariableMarshaller{Object: v}
	return marshaller.Marshal()
}

// VariableMarshaller variable converter
type VariableMarshaller struct {
	NameMarshalFuncs map[string]VariableMarshalFunc
	KindMarshalFuncs map[reflect.Kind]VariableMarshalFunc

	Object interface{}
}

// NameFuncs return name convert funcs
func (v *VariableMarshaller) NameFuncs() map[string]VariableMarshalFunc {
	if v.NameMarshalFuncs == nil {
		v.NameMarshalFuncs = DefaultNameMarshalFuncs
	} else {
		// Add the default marshal function, if user-defined, ignore it.
		for name, value := range DefaultNameMarshalFuncs {
			if _, ok := v.NameMarshalFuncs[name]; !ok {
				v.NameMarshalFuncs[name] = value
			}
		}
	}

	return v.NameMarshalFuncs
}

// KindFuncs return kind convert funcs
func (v *VariableMarshaller) KindFuncs() map[reflect.Kind]VariableMarshalFunc {
	if v.KindMarshalFuncs == nil {
		v.KindMarshalFuncs = DefaultKindMarshalFuncs
	}
	return DefaultKindMarshalFuncs
}

// Marshal convert object to variable list.
func (v *VariableMarshaller) Marshal() (v1alpha1.VariableList, error) {
	if v.Object == nil {
		return v1alpha1.VariableList{}, nil
	}

	st := reflect.TypeOf(v.Object)
	list, err := marshalType(st, nil, v)
	if err != nil {
		return v1alpha1.VariableList{}, err
	}

	return v1alpha1.VariableList{Items: list}, nil
}

// marshalType marshal type to variable list.
func marshalType(st reflect.Type, base *field.Path, marshalFuncs MarshalFuncManager) ([]v1alpha1.Variable, error) {
	// prioritize custom conversion based on the structure name.
	if f, ok := marshalFuncs.NameFuncs()[st.Name()]; ok {
		return f(st, base, marshalFuncs)
	}

	if !isValidVariableKind(st) {
		return nil, fmt.Errorf("unsupported type [%s]", st.Kind().String())
	}

	// marshal the structure according to the kind.
	if f, ok := marshalFuncs.KindFuncs()[st.Kind()]; ok {
		return f(st, base, marshalFuncs)
	}

	return []v1alpha1.Variable{}, nil
}

func isValidVariableKind(st reflect.Type) bool {
	switch st.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.String, reflect.Interface, reflect.Map, reflect.Slice, reflect.Array, reflect.Struct, reflect.Pointer:
		return true
	default:
		return false
	}
}

// MarshalStructToVariable marshal struct to variable list.
func MarshalStructToVariable(st reflect.Type, base *field.Path, marshalFuncs MarshalFuncManager) (list []v1alpha1.Variable, err error) {
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		nextBase := base
		if !field.Anonymous {
			variable := NewVariable(field, base)
			if variable != nil {
				list = append(list, *variable)
			}
			nextBase = base.Child(getJsonTagName(field))
		}

		subList, err := marshalType(field.Type, nextBase, marshalFuncs)
		if err != nil {
			return list, err
		}

		list = append(list, subList...)
	}
	return list, nil
}

// MarshalPointerToVariable marshal pointer to variable list.
func MarshalPointerToVariable(st reflect.Type, base *field.Path, marshalFuncs MarshalFuncManager) ([]v1alpha1.Variable, error) {
	return marshalType(st.Elem(), base, marshalFuncs)
}

// NewVariable return variable from StructField.
func NewVariable(field reflect.StructField, base *field.Path) *v1alpha1.Variable {
	switch field.Type.Kind() {
	// Ignoring the production of list types, usually the value of variables of list type is undefined.
	case reflect.Map, reflect.Slice, reflect.Array:
		return nil
	default:
		// nothing
	}

	jsonTag := getJsonTagName(field)
	if jsonTag == "" {
		return nil
	}

	variableTagStr := field.Tag.Get("variable")
	if variableTagStr == "-" {
		return nil
	}

	variableTagValues := parseVariableTag(variableTagStr)
	return &v1alpha1.Variable{
		Name:    base.Child(jsonTag).String(),
		Label:   variableTagValues["label"],
		Example: variableTagValues["example"],
	}
}

func parseVariableTag(variableTagStr string) map[string]string {
	result := map[string]string{}
	variableTags := strings.Split(variableTagStr, ";")
	for _, tagStr := range variableTags {
		tags := strings.Split(tagStr, "=")
		if len(tags) == 2 {
			result[tags[0]] = tags[1]
		}
	}
	return result
}

func getJsonTagName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "-" || jsonTag == "" {
		return ""
	}
	return strings.Split(jsonTag, ",")[0]
}
