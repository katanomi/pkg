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

package references

import (
	"testing"

	ktesting "github.com/AlaudaDevops/pkg/testing"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestFindOwnerReference(t *testing.T) {
	g := NewGomegaWithT(t)

	cmWithRefs := &corev1.ConfigMap{}
	cmWithoutRefs := &corev1.ConfigMap{}
	ktesting.MustLoadYaml("testdata/findOwnerReference.cm.withRefs.yaml", cmWithRefs)
	ktesting.MustLoadYaml("testdata/findOwnerReference.cm.withoutRefs.yaml", cmWithoutRefs)

	gvk := schema.GroupVersionKind{
		Version: "v1",
		Kind:    "ConfigMap",
	}
	notExistGvk := schema.GroupVersionKind{
		Group:   "not-exist",
		Version: "not-exist",
		Kind:    "not-exist",
	}
	g.Expect(FindOwnerReference(cmWithoutRefs.GetOwnerReferences(), gvk)).To(BeNil())
	g.Expect(FindOwnerReference(cmWithoutRefs.GetOwnerReferences(), notExistGvk)).To(BeNil())
	g.Expect(FindOwnerReference(cmWithRefs.GetOwnerReferences(), gvk)).NotTo(BeNil())
	g.Expect(FindOwnerReference(cmWithRefs.GetOwnerReferences(), notExistGvk)).To(BeNil())
}

func TestFindOwnerReferenceWithGroupKind(t *testing.T) {
	g := NewGomegaWithT(t)

	cmWithRefs := &corev1.ConfigMap{}
	cmWithoutRefs := &corev1.ConfigMap{}
	ktesting.MustLoadYaml("testdata/findOwnerReference.cm.withRefs.yaml", cmWithRefs)
	ktesting.MustLoadYaml("testdata/findOwnerReference.cm.withoutRefs.yaml", cmWithoutRefs)

	cmGroupKind := cmWithRefs.GroupVersionKind().GroupKind()
	noExistGroupKind := schema.GroupKind{
		Group: "not-exist",
		Kind:  "not-exist",
	}
	g.Expect(FindOwnerReferenceWithGroupKind(cmWithoutRefs.GetOwnerReferences(), cmGroupKind)).To(BeNil())
	g.Expect(FindOwnerReferenceWithGroupKind(cmWithoutRefs.GetOwnerReferences(), noExistGroupKind)).To(BeNil())
	g.Expect(FindOwnerReferenceWithGroupKind(cmWithRefs.GetOwnerReferences(), cmGroupKind)).NotTo(BeNil())
	g.Expect(FindOwnerReferenceWithGroupKind(cmWithRefs.GetOwnerReferences(), noExistGroupKind)).To(BeNil())
}
