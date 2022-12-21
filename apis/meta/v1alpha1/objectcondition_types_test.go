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

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	ktesting "github.com/katanomi/pkg/testing"
	mockmetav1alpha1 "github.com/katanomi/pkg/testing/mock/github.com/katanomi/pkg/apis/meta/v1alpha1"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"knative.dev/pkg/apis"
)

//go:generate ../../../bin/mockgen -package=apis -destination=../../../testing/mock/github.com/katanomi/pkg/apis/meta/v1alpha1/top_level_condition_object.go  github.com/katanomi/pkg/apis/meta/v1alpha1 TopLevelConditionObject

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

func TestFromTopLevelConditionObject(t *testing.T) {
	g := NewGomegaWithT(t)

	abcPod := &ObjectCondition{}
	g.Expect(ktesting.LoadYAML("testdata/objectcondition-pod-abc.yaml", abcPod)).To(Succeed())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	condition := &apis.Condition{
		Type:    apis.ConditionReady,
		Status:  corev1.ConditionTrue,
		Reason:  "Reason",
		Message: "message",
	}

	topLevelCondition := mockmetav1alpha1.NewMockTopLevelConditionObject(ctrl)
	topLevelCondition.EXPECT().GetTopLevelCondition().Return(condition).Times(2)
	topLevelCondition.EXPECT().GetAnnotations().Return(map[string]string{"abc": "def"}).Times(1)
	topLevelCondition.EXPECT().GetName().Return("some-name")
	topLevelCondition.EXPECT().GetNamespace().Return("default")
	uid := types.UID("1234")
	topLevelCondition.EXPECT().GetUID().Return(uid)

	abcPod.FromTopLevelConditionObject(topLevelCondition)

	g.Expect(abcPod.Condition).To(Equal(*condition))
	g.Expect(abcPod.Name).To(Equal("some-name"))
	g.Expect(abcPod.Namespace).To(Equal("default"))
	g.Expect(abcPod.UID).To(Equal(uid))
	g.Expect(abcPod.Annotations).To(Equal(map[string]string{"abc": "def"}))
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

func TestObjectConditionFromTopLevelConditionObject(t *testing.T) {
	g := NewGomegaWithT(t)

	objcs := &ObjectCondition{}
	g.Expect(objcs.FromTopLevelConditionObject(nil)).To(BeNil())
}
