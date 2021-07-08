package v1alpha1

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	ktesting "github.com/katanomi/pkg/testing"
	. "github.com/onsi/gomega"
)

func TestObjectConditions(t *testing.T) {
	g := NewGomegaWithT(t)

	// ctx := context.TODO()

	objcs := ObjectConditions{}
	g.Expect(ktesting.LoadYAML("testdata/objectconditions.yaml", &objcs)).To(Succeed())

	abcPod := ObjectCondition{}
	g.Expect(ktesting.LoadYAML("testdata/objectcondition-pod-abc.yaml", &abcPod)).To(Succeed())

	g.Expect(objcs).To(HaveLen(2))

	// Sets the same object again
	newobjcs := objcs.SetObjectCondition(abcPod)
	diff := cmp.Diff(newobjcs, objcs)
	g.Expect(diff).To(Equal(""))

	// removes the object
	newobjcs = objcs.RemoveObjectConditionByObjRef(abcPod.ObjectReference)
	g.Expect(newobjcs).To(HaveLen(1))
	g.Expect(objcs).To(HaveLen(2))

	// adds the object again
	abcobjcs := newobjcs.SetObjectCondition(abcPod)
	g.Expect(abcobjcs).To(HaveLen(2))
	diff = cmp.Diff(abcobjcs[1], abcPod)
	g.Expect(diff).To(BeEmpty())
}
