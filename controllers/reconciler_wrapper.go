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

// Importing necessary packages.
import (
	"context"
	"fmt"
	"reflect"
	"runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type WrapperOption func(*WrapperOptions)

type RequestFunc func(ctx context.Context, request reconcile.Request) (context.Context, error)

type ResultFunc func(ctx context.Context, request reconcile.Request, result *reconcile.Result) error

// scheduleReconciler is a wrapper for Reconciler interface from controller-runtime package.
type reconcilerWrapper struct {
	reconciler reconcile.Reconciler

	requestFuncs []RequestFunc
	resultFuncs  []ResultFunc
	Object       client.Object
	Client       client.Client
}

type WrapperOptions struct {
	// define the synchronization interval.
	RequestFuncs []RequestFunc
	ResultFuncs  []ResultFunc
	Object       client.Object
	Client       client.Client
}

// NewReconcilerWrapper creates a new scheduleReconciler with the provided Reconciler and ScheduleOptions.
func NewReconcilerWrapper(r reconcile.Reconciler, opts ...WrapperOption) reconcile.Reconciler {
	options := WrapperOptions{}
	for _, option := range opts {
		option(&options)
	}

	return &reconcilerWrapper{
		reconciler:   r,
		requestFuncs: options.RequestFuncs,
		resultFuncs:  options.ResultFuncs,
		Client:       options.Client,
		Object:       options.Object,
	}
}

// Reconcile is the method that will be called whenever an event occurs that the Reconciler should handle.
// It calls the Reconcile method of the embedded Reconciler and adjusts the RequeueAfter
func (s *reconcilerWrapper) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	if s.reconciler == nil {
		return reconcile.Result{}, fmt.Errorf("reconciler should not be empty")
	}

	if s.Object == nil {
		return reconcile.Result{}, fmt.Errorf("reconciler object should not be empty")
	}

	err := s.Client.Get(ctx, request.NamespacedName, s.Object)
	if err != nil {
		// if object does not exist, do nothing.
		// when the DeletionTimestamp is not 0, should not return here, as its
		// finalizer needs to be processed.
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	// modify request context
	for _, requestFunc := range s.requestFuncs {
		requestCtx, err := requestFunc(ctx, request)
		if err != nil {
			return reconcile.Result{}, funcError(requestFunc, err)
		}
		ctx = requestCtx
	}

	result, err := s.reconciler.Reconcile(ctx, request)
	if err != nil {
		return result, err
	}

	// modify request result
	for _, resultFunc := range s.resultFuncs {
		if err := resultFunc(ctx, request, &result); err != nil {
			return result, funcError(resultFunc, err)
		}
	}

	return result, nil
}

func funcError(f interface{}, err error) error {
	funcValue := reflect.ValueOf(f)
	funcPtr := funcValue.Pointer()
	name := runtime.FuncForPC(funcPtr).Name()

	return fmt.Errorf("%s func return error: %s", name, err.Error())
}
