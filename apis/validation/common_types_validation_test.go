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
	"testing"

	ktesting "github.com/katanomi/pkg/testing"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateName(t *testing.T) {
	g := NewGomegaWithT(t)

	table := map[string]struct {
		Object     metav1.Object
		FieldPath  *field.Path
		Evaluation func(g *WithT, errs field.ErrorList)
	}{
		"Invalid name with caps and underscore \"113_-Aabc\"": {
			ktesting.LoadObjectOrDie(g, "testdata/pod-1abc.yaml", &corev1.Pod{}, ktesting.SetName("113_-Aabc")),
			nil,
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(1))
			},
		},
		"Invalid name with space \"113 abc\"": {
			ktesting.LoadObjectOrDie(g, "testdata/pod-1abc.yaml", &corev1.Pod{}, ktesting.SetName("1113 abc")),
			nil,
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(1))
			},
		},
		"Invalid name with underscore \"abc_abc\"": {
			ktesting.LoadObjectOrDie(g, "testdata/pod-1abc.yaml", &corev1.Pod{}, ktesting.SetName("abc_abc")),
			nil,
			func(g *WithT, errs field.ErrorList) {
				g.Expect(errs).To(HaveLen(1))
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
			errs := ValidateName(test.Object, test.FieldPath)
			test.Evaluation(g, errs)
		})
	}

}
