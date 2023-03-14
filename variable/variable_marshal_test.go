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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	. "github.com/katanomi/pkg/testing"
	. "github.com/onsi/gomega"
	"k8s.io/utils/field"
)

func TestConvertToVariableList(t *testing.T) {
	t.Run("test BuildRunGitStatus", func(t *testing.T) {
		g := NewGomegaWithT(t)

		expected := v1alpha1.VariableList{}
		g.Expect(LoadYAML("testdata/converttovariablelist.buildrungitstatus.golden.json", &expected)).To(Succeed())
		marshaller := VariableMarshaller{Object: v1alpha1.BuildRunGitStatus{}}
		got, err := marshaller.Marshal()
		g.Expect(err).To(Succeed())
		diff := cmp.Diff(got, expected)
		g.Expect(diff).To(BeEmpty())
	})

	t.Run("test TriggeredBy", func(t *testing.T) {
		g := NewGomegaWithT(t)

		expected := v1alpha1.VariableList{}
		g.Expect(LoadYAML("testdata/converttovariablelist.triggeredby.golden.json", &expected)).To(Succeed())
		marshaller := VariableMarshaller{Object: v1alpha1.TriggeredBy{}}
		got, err := marshaller.Marshal()
		g.Expect(err).To(Succeed())
		diff := cmp.Diff(got, expected)
		g.Expect(diff).To(BeEmpty())
	})

	t.Run("test object is nil", func(t *testing.T) {
		g := NewGomegaWithT(t)

		expected := v1alpha1.VariableList{}
		marshaller := VariableMarshaller{}
		got, err := marshaller.Marshal()
		g.Expect(err).To(Succeed())
		g.Expect(expected).To(Equal(got))
	})

	t.Run("test BuildRunGitStatus with label filter", func(t *testing.T) {
		g := NewGomegaWithT(t)

		expected := v1alpha1.VariableList{}
		g.Expect(LoadYAML("testdata/converttovariablelist.buildrungitstatus.build.golden.json", &expected)).To(Succeed())
		marshaller := VariableMarshaller{Object: v1alpha1.BuildRunGitStatus{}}
		got, err := marshaller.Marshal()
		g.Expect(err).To(Succeed())

		got.Filter(v1alpha1.LabelFilter("default"))
		diff := cmp.Diff(got, expected)
		g.Expect(diff).To(BeEmpty())
	})

	t.Run("has unsupported type", func(t *testing.T) {
		g := NewGomegaWithT(t)

		obj := struct {
			Number int    `json:"number"`
			Func   func() `json:"channel"`
		}{
			Number: 1,
			Func:   func() {},
		}

		marshaller := VariableMarshaller{Object: obj}
		_, err := marshaller.Marshal()
		g.Expect(err).To(Equal(fmt.Errorf("unsupported type [%s]", reflect.Func.String())))
	})
}

func TestNewVariable(t *testing.T) {
	base := field.NewPath("test")
	tests := map[string]struct {
		field reflect.StructField
		base  *field.Path
		want  *v1alpha1.Variable
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
			want: &v1alpha1.Variable{
				Name:    base.Child("test").String(),
				Example: "12",
				Label:   "test",
			},
		},
		"parse tag without tag value": {
			base:  base,
			field: reflect.StructField{Tag: reflect.StructTag(`json:"test" variable:"example=12;label;"`), Type: reflect.TypeOf("test")},
			want: &v1alpha1.Variable{
				Name:    base.Child("test").String(),
				Example: "12",
			},
		},
		"parse tag without value": {
			base:  base,
			field: reflect.StructField{Tag: reflect.StructTag(`json:"test"`), Type: reflect.TypeOf("test")},
			want: &v1alpha1.Variable{
				Name: base.Child("test").String(),
			},
		},
		"parse tag with perfix key": {
			base:  base,
			field: reflect.StructField{Tag: reflect.StructTag(`json:"test" variable:"label2=test;example=12"`), Type: reflect.TypeOf("test")},
			want: &v1alpha1.Variable{
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
