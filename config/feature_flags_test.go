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

import "testing"

func TestFeatureFlags_FeatureValue(t *testing.T) {
	tests := map[string]struct {
		featureFlags *FeatureFlags
		flag         string
		want         FeatureValue
	}{
		"featureflag is empty": {
			flag: VersionEnabledFeatureKey,
			want: DefaultVersionEnabled,
		},
		"empty data": {
			featureFlags: &FeatureFlags{},
			flag:         VersionEnabledFeatureKey,
			want:         DefaultVersionEnabled,
		},
		"match feature with default": {
			featureFlags: &FeatureFlags{},
			flag:         VersionEnabledFeatureKey,
			want:         DefaultVersionEnabled,
		},
		"not match feature": {
			featureFlags: &FeatureFlags{},
			flag:         "notfound.flag",
			want:         "",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tt.featureFlags.FeatureValue(tt.flag); got != tt.want {
				t.Errorf("FeatureFlags.Config() = %v, want %v", got, tt.want)
			}
		})
	}
}
