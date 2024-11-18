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

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/AlaudaDevops/pkg/testing"
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

// convertToSetNamespace will set object namespace according TestContext
func convertToSetNamespace(namespace string) testing.ConvertRuntimeObjctToClientObjectFunc {
	return func(object runtime.Object) (ctrlclient.Object, error) {
		obj, ok := object.(ctrlclient.Object)
		if !ok {
			err := fmt.Errorf("Unsupported gvk: %s", object.GetObjectKind().GroupVersionKind())
			return nil, err
		}
		if obj.GetNamespace() == "" {
			obj.SetNamespace(namespace)
		}
		return obj, nil
	}
}

// LoadKubeResourcesToCtx will load kubernetes resources from file
// it will use client in ctx and set namespace as ctx.Namespace if namespace is empty in file
// if you set namespace in file, it will just use it and ignore ctx.Namespace
// you can set converts to change some extra fields before create in kubernetes
// but you must pay attention to that when you set more than one converts: just one converts will be executed when match
// if you set converts , you should set namespace to ctx.Namespace explicitly
// eg.
//
//	LoadKubeResources(ctx, file)// load resource from file and apply to kubernetes with ctx.Namespace and ctx.Client
//
//	// load resource from file and execute convert if resource type is configmap before apply to kubernetes with ctx.Namespace and ctx.Client
//
//	LoadKubeResources(ctx, file,func(object runtime.Object) (ctrlclient.Object, error) {
//		configmap, ok := object.(*corev1.ConfigMap)
//		if !ok {
//			return nil, fmt.Errorf("Unsupported gvk: %s", object.GetObjectKind().GroupVersionKind())
//		}
//		configmap.Namespace = "default-e2e"
//		configmap.Annotations = map[string]string{
//			"a": "1",
//		}
//		return configmap, nil
//	})
func LoadKubeResourcesToCtx(ctx *TestContext, file string, converts ...testing.ConvertRuntimeObjctToClientObjectFunc) (err error) {
	converts = append(converts, convertToSetNamespace(ctx.Namespace))
	return testing.LoadKubeResources(file, ctx.Client, converts...)
}

// MustLoadKubeResourcesToCtx similar with LoadKubeResources but panic if LoadKubeResources error
func MustLoadKubeResourcesToCtx(ctx *TestContext, file string, converts ...testing.ConvertRuntimeObjctToClientObjectFunc) {
	err := LoadKubeResourcesToCtx(ctx, file, converts...)
	if err != nil {
		panic(fmt.Sprintf("load yaml file failed, file path: %s, err: %s", file, err))
	}
}
