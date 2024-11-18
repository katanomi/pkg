/*
Copyright 2024 The AlaudaDevops Authors.

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

// Package assertions contains all the assertions that can be used in tests
// for cluster testing. All methods are Gomega compatible as async assertions
// or assertions that can be used directly with other Gomega methods
package assertions

import (
	"context"
	"fmt"

	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/errors"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"go.uber.org/zap"
)

// CheckObjectFunc is a function that checks an object and returns an error if it fails
// useful to be used with WaitResource AsyncAssertion for custom object checks
type CheckObjectFunc func(ctx TestContexter, obj client.Object) error

// TestContexter a basic text context that provides Client, Context and SugeredLogger
type TestContexter interface {
	GetClient() client.Client
	GetLogger() *zap.SugaredLogger
	GetContext() context.Context
}

// WaitObject will create an Eventually AsyncAssertion that
// checks for a resource status given the checkFunc function
// will print an error and return if getting the resource fails.
// Since it returns AsyncAssertion it is possible to use it with other
// Gomega methods
// Example:
//
//	checkPod := func(ctx TestContexter, obj client.Object) error { return check.Pod.Check(ctx, obj) }
//	pod := &corev1.Pod{ObjectMeta{Name: "pod", Namespace: "default"}}
//	WaitObject(ctx, pod, checkPod).To(Succeed())
//
//	// or with custom timeout and intervals
//	WaitObject(ctx, pod, checkPod).
//	  WithPolling(time.Second).
//	  WithTimeout(time.Minute*10).To(Succeed())
func WaitObject(ctx TestContexter, obj client.Object, checkFunc ...CheckObjectFunc) AsyncAssertion {
	key := client.ObjectKey{Name: obj.GetName(), Namespace: obj.GetNamespace()}
	return Eventually(func() error {
		if err := ctx.GetClient().Get(ctx.GetContext(), key, obj); err != nil {
			ctx.GetLogger().Warnw("error fetching object", "key", key, "err", err)
			return err
		}
		errs := make([]error, 0, len(checkFunc))
		for _, check := range checkFunc {
			if err := check(ctx, obj); err != nil {
				errs = append(errs, err)
			}
		}
		return errors.NewAggregate(errs)
	})
}

// IsDoneObject interface to have a client.Object
// that also implements a IsDone() bool function
type IsDoneObject interface {
	client.Object
	IsDone() bool
}

// IsDone will check if a object is done using IsDone() bool function
// and return nil when it is. Otherwise will return an error
// printing the object struct and current data.
// Recommended to be used with WaitResource method
// Example:
//
//	WaitObject(ctx, obj, IsDone).To(Succeed(), "object should be done: %#v"  obj)
func IsDone(ctx TestContexter, obj client.Object) error {
	if runnable, ok := obj.(IsDoneObject); ok {
		if runnable.IsDone() {
			return nil
		}
		return fmt.Errorf("object is not done yet obj.IsDone() == false: %#v", obj)
	}
	return fmt.Errorf("Object %v is not IsDoneObject and cannot be used with this function", obj)
}
