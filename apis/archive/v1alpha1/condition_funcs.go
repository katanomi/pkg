/*
Copyright 2022 The Katanomi Authors.

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
	"reflect"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Or generate condition with or operator
func Or(conditions ...Condition) Condition {
	return Condition{Key: "", Operator: ConditionOperatorOr, Value: conditions}
}

// And generate condition with and operator
func And(conditions ...Condition) Condition {
	return Condition{Key: "", Operator: ConditionOperatorAnd, Value: conditions}
}

// UID generate condition for uid field
func UID(v string) Condition {
	return Equal(UIDField, v)
}

// ParentUID generate condition for parent uid field
func ParentUID(v string) Condition {
	return Equal(ParentUIDField, v)
}

// TopUID generate condition for top uid field
func TopUID(v string) Condition {
	return Equal(TopUIDField, v)
}

// Cluster generate condition for cluster field
func Cluster(v string) Condition {
	return Equal(ClusterField, v)
}

// ParentCluster generate condition for parent cluster field
func ParentCluster(v string) Condition {
	return Equal(ParentClusterField, v)
}

// TopCluster generate condition for top cluster field
func TopCluster(v string) Condition {
	return Equal(TopClusterField, v)
}

// Namespace generate condition for namespace field
func Namespace(v string) Condition {
	return Equal(NamespaceField, v)
}

// ParentNamespace generate condition for parent namespace field
func ParentNamespace(v string) Condition {
	return Equal(ParentNamespaceField, v)
}

// TopNamespace generate condition for top namespace field
func TopNamespace(v string) Condition {
	return Equal(TopNamespaceField, v)
}

// Name generate condition for name field
func Name(v string) Condition {
	return Equal(NameField, v)
}

// ParentName generate condition for parent name field
func ParentName(v string) Condition {
	return Equal(ParentNameField, v)
}

// TopName generate condition for top name field
func TopName(v string) Condition {
	return Equal(TopNameField, v)
}

// Equal generate condition with equal operator
func Equal(key, value string) Condition {
	return Condition{
		Key:      key,
		Operator: ConditionOperatorEqual,
		Value:    value,
	}
}

// EqualColumn generate condition with equal column operator
func EqualColumn(key1, key2 string) Condition {
	return Condition{
		Key:      key1,
		Operator: ConditionOperatorEqualColumn,
		Value:    key2,
	}
}

// Gt generate condition with gt operator
func Gt(key, value string) Condition {
	return Condition{
		Key:      key,
		Operator: ConditionOperatorGt,
		Value:    value,
	}
}

// Gte generate condition with gte operator
func Gte(key, value string) Condition {
	return Condition{
		Key:      key,
		Operator: ConditionOperatorGte,
		Value:    value,
	}
}

// Lt generate condition with lt operator
func Lt(key, value string) Condition {
	return Condition{
		Key:      key,
		Operator: ConditionOperatorLt,
		Value:    value,
	}
}

// Lte generate condition with lte operator
func Lte(key, value string) Condition {
	return Condition{
		Key:      key,
		Operator: ConditionOperatorLte,
		Value:    value,
	}
}

// In generate condition with in operator
func In(key string, values ...string) Condition {
	return Condition{
		Key:      key,
		Operator: ConditionOperatorIn,
		Value:    values,
	}
}

// Like generate condition with like operator
func Like(key, value string) Condition {
	return Condition{
		Key:      key,
		Operator: ConditionOperatorLike,
		Value:    value,
	}
}

// NotEqual generate condition with not equal operator
func NotEqual(key, value string) Condition {
	return Condition{Key: key, Operator: ConditionOperatorNotEqual, Value: value}
}

// GVK generate condition for gvk field
func GVK(gvk schema.GroupVersionKind) []Condition {
	return []Condition{
		Equal(GroupField, gvk.Group),
		Equal(VersionField, gvk.Version),
		Equal(KindField, gvk.Kind),
	}
}

// NamespacedName generate condition for namespace and name field
func NamespacedName(ns, name string) []Condition {
	return []Condition{
		Equal(NamespaceField, ns),
		Equal(NameField, name),
	}
}

// CompletedStatus generate condition for completed status
func CompletedStatus() Condition {
	return Or(
		Equal(MetadataKey("status"), string(corev1.ConditionFalse)),
		Equal(MetadataKey("status"), string(corev1.ConditionTrue)),
	)
}

// ToInterfaceSlice convert slice to interface slice
func ToInterfaceSlice(value interface{}) []interface{} {
	reflectValue := reflect.ValueOf(value)
	valueLen := reflectValue.Len()
	values := make([]interface{}, valueLen)
	for i := 0; i < valueLen; i++ {
		values[i] = reflectValue.Index(i).Interface()
	}
	return values
}

// Exist generate condition with Exist operator
func Exist(key string) Condition {
	return Condition{
		Key:      key,
		Operator: ConditionOperatorExist,
		Value:    "",
	}
}
