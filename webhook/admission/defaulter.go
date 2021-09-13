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
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// Defaulter defines functions for setting defaults on resources
type Defaulter interface {
	runtime.Object
	metav1.Object
	Default(context.Context)
}

// DefaultingWebhookFor creates a new Webhook for Defaulting the provided type.
func DefaultingWebhookFor(ctx context.Context, defaulter Defaulter, transforms ...TransformFunc) *admission.Webhook {
	return &admission.Webhook{
		Handler: &mutatingHandler{
			defaulter:     defaulter,
			transforms:    transforms,
			SugaredLogger: logging.FromContext(ctx),
		},
	}
}

type mutatingHandler struct {
	defaulter  Defaulter
	decoder    *admission.Decoder
	transforms []TransformFunc

	*zap.SugaredLogger
}

var _ admission.DecoderInjector = &mutatingHandler{}

// InjectDecoder injects the decoder into a mutatingHandler.
func (h *mutatingHandler) InjectDecoder(d *admission.Decoder) error {
	h.decoder = d
	return nil
}

// Handle handles admission requests.
func (h *mutatingHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	if h.defaulter == nil {
		panic("defaulter should never be nil")
	}

	ctx = logging.WithLogger(ctx, h.SugaredLogger)
	ctx = WithAdmissionRequest(ctx, req)

	// Get the object in the request
	obj := h.defaulter.DeepCopyObject().(Defaulter)
	err := h.decoder.Decode(req, obj)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	switch req.Operation {
	case admissionv1.Create:
		ctx = apis.WithinCreate(ctx)
	case admissionv1.Update:
		old := h.defaulter.DeepCopyObject()
		err := h.decoder.DecodeRaw(req.OldObject, old)
		if err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}
		ctx = apis.WithinUpdate(ctx, old)
	case admissionv1.Delete:
		ctx = apis.WithinDelete(ctx)
	}

	// apply some common transformations before handling defaults
	for _, transform := range h.transforms {
		transform(ctx, obj, req)
	}

	// Default the object
	obj.Default(ctx)
	marshalled, err := json.Marshal(obj)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	// Create the patch
	return admission.PatchResponseFromRaw(req.Object.Raw, marshalled)
}
