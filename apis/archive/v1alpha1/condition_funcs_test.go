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

package v1alpha1

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	ktesting "github.com/katanomi/pkg/testing"
	"github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestFieldCondition(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	g.Expect(Gt("key", "value")).To(gomega.Equal(Condition{
		Key:      "key",
		Operator: ConditionOperatorGt,
		Value:    "value",
	}))
	g.Expect(Gte("key", "value")).To(gomega.Equal(Condition{
		Key:      "key",
		Operator: ConditionOperatorGte,
		Value:    "value",
	}))
	g.Expect(Lt("key", "value")).To(gomega.Equal(Condition{
		Key:      "key",
		Operator: ConditionOperatorLt,
		Value:    "value",
	}))
	g.Expect(Lte("key", "value")).To(gomega.Equal(Condition{
		Key:      "key",
		Operator: ConditionOperatorLte,
		Value:    "value",
	}))
	g.Expect(Like("key", "value")).To(gomega.Equal(Condition{
		Key:      "key",
		Operator: ConditionOperatorLike,
		Value:    "value",
	}))
	g.Expect(NotEqual("key", "value")).To(gomega.Equal(Condition{
		Key:      "key",
		Operator: ConditionOperatorNotEqual,
		Value:    "value",
	}))
}

func TestFieldCondition_equal(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	g.Expect(UID("test")).To(gomega.Equal(Condition{
		Key:      UIDField,
		Operator: ConditionOperatorEqual,
		Value:    "test",
	}))
	g.Expect(ParentUID("test")).To(gomega.Equal(Condition{
		Key:      ParentUIDField,
		Operator: ConditionOperatorEqual,
		Value:    "test",
	}))
	g.Expect(TopUID("test")).To(gomega.Equal(Condition{
		Key:      TopUIDField,
		Operator: ConditionOperatorEqual,
		Value:    "test",
	}))
	g.Expect(Cluster("test")).To(gomega.Equal(Condition{
		Key:      ClusterField,
		Operator: ConditionOperatorEqual,
		Value:    "test",
	}))
	g.Expect(ParentCluster("test")).To(gomega.Equal(Condition{
		Key:      ParentClusterField,
		Operator: ConditionOperatorEqual,
		Value:    "test",
	}))
	g.Expect(TopCluster("test")).To(gomega.Equal(Condition{
		Key:      TopClusterField,
		Operator: ConditionOperatorEqual,
		Value:    "test",
	}))
	g.Expect(Name("test")).To(gomega.Equal(Condition{
		Key:      NameField,
		Operator: ConditionOperatorEqual,
		Value:    "test",
	}))
	g.Expect(ParentName("test")).To(gomega.Equal(Condition{
		Key:      ParentNameField,
		Operator: ConditionOperatorEqual,
		Value:    "test",
	}))
	g.Expect(TopName("test")).To(gomega.Equal(Condition{
		Key:      TopNameField,
		Operator: ConditionOperatorEqual,
		Value:    "test",
	}))
	g.Expect(Namespace("test")).To(gomega.Equal(Condition{
		Key:      NamespaceField,
		Operator: ConditionOperatorEqual,
		Value:    "test",
	}))
	g.Expect(ParentNamespace("test")).To(gomega.Equal(Condition{
		Key:      ParentNamespaceField,
		Operator: ConditionOperatorEqual,
		Value:    "test",
	}))
	g.Expect(TopNamespace("test")).To(gomega.Equal(Condition{
		Key:      TopNamespaceField,
		Operator: ConditionOperatorEqual,
		Value:    "test",
	}))
	g.Expect(NamespacedName("ns", "name")).To(gomega.Equal([]Condition{
		Namespace("ns"),
		Name("name"),
	}))
	g.Expect(GVK(schema.GroupVersionKind{
		Group:   "group",
		Version: "version",
		Kind:    "kind",
	})).To(gomega.Equal([]Condition{
		Equal(GroupField, "group"),
		Equal(VersionField, "version"),
		Equal(KindField, "kind"),
	}))
}

func TestOr(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	expect := Condition{}
	ktesting.MustLoadYaml("testdata/condition.or.golden.yaml", &expect)
	cond := Or(
		UID("test-uid"),
		Name("test-name"),
	)
	g.Expect(cmp.Diff(cond, expect)).To(gomega.BeEmpty())
}

func TestAnd(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	expect := Condition{}
	ktesting.MustLoadYaml("testdata/condition.and.golden.yaml", &expect)
	cond := And(
		UID("test-uid"),
		Name("test-name"),
	)
	g.Expect(cmp.Diff(cond, expect)).To(gomega.BeEmpty())
}

func TestIn(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	expect := Condition{}
	ktesting.MustLoadYaml("testdata/condition.in.string.golden.yaml", &expect)
	cond := In("uid", "uid1", "uid2")
	cond.Value = ToInterfaceSlice(cond.Value)
	g.Expect(cmp.Diff(cond, expect)).To(gomega.BeEmpty())
}

func TestCompleteStatus(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	expect := Condition{}
	ktesting.MustLoadYaml("testdata/condition.completedStatus.golden.yaml", &expect)
	cond := CompletedStatus()
	g.Expect(cmp.Diff(cond, expect)).To(gomega.BeEmpty())
}

func TestEqualColumn(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	expect := Condition{}
	ktesting.MustLoadYaml("testdata/condition.eq.column.golden.yaml", &expect)
	cond := EqualColumn("uid", "uid2")
	g.Expect(cmp.Diff(cond, expect)).To(gomega.BeEmpty())
}
