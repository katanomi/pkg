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

	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

// RegisterDefaultWebhookFor registers a mutate webhook for the defaulter with transforms
func RegisterDefaultWebhookFor(ctx context.Context, mgr ctrl.Manager, defaulter Defaulter, transforms ...TransformFunc) (err error) {
	var gvk schema.GroupVersionKind
	if gvk, err = apiutil.GVKForObject(defaulter.DeepCopyObject(), mgr.GetScheme()); err != nil {
		return
	}
	mgr.GetWebhookServer().Register(
		generateMutatePath(gvk),
		DefaultingWebhookFor(ctx, defaulter, transforms...),
	)
	return
}

// RegisterValidateWebhookFor registers a mutate webhook for the defaulter with transforms
func RegisterValidateWebhookFor(ctx context.Context, mgr ctrl.Manager, validator Validator, validateCreateFuncs []ValidateCreateFunc, validateUpdateFuncs []ValidateUpdateFunc, validateDeleteFuncs []ValidateDeleteFunc) (err error) {
	var gvk schema.GroupVersionKind
	if gvk, err = apiutil.GVKForObject(validator.DeepCopyObject(), mgr.GetScheme()); err != nil {
		return
	}
	mgr.GetWebhookServer().Register(
		generateValidatePath(gvk),
		ValidatingWebhookFor(ctx, validator, validateCreateFuncs, validateUpdateFuncs, validateDeleteFuncs),
	)
	return
}
