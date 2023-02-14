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

	"k8s.io/utils/field"
)

// MarshalFuncManager defines an interface for getting and returning conversion functions.
type MarshalFuncManager interface {
	NameFuncs() map[string]ConvertFunc
	KindFuncs() map[reflect.Kind]ConvertFunc
}

// ConvertFunc define variable transformation function
type ConvertFunc func(reflect.Type, *field.Path, MarshalFuncManager) ([]Variable, error)

// DefaultKindConvertFuncs provides a default Kind conversion function.
var DefaultKindMarshalFuncs = map[reflect.Kind]ConvertFunc{
	reflect.Struct:  convertStruct,
	reflect.Pointer: convertPointer,
}

// Marshal returns a list of variables based on v.
// use the default parameters of VariableMarshaller to call Marshal for processing.
func Marshal(v interface{}) (VariableList, error) {
	marshaller := VariableMarshaller{Object: v}
	return marshaller.Marshal()
}

// VariableMarshaller variable converter
type VariableMarshaller struct {
	NameMarshalFuncs map[string]ConvertFunc
	KindMarshalFuncs map[reflect.Kind]ConvertFunc

	Object interface{}
}

// NameFuncs return name convert funcs
func (v *VariableMarshaller) NameFuncs() map[string]ConvertFunc {
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
func (v *VariableMarshaller) KindFuncs() map[reflect.Kind]ConvertFunc {
	if v.KindMarshalFuncs == nil {
		v.KindMarshalFuncs = DefaultKindMarshalFuncs
	}
	return DefaultKindMarshalFuncs
}

// ConvertToVariableList convert object to variable list.
func (v *VariableMarshaller) Marshal() (VariableList, error) {
	if v.Object == nil {
		return VariableList{}, nil
	}

	st := reflect.TypeOf(v.Object)
	list, err := convertType(st, nil, v)
	if err != nil {
		return VariableList{}, err
	}

	return VariableList{Items: list}, nil
}

// convertType convert type to variable list.
func convertType(st reflect.Type, base *field.Path, convertFuncs MarshalFuncManager) ([]Variable, error) {
	// prioritize custom conversion based on the structure name.
	if f, ok := convertFuncs.NameFuncs()[st.Name()]; ok {
		return f(st, base, convertFuncs)
	}

	if !isValidVariableKind(st) {
		return nil, fmt.Errorf("unsupported type [%s]", st.Kind().String())
	}

	// Convert the structure according to the kind.
	if f, ok := convertFuncs.KindFuncs()[st.Kind()]; ok {
		return f(st, base, convertFuncs)
	}

	return []Variable{}, nil
}

func isValidVariableKind(st reflect.Type) bool {
	switch st.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.String, reflect.Interface, reflect.Map, reflect.Slice, reflect.Array, reflect.Struct, reflect.Pointer:
		return true
	default:
		return false
	}
}

// ConvertStruct convert struct to variable list.
func convertStruct(st reflect.Type, base *field.Path, convertFuncs MarshalFuncManager) (list []Variable, err error) {
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

		subList, err := convertType(field.Type, nextBase, convertFuncs)
		if err != nil {
			return list, err
		}

		list = append(list, subList...)
	}
	return list, nil
}

// ConvertPointer convert pointer to variable list.
func convertPointer(st reflect.Type, base *field.Path, convertFuncs MarshalFuncManager) ([]Variable, error) {
	return convertType(st.Elem(), base, convertFuncs)
}
