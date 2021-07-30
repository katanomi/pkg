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
	"testing"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

func TestObjectReferenceIsTheSame(t *testing.T) {
	tests := map[string]struct {
		Object   *corev1.ObjectReference
		Compared *corev1.ObjectReference

		Expected bool
	}{
		"Both Nil": {
			Object:   nil,
			Compared: nil,
			Expected: true,
		},
		"Both Not Nil": {
			Object:   &corev1.ObjectReference{},
			Compared: &corev1.ObjectReference{},
			Expected: true,
		},
		"Mixed Nil": {
			Object:   nil,
			Compared: &corev1.ObjectReference{},
			Expected: false,
		},
		"Reversed Mixed Nil": {
			Object:   &corev1.ObjectReference{},
			Compared: nil,
			Expected: false,
		},
		"Same reference": {
			Object: &corev1.ObjectReference{
				Name:      "abc",
				Namespace: "default",
			},
			Compared: &corev1.ObjectReference{
				Name:      "abc",
				Namespace: "default",
			},
			Expected: true,
		},
		"Different reference": {
			Object: &corev1.ObjectReference{
				Name:      "def",
				Namespace: "default",
			},
			Compared: &corev1.ObjectReference{
				Name:      "abc",
				Namespace: "default",
			},
			Expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			g.Expect((IsTheSameObjectReference(test.Object, test.Compared))).To(Equal(test.Expected))
		})
	}
}
