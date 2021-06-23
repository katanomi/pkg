package main

import (
	"context"

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
	return "controller-test"
}

func (c *Controller) Setup(ctx context.Context, mgr manager.Manager, logger *zap.SugaredLogger) error {
	logger.Infow("info msg", "hello", "world")
	return nil
}
