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
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/utils/field"
)

func TestMarshalSubject(t *testing.T) {
	t.Run("base is empty", func(t *testing.T) {
		g := NewGomegaWithT(t)

		got, err := MarshalSubject(reflect.TypeOf(rbacv1.Subject{}), nil, nil)
		g.Expect(err).To(BeNil())

		want := []Variable{
			{Name: field.NewPath("kind").String(), Example: "User"},
			{Name: field.NewPath("apiGroup").String()},
			{Name: field.NewPath("name").String(), Example: "joedoe@example.com"},
			{Name: field.NewPath("namespace").String(), Example: "default"},
		}
		diff := cmp.Diff(got, want)
		g.Expect(diff).To(BeEmpty())
	})

	t.Run("base has set", func(t *testing.T) {
		g := NewGomegaWithT(t)

		base := field.NewPath("base")
		got, err := MarshalSubject(reflect.TypeOf(rbacv1.Subject{}), base, nil)
		g.Expect(err).To(BeNil())
		want := []Variable{
			{Name: base.Child("kind").String(), Example: "User"},
			{Name: base.Child("apiGroup").String()},
			{Name: base.Child("name").String(), Example: "joedoe@example.com"},
			{Name: base.Child("namespace").String(), Example: "default"},
		}
		diff := cmp.Diff(got, want)
		g.Expect(diff).To(BeEmpty())
	})

	t.Run("input type isn't match", func(t *testing.T) {
		g := NewGomegaWithT(t)

		type Subject struct{}
		base := field.NewPath("base")
		want := []Variable{}

		st := reflect.TypeOf(Subject{})
		got, err := MarshalSubject(st, base, nil)
		g.Expect(err).To(Equal(fmt.Errorf(
			"get marshal type[%s/%s] don't match %s/%s",
			st.PkgPath(), st.Name(), rbacv1SubjectPkgPath, rbacv1SubjectName)))

		diff := cmp.Diff(got, want)
		g.Expect(diff).To(BeEmpty())
	})
}

func TestMarshalObjectReference(t *testing.T) {
	t.Run("base is empty", func(t *testing.T) {
		g := NewGomegaWithT(t)

		got, err := MarshalObjectReference(reflect.TypeOf(corev1.ObjectReference{}), nil, nil)
		g.Expect(err).To(BeNil())

		want := []Variable{
			{Name: field.NewPath("kind").String(), Example: "DeliveryRun"},
			{Name: field.NewPath("namespace").String(), Example: "default"},
			{Name: field.NewPath("name").String(), Example: "delivery-run-abdexy"},
			{Name: field.NewPath("uid").String(), Example: "b2fab970-f672-4af0-a9cd-5ad9a8dbcc29"},
			{Name: field.NewPath("apiVersion").String(), Example: "deliveries.katanomi.dev/v1alpha1"},
		}
		diff := cmp.Diff(got, want)
		g.Expect(diff).To(BeEmpty())
	})

	t.Run("base has set", func(t *testing.T) {
		g := NewGomegaWithT(t)

		base := field.NewPath("base")
		got, err := MarshalObjectReference(reflect.TypeOf(corev1.ObjectReference{}), base, nil)
		g.Expect(err).To(BeNil())
		want := []Variable{
			{Name: base.Child("kind").String(), Example: "DeliveryRun"},
			{Name: base.Child("namespace").String(), Example: "default"},
			{Name: base.Child("name").String(), Example: "delivery-run-abdexy"},
			{Name: base.Child("uid").String(), Example: "b2fab970-f672-4af0-a9cd-5ad9a8dbcc29"},
			{Name: base.Child("apiVersion").String(), Example: "deliveries.katanomi.dev/v1alpha1"},
		}
		diff := cmp.Diff(got, want)
		g.Expect(diff).To(BeEmpty())
	})

	t.Run("input type isn't match", func(t *testing.T) {
		g := NewGomegaWithT(t)

		type ObjectReference struct{}
		base := field.NewPath("base")
		want := []Variable{}

		st := reflect.TypeOf(ObjectReference{})
		got, err := MarshalObjectReference(st, base, nil)
		g.Expect(err).To(Equal(fmt.Errorf(
			"get marshal type[%s/%s] don't match %s/%s",
			st.PkgPath(), st.Name(), corev1ObjectReferencePkgPath, corev1ObjectReferenceName)))

		diff := cmp.Diff(got, want)
		g.Expect(diff).To(BeEmpty())
	})
}

func TestMarshalBuildGitBranchStatus(t *testing.T) {
	t.Run("object is't macth", func(t *testing.T) {
		g := NewGomegaWithT(t)

		type BuildGitBranchStatus struct{}
		base := field.NewPath("base")
		want := []Variable{}

		st := reflect.TypeOf(BuildGitBranchStatus{})
		got, err := MarshalBuildGitBranchStatus(st, base, nil)
		g.Expect(err).To(Equal(fmt.Errorf(
			"get marshal type[%s/%s] don't match %s/%s",
			st.PkgPath(), st.Name(), buildGitBranchStatusPkgPath, buildGitBranchStatusName)))

		diff := cmp.Diff(got, want)
		g.Expect(diff).To(BeEmpty())
	})

	t.Run("Is branch base", func(t *testing.T) {
		g := NewGomegaWithT(t)

		base := field.NewPath("test", "branch")
		want := []Variable{
			{Name: "test.branch.name", Example: "main", Label: "default"},
			{Name: "test.branch.protected", Example: "true"},
			{Name: "test.branch.default", Example: "true"},
			{Name: "test.branch.webURL", Example: "https://github.com/repository/tree/main"},
		}

		st := reflect.TypeOf(v1alpha1.BuildGitBranchStatus{})
		got, err := MarshalBuildGitBranchStatus(st, base, &VariableMarshaller{Object: v1alpha1.BuildGitBranchStatus{}})

		g.Expect(err).To(BeNil())
		diff := cmp.Diff(got, want)
		g.Expect(diff).To(BeEmpty())
	})

	t.Run("Is't branch base", func(t *testing.T) {
		g := NewGomegaWithT(t)

		base := field.NewPath("test")
		want := []Variable{
			{Name: "test.name", Example: "main"},
			{Name: "test.protected", Example: "true"},
			{Name: "test.default", Example: "true"},
			{Name: "test.webURL", Example: "https://github.com/repository/tree/main"},
		}

		st := reflect.TypeOf(v1alpha1.BuildGitBranchStatus{})
		got, err := MarshalBuildGitBranchStatus(st, base, &VariableMarshaller{Object: v1alpha1.BuildGitBranchStatus{}})

		g.Expect(err).To(BeNil())
		diff := cmp.Diff(got, want)
		g.Expect(diff).To(BeEmpty())
	})
}
