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
	"time"

	"go.uber.org/zap"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// infomerHasAlreadyStartedError is returned when the informer has already started
type infomerHasAlreadyStartedError struct {
	err error
}

// InfomerHasAlreadyStartedError is returned when the informer has already started
func (i infomerHasAlreadyStartedError) Error() string {
	return "informer has already started"
}

// controllerLazyLoader implementation of LazyLoader
type controllerLazyLoader struct {
	ctx context.Context
	mgr manager.Manager
	*zap.SugaredLogger
	interval time.Duration
	pending  []lazyItem
	done     []lazyItem
}

// NewLazyLoader constructs new LazyLoader for controllers
func NewLazyLoader(ctx context.Context, interval time.Duration) LazyLoader {
	return &controllerLazyLoader{
		interval:      interval,
		ctx:           ctx,
		pending:       []lazyItem{},
		done:          []lazyItem{},
		SugaredLogger: logging.FromContext(ctx).Named("lazyloader"),
	}
}

type lazyItem struct {
	logger  *zap.SugaredLogger
	checker ControllerChecker
}

// LazyLoad loads items to lazy load if any error found
func (c *controllerLazyLoader) LazyLoad(ctx context.Context, mgr manager.Manager, logger *zap.SugaredLogger, checker ControllerChecker) error {
	c.ctx = ctx
	c.mgr = mgr
	item := lazyItem{
		logger:  logger,
		checker: checker,
	}

	ok, err := c.checkPending(item)
	if err != nil {
		return err
	}
	if !ok {
		c.pending = append(c.pending, item)
	} else {
		c.done = append(c.done, item)
	}

	return nil
}

func (c *controllerLazyLoader) checkPending(item lazyItem) (ok bool, err error) {

	checkCrdInstalled, err := item.checker.CheckCrdInstalled(c.ctx, c.SugaredLogger)
	if err != nil {
		c.Errorw("failed to check crds", "ctrl", item.checker.Name(), "err", err)
		return false, err
	}
	if !checkCrdInstalled {
		c.Debugw("controller setup is pending by crds", "ctrl", item.checker.Name(), "err", err)
		return false, nil
	}
	c.Infow("controller setup is not pending by crds", "ctrl", item.checker.Name())

	if err = item.checker.CheckSetup(c.ctx, c.mgr, item.logger); err != nil {
		c.Debugw("controller setup is pending", "ctrl", item.checker.Name(), "err", err)
		// errors returned by this function will cause an fatal error in the application
		// therefore here we set a nil to avoid crashing
		err = nil
	} else {
		c.Infow("controller setup started", "ctrl", item.checker.Name())
		if err = item.checker.Setup(c.ctx, c.mgr, item.logger); err != nil {
			c.Errorw("controller setup failed with error", "ctrl", item.checker.Name(), "err", err)
			alreadyStarted := infomerHasAlreadyStartedError{}
			if err.Error() == alreadyStarted.Error() {
				c.Warnw("controller setup already started", "ctrl", item.checker.Name())
				return false, nil
			}
		}
		ok = true
	}
	return
}

// Start starts to check and load controllers
// this method will block execution and should be runned in a goroutine
func (c *controllerLazyLoader) Start(ctx context.Context) error {
	ticker := time.NewTicker(c.interval)
	for {
		select {
		case <-ticker.C:
			if len(c.pending) > 0 {
				c.Infow("layloader controller setup check", "len(pending)", len(c.pending), "len(done)", len(c.done))
			}
			names := []string{}
			for i := 0; i < len(c.pending); i++ {
				item := c.pending[i]
				c.Debugw("checking controller", "ctrl", item.checker.Name())
				ok, err := c.checkPending(item)
				if err != nil {
					return err
				}
				if ok {
					c.pending = append(c.pending[:i], c.pending[i+1:]...)
					i--
					c.done = append(c.done, item)
				} else {
					names = append(names, item.checker.Name())
				}
			}
			if len(names) > 0 {
				c.Infow("still have pending controllers", "ctrls", names)
			}

		case <-ctx.Done():
			c.Infow("shutting down lazy loader")
			return nil
		}
	}
}
