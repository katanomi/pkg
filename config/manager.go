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
	"sigs.k8s.io/controller-runtime/pkg/client"

	"go.uber.org/zap"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kclient "github.com/katanomi/pkg/client"
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

	// source store mange config source object.
	configMapRef *corev1.ObjectReference

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

	manager.configMapRef = &corev1.ObjectReference{
		Name:      coreCM.Name,
		Namespace: coreCM.Namespace,
	}

	watcher := configmap.NewWatcher(cmName, manager.Informer)
	watcher.AddWatch(coreCM.GetName(), configmap.NewConfigConstructor(coreCM, func(cm *corev1.ConfigMap) {
		manager.applyConfig(cm)
	}))
	watcher.Run()
	return &manager
}

func (manager *Manager) isSameConfigMap(obj metav1.Object) bool {
	if manager == nil || manager.configMapRef == nil {
		return false
	}
	return obj.GetName() == manager.configMapRef.Name && obj.GetNamespace() == manager.configMapRef.Namespace
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

// GetFeatureFlagByClient get the config configuration by requesting configmap from the client.
// prioritize the use of GetFeatureFlag, and use the current function in scenarios that require high real-time data.
func (manager *Manager) GetFeatureFlagByClient(ctx context.Context, flag string) FeatureValue {
	if manager == nil || manager.configMapRef == nil {
		// returns the default value of flag.
		return getFeatureFlag(flag, nil)
	}

	clt := kclient.Client(ctx)
	if clt == nil {
		// When the client is not specified in the context, it behaves the same as GetFeatureFlag.
		return manager.GetFeatureFlag(flag)
	}

	cm := &corev1.ConfigMap{}
	err := clt.Get(ctx, client.ObjectKey{Name: manager.configMapRef.Name, Namespace: manager.configMapRef.Namespace}, cm)
	if err != nil {
		// When getting configmap and reporting an error, it behaves the same as GetFeatureFlag.
		return manager.GetFeatureFlag(flag)
	}
	return getFeatureFlag(flag, &Config{Data: cm.Data})
}

// GetFeatureFlag get the function switch data, if the function switch is not set,
// return the default value of the switch.
func (manager *Manager) GetFeatureFlag(flag string) FeatureValue {
	defaultValue := defaultFeatureValue[flag]
	if manager == nil {
		return defaultValue
	}

	manager.lock.Lock()
	defer manager.lock.Unlock()
	return getFeatureFlag(flag, manager.Config)
}

func getFeatureFlag(flag string, config *Config) FeatureValue {
	defaultValue := defaultFeatureValue[flag]
	if config == nil || config.Data == nil {
		return defaultValue
	}

	value, ok := config.Data[flag]
	if ok {
		return FeatureValue(value)
	}
	return defaultValue
}

func (manager *Manager) applyConfig(cm *corev1.ConfigMap) {
	// Almost never reach here since applyConfig will be called
	// only after the watched configmap has been transformed
	if cm == nil {
		return
	}
	if cm.Data == nil {
		manager.Logger.Errorw("config manager configmap data is nil", "configmap", fmt.Sprintf("%s/%s", cm.Namespace, cm.Name))
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
