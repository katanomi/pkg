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
	goerrors "errors"
	"net/http"

	kclient "github.com/katanomi/pkg/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"go.uber.org/zap"
	v1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// Validator defines functions for validating an operation
type Validator interface {
	runtime.Object
	metav1.Object
	ValidateCreate(ctx context.Context) error
	ValidateUpdate(ctx context.Context, old runtime.Object) error
	ValidateDelete(ctx context.Context) error
}

// ValidatingWebhookFor creates a new Webhook for Validating the provided type.
func ValidatingWebhookFor(ctx context.Context, validator Validator, creates []ValidateCreateFunc, updates []ValidateUpdateFunc, deletes []ValidateDeleteFunc) *admission.Webhook {
	client := kclient.Client(ctx)
	if client == nil {
		panic("validating webhook need client")
	}
	return &admission.Webhook{
		Handler: &validatingHandler{
			decoder:       admission.NewDecoder(client.Scheme()),
			validator:     validator,
			creates:       creates,
			updates:       updates,
			deletes:       deletes,
			SugaredLogger: logging.FromContext(ctx),
		},
	}
}

// a internal handler for an extended validation webhook methods
type validatingHandler struct {
	validator Validator
	decoder   *admission.Decoder
	creates   []ValidateCreateFunc
	updates   []ValidateUpdateFunc
	deletes   []ValidateDeleteFunc

	*zap.SugaredLogger
}

// Handle handles admission requests.
func (h *validatingHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	// the below panic was copied from the original controller-runtime validatingHandler
	// controller-runtime/pkg/webhook/admission/validator.go
	if h.validator == nil {
		panic("validator should never be nil")
	}
	ctx = logging.WithLogger(ctx, h.SugaredLogger)
	ctx = WithAdmissionRequest(ctx, req)

	if injector, ok := h.validator.(ContextInjector); ok {
		ctx = injector.InjectContext(ctx)
	}

	// Get the object in the request
	obj := h.validator.DeepCopyObject().(Validator)
	switch req.Operation {
	case v1.Create:
		ctx = apis.WithinCreate(ctx)
		err := h.decoder.Decode(req, obj)
		if err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}
		err = obj.ValidateCreate(ctx)
		if err != nil {
			return convertToResponse(err)
		}
		for _, val := range h.creates {
			if err = val(ctx, obj, req); err != nil {
				return convertToResponse(err)
			}
		}
	case v1.Update:
		oldObj := obj.DeepCopyObject()

		err := h.decoder.DecodeRaw(req.Object, obj)
		if err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}
		err = h.decoder.DecodeRaw(req.OldObject, oldObj)
		if err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}
		ctx = apis.WithinUpdate(ctx, oldObj)

		err = obj.ValidateUpdate(ctx, oldObj)
		if err != nil {
			return convertToResponse(err)
		}
		for _, val := range h.updates {
			if err = val(ctx, obj, oldObj, req); err != nil {
				return convertToResponse(err)
			}
		}
	case v1.Delete:
		ctx = apis.WithinDelete(ctx)
		// In reference to PR: https://github.com/kubernetes/kubernetes/pull/76346
		// OldObject contains the object being deleted
		err := h.decoder.DecodeRaw(req.OldObject, obj)
		if err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}

		err = obj.ValidateDelete(ctx)
		if err != nil {
			return convertToResponse(err)
		}
		for _, val := range h.deletes {
			if err = val(ctx, obj, req); err != nil {
				return convertToResponse(err)
			}
		}
	default:
		// do nothing
	}
	return admission.Allowed("")
}

func convertToResponse(err error) admission.Response {
	var apiStatus errors.APIStatus
	if goerrors.As(err, &apiStatus) {
		return validationResponseFromStatus(apiStatus.Status())
	}
	return admission.Denied(err.Error())
}

// validationResponseFromStatus returns a response for admitting a request with provided Status object.
func validationResponseFromStatus(status metav1.Status) admission.Response {
	// this method was copied from controller-runtime/pkg/webhook/admission/response.go
	resp := admission.Response{
		AdmissionResponse: v1.AdmissionResponse{
			Allowed: false,
			Result:  &status,
		},
	}
	return resp
}
