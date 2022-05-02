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

	ktesting "github.com/katanomi/pkg/testing"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
)

type SomeStatusObj struct {
	conditions ObjectConditions
}

func (c *SomeStatusObj) GetObjectConditions() ObjectConditions {
	return c.conditions
}

func (c *SomeStatusObj) SetObjectConditions(objs ObjectConditions) {
	c.conditions = objs
}

func TestObjectConditionsSet(t *testing.T) {
	g := NewGomegaWithT(t)

	abcPod := ObjectCondition{}
	g.Expect(ktesting.LoadYAML("testdata/objectcondition-pod-abc.yaml", &abcPod)).To(Succeed())

	t.Run("MarkTrue", func(t *testing.T) {
		g := NewGomegaWithT(t)
		objcs := ObjectConditions{}
		g.Expect(ktesting.LoadYAML("testdata/objectconditions.yaml", &objcs)).To(Succeed())

		obj := &SomeStatusObj{conditions: objcs}
		mgr := objcs.Manage(obj)
		mgr.MarkTrue(abcPod.ObjectReference)

		gotObj := mgr.GetObjectConditionByObjRef(abcPod.ObjectReference)
		g.Expect(gotObj.Status).To(Equal(corev1.ConditionTrue))
		g.Expect(gotObj.Reason).To(Equal(""))
		g.Expect(gotObj.Message).To(Equal(""))
		g.Expect(gotObj.LastTransitionTime.Inner.IsZero()).ToNot(BeTrue())
	})

	t.Run("MarkTrueWithReason", func(t *testing.T) {
		g := NewGomegaWithT(t)
		objcs := ObjectConditions{}
		g.Expect(ktesting.LoadYAML("testdata/objectconditions.yaml", &objcs)).To(Succeed())

		obj := &SomeStatusObj{conditions: objcs}
		mgr := objcs.Manage(obj)
		mgr.MarkTrueWithReason(abcPod.ObjectReference, "MyReason", "some error %s abc", "error message here")
		gotObj := mgr.GetObjectConditionByObjRef(abcPod.ObjectReference)
		g.Expect(gotObj.Status).To(Equal(corev1.ConditionTrue))
		g.Expect(gotObj.Reason).To(Equal("MyReason"))
		g.Expect(gotObj.Message).To(Equal("some error error message here abc"))
		g.Expect(gotObj.LastTransitionTime.Inner.IsZero()).ToNot(BeTrue())
	})

	t.Run("MarkUnknown", func(t *testing.T) {
		g := NewGomegaWithT(t)
		objcs := ObjectConditions{}
		g.Expect(ktesting.LoadYAML("testdata/objectconditions.yaml", &objcs)).To(Succeed())

		obj := &SomeStatusObj{conditions: objcs}
		mgr := objcs.Manage(obj)
		mgr.MarkUnknown(abcPod.ObjectReference, "ABCMyReason", "abc error %s abc", "error message here")
		gotObj := mgr.GetObjectConditionByObjRef(abcPod.ObjectReference)
		g.Expect(gotObj.Status).To(Equal(corev1.ConditionUnknown))
		g.Expect(gotObj.Reason).To(Equal("ABCMyReason"))
		g.Expect(gotObj.Message).To(Equal("abc error error message here abc"))
		g.Expect(gotObj.LastTransitionTime.Inner.IsZero()).ToNot(BeTrue())
	})

	t.Run("MarkFalse", func(t *testing.T) {
		g := NewGomegaWithT(t)
		objcs := ObjectConditions{}
		g.Expect(ktesting.LoadYAML("testdata/objectconditions.yaml", &objcs)).To(Succeed())

		obj := &SomeStatusObj{conditions: objcs}
		mgr := objcs.Manage(obj)
		mgr.MarkFalse(abcPod.ObjectReference, "FalseReason", "My error message")
		mgr.SetConditionType(abcPod.ObjectReference, apis.ConditionReady)
		mgr.SetSeverity(abcPod.ObjectReference, apis.ConditionSeverityInfo)
		gotObj := mgr.GetObjectConditionByObjRef(abcPod.ObjectReference)
		g.Expect(gotObj.Status).To(Equal(corev1.ConditionFalse))
		g.Expect(gotObj.Reason).To(Equal("FalseReason"))
		g.Expect(gotObj.Message).To(Equal("My error message"))
		g.Expect(gotObj.LastTransitionTime.Inner.IsZero()).ToNot(BeTrue())
		g.Expect(gotObj.Type).To(Equal(apis.ConditionReady))
		g.Expect(gotObj.Severity).To(Equal(apis.ConditionSeverityInfo))
	})

	t.Run("MarkTrueWithReason for new condition", func(t *testing.T) {
		g := NewGomegaWithT(t)
		objcs := ObjectConditions{}
		g.Expect(ktesting.LoadYAML("testdata/objectconditions.yaml", &objcs)).To(Succeed())

		t.Logf("object conditions start: %d", len(objcs))

		obj := &SomeStatusObj{conditions: objcs}
		mgr := objcs.Manage(obj)
		ref := corev1.ObjectReference{Name: "new-object-ref", APIVersion: "v1", Kind: "ConfigMap"}
		mgr.MarkTrueWithReason(ref, "OK", "All good")

		t.Logf("object conditions after set new: %d \t %d", len(objcs), len(mgr.GetObjectConditions()))
		t.Logf("%#v", mgr.GetObjectConditions())

		gotObj := mgr.GetObjectConditionByObjRef(ref)
		g.Expect(gotObj).ToNot(BeNil())
		g.Expect(gotObj.Status).To(Equal(corev1.ConditionTrue))
		g.Expect(gotObj.Reason).To(Equal("OK"))
		g.Expect(gotObj.Message).To(Equal("All good"))
		g.Expect(gotObj.LastTransitionTime.Inner.IsZero()).ToNot(BeTrue())
	})

	t.Run("remove and add again", func(t *testing.T) {
		g := NewGomegaWithT(t)

		objcs := ObjectConditions{}
		g.Expect(ktesting.LoadYAML("testdata/objectconditions.yaml", &objcs)).To(Succeed())

		obj := &SomeStatusObj{conditions: objcs}
		mgr := objcs.Manage(obj)

		mgr.RemoveObjectConditionByObjRef(abcPod.ObjectReference)
		g.Expect(mgr.GetObjectConditions()).To(HaveLen(1))
		g.Expect(mgr.GetObjectConditionByObjRef(abcPod.ObjectReference)).To(BeNil())

		// check if it still have 2 items
		mgr.SetObjectConditions(objcs)
		g.Expect(mgr.GetObjectConditions()).To(HaveLen(2))
	})

}
