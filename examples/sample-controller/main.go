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

	"github.com/katanomi/pkg/sharedmain"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func main() {
	sharedmain.App("test").
		Scheme(scheme.Scheme).
		Log().
		Profiling().
		Controllers(&Controller{}, &Controller2{}).
		Run()
}

type Controller struct {
}

func (c *Controller) Name() string {
	return "controller-test"
}

func (c *Controller) Setup(ctx context.Context, mgr manager.Manager, logger *zap.SugaredLogger) error {
	logger.Errorf("error msg", "hello", "001")
	logger.Warnw("warn msg", "hello", "001")
	logger.Infow("info msg", "hello", "001")
	logger.Debugf("debug msg", "hello", "001")
	return nil
}

type Controller2 struct {
}

func (c *Controller2) Name() string {
	return "controller-test-bak"
}

func (c *Controller2) Setup(ctx context.Context, mgr manager.Manager, logger *zap.SugaredLogger) error {
	logger.Errorf("error msg", "hello", "002")
	logger.Warnw("warn msg", "hello", "002")
	logger.Infow("info msg", "hello", "002")
	logger.Debugf("debug msg", "hello", "002")
	return nil
}
