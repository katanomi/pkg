/*
Copyright 2024 The Katanomi Authors.

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

package tekton

import (
	"context"

	"github.com/tektoncd/pipeline/pkg/apis/config"
)

// ConfigOptions is a function that can modify the config
type ConfigOptions func(*config.Config)

// ensureConfigDefaults ensures that the config has all the default values set
func ensureConfigDefaults(cfg *config.Config) *config.Config {
	defaultConfigs := config.FromContextOrDefaults(context.TODO())

	if cfg == nil {
		cfg = &config.Config{}
	}
	if cfg.Defaults == nil {
		cfg.Defaults = defaultConfigs.Defaults
	}
	if cfg.FeatureFlags == nil {
		cfg.FeatureFlags = defaultConfigs.FeatureFlags
	}
	if cfg.Metrics == nil {
		cfg.Metrics = defaultConfigs.Metrics
	}
	if cfg.SpireConfig == nil {
		cfg.SpireConfig = defaultConfigs.SpireConfig
	}
	if cfg.Events == nil {
		cfg.Events = defaultConfigs.Events
	}
	if cfg.Tracing == nil {
		cfg.Tracing = defaultConfigs.Tracing
	}
	return cfg
}

// WithEnableAPIFields enables the alpha API fields
func WithEnableAPIFields() ConfigOptions {
	return func(cfg *config.Config) {
		cfg.FeatureFlags.EnableAPIFields = config.AlphaAPIFields
	}
}

// WithEnableParamEnum enables the enum for params
func WithEnableParamEnum(enabled bool) ConfigOptions {
	return func(cfg *config.Config) {
		cfg.FeatureFlags.EnableParamEnum = enabled
	}
}

// WithEnableCELInWhenExpression enables the CEL in when expression
func WithEnableCELInWhenExpression(enabled bool) ConfigOptions {
	return func(cfg *config.Config) {
		cfg.FeatureFlags.EnableCELInWhenExpression = enabled
	}
}

// WithDefaultConfig sets the default config
func WithDefaultConfig(ctx context.Context, options ...ConfigOptions) context.Context {
	cfg := config.FromContextOrDefaults(ctx)
	cfg = ensureConfigDefaults(cfg)

	// set the default options
	definedOptions := []ConfigOptions{
		WithEnableAPIFields(),
		WithEnableParamEnum(true),
		WithEnableCELInWhenExpression(true),
	}

	for _, option := range definedOptions {
		option(cfg)
	}
	for _, option := range options {
		option(cfg)
	}

	return config.ToContext(ctx, cfg)
}
