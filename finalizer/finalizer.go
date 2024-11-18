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

// Package finalizer provides a set of functions to manage finalizers
package finalizer

import (
	"context"

	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// AddFinalizer adds a finalizer to the object
func AddFinalizer(ctx context.Context, clt client.Client, o client.Object, finalizerKey string) error {
	finalizers := o.GetFinalizers()
	if controllerutil.ContainsFinalizer(o, finalizerKey) {
		return nil
	}

	toUpdate := o.DeepCopyObject().(client.Object)
	toUpdate.SetFinalizers(append(finalizers, finalizerKey))
	err := clt.Patch(ctx, toUpdate, client.MergeFrom(o))
	if err != nil {
		logging.FromContext(ctx).Errorw("failed to add finalizer", "err", err,
			"namespacedName", client.ObjectKeyFromObject(o),
			"finalizerKey", finalizerKey,
		)
		return err
	}
	o.SetFinalizers(append(finalizers, finalizerKey))
	o.SetResourceVersion(toUpdate.GetResourceVersion())
	return nil
}

// PrependFinalizer adds a finalizer to the object and prepends it to the list
// NOTE: do not use this method as much as possible, use AddFinalizer instead
func PrependFinalizer(ctx context.Context, clt client.Client, o client.Object, finalizerKey string) error {
	finalizers := o.GetFinalizers()
	if controllerutil.ContainsFinalizer(o, finalizerKey) {
		return nil

	}
	toUpdate := o.DeepCopyObject().(client.Object)
	toUpdate.SetFinalizers(append([]string{finalizerKey}, finalizers...))
	err := clt.Patch(ctx, toUpdate, client.MergeFrom(o))
	if err != nil {
		logging.FromContext(ctx).Errorw("failed to append finalizer", "err", err,
			"namespacedName", client.ObjectKeyFromObject(o),
			"finalizerKey", finalizerKey,
		)
		return err
	}
	o.SetFinalizers(append([]string{finalizerKey}, finalizers...))
	o.SetResourceVersion(toUpdate.GetResourceVersion())
	return nil
}

// RemoveFinalizer removes a finalizer from the object
func RemoveFinalizer(ctx context.Context, clt client.Client, o client.Object, finalizerKey string,
	callback func(context.Context) error) error {
	if !controllerutil.ContainsFinalizer(o, finalizerKey) {
		return nil
	}

	if callback != nil {
		if err := callback(ctx); err != nil {
			return err
		}
	}

	toUpdate := o.DeepCopyObject().(client.Object)
	controllerutil.RemoveFinalizer(toUpdate, finalizerKey)
	err := clt.Patch(ctx, toUpdate, client.MergeFrom(o))
	if client.IgnoreNotFound(err) != nil {
		logging.FromContext(ctx).Errorw("failed to remove finalizer", "err", err,
			"namespacedName", client.ObjectKeyFromObject(o),
			"finalizerKey", finalizerKey,
		)
		return err
	}
	controllerutil.RemoveFinalizer(o, finalizerKey)
	o.SetResourceVersion(toUpdate.GetResourceVersion())
	return nil
}
