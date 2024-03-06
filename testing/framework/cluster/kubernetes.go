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

package cluster

import (
	"fmt"

	"github.com/katanomi/pkg/testing"
	. "github.com/onsi/gomega"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// ExpectKubeObject expect kubernetes resource exists by object.name and object.namespace, clean some fields in object and return gomega.Expect for this object.
// the param object should contains Namespace and name.
// you can use cleanFuncs to clean some fileds for object that load from kubernetes
// eg.
//
//	ExpectKubeObject(&TestContext{Client: clt, Namespace: "default"}, &corev1.ConfigMap{
//		ObjectMeta: metav1.ObjectMeta{
//			Namespace: "default",
//			Name:      "config-1",
//		},
//	}).Should(testing.DiffEqualTo(&corev1.ConfigMap{
//		ObjectMeta: metav1.ObjectMeta{
//			Namespace: "default",
//			Name:      "config-1",
//		},
//		Data: map[string]string{
//			"a": "a",
//			"b": "b",
//		},
//	}))
func ExpectKubeObject(ctx *TestContext, object ctrlclient.Object, cleanFuncs ...func(object interface{}) interface{}) Assertion {

	gvk := object.GetObjectKind().GroupVersionKind()
	err := ctx.Client.Get(ctx.Context, ctrlclient.ObjectKey{
		Namespace: object.GetNamespace(),
		Name:      object.GetName(),
	}, object)

	Expect(err).Should(BeNil())
	object.GetObjectKind().SetGroupVersionKind(gvk)
	fmt.Printf("===> %#v", object)
	if len(cleanFuncs) == 0 {
		cleanFuncs = []func(object interface{}) interface{}{
			testing.KubeObjectDiffClean,
		}
	}

	for _, clean := range cleanFuncs {
		object = clean(object).(ctrlclient.Object)
	}
	return Expect(object)
}
