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

// Package fake provide fake manager for config
package fake

import "github.com/katanomi/pkg/config"

var fake config.ManagerInterface = &FakeManager{}

type FakeManager struct {
	Data         map[string]string
	FeatureFlags map[string]config.FeatureValue
}

func (m *FakeManager) GetConfig() *config.Config {
	return &config.Config{Data: m.Data}
}

func (m *FakeManager) GetFeatureFlag(flag string) config.FeatureValue {
	return m.FeatureFlags[flag]
}
