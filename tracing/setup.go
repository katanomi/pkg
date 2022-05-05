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

package tracing

import (
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"knative.dev/pkg/configmap"
	cminformer "knative.dev/pkg/configmap/informer"
)

func defaultConfigMap(name string) *v1.ConfigMap {
	cm := &v1.ConfigMap{}
	cm.Name = name
	return cm
}

// SetupDynamicPublishing sets up trace publishing for the process, by watching a
// ConfigMap for the configuration. Note that other pieces still need to generate the traces, this
// just ensures that if generated, they are collected appropriately. This is normally done by using
// tracing.HTTPSpanMiddleware as a middleware HTTP handler. The configuration will be dynamically
// updated when the ConfigMap is updated.
func SetupDynamicPublishing(tracing *Tracing, configMapWatcher *cminformer.InformedWatcher) error {
	tracerUpdater := func(name string, value interface{}) {
		if name != tracing.ConfigMapName {
			return
		}
		cfg := value.(*Config)
		tracing.ApplyConfig(cfg)
	}

	// Set up our config store.
	w := NewDftConfigMapWatcher("config-tracing-store", tracing.logger, configMapWatcher)
	w.AddWatch(tracing.ConfigMapName, newTracingConfigFromConfigMap, defaultConfigMap(tracing.ConfigMapName))
	w.Run(tracerUpdater)

	return nil
}

// NewDftConfigMapWatcher constructs new dftConfigMapWatcher
func NewDftConfigMapWatcher(name string, logger *zap.SugaredLogger, informedWatcher *cminformer.InformedWatcher) *dftConfigMapWatcher {
	return &dftConfigMapWatcher{
		name:            name,
		InformedWatcher: informedWatcher,
		constructors:    map[string]interface{}{},
		defaults:        map[string]v1.ConfigMap{},
		logger:          logger,
	}
}

// dftConfigMapWatcher describe configmap watcher
type dftConfigMapWatcher struct {
	*cminformer.InformedWatcher

	name         string
	constructors map[string]interface{}
	defaults     map[string]v1.ConfigMap

	logger *zap.SugaredLogger
}

// AddWatch watch the configmap based on the given name
// If there is no configmap specified in the current namespace,
// the default configuration you provide will be used, if not, `InformedWatcher` will exit
func (d *dftConfigMapWatcher) AddWatch(name string, constructor interface{}, defaultCM *v1.ConfigMap) {
	d.constructors[name] = constructor
	if defaultCM != nil {
		d.defaults[name] = *defaultCM
	}
}

// Run register the watcher to configStore
func (d *dftConfigMapWatcher) Run(onAfterStore ...func(name string, value interface{})) {
	configStore := configmap.NewUntypedStore(d.name, d.logger, d.constructors, onAfterStore...)
	configStore.WatchConfigs(d)
}

// Watch register the watcher to InformedWatcher
func (d *dftConfigMapWatcher) Watch(name string, obs ...configmap.Observer) {
	if cm, ok := d.defaults[name]; ok {
		d.InformedWatcher.WatchWithDefault(cm, obs...)
	} else {
		d.InformedWatcher.Watch(name, obs...)
	}
}
