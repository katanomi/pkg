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
)

func TestConvertToVariableList(t *testing.T) {
	t.Run("test BuildRunGitStatus", func(t *testing.T) {
		g := NewGomegaWithT(t)

		expected := VariableList{}
		g.Expect(LoadYAML("testdata/converttovariablelist.buildrungitstatus.golden.json", &expected)).To(Succeed())
		convertor := VariableConverter{}
		got, err := convertor.ConvertToVariableList(v1alpha1.BuildRunGitStatus{})
		g.Expect(err).To(Succeed())
		diff := cmp.Diff(got, expected)
		g.Expect(diff).To(BeEmpty())
	})

	t.Run("test TriggeredBy", func(t *testing.T) {
		g := NewGomegaWithT(t)

		expected := VariableList{}
		g.Expect(LoadYAML("testdata/converttovariablelist.triggeredby.golden.json", &expected)).To(Succeed())
		convertor := VariableConverter{}
		got, err := convertor.ConvertToVariableList(v1alpha1.TriggeredBy{})
		g.Expect(err).To(Succeed())
		diff := cmp.Diff(got, expected)
		g.Expect(diff).To(BeEmpty())
	})

	t.Run("test BuildRunGitStatus with label", func(t *testing.T) {
		g := NewGomegaWithT(t)

		expected := VariableList{}
		g.Expect(LoadYAML("testdata/converttovariablelist.buildrungitstatus.build.golden.json", &expected)).To(Succeed())
		convertor := VariableConverter{}
		got, err := convertor.ConvertToVariableList(v1alpha1.BuildRunGitStatus{}, LabelFilter("common"))
		g.Expect(err).To(Succeed())
		diff := cmp.Diff(got, expected)
		g.Expect(diff).To(BeEmpty())
	})

	t.Run("has unsupported type", func(t *testing.T) {
		g := NewGomegaWithT(t)

		convertor := VariableConverter{}
		obj := struct {
			Number int    `json:"number"`
			Func   func() `json:"channel"`
		}{
			Number: 1,
			Func:   func() {},
		}

		_, err := convertor.ConvertToVariableList(obj)
		g.Expect(err).To(Equal(fmt.Errorf("unsupported type [%s]", reflect.Func.String())))
	})
}
