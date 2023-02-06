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

// ConvertFuncManager defines an interface for getting and returning conversion functions.
type ConvertFuncManager interface {
	NameFuncs() map[string]ConvertFunc
	KindFuncs() map[reflect.Kind]ConvertFunc
}

// ConvertFunc define variable transformation function
type ConvertFunc func(reflect.Type, *field.Path, ConvertFuncManager) ([]Variable, error)

// DefaultKindConvertFuncs provides a default Kind conversion function.
var DefaultKindConvertFuncs = map[reflect.Kind]ConvertFunc{
	reflect.Struct:  ConvertStruct,
	reflect.Pointer: ConvertPointer,
}

// FilterFunc define a filter function for a variable
type FilterFunc func(*Variable) bool

// VariableConverter variable converter
type VariableConverter struct {
	NameConvertFuncs map[string]ConvertFunc
	KindConvertFuncs map[reflect.Kind]ConvertFunc

	Filter []FilterFunc
}

// NameFuncs return name convert funcs
func (v *VariableConverter) NameFuncs() map[string]ConvertFunc {
	if v.NameConvertFuncs == nil {
		v.NameConvertFuncs = DefaultNameConvertFuncs
	}
	return v.NameConvertFuncs
}

// KindFuncs return kind convert funcs
func (v *VariableConverter) KindFuncs() map[reflect.Kind]ConvertFunc {
	if v.KindConvertFuncs == nil {
		v.KindConvertFuncs = DefaultKindConvertFuncs
	}
	return v.KindConvertFuncs
}

// ConvertToVariableList convert object to variable list.
func (v *VariableConverter) ConvertToVariableList(obj interface{}, filters ...FilterFunc) (VariableList, error) {
	if obj == nil {
		return VariableList{}, nil
	}

	st := reflect.TypeOf(obj)
	list, err := ConvertType(st, nil, v)
	if err != nil {
		return VariableList{}, err
	}

	varList := VariableList{}
	for i := range list {
		if filtVariable(&list[i], filters...) {
			varList.Items = append(varList.Items, list[i])
		}
	}
	return varList, nil
}

// ConvertType convert type to variable list.
func ConvertType(st reflect.Type, base *field.Path, convertFuncs ConvertFuncManager) ([]Variable, error) {
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
func ConvertStruct(st reflect.Type, base *field.Path, convertFuncs ConvertFuncManager) (list []Variable, err error) {
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

		subList, err := ConvertType(field.Type, nextBase, convertFuncs)
		if err != nil {
			return list, err
		}

		list = append(list, subList...)
	}
	return list, nil
}

// ConvertPointer convert pointer to variable list.
func ConvertPointer(st reflect.Type, base *field.Path, convertFuncs ConvertFuncManager) ([]Variable, error) {
	return ConvertType(st.Elem(), base, convertFuncs)
}
