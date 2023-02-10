/*
Copyright 2023 The Katanomi Authors.

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
	"fmt"
	"strconv"
	"time"
)

// FeatureFlags holds the features configurations
type FeatureFlags struct {
	// VersionEnabled whether the version policy is enabled, the default is false
	VersionEnabled bool
	// InitializeAllowLocalRequests indicates gitlab-plugin allow local requests
	InitializeAllowLocalRequests bool
	// PrunerDelayAfterCompleted represent default duration for delay taskRun
	PrunerDelayAfterCompleted time.Duration
	// PrunerKeep represent default keep number for taskRun
	PrunerKeep int
}

const (
	// VersionEnabledFeatureKey indicates the configuration key of the version feature gate.
	// If the value is true, the feature is enabled cluster-wide.
	VersionEnabledFeatureKey = "version.enabled"

	// InitializeAllowLocalRequestsFeatureKey indicates the configuration key of.
	// If the value is true, the feature is enabled cluster-wide.
	InitializeAllowLocalRequestsFeatureKey = "plugin.gitlab.allow-local-requests"

	// PrunerDelayAfterCompletedFeatureKey represent taskRun delay configuration key
	PrunerDelayAfterCompletedFeatureKey = "taskRunPruner.delayAfterCompleted"

	// PrunerKeepFeatureKey represent taskRun keep configuration key
	PrunerKeepFeatureKey = "taskRunPruner.keep"
)

const (
	// DefaultVersionEnabled indicates the default value of the version feature gate.
	// If the corresponding key does not exist, the default value is returned.
	DefaultVersionEnabled = false

	// DefaultInitializeAllowLocalRequests indicates the configuration key of.
	// If the corresponding key does not exist, the default value is returned.
	DefaultInitializeAllowLocalRequests = true

	// DefaultPrunerDelayAfterCompleted represent default duration for delay taskRun
	// If the corresponding key does not exist, the default value is returned.
	DefaultPrunerDelayAfterCompleted = time.Hour

	// DefaultPrunerKeep represent default keep number for taskRun
	// If the corresponding key does not exist, the default value is returned.
	DefaultPrunerKeep = 10000
)

// defaultFeatureValue defines the default value for the feature switch.
var defaultFeatureValue = map[string]interface{}{
	VersionEnabledFeatureKey:               DefaultVersionEnabled,
	InitializeAllowLocalRequestsFeatureKey: DefaultInitializeAllowLocalRequests,
	PrunerDelayAfterCompletedFeatureKey:    DefaultPrunerDelayAfterCompleted,
	PrunerKeepFeatureKey:                   DefaultPrunerKeep,
}

// NewFeatureFlagsFromMap returns a Config given a map corresponding to a ConfigMap
func NewFeatureFlagsFromMap(cfgMap map[string]string) (*FeatureFlags, error) {
	tc := FeatureFlags{}
	if err := parseBoolFeature(cfgMap, VersionEnabledFeatureKey, DefaultVersionEnabled, &tc.VersionEnabled); err != nil {
		return nil, err
	}

	if err := parseBoolFeature(cfgMap, InitializeAllowLocalRequestsFeatureKey, DefaultInitializeAllowLocalRequests, &tc.InitializeAllowLocalRequests); err != nil {
		return nil, err
	}

	if err := parseDurationFeature(cfgMap, PrunerDelayAfterCompletedFeatureKey, DefaultPrunerDelayAfterCompleted, &tc.PrunerDelayAfterCompleted); err != nil {
		return nil, err
	}

	if err := parseIntFeature(cfgMap, PrunerKeepFeatureKey, DefaultPrunerKeep, &tc.PrunerKeep); err != nil {
		return nil, err
	}

	return &tc, nil
}

func parseBoolFeature(cfgMap map[string]string, key string, defaultValue bool, feature *bool) error {
	value := defaultValue
	if cfg, ok := cfgMap[key]; ok {
		v, err := strconv.ParseBool(cfg)
		if err != nil {
			return fmt.Errorf("failed parsing feature flags config %q: %v", cfg, err)
		}
		value = v
	}
	*feature = value
	return nil
}

func parseDurationFeature(cfgMap map[string]string, key string, defaultValue time.Duration, feature *time.Duration) error {
	value := defaultValue
	if cfg, ok := cfgMap[key]; ok {
		v, err := time.ParseDuration(cfg)
		if err != nil {
			return fmt.Errorf("failed parsing feature flags config %q: %v", cfg, err)
		}
		value = v
	}
	*feature = value
	return nil
}

func parseIntFeature(cfgMap map[string]string, key string, defaultValue int, feature *int) error {
	value := defaultValue
	if cfg, ok := cfgMap[key]; ok {
		v, err := strconv.Atoi(cfg)
		if err != nil {
			return fmt.Errorf("failed parsing feature flags config %q: %v", cfg, err)
		}
		value = v
	}
	*feature = value
	return nil
}
