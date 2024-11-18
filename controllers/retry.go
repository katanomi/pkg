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

package controllers

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/util/retry"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CreateOrGetWithRetry will retry to create source multi times if encounter error
// however if error is alreadyExist then will get the resource and return it
func CreateOrGetWithRetry(ctx context.Context, clt client.Client, obj client.Object) error {
	logger := logging.FromContext(ctx)
	if clt == nil || obj == nil {
		return fmt.Errorf("client or obj is nil")
	}

	createObj := func() error {
		err := clt.Create(ctx, obj)
		if errors.IsAlreadyExists(err) {
			logger.Warnw("obj %s already exists, try to get it", "object", fmt.Sprintf("%s/%s/%s", obj.GetObjectKind(), obj.GetNamespace(), obj.GetName()))
			return clt.Get(ctx, client.ObjectKeyFromObject(obj), obj)
		}
		return err
	}
	retriable := func(err error) bool {
		return err != nil
	}

	return retry.OnError(retry.DefaultRetry, retriable, createObj)
}
