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
	"reflect"
	"testing"

	"k8s.io/utils/field"
)

func TestNewVariable(t *testing.T) {
	base := field.NewPath("test")
	tests := map[string]struct {
		field reflect.StructField
		base  *field.Path
		want  *Variable
	}{
		"input map": {
			field: reflect.StructField{Type: reflect.TypeOf(map[string]string{})},
			base:  base,
		},
		"inline struct": {
			base:  base,
			field: reflect.StructField{Type: reflect.TypeOf("test")},
		},
		"ignore json tag": {
			base:  base,
			field: reflect.StructField{Tag: reflect.StructTag("-"), Type: reflect.TypeOf("test")},
		},
		"parse tag success": {
			base:  base,
			field: reflect.StructField{Tag: reflect.StructTag(`json:"test" variable:"label=test;other=;example=12"`), Type: reflect.TypeOf("test")},
			want: &Variable{
				Name:    base.Child("test").String(),
				Example: "12",
				Label:   "test",
			},
		},
		"parse tag without tag value": {
			base:  base,
			field: reflect.StructField{Tag: reflect.StructTag(`json:"test" variable:"example=12;label;"`), Type: reflect.TypeOf("test")},
			want: &Variable{
				Name:    base.Child("test").String(),
				Example: "12",
			},
		},
		"parse tag without value": {
			base:  base,
			field: reflect.StructField{Tag: reflect.StructTag(`json:"test"`), Type: reflect.TypeOf("test")},
			want: &Variable{
				Name: base.Child("test").String(),
			},
		},
		"parse tag with perfix key": {
			base:  base,
			field: reflect.StructField{Tag: reflect.StructTag(`json:"test" variable:"label2=test;example=12"`), Type: reflect.TypeOf("test")},
			want: &Variable{
				Name:    base.Child("test").String(),
				Example: "12",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := NewVariable(tt.field, tt.base); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}
