/*
Copyright 2024 The Katanomi Authors.

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

package testing

import (
	"time"

	"github.com/AlaudaDevops/pkg/testing/assertions"

	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/onsi/gomega"
)

// ExpectDiff will using github.com/google/go-cmp/cmp.Diff to compare two data
// as default, it will not compare fields that have some undetermined value like uid, timestamp when compare kubernetes object
// you can use diffCleanFuncs to change object in order to clean some fields that you don't want to compare
// eg: clean status when compare
//
//	ExpectDiff(actual, expected, KubeObjectDiffClean, func(object ctrlclient.Object) {
//		object.(*corev1.Pod).Status = corev1.PodStatus{} // do not compare Status
//		return object
//	}).Should(BeEmpty())
//
// Deprecated: should use regular Expect with assertions.DiffEqual instead
// Example:
//
//	Expect(obj).To(DiffEqual(expected, cmpopts.IgnoreTypes(v1.TypeMeta{}, v1.PodStatus{})))
func ExpectDiff(actual interface{}, expected interface{}, diffCleanFuncs ...func(object interface{}) interface{}) gomega.Assertion {
	cleanFuncs := []func(object interface{}) interface{}{}
	if len(diffCleanFuncs) == 0 {
		cleanFuncs = append(cleanFuncs, KubeObjectDiffClean)
	}
	cleanFuncs = append(cleanFuncs, diffCleanFuncs...)

	for _, cleanF := range cleanFuncs {
		expected = cleanF(expected)
		actual = cleanF(actual)
	}

	return gomega.Expect(cmp.Diff(expected, actual))
}

// KubeObjectDiffClean will clean these fields in order to clean when diff kubernetes objects
//   - CreationTimestamp
//   - ManagedFields
//   - UID
//   - ResourceVersion
//   - Generation
//   - SelfLink
//
// Deprecated: use assertions.DiffEqual and assertions.IgnoreObjectMetaFields instead
func KubeObjectDiffClean(object interface{}) interface{} {

	k8sObject, ok := object.(ctrlclient.Object)
	if !ok {
		return object
	}
	k8sObject = k8sObject.DeepCopyObject().(ctrlclient.Object)

	k8sObject.SetCreationTimestamp(metav1.NewTime(time.Time{}))
	k8sObject.SetManagedFields(nil)
	k8sObject.SetUID("")
	k8sObject.SetResourceVersion("")
	k8sObject.SetGeneration(0)
	k8sObject.SetSelfLink("")

	return k8sObject
}

// DiffEqualTo will use github.com/google/go-cmp/cmp.Diff to compare
// you can use diffCleanFuncs to change object in order to clean some fields that you don't want to compare
// Deprecated: use assertions.DiffEqual instead
func DiffEqualTo(expected interface{}, diffCleanFuncs ...func(object interface{}) interface{}) gomega.OmegaMatcher {
	if len(diffCleanFuncs) == 0 {
		diffCleanFuncs = append(diffCleanFuncs, KubeObjectDiffClean)
	}
	return &assertions.DiffEqualMatcher{Expected: expected, DiffCleanFunc: diffCleanFuncs}
}
