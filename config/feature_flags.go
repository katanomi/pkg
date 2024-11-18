/*
Copyright 2023 The AlaudaDevops Authors.

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
	// PprofEnabledKey indicates the configuration key of the /debug/pprof debugging api/
	PprofEnabledKey = "pprof.enabled"
)

const (
	// True represents the value "true" for the feature switch.
	True FeatureValue = "true"

	// False represents the value "false" for the feature switch.
	False FeatureValue = "false"

	// DefaultPprofEnabled stores the default value "false" for the "pprof.enabled" /debug/pprof debugging api.
	// If the corresponding key does not exist, the default value is returned.
	DefaultPprofEnabled FeatureValue = False
)

// defaultFeatureValue defines the default value for the feature switch.
var defaultFeatureValue = map[string]FeatureValue{
	PprofEnabledKey: DefaultPprofEnabled,
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
