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

package multicluster

import (
	"context"
	"testing"

	ktesting "github.com/katanomi/pkg/testing"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes/scheme"
)

func TestClusterRegistryClientGetConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	objs, err := ktesting.LoadKubeResourcesAsUnstructured("testdata/clusterregistry_client.GetConfig.cluster.yaml")
	g.Expect(err).To(BeNil())

	ctx := context.TODO()

	runtimeObjs := make([]runtime.Object, len(objs))
	for i, o := range objs {
		obj := o
		runtimeObjs[i] = &obj
	}
	clusterClient := &ClusterRegistryClient{
		Interface: fake.NewSimpleDynamicClient(scheme.Scheme, runtimeObjs...),
	}

	ref := &corev1.ObjectReference{
		Namespace: "default",
		Name:      "my-cluster",
	}

	t.Run("get config", func(t *testing.T) {
		g := NewGomegaWithT(t)
		config, err := clusterClient.GetConfig(ctx, ref)

		g.Expect(err).To(BeNil())
		g.Expect(config).ToNot(BeNil())
		g.Expect(config.Host).To(Equal("https://127.0.0.1:1111"))
		g.Expect(config.BearerToken).To(Equal("abctoken"))
	})

	t.Run("get client", func(t *testing.T) {
		client, err := clusterClient.GetClient(ctx, ref, scheme.Scheme)

		// The create client method will do a request
		// and in this case it should fail
		// so there is nothing we can do actually
		g.Expect(err).ToNot(BeNil())
		g.Expect(client).To(BeNil())
	})

	t.Run("get dynamic", func(t *testing.T) {
		client, err := clusterClient.GetDynamic(ctx, ref)

		g.Expect(err).To(BeNil())
		g.Expect(client).ToNot(BeNil())
	})

	t.Run("list clusters namespaces", func(t *testing.T) {
		clusterNamespaces, err := clusterClient.ListClustersNamespaces(ctx, "namespace")

		g.Expect(err).To(BeNil())
		g.Expect(clusterNamespaces).To(HaveLen(0))
	})

	t.Run("is namespace in this project", func(t *testing.T) {
		exist, err := clusterClient.IsNamespaceInProject(ctx, "projectName", nil, "namespace")

		g.Expect(err).To(BeNil())
		g.Expect(exist).To(BeTrue())
	})

}
