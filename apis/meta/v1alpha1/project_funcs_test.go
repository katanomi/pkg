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

package v1alpha1

import (
	"strconv"
	"testing"

	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
)

func TestProjectAddNamespaceRef(t *testing.T) {
	cases := []struct {
		Refs     []v1.ObjectReference
		Expected int
	}{
		{
			Refs: []v1.ObjectReference{
				{
					Name: "test-1",
				},
				{
					Name: "test-1",
				},
			},
			Expected: 1,
		},
		{
			Refs: []v1.ObjectReference{
				{
					Name: "test-1",
				},
				{
					Name: "test-2",
				},
			},
			Expected: 2,
		},
		{
			Refs: []v1.ObjectReference{
				{
					Name: "test-1",
				},
				{
					Name: "test-2",
				},
				{
					Name: "test-1",
				},
			},
			Expected: 2,
		},
	}

	for i := range cases {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			g := NewGomegaWithT(t)
			p := &Project{}
			p.AddNamespaceRef(cases[i].Refs...)
			g.Expect(p.Spec.NamespaceRefs).To(HaveLen(cases[i].Expected))
		})
	}
}
