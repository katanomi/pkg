/*
Copyright 2022 The AlaudaDevops Authors.

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
	"testing"

	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateRegExp(t *testing.T) {
	table := map[string]struct {
		Pattern    string
		FieldPath  *field.Path
		Evaluation func(g *WithT, errs field.ErrorList)
	}{
		"empty string": {
			``,
			field.NewPath("abc"),
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
		"normal string": {
			`abc`,
			field.NewPath("abc"),
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
		"correct regular expressions": {
			`^abc.*z$`,
			field.NewPath("abc"),
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
		"complex regular expressions": {
			`v(?P<major>\d+)\.(?P<minor>\d+).*(?P<patch>\d+)*`,
			field.NewPath("abc"),
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(0))
			},
		},
		"incorrect regular expressions": {
			`^(1234\$`,
			field.NewPath("abc"),
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(1))
			},
		},
	}

	for i, test := range table {
		t.Run(i, func(t *testing.T) {
			g := NewGomegaWithT(t)
			errs := ValidateRegExp(test.Pattern, test.FieldPath)
			test.Evaluation(g, errs)
		})
	}
}
