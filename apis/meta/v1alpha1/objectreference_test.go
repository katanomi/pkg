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

	"context"

	. "github.com/katanomi/pkg/testing"
	. "github.com/onsi/ginkgo/v2"

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
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

func TestGetObjectReferenceFromObject(t *testing.T) {
	g := NewGomegaWithT(t)
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod",
			Namespace: "default",
			UID:       types.UID("abc"),
		},
	}

	ref := GetObjectReferenceFromObject(pod)
	g.Expect(ref.Name).To(Equal("pod"))
	g.Expect(ref.Namespace, ref.APIVersion, ref.Kind, ref.UID).To(BeEmpty())

	ref = GetObjectReferenceFromObject(pod, ObjectRefWithTypeMeta(), ObjectRefWithNamespace(), ObjectRefWithUID())
	g.Expect(ref.Name).To(Equal("pod"))
	g.Expect(ref.Namespace).To(Equal("default"))
	g.Expect(ref.APIVersion).To(Equal("v1"))
	g.Expect(ref.Kind).To(Equal("Pod"))
	g.Expect(ref.UID).To(Equal(types.UID("abc")))
}

func TestGetNamespacedNameFromRef(t *testing.T) {
	table := map[string]struct {
		Object *corev1.ObjectReference
		Result types.NamespacedName
	}{
		"Simple secret object": {
			Object: &corev1.ObjectReference{Name: "secret", Namespace: "default"},
			Result: types.NamespacedName{Name: "secret", Namespace: "default"},
		},
		"Nil object": {
			Object: nil,
			Result: types.NamespacedName{},
		},
	}

	for name, item := range table {
		test := item
		t.Run(name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			result := GetNamespacedNameFromRef(test.Object)
			diff := cmp.Diff(test.Result, result)

			g.Expect(diff).To(BeEmpty())
		})
	}
}

var _ = Describe("ObjectReferenceValGetter.GetValWithKey", func() {
	var (
		ctx      context.Context
		path     *field.Path
		obj      *corev1.ObjectReference
		values   map[string]string
		expected map[string]string
	)
	BeforeEach(func() {
		ctx = context.TODO()
		path = field.NewPath("ref")
		obj = &corev1.ObjectReference{}
		expected = map[string]string{}
	})
	JustBeforeEach(func() {
		values = ObjectReferenceValGetter(obj)(ctx, path)
	})
	Context("corev1.ObjectReference with all variables", func() {
		BeforeEach(func() {
			Expect(LoadYAML("testdata/objectreference_vars.all.yaml", obj)).To(Succeed())
			Expect(LoadYAML("testdata/objectreference_vars.all.golden.yaml", &expected)).To(Succeed())
			Expect(expected).ToNot(BeEmpty())
		})
		It("should return the same amount of data", func() {
			Expect(values).To(Equal(expected))
		})
	})
})
