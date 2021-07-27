/*
Copyright 2021 The Katanomi Authors.

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

package tracing

import (
	"encoding/json"
	"fmt"
	"strings"
)

type TraceConfig struct {
	JaegerUrl       string  `json:"jaegerUrl"`
	SampleType      string  `json:"sampleType"`
	SampleParam     float64 `json:"sampleParam"`
	SampleServerURL string  `json:"sampleServerUrl"`
	LogSpan         bool    `json:"logSpan"`
}

func (t *TraceConfig) Copy() TraceConfig {
	return TraceConfig{
		JaegerUrl:       t.JaegerUrl,
		SampleType:      t.SampleType,
		SampleParam:     t.SampleParam,
		SampleServerURL: t.SampleServerURL,
		LogSpan:         t.LogSpan,
	}
}

const (
	tracerConfigKey = "tracer-config-key"
)

type configManager struct {
	global  *TraceConfig
	mapping map[string]*TraceConfig
}

func ParseConfig(configMap map[string]string) (*configManager, error) {
	manager := &configManager{
		mapping: make(map[string]*TraceConfig),
	}
	if v, exist := configMap[tracerConfigKey]; exist {
		global, err := manager.parse(v)
		if err != nil {
			return nil, err
		}
		manager.global = global
	}

	for k, v := range configMap {
		if name := strings.TrimPrefix(k, "sampler."); name != k && name != "" {
			if len(v) > 0 {
				config, err := manager.mergeWithGlobal(v)
				if err != nil {
					return nil, fmt.Errorf("parse %s error: %s", name, err.Error())
				}

				manager.mapping[name] = config
			}
		}
	}

	return manager, nil
}

func (c *configManager) Get(name string) *TraceConfig {
	if v, exist := c.mapping[name]; exist {
		return v
	}

	return c.global
}

func (c *configManager) mergeWithGlobal(v string) (*TraceConfig, error) {
	config, err := c.parse(v)
	if err != nil {
		return nil, err
	}

	if c.global != nil {
		cp := c.global.Copy()
		cp.SampleType = config.SampleType
		cp.SampleParam = config.SampleParam
		cp.SampleServerURL = config.SampleServerURL

		return &cp, nil
	}

	return config, nil
}

func (c *configManager) parse(v string) (*TraceConfig, error) {
	config := &TraceConfig{}
	if err := json.Unmarshal([]byte(v), config); err != nil {
		return nil, err
	}
	return config, nil
}
