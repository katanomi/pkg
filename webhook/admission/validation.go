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

	"k8s.io/apimachinery/pkg/runtime"

	// "knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// ValidateCreateFunc function to add validation functions when operation is create
// using a context, an object and a request
type ValidateCreateFunc func(ctx context.Context, obj runtime.Object, req admission.Request) error

// ValidateUpdateFunc function to add validation functions when operation is update
// using a context, the current object, the old object and a request
type ValidateUpdateFunc func(ctx context.Context, obj runtime.Object, old runtime.Object, req admission.Request) error

// ValidateDeleteFunc function to add validation functions when operation is delete
// using a context, an object and a request
type ValidateDeleteFunc func(ctx context.Context, obj runtime.Object, req admission.Request) error
