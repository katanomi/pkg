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

package logging

import (
	"encoding/json"
	"errors"
	"sync"

	"knative.dev/pkg/logging"

	corev1 "k8s.io/api/core/v1"

	"github.com/blendle/zapdriver"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	errEmptyLoggerConfig = errors.New("empty logger configuration")
)

// ControllerLevel level for Controller
type ControllerLevel struct {
	Inherit bool
	Level   zap.AtomicLevel
}

// LevelManager a manager for level. BaseLevel will be used when can't find level in ControllerLevelMap
type LevelManager struct {
	BaseLevel          zap.AtomicLevel
	ControllerLevelMap map[string]ControllerLevel
	Locker             *sync.Mutex
}

// NewLevelManager create a new levelManager object
func NewLevelManager() LevelManager {
	return LevelManager{
		BaseLevel:          zap.NewAtomicLevel(),
		ControllerLevelMap: map[string]ControllerLevel{},
		Locker:             &sync.Mutex{},
	}
}

func stackdriverConfig() zap.Config {
	cfg := zapdriver.NewProductionConfig()
	cfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	return cfg
}

// ZapConfigFromJSON read zap.Config from json
func ZapConfigFromJSON(configJSON string) (*zap.Config, error) {
	loggingCfg := stackdriverConfig()

	if configJSON != "" {
		if err := json.Unmarshal([]byte(configJSON), &loggingCfg); err != nil {
			return nil, err
		}
	}
	return &loggingCfg, nil
}

// Update read a configmap to update LevelManager. Set BaseLevel with zap.Config, set ControllerLevels with LoggingLevel
func (l *LevelManager) Update() func(configMap *corev1.ConfigMap) {
	return func(configMap *corev1.ConfigMap) {
		l.Locker.Lock()
		defer l.Locker.Unlock()
		logger := zap.NewExample()
		config, err := logging.NewConfigFromConfigMap(configMap)
		if err != nil {
			logger.Error("Failed to parse the logging configmap. Previous config map will be used.", zap.Error(err))
			return
		}
		loggingCfg, err := ZapConfigFromJSON(config.LoggingConfig)
		var level zapcore.Level
		switch {
		case errors.Is(err, errEmptyLoggerConfig):
			level = zap.NewAtomicLevel().Level()
		case err != nil:
			logger.Error("Failed to parse logger configuration.", zap.Error(err))
			return
		default:
			level = loggingCfg.Level.Level()
		}
		l.BaseLevel.SetLevel(level)

		for k, v := range config.LoggingLevel {
			if controllerLevel, ok := l.ControllerLevelMap[k]; !ok {
				l.ControllerLevelMap[k] = ControllerLevel{Inherit: false, Level: zap.NewAtomicLevelAt(v)}
			} else {
				controllerLevel.Inherit = false
				controllerLevel.Level.SetLevel(v)
			}
		}

		for _, v := range l.ControllerLevelMap {
			if v.Inherit {
				v.Level.SetLevel(level)
			}
		}
	}
}

// Get find atomiclevel by name.If not found than insert one which inherit baselevel.
func (l *LevelManager) Get(name string) zap.AtomicLevel {
	l.Locker.Lock()
	defer l.Locker.Unlock()
	if _, ok := l.ControllerLevelMap[name]; !ok {
		l.ControllerLevelMap[name] = ControllerLevel{
			Inherit: true,
			Level:   zap.NewAtomicLevelAt(l.BaseLevel.Level()),
		}
	}
	return l.ControllerLevelMap[name].Level
}
