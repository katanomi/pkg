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

	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"

	"k8s.io/client-go/dynamic/fake"
)

func TestContext(t *testing.T) {
	// g := NewGomegaWithT(t)
	// ctx := context.TODO()
	clusterClient := &ClusterRegistryClient{
		Interface: fake.NewSimpleDynamicClient(scheme.Scheme),
	}

	t.Run("get from context nil", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctx := context.TODO()

		g.Expect(MultiCluster(ctx)).To(BeNil())
	})

	t.Run("add and get", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctx := context.TODO()

		ctx = WithMultiCluster(ctx, clusterClient)

		g.Expect(MultiCluster(ctx)).To(Equal(clusterClient))

	})
}
