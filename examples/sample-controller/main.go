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

package main

import (
	"context"

	"github.com/katanomi/pkg/multicluster"
	"github.com/katanomi/pkg/sharedmain"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func main() {

	sharedmain.App("test").
		Scheme(scheme.Scheme).
		// Overrides a MultiClusterClient
		// by default will load the multicluster.ClusterRegistryClient
		// MultiClusterClient(Interface).
		Log().
		Profiling().
		Controllers(&Controller{}, &Controller2{}).
		APIDocs().
		Run()
}

type Controller struct {
	*zap.SugaredLogger

	MultiCluster multicluster.Interface
}

func (c *Controller) Name() string {
	return "controller-test"
}

func (c *Controller) Setup(ctx context.Context, mgr manager.Manager, logger *zap.SugaredLogger) error {
	// logger.Errorw("error msg", "hello", "001")
	// logger.Warnw("warn msg", "hello", "001")
	// logger.Infow("info msg", "hello", "001")
	// logger.Debugw("debug msg", "hello", "001")
	c.SugaredLogger = logger
	c.MultiCluster = multicluster.MultiCluster(ctx)

	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.ConfigMap{}).
		Complete(c)
}

func (c *Controller) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	log := c.With("key", req)

	client, err := c.MultiCluster.GetClient(ctx, &corev1.ObjectReference{
		Name:      "my-cluster",
		Namespace: "default",
	}, scheme.Scheme)
	log.Infow("info msg", "client", client, "err", err)
	if client != nil {
		secretList := &corev1.SecretList{}
		err = client.List(ctx, secretList)
		log.Infow("secret list", "err", err, "list(len)", len(secretList.Items))
	}

	// log.Errorw("error msg", "hello", "001")
	// log.Warnw("warn msg", "hello", "001")

	// log.Debugw("debug msg", "hello", "001")
	return
}

type Controller2 struct {
}

func (c *Controller2) Name() string {
	return "controller-test-bak"
}

func (c *Controller2) Setup(ctx context.Context, mgr manager.Manager, logger *zap.SugaredLogger) error {
	// logger.Errorw("error msg", "hello", "002")
	// logger.Warnw("warn msg", "hello", "002")
	// logger.Infow("info msg", "hello", "002")
	// logger.Debugw("debug msg", "hello", "002")
	return nil
}
