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

package v1alpha1

import (
	"testing"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCopyLabelsAnnotations(t *testing.T) {

	t.Run("both nil", func(t *testing.T) {
		g := NewGomegaWithT(t)

		left := &corev1.Secret{}
		right := &corev1.Secret{}
		CopyLabels(left, right)
		CopyAnnotations(left, right)
		g.Expect(right.Labels).To(BeNil())
		g.Expect(right.Annotations).To(BeNil())
	})

	t.Run("left nil, right with value", func(t *testing.T) {
		g := NewGomegaWithT(t)

		left := &corev1.Secret{}
		right := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Labels:      map[string]string{"abc": "def"},
				Annotations: map[string]string{"abc": "def"},
			},
		}
		CopyLabels(left, right)
		CopyAnnotations(left, right)
		g.Expect(right.Labels).To(Equal(map[string]string{"abc": "def"}))
		g.Expect(right.Annotations).To(Equal(map[string]string{"abc": "def"}))
	})

	t.Run("left with value, right nil", func(t *testing.T) {
		g := NewGomegaWithT(t)

		left := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Labels:      map[string]string{"xyz": "qwe"},
				Annotations: map[string]string{"cxz": "dsa"},
			},
		}
		right := &corev1.Secret{}
		CopyLabels(left, right)
		CopyAnnotations(left, right)
		g.Expect(right.Labels).To(Equal(map[string]string{"xyz": "qwe"}))
		g.Expect(right.Annotations).To(Equal(map[string]string{"cxz": "dsa"}))
	})
}

func TestHasValueLabelsAnnotations(t *testing.T) {

	t.Run("HasAnnotations", func(t *testing.T) {
		g := NewGomegaWithT(t)

		obj := &corev1.Secret{}
		g.Expect(HasAnnotation(obj, "abc", "def")).To(BeFalse())

		obj.SetAnnotations(map[string]string{"abc": "def"})

		g.Expect(HasAnnotation(obj, "abc", "def")).To(BeTrue())
		g.Expect(HasAnnotation(obj, "xyz", "qwe")).To(BeFalse())
	})

	t.Run("HasLabels", func(t *testing.T) {
		g := NewGomegaWithT(t)

		obj := &corev1.Secret{}
		g.Expect(HasLabel(obj, "abc", "def")).To(BeFalse())

		obj.SetLabels(map[string]string{"abc": "def"})

		g.Expect(HasLabel(obj, "abc", "def")).To(BeTrue())
		g.Expect(HasLabel(obj, "xyz", "qwe")).To(BeFalse())
	})

}

func TestFilterMapKeys(t *testing.T) {
	g := NewGomegaWithT(t)
	tests := []struct {
		mapObj         map[string]string
		ignorePrefixes []string
		want           map[string]string
	}{
		{
			mapObj:         nil,
			ignorePrefixes: nil,
			want:           nil,
		},
		{
			mapObj:         map[string]string{"a": "b"},
			ignorePrefixes: []string{"a"},
			want:           map[string]string{},
		},
		{
			mapObj:         map[string]string{"a": "b"},
			ignorePrefixes: []string{"b"},
			want:           map[string]string{"a": "b"},
		},
		{
			mapObj:         map[string]string{"abc": "b"},
			ignorePrefixes: []string{"ab"},
			want:           map[string]string{},
		},
		{
			mapObj:         map[string]string{"abc": "b"},
			ignorePrefixes: []string{"ac", "ab"},
			want:           map[string]string{},
		},
	}
	for _, tt := range tests {
		g.Expect(FilterMapKeys(tt.mapObj, tt.ignorePrefixes...)).To(Equal(tt.want))
	}
}

func TestHasLabelsAnnotationsKey(t *testing.T) {

	t.Run("HasAnnotationKey", func(t *testing.T) {
		g := NewGomegaWithT(t)

		obj := &corev1.Secret{}
		g.Expect(HasAnnotationKey(obj, "abc")).To(BeFalse())

		obj.SetAnnotations(map[string]string{"abc": "def", "empty": ""})

		g.Expect(HasAnnotationKey(obj, "abc")).To(BeTrue())
		g.Expect(HasAnnotationKey(obj, "xyz")).To(BeFalse())
		g.Expect(HasAnnotationKey(obj, "empty")).To(BeTrue())
	})

	t.Run("HasLabelKey", func(t *testing.T) {
		g := NewGomegaWithT(t)

		obj := &corev1.Secret{}
		g.Expect(HasLabelKey(obj, "abc")).To(BeFalse())

		obj.SetLabels(map[string]string{"abc": "def", "empty": ""})

		g.Expect(HasLabelKey(obj, "abc")).To(BeTrue())
		g.Expect(HasLabelKey(obj, "xyz")).To(BeFalse())
		g.Expect(HasLabelKey(obj, "empty")).To(BeTrue())
	})

}
