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

// Package admission contains functions to add and retrieve admission request from context
package admission

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type admissionReqKey struct{}

// WithAdmissionRequest adds an admission request to the context
func WithAdmissionRequest(ctx context.Context, req admission.Request) context.Context {
	return context.WithValue(ctx, admissionReqKey{}, req)
}

// AdmissionRequest returns admission request from context
func AdmissionRequest(ctx context.Context) admission.Request {
	val := ctx.Value(admissionReqKey{})
	if val == nil {
		return admission.Request{}
	}
	return val.(admission.Request)
}
