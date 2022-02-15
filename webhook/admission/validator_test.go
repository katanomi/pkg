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

package admission

import (
	"context"
	"fmt"
	"net/http"
	"testing"

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
	table := map[string]struct {
		Validator *admission.Webhook
		Context   context.Context
		Request   admission.Request
		Response  admission.Response
	}{
		"simple ok create validation": {
			Validator: ValidatingWebhookFor(context.TODO(), &MyObject{}, []ValidateCreateFunc{validateCreateFunc(nil)}, nil, nil),
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
			Validator: ValidatingWebhookFor(context.TODO(), &MyObject{}, nil, nil, nil),
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
			Validator: ValidatingWebhookFor(context.TODO(), &MyObject{}, []ValidateCreateFunc{validateCreateFunc(fmt.Errorf("this is an extra error"))}, nil, nil),
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

	decoder, _ := admission.NewDecoder(scheme.Scheme)
	t.Logf("decoder: %#v", decoder)
	for name, test := range table {
		if inject, ok := test.Validator.Handler.(admission.DecoderInjector); ok {
			inject.InjectDecoder(decoder)
		}
		t.Run(name, func(t *testing.T) {

			returned := test.Validator.Handle(test.Context, test.Request)
			diff := cmp.Diff(returned, test.Response)
			t.Logf("diff is: \n%s\n %#v == %#v", diff, test.Response, returned)
			if diff != "" {
				t.Failed()
			}
			if returned.Allowed != test.Response.Allowed {
				t.Failed()
			}
		})
	}
}

func TestValidatorUpdate(t *testing.T) {
	table := map[string]struct {
		Validator *admission.Webhook
		Context   context.Context
		Request   admission.Request
		Response  admission.Response
	}{
		"simple ok update validation": {
			Validator: ValidatingWebhookFor(context.TODO(), &MyObject{}, nil, []ValidateUpdateFunc{validateUpdateFunc(nil)}, nil),
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
			Validator: ValidatingWebhookFor(context.TODO(), &MyObject{}, nil, nil, nil),
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
			Validator: ValidatingWebhookFor(context.TODO(), &MyObject{}, nil, []ValidateUpdateFunc{validateUpdateFunc(fmt.Errorf("this is an extra error"))}, nil),
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

	decoder, _ := admission.NewDecoder(scheme.Scheme)
	t.Logf("decoder: %#v", decoder)
	for name, test := range table {
		if inject, ok := test.Validator.Handler.(admission.DecoderInjector); ok {
			inject.InjectDecoder(decoder)
		}
		t.Run(name, func(t *testing.T) {

			returned := test.Validator.Handle(test.Context, test.Request)
			diff := cmp.Diff(returned, test.Response)
			t.Logf("diff is: \n%s\n %#v == %#v", diff, test.Response, returned)
			if diff != "" {
				t.Failed()
			}
			if returned.Allowed != test.Response.Allowed {
				t.Failed()
			}
		})
	}
}

func TestValidatorDelete(t *testing.T) {
	table := map[string]struct {
		Validator *admission.Webhook
		Context   context.Context
		Request   admission.Request
		Response  admission.Response
	}{
		"simple ok delete validation": {
			Validator: ValidatingWebhookFor(context.TODO(), &MyObject{}, nil, nil, []ValidateDeleteFunc{validateDeleteFunc(nil)}),
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
			Validator: ValidatingWebhookFor(context.TODO(), &MyObject{}, nil, nil, nil),
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
			Response: validationResponseFromStatus(metav1.Status{
				Status:  metav1.StatusFailure,
				Code:    http.StatusBadRequest,
				Reason:  metav1.StatusReasonBadRequest,
				Message: "needs to have a namespace",
			}),
		},
		"returns error from extra added validation": {
			Validator: ValidatingWebhookFor(context.TODO(), &MyObject{}, nil, nil, []ValidateDeleteFunc{validateDeleteFunc(fmt.Errorf("this is an extra error"))}),
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

	decoder, _ := admission.NewDecoder(scheme.Scheme)
	t.Logf("decoder: %#v", decoder)
	for name, test := range table {
		if inject, ok := test.Validator.Handler.(admission.DecoderInjector); ok {
			inject.InjectDecoder(decoder)
		}
		t.Run(name, func(t *testing.T) {

			returned := test.Validator.Handle(test.Context, test.Request)
			diff := cmp.Diff(returned, test.Response)
			t.Logf("diff is: \n%s\n %#v == %#v", diff, test.Response, returned)
			if diff != "" {
				t.Failed()
			}
			if returned.Allowed != test.Response.Allowed {
				t.Failed()
			}
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
