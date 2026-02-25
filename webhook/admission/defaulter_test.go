/*
Copyright 2023 The AlaudaDevops Authors.

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

package admission

import (
	"context"
	"testing"

	kclient "github.com/AlaudaDevops/pkg/client"
	kconfig "github.com/AlaudaDevops/pkg/config"
	kscheme "github.com/AlaudaDevops/pkg/scheme"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/onsi/gomega"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	ctrlclientfake "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type testDefaulterObj struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

func (m *testDefaulterObj) Default(ctx context.Context) {
}

func (m *testDefaulterObj) DeepCopyObject() runtime.Object {
	return &testDefaulterObj{
		TypeMeta:   m.TypeMeta,
		ObjectMeta: m.ObjectMeta,
	}
}

type testContextInjectorObject struct {
	testDefaulterObj
}

func (m *testContextInjectorObject) InjectContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "foo", "bar")
}

func (m *testContextInjectorObject) Default(ctx context.Context) {
	value := ctx.Value("foo")
	if value.(string) != "bar" {
		panic("context not injected")
	}
}

func (m *testContextInjectorObject) DeepCopyObject() runtime.Object {
	return &testContextInjectorObject{
		testDefaulterObj: m.testDefaulterObj,
	}
}

func TestDefaulterContextInjector(t *testing.T) {
	ctx := context.Background()
	ctx = kscheme.WithScheme(ctx, scheme.Scheme)

	g := gomega.NewGomegaWithT(t)
	req := admission.Request{
		AdmissionRequest: admissionv1.AdmissionRequest{
			Operation: admissionv1.Delete,
			OldObject: runtime.RawExtension{
				Raw:    []byte(`{"metadata":{"name":"def","namespace":"default"}}`),
				Object: &testDefaulterObj{},
			},
		},
	}
	table := map[string]struct {
		Validator *admission.Webhook
	}{
		"object not implementing context injector interface": {
			Validator: DefaultingWebhookFor(ctx, &testDefaulterObj{}),
		},
		"object implement the context injector interface": {
			Validator: DefaultingWebhookFor(ctx, &testContextInjectorObject{}),
		},
	}

	for name, test := range table {
		t.Run(name, func(t *testing.T) {
			g.Expect(func() {
				test.Validator.Handle(ctx, req)
			}).NotTo(gomega.Panic())
		})
	}
}

// testConfigManagerDefaulterObj is a Defaulter that captures the ConfigManager and Client from context
// during Default() to allow verification in tests.
// The state pointer is shared across DeepCopy instances so captures are visible to the original.
type testConfigManagerDefaulterObj struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// state holds captured values and is shared across DeepCopy instances.
	state *testConfigManagerDefaulterState
}

type testConfigManagerDefaulterState struct {
	CapturedConfigManager *kconfig.Manager
	CapturedClient        interface{}
}

func (m *testConfigManagerDefaulterObj) Default(ctx context.Context) {
	if m.state != nil {
		m.state.CapturedConfigManager = kconfig.ConfigManager(ctx)
		m.state.CapturedClient = kclient.Client(ctx)
	}
}

func (m *testConfigManagerDefaulterObj) DeepCopyObject() runtime.Object {
	return &testConfigManagerDefaulterObj{
		TypeMeta:   m.TypeMeta,
		ObjectMeta: m.ObjectMeta,
		state:      m.state, // share state pointer so Default() captures are observable
	}
}

// TestDefaulterInjectsConfigManagerAndClientFromHandlerCtx verifies that
// the mutatingHandler.Handle method propagates the ConfigManager and Client
// stored in h.ctx (set at webhook construction time) into the per-request ctx
// so that Defaulter implementations can access them.
func TestDefaulterInjectsConfigManagerAndClientFromHandlerCtx(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	// Build a base context with scheme (required for decoder).
	baseCtx := context.Background()
	baseCtx = kscheme.WithScheme(baseCtx, scheme.Scheme)

	// Prepare a non-nil ConfigManager and a fake client to inject via h.ctx.
	configManager := &kconfig.Manager{}
	fakeClient := ctrlclientfake.NewClientBuilder().Build()

	// handlerCtx simulates the context passed to DefaultingWebhookFor (stored as h.ctx).
	handlerCtx := kconfig.WithConfigManager(baseCtx, configManager)
	handlerCtx = kclient.WithClient(handlerCtx, fakeClient)

	table := map[string]struct {
		handlerCtx        context.Context
		wantConfigManager *kconfig.Manager
		wantClientNil     bool
	}{
		"ConfigManager and Client are injected when set in handler ctx": {
			handlerCtx:        handlerCtx,
			wantConfigManager: configManager,
			wantClientNil:     false,
		},
		"ConfigManager is nil when not set in handler ctx": {
			handlerCtx:        baseCtx,
			wantConfigManager: nil,
			wantClientNil:     true,
		},
	}

	req := admission.Request{
		AdmissionRequest: admissionv1.AdmissionRequest{
			Operation: admissionv1.Create,
			Object: runtime.RawExtension{
				Raw: []byte(`{"metadata":{"name":"test","namespace":"default"}}`),
			},
		},
	}

	for name, tc := range table {
		t.Run(name, func(t *testing.T) {
			state := &testConfigManagerDefaulterState{}
			obj := &testConfigManagerDefaulterObj{state: state}
			wh := DefaultingWebhookFor(tc.handlerCtx, obj)
			response := wh.Handle(baseCtx, req)
			g.Expect(response.Allowed).To(gomega.BeTrue())
			g.Expect(state.CapturedConfigManager).To(gomega.Equal(tc.wantConfigManager))
			if tc.wantClientNil {
				g.Expect(state.CapturedClient).To(gomega.BeNil())
			} else {
				g.Expect(state.CapturedClient).NotTo(gomega.BeNil())
			}
		})
	}
}
