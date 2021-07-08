package v1alpha1

import (
	"testing"

	ktesting "github.com/katanomi/pkg/testing"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
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

	// ctx := context.TODO()

	objcs := ObjectConditions{}
	g.Expect(ktesting.LoadYAML("testdata/objectconditions.yaml", &objcs)).To(Succeed())

	obj := &SomeStatusObj{conditions: objcs}
	mgr := objcs.Manage(obj)

	abcPod := ObjectCondition{}
	g.Expect(ktesting.LoadYAML("testdata/objectcondition-pod-abc.yaml", &abcPod)).To(Succeed())

	// g.Expect(gotObj).ToNot(BeNil())
	// diff := cmp.Diff(*gotObj, abcPod)
	// g.Expect(diff).To(Equal(""))

	t.Run("MarkTrue", func(t *testing.T) {
		g := NewGomegaWithT(t)
		mgr.MarkTrue(abcPod.ObjectReference)

		gotObj := mgr.GetObjectConditionByObjRef(abcPod.ObjectReference)
		g.Expect(gotObj.Status).To(Equal(corev1.ConditionTrue))
	})

	t.Run("MarkTrueWithReason", func(t *testing.T) {
		g := NewGomegaWithT(t)
		mgr.MarkTrueWithReason(abcPod.ObjectReference, "MyReason", "some error %s abc", "error message here")
		gotObj := mgr.GetObjectConditionByObjRef(abcPod.ObjectReference)
		g.Expect(gotObj.Status).To(Equal(corev1.ConditionTrue))
		g.Expect(gotObj.Reason).To(Equal("MyReason"))
		g.Expect(gotObj.Message).To(Equal("some error error message here abc"))
	})

	t.Run("MarkUnknown", func(t *testing.T) {
		g := NewGomegaWithT(t)
		mgr.MarkUnknown(abcPod.ObjectReference, "ABCMyReason", "abc error %s abc", "error message here")
		gotObj := mgr.GetObjectConditionByObjRef(abcPod.ObjectReference)
		g.Expect(gotObj.Status).To(Equal(corev1.ConditionUnknown))
		g.Expect(gotObj.Reason).To(Equal("ABCMyReason"))
		g.Expect(gotObj.Message).To(Equal("abc error error message here abc"))
	})

	t.Run("MarkFalse", func(t *testing.T) {
		g := NewGomegaWithT(t)
		mgr.MarkFalse(abcPod.ObjectReference, "FalseReason", "My errror message")
		gotObj := mgr.GetObjectConditionByObjRef(abcPod.ObjectReference)
		g.Expect(gotObj.Status).To(Equal(corev1.ConditionFalse))
		g.Expect(gotObj.Reason).To(Equal("FalseReason"))
		g.Expect(gotObj.Message).To(Equal("My errror message"))
	})

	mgr.RemoveObjectConditionByObjRef(abcPod.ObjectReference)
	g.Expect(mgr.GetObjectConditions()).To(HaveLen(1))
	g.Expect(mgr.GetObjectConditionByObjRef(abcPod.ObjectReference)).To(BeNil())

	// check if it still have 2 items
	mgr.SetObjectConditions(objcs)
	g.Expect(mgr.GetObjectConditions()).To(HaveLen(2))

}
