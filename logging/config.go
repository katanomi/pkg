/*
Copyright 2021 The AlaudaDevops Authors.

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

// Package logging contains useful functionality
package logging

import (
	"context"
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
	ControllerLevelMap map[string]*ControllerLevel
	Locker             *sync.Mutex
	Name               string

	// customObservers is a list of custom observers that will be invoked when the log config is updated
	customObservers []func(map[string]string)

	*zap.SugaredLogger
}

// NewLevelManager create a new levelManager object
func NewLevelManager(ctx context.Context, name string) LevelManager {
	return LevelManager{
		BaseLevel:          zap.NewAtomicLevel(),
		ControllerLevelMap: map[string]*ControllerLevel{},
		Locker:             &sync.Mutex{},
		Name:               name,
		// most probably will init with a fallback logger
		// which is production configuration
		// so it is recommended to change the logger once started
		SugaredLogger: logging.FromContext(ctx),
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

// SetLogger overwrites the previous logger
func (l *LevelManager) SetLogger(logger *zap.SugaredLogger) {
	l.SugaredLogger = logger
}

// Update read a configmap to update LevelManager. Set BaseLevel with zap.Config, set ControllerLevels with LoggingLevel
func (l *LevelManager) Update() func(configMap *corev1.ConfigMap) {
	return func(configMap *corev1.ConfigMap) {
		l.Locker.Lock()
		defer l.Locker.Unlock()
		l.Infow("reloading log configuration after change")

		config, err := logging.NewConfigFromConfigMap(configMap)
		if err != nil {
			l.Error("Failed to parse the logging configmap. Previous config map will be used.", "err", err)
			return
		}
		loggingCfg, err := ZapConfigFromJSON(config.LoggingConfig)
		var level zapcore.Level
		switch {
		case errors.Is(err, errEmptyLoggerConfig):
			level = zap.NewAtomicLevel().Level()
		case err != nil:
			l.Error("Failed to parse logger configuration.", "err", err)
			return
		default:
			// the component name should be the base for all logs in the same
			// component, and fallbacks to the level inside "zap-logger-config" if not existing
			namedLevel, ok := config.LoggingLevel[l.Name]
			if ok {
				level = namedLevel
			} else {
				level = loggingCfg.Level.Level()
			}
		}
		oldLevel := l.BaseLevel.Level().String()
		l.BaseLevel.SetLevel(level)
		if oldLevel != level.String() {
			l.Infow("logging base level changed", "old", oldLevel, "current", level)
		}

		for k, v := range config.LoggingLevel {
			if controllerLevel, ok := l.ControllerLevelMap[k]; !ok {
				l.Infow("adding new log level. Obs: This change only takes effect after restaring the pod", "name", k, "level", v)
				l.ControllerLevelMap[k] = &ControllerLevel{Inherit: false, Level: zap.NewAtomicLevelAt(v)}
			} else {
				controllerLevel.Inherit = false
				if controllerLevel.Level.String() != v.String() {
					l.Infow("updating log level", "name", k, "old", controllerLevel.Level, "current", v)
				}
				controllerLevel.Level.SetLevel(v)
			}
		}

		for k, v := range l.ControllerLevelMap {
			if v.Inherit {
				if v.Level.String() != level.String() {
					l.Infow("updating log level", "name", k, "old", v, "current", level)
				}
				v.Level.SetLevel(level)
			}
		}

		for _, observer := range l.customObservers {
			func() {
				defer func() {
					if invokeErr := recover(); invokeErr != nil {
						l.Errorw("failed to invoke log config updater", "err", err)
					}
				}()
				observer(configMap.Data)
			}()
		}
	}
}

// Get find atomiclevel by name.If not found than insert one which inherit baselevel.
func (l *LevelManager) Get(name string) zap.AtomicLevel {
	l.Locker.Lock()
	defer l.Locker.Unlock()
	if _, ok := l.ControllerLevelMap[name]; !ok {
		l.ControllerLevelMap[name] = &ControllerLevel{
			Inherit: true,
			Level:   zap.NewAtomicLevelAt(l.BaseLevel.Level()),
		}
	}
	return l.ControllerLevelMap[name].Level
}

// AddCustomObserver add a custom observer to the level manager
func (l *LevelManager) AddCustomObserver(f func(map[string]string)) {
	l.Locker.Lock()
	defer l.Locker.Unlock()
	l.customObservers = append(l.customObservers, f)
}
