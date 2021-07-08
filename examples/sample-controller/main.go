package main

import (
	"context"

	"github.com/katanomi/pkg/sharedmain"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func main() {
	sharedmain.Main("test", scheme.Scheme, &Controller{}, &Controller2{})
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
