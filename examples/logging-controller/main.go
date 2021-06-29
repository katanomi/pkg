package main

import (
	"context"
	"github.com/katanomi/pkg/encoder"

	"github.com/katanomi/pkg/sharedmain"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func main() {
	sharedmain.Main("test", scheme.Scheme, &Controller{})
}

type Controller struct {
}

func (c *Controller) Name() string {
	return "controller-test-logging"
}

func (c *Controller) Setup(ctx context.Context, mgr manager.Manager, logger *zap.SugaredLogger) error {
	logger.Infow("info msg", "hello", "world")
	return nil
}

type ControllerFortest struct {
}

func (c *ControllerFortest) Name() string {
	return "controller-test-logging"
}

func (c *ControllerFortest) Setup(ctx context.Context, mgr manager.Manager, logger *zap.SugaredLogger) error {
	logger.Infow("info msg", "hello", "katanomi")
	return nil
}

func init() {
	encoder.SetFilter("filter", []string{"controller-test-logging"})
}
