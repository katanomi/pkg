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

package accessmode

import (
	"encoding/json"
	"os"
	"sync"

	kconfigmap "github.com/katanomi/pkg/storage/configmap"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/system"
)

const (
	configMapNameEnv = "CONFIG_STORAGE_ACCESSMODE_RULE"

	defaultConfigMapName = "katanomi-storage-accessmode-rule"

	rulesKey = "rules"
)

// ConfigMapName gets the name of the tracing ConfigMap
func ConfigMapName() string {
	if name := os.Getenv(configMapNameEnv); name != "" {
		return name
	}
	return defaultConfigMapName
}

// dftCm defatule configmap contains accessMode rules
func dftCm() *corev1.ConfigMap {
	cm := &corev1.ConfigMap{}
	cm.Name = ConfigMapName()
	cm.Namespace = system.Namespace()
	return cm
}

// AccessModeManager accessMode manager supported by storageClass
type AccessModeManager interface {
	// SupportedAccessModes return the accessModes supported by specify storageClass
	SupportedAccessModes(sc *storagev1.StorageClass) []corev1.PersistentVolumeAccessMode
}

// NewDynamicAccessModeManager helper function for constructing AccessModeManager
func NewDynamicAccessModeManager(logger *zap.SugaredLogger, informer configmap.DefaultingWatcher) AccessModeManager {
	dftRules := dftAccessModeRules()

	d := &dynamicAccessModeManager{}
	d.logger = logger
	d.rules = dftRules
	d.dftRules = dftRules

	dftCM := dftCm()
	watcher := kconfigmap.NewWatcher("accessmode-rules-cm", informer)
	watcher.AddWatch(dftCM.GetName(), kconfigmap.NewConfigConstructor(dftCM, func(cm *corev1.ConfigMap) {
		d.applyConfig(cm)
	}))
	watcher.Run()

	return d
}

// dynamicAccessModeManager Users can define the mapping relationship
// between storageClass and accessmode through configMap
type dynamicAccessModeManager struct {
	logger *zap.SugaredLogger

	rules    map[string][]corev1.PersistentVolumeAccessMode
	dftRules map[string][]corev1.PersistentVolumeAccessMode

	lock sync.RWMutex
}

func (d *dynamicAccessModeManager) applyConfig(cm *corev1.ConfigMap) {
	if cm == nil || cm.Data == nil || cm.Data[rulesKey] == "" {
		return
	}
	rules := make(map[string][]corev1.PersistentVolumeAccessMode)
	if err := json.Unmarshal([]byte(cm.Data[rulesKey]), &rules); err != nil {
		d.logger.Errorw("unmarshal rules data err",
			"rules_data", cm.Data["rules"],
			"err", err,
		)
		return
	}

	d.lock.Lock()
	defer d.lock.Unlock()
	newRules := make(map[string][]corev1.PersistentVolumeAccessMode)

	for k, v := range d.dftRules {
		newRules[k] = v
	}

	for k, v := range rules {
		newRules[k] = v
	}

	d.rules = newRules
}

// SupportedAccessModes return the accessModes supported by specify storageClass
func (d *dynamicAccessModeManager) SupportedAccessModes(sc *storagev1.StorageClass) []corev1.PersistentVolumeAccessMode {
	if sc == nil || sc.Provisioner == "" {
		return nil
	}

	d.lock.RLock()
	defer d.lock.RUnlock()

	if d.rules != nil {
		return d.rules[sc.Provisioner]
	}

	return nil
}
