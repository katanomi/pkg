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
	"context"
	"testing"
	"time"

	. "github.com/katanomi/pkg/testing"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo/v2"
	"k8s.io/apimachinery/pkg/util/validation/field"

	. "github.com/onsi/gomega"
)

func TestFromAnnotation(t *testing.T) {
	t.Run("annotations has no triggeredby key", func(t *testing.T) {
		g := NewGomegaWithT(t)

		annotations := map[string]string{}
		by, err := (&TriggeredBy{}).FromAnnotation(annotations)
		g.Expect(err).Should(BeNil())
		g.Expect(by).Should(BeNil())
	})

	t.Run("annotations has triggededby key with invalid json", func(t *testing.T) {
		g := NewGomegaWithT(t)

		annotations := map[string]string{
			TriggeredByAnnotationKey: "{",
		}
		by, err := (&TriggeredBy{}).FromAnnotation(annotations)
		g.Expect(err).ShouldNot(BeNil())
		g.Expect(by).Should(BeNil())
	})

	t.Run("annotations has triggeredby key with valid json", func(t *testing.T) {
		g := NewGomegaWithT(t)

		annotations := map[string]string{
			TriggeredByAnnotationKey: `{
	"ref": { "kind": "Trigger", "apiVersion": "core.katanomi.dev/v1alpha1", "name": "git-trigger", "namespace": "default"},
	"cloudEvent": {"id":"11", "type":"Tag Hook", "source": "https://example.com/demo/demo"}
}`,
		}
		by, err := (&TriggeredBy{}).FromAnnotation(annotations)
		g.Expect(err).Should(BeNil())
		g.Expect(by).ShouldNot(BeNil())
		g.Expect(by.CloudEvent.Type).Should(BeEquivalentTo("Tag Hook"))
		g.Expect(by.CloudEvent.ID).Should(BeEquivalentTo("11"))
		g.Expect(by.CloudEvent.Source).Should(BeEquivalentTo("https://example.com/demo/demo"))
		g.Expect(by.Ref.Kind).Should(BeEquivalentTo("Trigger"))
	})

	t.Run("annotations has triggeredby key with user and triggeredTimestamp", func(t *testing.T) {
		g := NewGomegaWithT(t)

		annotations := map[string]string{
			TriggeredByAnnotationKey: `{"user":{"kind":"User","name":"admin"},"triggeredTimestamp":"2022-02-09T05:34:22Z"}`,
		}
		by, err := (&TriggeredBy{}).FromAnnotation(annotations)
		g.Expect(err).Should(BeNil())
		g.Expect(by).ShouldNot(BeNil())
		g.Expect(by.User.Name).Should(BeEquivalentTo("admin"))
		g.Expect(by.TriggeredTimestamp).ShouldNot(BeNil())
	})

	t.Run("by is nil", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var by *TriggeredBy = nil
		annotations := map[string]string{
			TriggeredByAnnotationKey: `{
	"ref": { "kind": "Trigger", "apiVersion": "core.katanomi.dev/v1alpha1", "name": "git-trigger", "namespace": "default"},
	"cloudEvent": {"id":"11", "type":"Tag Hook", "source": "https://example.com/demo/demo"}
}`,
		}

		by, err := by.FromAnnotation(annotations)
		g.Expect(err).Should(BeNil())
		g.Expect(by).ShouldNot(BeNil())
		g.Expect(by.CloudEvent.Type).Should(BeEquivalentTo("Tag Hook"))
		g.Expect(by.CloudEvent.ID).Should(BeEquivalentTo("11"))
		g.Expect(by.CloudEvent.Source).Should(BeEquivalentTo("https://example.com/demo/demo"))
		g.Expect(by.Ref.Kind).Should(BeEquivalentTo("Trigger"))

	})
}

