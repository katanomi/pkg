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

package test

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"k8s.io/client-go/tools/record"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	pkgctrl "github.com/AlaudaDevops/pkg/controllers"
	testv1alpha1 "github.com/AlaudaDevops/pkg/examples/sample-controller/apis/test/v1alpha1"
	"github.com/AlaudaDevops/pkg/sharedmain"
	"go.uber.org/zap"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ sharedmain.Controller = &FooBarReconciler{}

type FooBarReconciler struct {
	Client client.Client
	// to avoid race conditions with the cache, listing matching webhook items
	// needs to use a direct client
	DirectClient  client.Client
	Log           *zap.SugaredLogger
	Scheme        *runtime.Scheme
	EventRecorder record.EventRecorder
	crds          []string
}

var _ pkgctrl.ControllerChecker = &FooBarReconciler{}

// return trigger reconciler name
func (FooBarReconciler) Name() string {
	return "sugar"
}

// Setup will setup webhook and controller
func (r *FooBarReconciler) Setup(ctx context.Context, mgr ctrl.Manager, logger *zap.SugaredLogger) error {
	r.Log = logger
	r.Scheme = mgr.GetScheme()
	r.Client = mgr.GetClient()
	r.EventRecorder = mgr.GetEventRecorderFor(r.Name())
	var err error
	if r.DirectClient, err = client.New(mgr.GetConfig(), client.Options{Scheme: mgr.GetScheme()}); err != nil {
		logger.Errorw("error initializing direct client")
		return err
	}

	return r.SetupWithManager(mgr)
}

func (r *FooBarReconciler) CheckSetup(ctx context.Context, mgr ctrl.Manager, logger *zap.SugaredLogger) error {
	// no dependencies for this reconciler
	return nil
}

func (m *FooBarReconciler) DependentCrdInstalled(ctx context.Context, logger *zap.SugaredLogger) (bool, error) {
	return true, nil
}

//+kubebuilder:rbac:groups=test.katanomi.dev,resources=*,verbs=get;list;watch;update;patch;create;delete
//+kubebuilder:rbac:groups=``,resources=namespaces;configmaps;secrets,verbs=get;list;watch;update;patch

// Reconcile for Trigger should be implemented in each plugin
func (r *FooBarReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	log := r.Log.With("key", req)

	defer func() {
		log.Infow("reconcile finished", "err", err, "res", result)
	}()

	fb := &testv1alpha1.FooBar{}
	err = r.DirectClient.Get(ctx, req.NamespacedName, fb)
	if err != nil {
		err = client.IgnoreNotFound(err)
		return
	}
	return
}

// SetupTrigger sets up the controller with the Manager.
func (r *FooBarReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&testv1alpha1.FooBar{}).
		WithOptions(controller.Options{
			RateLimiter: pkgctrl.DefaultTypedRateLimiter[reconcile.Request](),
		}).Complete(r)
}
