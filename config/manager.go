/*
Copyright 2022 The Katanomi Authors.

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

package config

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/katanomi/pkg/watcher"

	"go.uber.org/zap"

	corev1 "k8s.io/api/core/v1"

	"github.com/katanomi/pkg/storage/configmap"
	"knative.dev/pkg/system"
)

const (
	defaultConfig = "katanomi-config"
	configNameEnv = "KATANOMI_CONFIG_NAME"
)

// Config store katanomi configuration
type Config struct {
	Data map[string]string
}

// Manager will manage katanomi configuration and store in Config
type Manager struct {
	Informer watcher.DefaultingWatcherWithOnChange
	Logger   *zap.SugaredLogger
	lock     sync.RWMutex
	*Config
}

// NewManager will instantiate a manager that watch configmap for core component configuration
func NewManager(informer watcher.DefaultingWatcherWithOnChange, logger *zap.SugaredLogger, cmName string) *Manager {
	manager := Manager{
		Informer: informer,
		Logger:   logger,
		Config: &Config{
			Data: make(map[string]string),
		},
	}
	coreCM := &corev1.ConfigMap{}
	coreCM.Namespace = system.Namespace()
	coreCM.Name = cmName

	watcher := configmap.NewWatcher(cmName, manager.Informer)
	watcher.AddWatch(coreCM.GetName(), configmap.NewConfigConstructor(coreCM, func(cm *corev1.ConfigMap) {
		manager.applyConfig(cm)
	}))
	watcher.Run()
	return &manager
}

// GetConfig will return the config of manager
func (manager *Manager) GetConfig() *Config {
	if manager == nil {
		return nil
	}
	manager.lock.Lock()
	defer manager.lock.Unlock()
	return manager.Config
}

func (manager *Manager) applyConfig(cm *corev1.ConfigMap) {
	if cm == nil || cm.Data == nil {
		manager.Logger.Errorw("configmap or configmap.data is null", "configmap", fmt.Sprintf("%s/%s", cm.Namespace, cm.Name))
		return
	}
	manager.lock.Lock()
	defer manager.lock.Unlock()
	// whole replacement
	manager.Config = &Config{
		Data: cm.Data,
	}
}

// Name return config name for configuration
func Name() string {
	if name := os.Getenv(configNameEnv); name != "" {
		return name
	}
	return defaultConfig
}

type katanomiConfigKey struct{}

// WithKatanomiConfigManager sets a Config Manager instance into a context
func WithKatanomiConfigManager(ctx context.Context, manager *Manager) context.Context {
	return context.WithValue(ctx, katanomiConfigKey{}, manager)
}

// KatanomiConfigManager returns a Config Manager, returns nil if not found
func KatanomiConfigManager(ctx context.Context) *Manager {
	val := ctx.Value(katanomiConfigKey{})
	if val == nil {
		return nil
	}
	return val.(*Manager)
}
