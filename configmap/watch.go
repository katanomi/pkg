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

package configmap

import (
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/configmap"
)

func NewWatcher(watcherName string, informer configmap.DefaultingWatcher) *watcher {
	logger, _ := zap.NewDevelopment()

	w := &watcher{}
	w.name = watcherName
	w.DefaultingWatcher = informer
	w.constructors = make(map[string]interface{})
	w.defaults = make(map[string]corev1.ConfigMap)
	w.logger = logger.Sugar()
	return w
}

// watcher describe configmap watcher
type watcher struct {
	configmap.DefaultingWatcher

	name         string
	constructors map[string]interface{}
	defaults     map[string]corev1.ConfigMap
	onAfterStore []func(name string, value interface{})

	logger *zap.SugaredLogger
}

// WithLogger replace with customize logger
func (d *watcher) WithLogger(logger *zap.SugaredLogger) *watcher {
	if logger != nil {
		d.logger = logger
	}
	return d
}

// AddWatch watch the configmap based on the given name
// If there is no configmap specified in the current namespace,
// the default configuration you provide will be used, if not, `DefaultingWatcher` will exit
func (d *watcher) AddWatch(cmName string, c ConfigConstructor) {
	if c.CmName() != "" {
		cmName = c.CmName()
	}
	handle := c.Handle
	dft := c.Default()

	if handle == nil || cmName == "" {
		return
	}

	d.constructors[cmName] = func(config *corev1.ConfigMap) (*corev1.ConfigMap, error) {
		return config, nil
	}

	d.onAfterStore = append(d.onAfterStore, func(name string, value interface{}) {
		if name != cmName {
			return
		}

		handle(value.(*corev1.ConfigMap))
	})

	if dft != nil {
		d.defaults[cmName] = *dft
	}
}

// Run register the watcher to configStore
func (d *watcher) Run() {
	configStore := configmap.NewUntypedStore(d.name, d.logger, d.constructors, d.onAfterStore...)
	configStore.WatchConfigs(d)
}

// Watch register the watcher to DefaultingWatcher
func (d *watcher) Watch(name string, obs ...configmap.Observer) {
	if cm, ok := d.defaults[name]; ok {
		d.DefaultingWatcher.WatchWithDefault(cm, obs...)
	} else {
		d.DefaultingWatcher.Watch(name, obs...)
	}
}
