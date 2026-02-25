/*
Copyright 2021 The AlaudaDevops Authors.

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
	"fmt"
	"net/http"
	"testing"

	kclient "github.com/AlaudaDevops/pkg/client"
	kconfig "github.com/AlaudaDevops/pkg/config"
	kscheme "github.com/AlaudaDevops/pkg/scheme"

	"github.com/onsi/gomega"

	"github.com/google/go-cmp/cmp"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type MyObject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

func (m *MyObject) DeepCopyObject() runtime.Object {
	return m.DeepCopy()
}

func (m *MyObject) DeepCopy() *MyObject {
	return &MyObject{
		TypeMeta:   m.TypeMeta,
		ObjectMeta: m.ObjectMeta,
	}
}

func (m *MyObject) ValidateCreate(ctx context.Context) error {
	if m.Name == "" {
		return fmt.Errorf("some regular error")
	}
	if m.Namespace == "" {
		return errors.NewBadRequest("needs to have a namespace")
	}
	return nil
}

func (m *MyObject) ValidateDelete(ctx context.Context) error {
	return m.ValidateCreate(ctx)
}

func (m *MyObject) ValidateUpdate(ctx context.Context, old runtime.Object) error {
	return m.ValidateCreate(ctx)
}

func validateCreateFunc(err error) ValidateCreateFunc {
	return func(_ context.Context, _ runtime.Object, _ admission.Request) error {
		// no-op validation func, returns error
		return err
	}
}

func validateUpdateFunc(err error) ValidateUpdateFunc {
	return func(_ context.Context, _, _ runtime.Object, _ admission.Request) error {
		// no-op validation func, returns error
		return err
	}
}

func validateDeleteFunc(err error) ValidateDeleteFunc {
	return func(_ context.Context, _ runtime.Object, _ admission.Request) error {
		// no-op validation func, returns error
		return err
	}
}

func TestValidatorCreate(t *testing.T) {
	ctx := context.Background()
	ctx = kscheme.WithScheme(ctx, scheme.Scheme)

	table := map[string]struct {
		Validator *admission.Webhook
		Context   context.Context
		Request   admission.Request
		Response  admission.Response
	}{
		"simple ok create validation": {
			Validator: ValidatingWebhookFor(ctx, &MyObject{}, []ValidateCreateFunc{validateCreateFunc(nil)}, nil, nil),
			Context:   context.TODO(),
			Request: admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					Operation: admissionv1.Create,
					Object: runtime.RawExtension{
						Raw:    []byte(`{"metadata":{"name":"abc","namespace":"default"}}`),
						Object: &MyObject{},
					},
				},
			},
			Response: admission.Allowed(""),
		},
		"error create validation": {
			Validator: ValidatingWebhookFor(ctx, &MyObject{}, nil, nil, nil),
			Context:   context.TODO(),
			Request: admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					Operation: admissionv1.Create,
					Object: runtime.RawExtension{
						Raw:    []byte(`{"metadata":{"name":"abc","namespace":""}}`),
						Object: &MyObject{},
					},
				},
			},
			Response: validationResponseFromStatus(metav1.Status{
				Status:  metav1.StatusFailure,
				Code:    http.StatusBadRequest,
				Reason:  metav1.StatusReasonBadRequest,
				Message: "needs to have a namespace",
			}),
		},
		"returns error from extra added validation": {
			Validator: ValidatingWebhookFor(ctx, &MyObject{}, []ValidateCreateFunc{validateCreateFunc(fmt.Errorf("this is an extra error"))}, nil, nil),
			Context:   context.TODO(),
			Request: admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					Operation: admissionv1.Create,
					Object: runtime.RawExtension{
						Raw:    []byte(`{"metadata":{"name":"abc","namespace":"default"}}`),
						Object: &MyObject{},
					},
				},
			},
			Response: admission.Denied(fmt.Errorf("this is an extra error").Error()),
		},
	}

	for name, test := range table {
		t.Run(name, func(t *testing.T) {

			returned := test.Validator.Handle(test.Context, test.Request)
			diff := cmp.Diff(returned, test.Response)
			t.Logf("diff is: \n%s\n %#v == %#v", diff, test.Response, returned)
			if diff != "" {
				t.Fail()
			}
			if returned.Allowed != test.Response.Allowed {
				t.Fail()
			}
		})
	}
}

func TestValidatorUpdate(t *testing.T) {
	ctx := context.Background()
	ctx = kscheme.WithScheme(ctx, scheme.Scheme)

	table := map[string]struct {
		Validator *admission.Webhook
		Context   context.Context
		Request   admission.Request
		Response  admission.Response
	}{
		"simple ok update validation": {
			Validator: ValidatingWebhookFor(ctx, &MyObject{}, nil, []ValidateUpdateFunc{validateUpdateFunc(nil)}, nil),
			Context:   context.TODO(),
			Request: admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					Operation: admissionv1.Update,
					Object: runtime.RawExtension{
						Raw:    []byte(`{"metadata":{"name":"abc","namespace":"default"}}`),
						Object: &MyObject{},
					},
					OldObject: runtime.RawExtension{
						Raw:    []byte(`{"metadata":{"name":"def","namespace":"default"}}`),
						Object: &MyObject{},
					},
				},
			},
			Response: admission.Allowed(""),
		},
		"error update validation": {
			Validator: ValidatingWebhookFor(ctx, &MyObject{}, nil, nil, nil),
			Context:   context.TODO(),
			Request: admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					Operation: admissionv1.Update,
					Object: runtime.RawExtension{
						Raw:    []byte(`{"metadata":{"name":"abc","namespace":""}}`),
						Object: &MyObject{},
					},
					OldObject: runtime.RawExtension{
						Raw:    []byte(`{"metadata":{"name":"def","namespace":"default"}}`),
						Object: &MyObject{},
					},
				},
			},
			Response: validationResponseFromStatus(metav1.Status{
				Status:  metav1.StatusFailure,
				Code:    http.StatusBadRequest,
				Reason:  metav1.StatusReasonBadRequest,
				Message: "needs to have a namespace",
			}),
		},
		"returns error from extra added validation": {
			Validator: ValidatingWebhookFor(ctx, &MyObject{}, nil, []ValidateUpdateFunc{validateUpdateFunc(fmt.Errorf("this is an extra error"))}, nil),
			Context:   context.TODO(),
			Request: admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					Operation: admissionv1.Update,
					Object: runtime.RawExtension{
						Raw:    []byte(`{"metadata":{"name":"abc","namespace":"default"}}`),
						Object: &MyObject{},
					},
					OldObject: runtime.RawExtension{
						Raw:    []byte(`{}`),
						Object: &MyObject{},
					},
				},
			},
			Response: admission.Denied(fmt.Errorf("this is an extra error").Error()),
		},
	}

	for name, test := range table {
		t.Run(name, func(t *testing.T) {

			returned := test.Validator.Handle(test.Context, test.Request)
			diff := cmp.Diff(returned, test.Response)
			t.Logf("diff is: \n%s\n %#v == %#v", diff, test.Response, returned)
			if diff != "" {
				t.Fail()
			}
			if returned.Allowed != test.Response.Allowed {
				t.Fail()
			}
		})
	}
}

func TestValidatorDelete(t *testing.T) {
	ctx := context.Background()
	ctx = kscheme.WithScheme(ctx, scheme.Scheme)

	table := map[string]struct {
		Validator *admission.Webhook
		Context   context.Context
		Request   admission.Request
		Response  admission.Response
	}{
		"simple ok delete validation": {
			Validator: ValidatingWebhookFor(ctx, &MyObject{}, nil, nil, []ValidateDeleteFunc{validateDeleteFunc(nil)}),
			Context:   context.TODO(),
			Request: admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					Operation: admissionv1.Delete,
					OldObject: runtime.RawExtension{
						Raw:    []byte(`{"metadata":{"name":"def","namespace":"default"}}`),
						Object: &MyObject{},
					},
				},
			},
			Response: admission.Allowed(""),
		},
		"error delete validation": {
			Validator: ValidatingWebhookFor(ctx, &MyObject{}, nil, nil, nil),
			Context:   context.TODO(),
			Request: admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					Operation: admissionv1.Delete,
					OldObject: runtime.RawExtension{
						Raw:    []byte(`{"metadata":{"name":"def"}}`),
						Object: &MyObject{},
					},
				},
			},
			Response: validationResponseFromStatus(metav1.Status{
				Status:  metav1.StatusFailure,
				Code:    http.StatusBadRequest,
				Reason:  metav1.StatusReasonBadRequest,
				Message: "needs to have a namespace",
			}),
		},
		"returns error from extra added validation": {
			Validator: ValidatingWebhookFor(ctx, &MyObject{}, nil, nil, []ValidateDeleteFunc{validateDeleteFunc(fmt.Errorf("this is an extra error"))}),
			Context:   context.TODO(),
			Request: admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					Operation: admissionv1.Delete,
					OldObject: runtime.RawExtension{
						Raw:    []byte(`{"metadata":{"name":"def","namespace":"default"}}`),
						Object: &MyObject{},
					},
				},
			},
			Response: admission.Denied(fmt.Errorf("this is an extra error").Error()),
		},
	}

	for name, test := range table {
		t.Run(name, func(t *testing.T) {
			returned := test.Validator.Handle(test.Context, test.Request)
			diff := cmp.Diff(returned, test.Response)
			t.Logf("diff is: \n%s\n %#v == %#v", diff, test.Response, returned)
			if diff != "" {
				t.Fail()
			}
			if returned.Allowed != test.Response.Allowed {
				t.Fail()
			}
		})
	}
}

type injectContextObject struct {
	MyObject
}

func (m *injectContextObject) DeepCopyObject() runtime.Object {
	return &injectContextObject{
		MyObject: m.MyObject,
	}
}

func (i *injectContextObject) InjectContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "foo", "bar")
}

func (i *injectContextObject) ValidateDelete(ctx context.Context) error {
	value := ctx.Value("foo")
	if value.(string) != "bar" {
		panic("context not injected")
	}
	return nil
}

func TestValidatorContextInjector(t *testing.T) {
	ctx := context.Background()
	ctx = kscheme.WithScheme(ctx, scheme.Scheme)

	g := gomega.NewGomegaWithT(t)
	req := admission.Request{
		AdmissionRequest: admissionv1.AdmissionRequest{
			Operation: admissionv1.Delete,
			OldObject: runtime.RawExtension{
				Raw:    []byte(`{"metadata":{"name":"def","namespace":"default"}}`),
				Object: &injectContextObject{},
			},
		},
	}
	table := map[string]struct {
		Validator *admission.Webhook
	}{
		"object not implementing context injector interface": {
			Validator: ValidatingWebhookFor(ctx, &MyObject{}, nil, nil, nil),
		},
		"object implement the context injector interface": {
			Validator: ValidatingWebhookFor(ctx, &injectContextObject{}, nil, nil, nil),
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

type fakeValidator struct {
	ErrorToReturn error `json:"ErrorToReturn,omitempty"`
	metav1.ObjectMeta
}

var _ Validator = &fakeValidator{}

var fakeValidatorVK = schema.GroupVersionKind{Group: "foo.test.org", Version: "v1", Kind: "fakeValidator"}

func (v *fakeValidator) ValidateCreate(ctx context.Context) error {
	return v.ErrorToReturn
}

func (v *fakeValidator) ValidateUpdate(ctx context.Context, old runtime.Object) error {
	return v.ErrorToReturn
}

func (v *fakeValidator) ValidateDelete(ctx context.Context) error {
	return v.ErrorToReturn
}

func (v *fakeValidator) GetObjectKind() schema.ObjectKind { return v }

func (v *fakeValidator) DeepCopyObject() runtime.Object {
	return &fakeValidator{ErrorToReturn: v.ErrorToReturn}
}

func (v *fakeValidator) GroupVersionKind() schema.GroupVersionKind {
	return fakeValidatorVK
}

func (v *fakeValidator) SetGroupVersionKind(gvk schema.GroupVersionKind) {}

// testConfigManagerValidatorState holds captured values shared across DeepCopy instances.
type testConfigManagerValidatorState struct {
	CapturedConfigManager *kconfig.Manager
	CapturedClient        interface{}
}

// testConfigManagerValidatorObj is a Validator that captures the ConfigManager from context
// during Validate*() to allow verification in tests.
// The state pointer is shared across DeepCopy instances so captures are visible to the original.
type testConfigManagerValidatorObj struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// state holds captured values shared across DeepCopy instances.
	state *testConfigManagerValidatorState
}

func (v *testConfigManagerValidatorObj) DeepCopyObject() runtime.Object {
	return &testConfigManagerValidatorObj{
		TypeMeta:   v.TypeMeta,
		ObjectMeta: v.ObjectMeta,
		state:      v.state, // share state pointer so Validate*() captures are observable
	}
}

func (v *testConfigManagerValidatorObj) ValidateCreate(ctx context.Context) error {
	if v.state != nil {
		v.state.CapturedConfigManager = kconfig.ConfigManager(ctx)
		v.state.CapturedClient = kclient.Client(ctx)
	}
	return nil
}

func (v *testConfigManagerValidatorObj) ValidateUpdate(ctx context.Context, old runtime.Object) error {
	if v.state != nil {
		v.state.CapturedConfigManager = kconfig.ConfigManager(ctx)
		v.state.CapturedClient = kclient.Client(ctx)
	}
	return nil
}

func (v *testConfigManagerValidatorObj) ValidateDelete(ctx context.Context) error {
	if v.state != nil {
		v.state.CapturedConfigManager = kconfig.ConfigManager(ctx)
		v.state.CapturedClient = kclient.Client(ctx)
	}
	return nil
}

// TestValidatorInjectsConfigManagerFromHandlerCtx verifies that
// validatingHandler.Handle propagates the ConfigManager stored in h.ctx
// (set at webhook construction time) into the per-request ctx so that
// Validator implementations can access it.
func TestValidatorInjectsConfigManagerFromHandlerCtx(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	// Build a base context with scheme (required for decoder).
	baseCtx := context.Background()
	baseCtx = kscheme.WithScheme(baseCtx, scheme.Scheme)

	// Prepare a non-nil ConfigManager to inject via h.ctx.
	configManager := &kconfig.Manager{}

	// handlerCtx simulates the context passed to ValidatingWebhookFor (stored as h.ctx).
	handlerCtx := kconfig.WithConfigManager(baseCtx, configManager)

	table := map[string]struct {
		description       string
		handlerCtx        context.Context
		req               admission.Request
		wantConfigManager *kconfig.Manager
	}{
		"ConfigManager is injected on Create when set in handler ctx": {
			handlerCtx:        handlerCtx,
			wantConfigManager: configManager,
			req: admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					Operation: admissionv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(`{"metadata":{"name":"test","namespace":"default"}}`),
					},
				},
			},
		},
		"ConfigManager is nil when not set in handler ctx on Create": {
			handlerCtx:        baseCtx,
			wantConfigManager: nil,
			req: admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					Operation: admissionv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(`{"metadata":{"name":"test","namespace":"default"}}`),
					},
				},
			},
		},
		"ConfigManager is injected on Delete when set in handler ctx": {
			handlerCtx:        handlerCtx,
			wantConfigManager: configManager,
			req: admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					Operation: admissionv1.Delete,
					OldObject: runtime.RawExtension{
						Raw: []byte(`{"metadata":{"name":"test","namespace":"default"}}`),
					},
				},
			},
		},
	}

	for name, tc := range table {
		t.Run(name, func(t *testing.T) {
			state := &testConfigManagerValidatorState{}
			obj := &testConfigManagerValidatorObj{state: state}
			wh := ValidatingWebhookFor(tc.handlerCtx, obj, nil, nil, nil)
			response := wh.Handle(baseCtx, tc.req)
			g.Expect(response.Allowed).To(gomega.BeTrue())
			g.Expect(state.CapturedConfigManager).To(gomega.Equal(tc.wantConfigManager))
		})
	}
}
