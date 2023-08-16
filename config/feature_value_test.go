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
	"context"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestFeatureValue_AsInt(t *testing.T) {
	tests := map[string]struct {
		f       FeatureValue
		want    int
		wantErr bool
	}{
		"feature is int value":     {f: FeatureValue("1000"), want: 1000},
		"feature is not int value": {f: FeatureValue("failed"), want: -1, wantErr: true},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tt.f.AsInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("FeatureValue.AsInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FeatureValue.AsInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFeatureValue_AsDuration(t *testing.T) {
	tests := map[string]struct {
		f       FeatureValue
		want    time.Duration
		wantErr bool
	}{
		"feature is string Duration value": {f: FeatureValue("1h"), want: time.Hour},
		"feature is not Duration value":    {f: FeatureValue("failed"), want: 0, wantErr: true},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tt.f.AsDuration()
			if (err != nil) != tt.wantErr {
				t.Errorf("FeatureValue.AsInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FeatureValue.AsInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFeatureValue_AsBool(t *testing.T) {
	tests := map[string]struct {
		f       FeatureValue
		want    bool
		wantErr bool
	}{
		"feature is int value":     {f: FeatureValue("true"), want: true},
		"feature is not int value": {f: FeatureValue("failed"), want: false, wantErr: true},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tt.f.AsBool()
			if (err != nil) != tt.wantErr {
				t.Errorf("FeatureValue.AsInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FeatureValue.AsInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetDurationConfig(t *testing.T) {
	RegisterTestingT(t)

	for _, c := range []struct {
		description string
		context     context.Context
		key         string
		defaultTime time.Duration
		expect      time.Duration
	}{
		{
			description: "return default duration",
			context:     context.Background(),
			key:         "not exist",
			defaultTime: time.Hour,
			expect:      time.Hour,
		},
		{
			description: "return default config duration",
			context:     WithKatanomiConfigManager(context.Background(), &Manager{}),
			key:         TemplateRenderRetentionTimeKey,
			defaultTime: time.Hour,
			expect:      30 * time.Minute,
		},
	} {
		t.Logf("<=== starting %s...", c.description)
		Expect(GetDurationConfig(c.context, c.key, c.defaultTime)).To(Equal(c.expect))
		t.Logf("===> passed %s...", c.description)
	}

}
