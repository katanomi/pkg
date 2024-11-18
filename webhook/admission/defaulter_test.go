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

	kscheme "github.com/AlaudaDevops/pkg/scheme"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/onsi/gomega"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
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