func TestSetIntoAnnotation(t *testing.T) {
	t.Run("triggeredBy is ok", func(t *testing.T) {
		g := NewGomegaWithT(t)

		evt := &CloudEvent{
			ID:          "11",
			Source:      "https://example.com/demo/demo",
			Subject:     "",
			Type:        "Tag Hook",
			Data:        `{"a":"b"}`,
			Time:        metav1.Time{},
			SpecVersion: "v1",
		}

		by := &TriggeredBy{
			Ref: &corev1.ObjectReference{
				Kind:      "Trigger",
				Namespace: "default",
				Name:      "push-trigger",
			},
			CloudEvent: evt,
		}
		annotations := map[string]string{}
		annotations, err := by.SetIntoAnnotation(annotations)
		g.Expect(err).Should(BeNil())
		g.Expect(len(annotations)).ShouldNot(BeZero())

		by, err = (&TriggeredBy{}).FromAnnotation(annotations)
		g.Expect(err).Should(BeNil())
		g.Expect(by.Ref.Name).Should(BeEquivalentTo("push-trigger"))
		g.Expect(by.CloudEvent.Source).Should(BeEquivalentTo("https://example.com/demo/demo"))
		g.Expect(by.CloudEvent.Data).Should(BeEmpty())
	})

	t.Run("nil annotations", func(t *testing.T) {
		g := NewGomegaWithT(t)

		by := &TriggeredBy{
			Ref: &corev1.ObjectReference{
				Kind:      "Trigger",
				Namespace: "default",
				Name:      "push-trigger",
			},
			CloudEvent: nil,
		}
		var annotations map[string]string
		annotations, err := by.SetIntoAnnotation(annotations)
		g.Expect(err).Should(BeNil())
		g.Expect(len(annotations)).ShouldNot(BeZero())
	})
}

var _ = Describe("TriggeredBy.GetValWithKey", func() {
	var (
		ctx      context.Context
		path     *field.Path
		by       *TriggeredBy
		values   map[string]string
		expected map[string]string
	)
	BeforeEach(func() {
		ctx = context.TODO()
		path = field.NewPath("triggeredBy")
		by = &TriggeredBy{}
		expected = map[string]string{}
	})
	JustBeforeEach(func() {
		values = by.GetValWithKey(ctx, path)
	})
	Context("triggeredby with all variables", func() {
		BeforeEach(func() {
			Expect(LoadYAML("testdata/triggeredby_types_vars.all.yaml", by)).To(Succeed())
			Expect(LoadYAML("testdata/triggeredby_types_vars.all.golden.yaml", &expected)).To(Succeed())
			Expect(expected).ToNot(BeEmpty())
		})
		It("should return the same amount of data", func() {
			Expect(values).To(Equal(expected))
		})
	})
	Context("triggeredby nil", func() {
		BeforeEach(func() {
			by = nil
			Expect(LoadYAML("testdata/triggeredby_types_vars.nil.golden.yaml", &expected)).To(Succeed())
			Expect(expected).ToNot(BeEmpty())
		})
		It("should return the same amount of data", func() {
			Expect(values).To(Equal(expected))
		})
	})
})

var _ = Describe("TriggeredBy.IsZero", func() {
	DescribeTable("TriggeredBy.IsZero",
		func(input *TriggeredBy, expected bool) {
			actual := input.IsZero()
			Expect(actual).To(Equal(expected))
		},
		Entry("nil pointer",
			nil,
			true,
		),
		Entry("no user",
			&TriggeredBy{},
			true,
		),
		Entry("user is not nil",
			&TriggeredBy{
				User: &rbacv1.Subject{
					Name: "user1",
				},
			},
			false,
		),
		Entry("cloudEvent is not nil",
			&TriggeredBy{
				CloudEvent: &CloudEvent{},
			},
			false,
		),
		Entry("ref is not nil",
			&TriggeredBy{
				Ref: &corev1.ObjectReference{},
			},
			false,
		),
		Entry("triggered timestamp is not zero",
			&TriggeredBy{
				TriggeredTimestamp: &metav1.Time{Time: time.Now()},
			},
			false,
		),
		Entry("triggered type is not empty",
			&TriggeredBy{
				TriggeredType: DefinitionTriggeredTypeValues.Manual,
			},
			false,
		),
	)
})
