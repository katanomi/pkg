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

	"github.com/google/go-cmp/cmp"
	ktesting "github.com/katanomi/pkg/testing"
	. "github.com/onsi/gomega"
)

func TestObjectConditions(t *testing.T) {
	g := NewGomegaWithT(t)

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

func TestDeepCopy(t *testing.T) {
	g := NewGomegaWithT(t)

	objcs := ObjectConditions{}
	g.Expect(ktesting.LoadYAML("testdata/objectconditions.yaml", &objcs)).To(Succeed())

	g.Expect(cmp.Diff(objcs, objcs.DeepCopy())).To(BeEmpty())

	abcPod := ObjectCondition{}
	g.Expect(ktesting.LoadYAML("testdata/objectcondition-pod-abc.yaml", &abcPod)).To(Succeed())

	g.Expect(cmp.Diff(&abcPod, abcPod.DeepCopy())).To(BeEmpty())
}
