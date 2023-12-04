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

	// BuildMRCheckTimeoutKey represent build merge request status check timeout key
	BuildMRCheckTimeoutKey = "build.mergerequest.checkTimeout"

	// PprofEnabledKey indicates the configuration key of the /debug/pprof debugging api/
	PprofEnabledKey = "pprof.enabled"
)

const (
	// DefaultVersionEnabled indicates the default value of the version feature gate.
	// If the corresponding key does not exist, the default value is returned.
	DefaultVersionEnabled FeatureValue = "false"

	// DefaultInitializeAllowLocalRequests indicates the configuration key of.
	// If the corresponding key does not exist, the default value is returned.
	DefaultInitializeAllowLocalRequests FeatureValue = "true"

	// DefaultPrunerDelayAfterCompleted represent default duration for delay taskRun
	// If the corresponding key does not exist, the default value is returned.
	DefaultPrunerDelayAfterCompleted FeatureValue = "time.Hour"

	// DefaultPrunerKeep represent default keep number for taskRun
	// If the corresponding key does not exist, the default value is returned.
	DefaultPrunerKeep FeatureValue = "10000"

	// DefaultMRCheckTimeout represent default timeout for merge request status check
	DefaultMRCheckTimeout FeatureValue = "10m"

	// DefaultPprofEnabled stores the default value "false" for the "pprof.enabled" /debug/pprof debugging api.
	// If the corresponding key does not exist, the default value is returned.
	DefaultPprofEnabled FeatureValue = "false"
)

// defaultFeatureValue defines the default value for the feature switch.
var defaultFeatureValue = map[string]FeatureValue{
	VersionEnabledFeatureKey:               DefaultVersionEnabled,
	InitializeAllowLocalRequestsFeatureKey: DefaultInitializeAllowLocalRequests,
	PrunerDelayAfterCompletedFeatureKey:    DefaultPrunerDelayAfterCompleted,
	PrunerKeepFeatureKey:                   DefaultPrunerKeep,
	BuildMRCheckTimeoutKey:                 DefaultMRCheckTimeout,
	PprofEnabledKey:                        DefaultPprofEnabled,
}

// FeatureFlags holds the features configurations
type FeatureFlags struct {
	Data map[string]string
}

// FeatureValue returns the value of the implemented feature flag, or the default if not found.
func (f *FeatureFlags) FeatureValue(flag string) FeatureValue {
	defaultValue := defaultFeatureValue[flag]
	if f == nil || f.Data == nil {
		return defaultValue
	}

	if value, ok := f.Data[flag]; ok {
		return FeatureValue(value)
	}
	return defaultValue
}
