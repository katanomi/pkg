package sharedmain

import (
	"context"

	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// Controller is a basic interface that every reconciler should implement to create
// a new controller and startup in the controller manager
type Controller interface {
	Name() string
	Setup(context.Context, manager.Manager, *zap.SugaredLogger) error
}
