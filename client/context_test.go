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

package client

import (
	"context"
	"testing"

	"k8s.io/client-go/rest"

	dynamicFake "k8s.io/client-go/dynamic/fake"

	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestClientContext(t *testing.T) {
	g := NewGomegaWithT(t)

	ctx := context.TODO()

	clt := Client(ctx)
	g.Expect(clt).To(BeNil())

	fakeClt := fake.NewClientBuilder().Build()
	ctx = WithClient(ctx, fakeClt)
	g.Expect(Client(ctx)).To(Equal(fakeClt))
	g.Expect(DirectClient(ctx)).To(BeNil())
}

func TestDirectClientContext(t *testing.T) {
	g := NewGomegaWithT(t)

	ctx := context.TODO()

	clt := DirectClient(ctx)
	g.Expect(clt).To(BeNil())

	fakeClt := fake.NewClientBuilder().Build()
	ctx = WithDirectClient(ctx, fakeClt)
	g.Expect(DirectClient(ctx)).To(Equal(fakeClt))
	g.Expect(Client(ctx)).To(BeNil())
}

func TestManagerContext(t *testing.T) {
	g := NewGomegaWithT(t)

	ctx := context.TODO()

	clt := ManagerCtx(ctx)
	g.Expect(clt).To(BeNil())

	mgr := NewManager(ctx, nil, nil)
	ctx = WithManager(ctx, mgr)
	g.Expect(ManagerCtx(ctx)).To(Equal(mgr))
}

func TestDynamicContext(t *testing.T) {
	g := NewGomegaWithT(t)

	ctx := context.TODO()

	clt, _ := DynamicClient(ctx)
	g.Expect(clt).To(BeNil())

	client := &dynamicFake.FakeDynamicClient{}
	ctx = WithDynamicClient(ctx, client)
	g.Expect(DynamicClient(ctx)).To(Equal(client))
}

func TestAppConfigContext(t *testing.T) {
	g := NewGomegaWithT(t)

	ctx := context.TODO()

	cfg := GetAppConfig(ctx)
	g.Expect(cfg).To(BeNil())

	cfg = &rest.Config{}
	ctx = WithAppConfig(ctx, cfg)
	g.Expect(GetAppConfig(ctx)).To(Equal(cfg))
}
