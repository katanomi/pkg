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
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func Test_parseBoolFeature(t *testing.T) {
	tests := map[string]struct {
		cfgMap       map[string]string
		key          string
		defaultValue bool
		feature      bool
		want         bool
		wantErr      bool
	}{
		"include feature, and parse success": {
			cfgMap:       map[string]string{VersionEnabledFeatureKey: "true"},
			key:          VersionEnabledFeatureKey,
			defaultValue: DefaultVersionEnabled,
			want:         true,
		},
		"include feature, and parse failed": {
			cfgMap:       map[string]string{VersionEnabledFeatureKey: "true1"},
			key:          VersionEnabledFeatureKey,
			defaultValue: DefaultVersionEnabled,
			want:         false,
			wantErr:      true,
		},
		"feature not set": {
			cfgMap:       map[string]string{"not.found": "true1"},
			key:          VersionEnabledFeatureKey,
			defaultValue: DefaultVersionEnabled,
			want:         DefaultVersionEnabled,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if err := parseBoolFeature(tt.cfgMap, tt.key, tt.defaultValue, &tt.feature); (err != nil) != tt.wantErr {
				t.Errorf("parseBoolFeature() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != tt.feature {
				t.Errorf("parseBoolFeature() get %v, want %v", tt.feature, tt.want)
			}
		})
	}
}

func Test_parseDurationFeature(t *testing.T) {
	tests := map[string]struct {
		cfgMap       map[string]string
		key          string
		defaultValue time.Duration
		feature      time.Duration
		want         time.Duration
		wantErr      bool
	}{
		"include feature, and parse success": {
			cfgMap:       map[string]string{PrunerDelayAfterCompletedFeatureKey: "10s"},
			key:          PrunerDelayAfterCompletedFeatureKey,
			defaultValue: DefaultPrunerDelayAfterCompleted,
			want:         time.Second * 10,
		},
		"include feature, and parse failed": {
			cfgMap:       map[string]string{PrunerDelayAfterCompletedFeatureKey: "10s89"},
			key:          PrunerDelayAfterCompletedFeatureKey,
			defaultValue: DefaultPrunerDelayAfterCompleted,
			want:         0,
			wantErr:      true,
		},
		"feature not set": {
			cfgMap:       map[string]string{"not.found": "10s"},
			key:          PrunerDelayAfterCompletedFeatureKey,
			defaultValue: DefaultPrunerDelayAfterCompleted,
			want:         DefaultPrunerDelayAfterCompleted,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if err := parseDurationFeature(tt.cfgMap, tt.key, tt.defaultValue, &tt.feature); (err != nil) != tt.wantErr {
				t.Errorf("parseDurationFeature() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != tt.feature {
				t.Errorf("parseDurationFeature() get %v, want %v", tt.feature, tt.want)
			}
		})
	}
}

func Test_parseIntFeature(t *testing.T) {
	tests := map[string]struct {
		cfgMap       map[string]string
		key          string
		defaultValue int
		feature      int
		want         int
		wantErr      bool
	}{
		"include feature, and parse success": {
			cfgMap:       map[string]string{PrunerKeepFeatureKey: "1089"},
			key:          PrunerKeepFeatureKey,
			defaultValue: DefaultPrunerKeep,
			want:         1089,
		},
		"include feature, and parse failed": {
			cfgMap:       map[string]string{PrunerKeepFeatureKey: "0fd8796"},
			key:          PrunerKeepFeatureKey,
			defaultValue: DefaultPrunerKeep,
			want:         0,
			wantErr:      true,
		},
		"feature not set": {
			cfgMap:       map[string]string{"not.found": "2147"},
			key:          PrunerKeepFeatureKey,
			defaultValue: DefaultPrunerKeep,
			want:         DefaultPrunerKeep,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if err := parseIntFeature(tt.cfgMap, tt.key, tt.defaultValue, &tt.feature); (err != nil) != tt.wantErr {
				t.Errorf("parseIntFeature() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != tt.feature {
				t.Errorf("parseIntFeature() get %v, want %v", tt.feature, tt.want)
			}
		})
	}
}

func TestNewFeatureFlagsFromMap(t *testing.T) {
	tests := map[string]struct {
		cfgMap  map[string]string
		want    *FeatureFlags
		wantErr bool
	}{
		"input nil, return default feature flags": {
			want: &FeatureFlags{
				VersionEnabled:               DefaultVersionEnabled,
				InitializeAllowLocalRequests: DefaultInitializeAllowLocalRequests,
				PrunerDelayAfterCompleted:    DefaultPrunerDelayAfterCompleted,
				PrunerKeep:                   DefaultPrunerKeep,
			},
		},
		"input part feature": {
			cfgMap: map[string]string{
				VersionEnabledFeatureKey: "true",
				PrunerKeepFeatureKey:     "50",
			},
			want: &FeatureFlags{
				VersionEnabled:               true,
				InitializeAllowLocalRequests: DefaultInitializeAllowLocalRequests,
				PrunerDelayAfterCompleted:    DefaultPrunerDelayAfterCompleted,
				PrunerKeep:                   50,
			},
		},
		"inclue parse failed feature": {
			cfgMap: map[string]string{
				VersionEnabledFeatureKey: "true",
				PrunerKeepFeatureKey:     "50789.",
			},
			wantErr: true,
		},
		"inclue other feature": {
			cfgMap: map[string]string{
				VersionEnabledFeatureKey: "true",
				PrunerKeepFeatureKey:     "50",
				"other.key":              "ok",
			},
			want: &FeatureFlags{
				VersionEnabled:               true,
				InitializeAllowLocalRequests: DefaultInitializeAllowLocalRequests,
				PrunerDelayAfterCompleted:    DefaultPrunerDelayAfterCompleted,
				PrunerKeep:                   50,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := NewFeatureFlagsFromMap(tt.cfgMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFeatureFlagsFromMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("NewFeatureFlagsFromMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
