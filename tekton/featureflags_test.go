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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tektoncd/pipeline/pkg/apis/config"
)

var _ = Describe("ConfigOptions", func() {
	var (
		ctx context.Context
		cfg *config.Config
	)

	BeforeEach(func() {
		// create a new config
		cfg = ensureConfigDefaults(nil)

		// set the config in the context
		ctx = config.ToContext(context.TODO(), cfg)
	})

	Describe("WithEnableAPIFields", func() {
		It("should enable the alpha API fields", func() {
			// create a new config option
			option := WithEnableAPIFields()

			// apply the option to the config
			option(cfg)

			// assert that the flag is enabled
			Expect(cfg.FeatureFlags.EnableAPIFields).To(Equal(config.AlphaAPIFields))
		})
	})

	Describe("WithEnableParamEnum", func() {
		It("should enable the enum for params", func() {
			// create a new config option
			option := WithEnableParamEnum(true)

			// apply the option to the config
			option(cfg)

			// assert that the flag is enabled
			Expect(cfg.FeatureFlags.EnableParamEnum).To(Equal(true))
		})
	})

	Describe("WithEnableCELInWhenExpression", func() {
		It("should enable the CEL in when expression", func() {
			// create a new config option
			option := WithEnableCELInWhenExpression(true)

			// apply the option to the config
			option(cfg)

			// assert that the flag is enabled
			Expect(cfg.FeatureFlags.EnableCELInWhenExpression).To(Equal(true))
		})
	})

	Describe("WithDefaultConfig", func() {
		It("should set the default config", func() {
			// create a new config option
			option := WithDefaultConfig(ctx)

			// apply the option to the config
			ctx = option

			// assert that the flags are set to the default values
			Expect(cfg.FeatureFlags.EnableAPIFields).To(Equal(config.AlphaAPIFields))
			Expect(cfg.FeatureFlags.EnableParamEnum).To(Equal(true))
			Expect(cfg.FeatureFlags.EnableCELInWhenExpression).To(Equal(true))
		})

		It("should overwrite the default config with provided options", func() {
			// create a new config option
			option := WithDefaultConfig(ctx, WithEnableParamEnum(false))

			// apply the option to the config
			ctx = option

			// assert that the flags are set according to the provided options
			Expect(cfg.FeatureFlags.EnableAPIFields).To(Equal(config.AlphaAPIFields))
			Expect(cfg.FeatureFlags.EnableParamEnum).To(Equal(false))
			Expect(cfg.FeatureFlags.EnableCELInWhenExpression).To(Equal(true))
		})
	})
})
