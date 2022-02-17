/*
Copyright 2021 The Katanomi Authors.

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

	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateDuplicatedName(t *testing.T) {
	// g := NewGomegaWithT(t)

	table := map[string]struct {
		Name       string
		FieldPath  *field.Path
		Set        sets.String
		Evaluation func(g *WithT, errs field.ErrorList)
	}{
		"Already added, should error": {
			"abc",
			field.NewPath("abc"),
			sets.NewString("abc"),
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(1))
			},
		},
		"Empty set, should be ok": {
			"abc",
			field.NewPath("abc"),
			sets.NewString(),
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
	}

	for i, test := range table {
		t.Run(i, func(t *testing.T) {
			g := NewGomegaWithT(t)
			errs := ValidateDuplicatedName(test.FieldPath, test.Name, test.Set)
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
	}

	for i, test := range table {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			g := NewGomegaWithT(t)
			errs := ValidateGenericResourceName(test.Name, test.FieldPath)
			test.Evaluation(g, errs)
		})
	}

}
