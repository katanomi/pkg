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

package base

import (
	"context"

	"github.com/onsi/ginkgo/v2/types"
)

// SharedExtension describe interface for shared extension
type SharedExtension interface {
	// SetShardInfo set shard info into context
	SetShardInfo(ctx context.Context) context.Context
}

// SharedExtensionFunc helper function to generate a new shared extension
type SharedExtensionFunc func(ctx context.Context) context.Context

// SetShardInfo set shard info into context
func (p SharedExtensionFunc) SetShardInfo(ctx context.Context) context.Context {
	return p(ctx)
}

// Configure describe interface for configure framework
type Configure interface {
	// Config mutate suite config and reporter config
	Config(suiteConfig *types.SuiteConfig, reporterConfig *types.ReporterConfig)
}

// ConfigureFunc helper function to generate a implementation of configure interface
type ConfigureFunc func(suiteConfig *types.SuiteConfig, reporterConfig *types.ReporterConfig)

// Config mutate suite config and reporter config
func (p ConfigureFunc) Config(suiteConfig *types.SuiteConfig, reporterConfig *types.ReporterConfig) {
	p(suiteConfig, reporterConfig)
}
