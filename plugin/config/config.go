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

package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

// Config global config
type Config struct {
	Trace   TraceConfig
	Log     LogConfig
	Server  ServerConfig
	Service AccessConfig
}

type TraceConfig struct {
	Enable          bool    `env:"TRACE_ENABLE"`
	JaegerUrl       string  `env:"TRACE_JAEGER_URL"`
	SampleType      string  `env:"TRACE_SAMPLE_TYPE"`
	SampleParam     float64 `env:"TRACE_SAMPLE_PARAM"`
	SampleServerURL string  `env:"TRACE_SAMPLE_SERVER_URL"`
}

type ServerConfig struct {
	Port  int `env:"SERVER_PORT" envDefault:"8080"`
	Debug int `env:"SERVER_DEBUG"`
}

type LogConfig struct {
	Level string `env:"LOG_LEVEL" envDefault:"info"`
	Path  string `env:"LOG_PATH" envDefault:"stderr"`
}

type AccessConfig struct {
	ServiceName     string `env:"SERVICE_NAME"`
	SystemNamespace string `env:"SYSTEM_NAMESPACE"`
	WebhookAddress  string `env:"WEBHOOK_ADDRESS"`
}

func NewConfig() *Config {
	config := &Config{}
	if err := env.Parse(config); err != nil {
		panic(fmt.Sprintf("parse config error: %s", err.Error()))
	}

	return config
}
