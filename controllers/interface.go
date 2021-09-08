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

package controllers

import (
	"context"

	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// LazyLoader loads whenever dependencies are ready
type LazyLoader interface {
	LazyLoad(context.Context, manager.Manager, *zap.SugaredLogger, SetupChecker) error
	Start(context.Context) error
}

// Interface is a basic interface that every reconciler should implement to create
// a new controller and startup in the controller manager
type Interface interface {
	Name() string
	Setup(context.Context, manager.Manager, *zap.SugaredLogger) error
}

// SetupChecker controllers with dependencies on other resources will need to implement
// this interface to allow lazy loading
type SetupChecker interface {
	Interface
	CheckSetup(context.Context, manager.Manager, *zap.SugaredLogger) error
}
