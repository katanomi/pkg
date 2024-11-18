/*
Copyright 2021 The AlaudaDevops Authors.

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

package validation

import (
	"fmt"
	"testing"

	ktesting "github.com/AlaudaDevops/pkg/testing"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateItemNameUnderscore(t *testing.T) {
	g := NewGomegaWithT(t)

	table := map[string]struct {
		Object     metav1.Object
		FieldPath  *field.Path
		Evaluation func(g *WithT, errs field.ErrorList)
	}{
		"Valid name with caps and underscore \"113_-Aabc\"": {
			ktesting.LoadObjectOrDie(g, "testdata/pod-1abc.yaml", &corev1.Pod{}, ktesting.SetName("113_-Aabc")),
			nil,
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
		"Valid name with caps slash and underscore \"My-Name113_ISAabc\"": {
			ktesting.LoadObjectOrDie(g, "testdata/pod-1abc.yaml", &corev1.Pod{}, ktesting.SetName("My-Name113_ISAabc")),
			nil,
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
		"Invalid name with space \"113 abc\"": {
			ktesting.LoadObjectOrDie(g, "testdata/pod-1abc.yaml", &corev1.Pod{}, ktesting.SetName("1113 abc")),
			nil,
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(1))
			},
		},
		"Invalid name over 63 characters": {
			ktesting.LoadObjectOrDie(g, "testdata/pod-1abc.yaml", &corev1.Pod{}, ktesting.SetName("0123456789001234567890012345678900123456789001234567890012345678901234")),
			nil,
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(1))
			},
		},
		"Valid name with underscore \"abc_abc\"": {
			ktesting.LoadObjectOrDie(g, "testdata/pod-1abc.yaml", &corev1.Pod{}, ktesting.SetName("abc_abc")),
			nil,
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
		"Valid name \"123-abc\"": {
			ktesting.LoadObjectOrDie(g, "testdata/pod-1abc.yaml", &corev1.Pod{}, ktesting.SetName("123-abc")),
			nil,
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
	}

	for i, test := range table {
		t.Run(i, func(t *testing.T) {
			g := NewGomegaWithT(t)
			errs := ValidateItemNameUnderscore(test.Object.GetName(), field.NewPath("metadata", "name"))
			test.Evaluation(g, errs)
		})
	}
}

func TestValidateGenericResourceName(t *testing.T) {
	// g := NewGomegaWithT(t)
	table := []struct {
		Name       string
		FieldPath  *field.Path
		Evaluation func(g *WithT, errs field.ErrorList)
	}{

		{
			"",
			field.NewPath("a"),
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(1))
			},
		},
		{
			"abc",
			field.NewPath("abc"),
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
		{
			"abc-def",
			field.NewPath("def"),
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
		{
			"abc-def/gh",
			field.NewPath("gh"),
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
		{
			"abc-def/g-h",
			field.NewPath("g-h"),
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
		{
			"abc-def/g-h_i",
			field.NewPath("ghi"),
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
		{
			"abc-def/g-h_i.j",
			field.NewPath("j"),
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
		{
			"Abc-def/g-h_i.j",
			field.NewPath("j2"),
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
		{
			"Abc-def/g-h_i:j",
			field.NewPath("j3"),
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(1))
			},
		},
		{
			"测试",
			field.NewPath("j3"),
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(1))
			},
		},
	}

	for i, test := range table {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			g := NewGomegaWithT(t)
			errs := ValidateGenericResourceName(test.Name, test.FieldPath)
			test.Evaluation(g, errs)
		})
	}
}

func TestValidateResourceNameWithChinese(t *testing.T) {
	g := NewGomegaWithT(t)

	table := map[string]struct {
		Object     metav1.Object
		FieldPath  *field.Path
		Evaluation func(g *WithT, errs field.ErrorList)
	}{
		"Valid name with caps and underscore \"113_-Aabc\"": {
			ktesting.LoadObjectOrDie(g, "testdata/pod-1abc.yaml", &corev1.Pod{}, ktesting.SetName("113_-Aabc")),
			nil,
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
		"Valid name with caps slash and underscore \"My-Name113_ISAabc\"": {
			ktesting.LoadObjectOrDie(g, "testdata/pod-1abc.yaml", &corev1.Pod{}, ktesting.SetName("My-Name113_ISAabc")),
			nil,
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
		"Valid name with underscore \"abc_abc\"": {
			ktesting.LoadObjectOrDie(g, "testdata/pod-1abc.yaml", &corev1.Pod{}, ktesting.SetName("abc_abc")),
			nil,
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
		"Valid name with slash \"123/abc\"": {
			ktesting.LoadObjectOrDie(g, "testdata/pod-1abc.yaml", &corev1.Pod{}, ktesting.SetName("123/abc")),
			nil,
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
		"Valid name with dot \"123.abc\"": {
			ktesting.LoadObjectOrDie(g, "testdata/pod-1abc.yaml", &corev1.Pod{}, ktesting.SetName("123.abc")),
			nil,
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
		"Valid name \"123-abc\"": {
			ktesting.LoadObjectOrDie(g, "testdata/pod-1abc.yaml", &corev1.Pod{}, ktesting.SetName("123-abc")),
			nil,
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
		"Valid name with single char n": {
			ktesting.LoadObjectOrDie(g, "testdata/pod-1abc.yaml", &corev1.Pod{}, ktesting.SetName("n")),
			nil,
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
		"Valid name with space \"123 abc\"": {
			ktesting.LoadObjectOrDie(g, "testdata/pod-1abc.yaml", &corev1.Pod{}, ktesting.SetName("1113 abc")),
			nil,
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
		"Invalid name starts with space \" 123abc\"": {
			ktesting.LoadObjectOrDie(g, "testdata/pod-1abc.yaml", &corev1.Pod{}, ktesting.SetName(" 1113abc")),
			nil,
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(1))
				g.Expect(errs[0].Error()).To(ContainSubstring("Chinese"))
			},
		},
		"Invalid name with ? \"123?abc\"": {
			ktesting.LoadObjectOrDie(g, "testdata/pod-1abc.yaml", &corev1.Pod{}, ktesting.SetName("123?abc")),
			nil,
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(1))
				g.Expect(errs[0].Error()).To(ContainSubstring("Chinese"))
			},
		},
		"Invalid name with 。 \"123abc。\"": {
			ktesting.LoadObjectOrDie(g, "testdata/pod-1abc.yaml", &corev1.Pod{}, ktesting.SetName("123abc。")),
			nil,
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(1))
				g.Expect(errs[0].Error()).To(ContainSubstring("Chinese"))
			},
		},
		"Invalid name with emoji \"123abc😀😃😄😁\"": {
			ktesting.LoadObjectOrDie(g, "testdata/pod-1abc.yaml", &corev1.Pod{}, ktesting.SetName("123abc😀😃😄😁")),
			nil,
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(1))
				g.Expect(errs[0].Error()).To(ContainSubstring("Chinese"))
			},
		},
		"Valid name with chinese \"123-abc_测试\"": {
			ktesting.LoadObjectOrDie(g, "testdata/pod-1abc.yaml", &corev1.Pod{}, ktesting.SetName("123-abc_测试")),
			nil,
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
		"Valid name with chinese \"测试苼膍茊\"": {
			ktesting.LoadObjectOrDie(g, "testdata/pod-1abc.yaml", &corev1.Pod{}, ktesting.SetName("测试苼膍茊")),
			nil,
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
	}

	for i, test := range table {
		t.Run(i, func(t *testing.T) {
			g := NewGomegaWithT(t)
			errs := ValidateResourceNameWithChinese(test.Object.GetName(), field.NewPath("metadata", "name"))
			test.Evaluation(g, errs)
		})
	}
}
