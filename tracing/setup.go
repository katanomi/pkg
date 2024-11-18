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
	v1 "k8s.io/api/core/v1"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/system"

	kconfigmap "github.com/AlaudaDevops/pkg/configmap"
)

func defaultConfigMap(name string) *v1.ConfigMap {
	cm := &v1.ConfigMap{}
	cm.Name = name
	cm.Namespace = system.Namespace()
	return cm
}

// SetupDynamicPublishing sets up trace publishing for the process, by watching a
// ConfigMap for the configuration. Note that other pieces still need to generate the traces, this
// just ensures that if generated, they are collected appropriately. This is normally done by using
// tracing.HTTPSpanMiddleware as a middleware HTTP handler. The configuration will be dynamically
// updated when the ConfigMap is updated.
func SetupDynamicPublishing(tracing *Tracing, configMapWatcher configmap.DefaultingWatcher) error {
	tracerUpdater := func(cm *v1.ConfigMap) {
		cfg, err := newTracingConfigFromConfigMap(cm)
		if err != nil {
			return
		}
		tracing.ApplyConfig(cfg)
	}

	// Set up our config store.
	dftCm := defaultConfigMap(tracing.ConfigMapName)
	w := kconfigmap.NewWatcher("config-tracing-store", configMapWatcher).WithLogger(tracing.logger)
	w.AddWatch(dftCm.GetName(), kconfigmap.NewConfigConstructor(dftCm, func(cm *v1.ConfigMap) {
		tracerUpdater(cm)
	}))
	w.Run()

	return nil
}
