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
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	cm "knative.dev/pkg/configmap"
)

const (
	requestIDHeaderKey = "X-Request-ID"

	configMapNameEnv     = "CONFIG_TRACING_NAME"
	defaultServiceName   = "katanomi"
	defaultConfigMapName = "katanomi-config-tracing"

	enableKey        = "enable"
	backendKey       = "backend"
	samplingRatioKey = "sampling-ratio"
	jaegerConfigKey  = "jaeger-config"
	zipkinConfigKey  = "zipkin-config"
	customConfigKey  = "custom-config"
)

// ExporterBackend Built-in supported export types
type ExporterBackend string

var (
	ExporterBackendJaeger ExporterBackend = "jaeger"
	ExporterBackendZipkin ExporterBackend = "zipkin"
	ExporterBackendCustom ExporterBackend = "custom"
)

// Config tracing config
type Config struct {
	// Enable Controls whether to enable tracing.
	// default false.
	Enable bool `json:"enable" yaml:"enable"`

	// SamplingRatio Control the rate of sampling.
	// SamplingRatio >= 1 will always sample.
	// SamplingRatio <= 0 will never sample.
	// default 0.
	SamplingRatio float64 `json:"sampling_ratio" yaml:"samplingRatio"`

	// Backend The type of exporter backend
	Backend ExporterBackend `json:"backend" yaml:"backend"`

	// Jaeger The configuration used by jaeger backend
	Jaeger JaegerConfig `json:"jaeger" yaml:"jaeger"`

	// Zipkin The configuration used by zipkin backend
	Zipkin ZipkinConfig `json:"zipkin" yaml:"zipkin"`

	// Custom The configuration used by custom backend
	Custom string `json:"custom" yaml:"custom"`
	// todo support OTLP
}

// ZipkinConfig The configuration used by zipkin backend
type ZipkinConfig struct {
	// Url The collector url of zipkin backend
	Url string `json:"url" yaml:"url"`
}

// JaegerConfig The configuration used by Jaeger backend
type JaegerConfig struct {
	// Host The host of jaeger backend.
	Host string `json:"host" yaml:"host"`

	// Port the port of jaeger backend.
	Port string `json:"port" yaml:"port"`

	// MaxPacketSize The maximum UDP packet size for transport to the Jaeger agent.
	MaxPacketSize int `json:"max_packet_size" yaml:"maxPacketSize"`

	// DisableAttemptReconnecting Disable reconnecting udp client.
	DisableAttemptReconnecting bool `json:"disable_attempt_reconnecting" yaml:"disableAttemptReconnecting"`

	// AttemptReconnectInterval The interval between attempts to re resolve agent endpoint.
	AttemptReconnectInterval time.Duration `json:"attempt_reconnect_interval" yaml:"attemptReconnectInterval"`
}

// newTracingConfigFromConfigMap returns a Config for the given configmap
func newTracingConfigFromConfigMap(config *corev1.ConfigMap) (*Config, error) {
	if config == nil {
		return nil, nil
	}
	c := &Config{}
	backend := ""
	err := cm.Parse(config.Data,
		cm.AsBool(enableKey, &c.Enable),
		cm.AsString(backendKey, &backend),
		cm.AsString(customConfigKey, &c.Custom),
		cm.AsFloat64(samplingRatioKey, &c.SamplingRatio),
	)
	if err != nil {
		return nil, err
	}
	c.Backend = ExporterBackend(backend)
	if !c.Enable || c.SamplingRatio <= 0 {
		return c, nil
	}
	switch c.Backend {
	case ExporterBackendJaeger:
		if s, ok := config.Data[jaegerConfigKey]; ok && s != "" {
			if err := json.Unmarshal([]byte(s), &c.Jaeger); err != nil {
				return nil, err
			}
		}
	case ExporterBackendZipkin:
		if s, ok := config.Data[zipkinConfigKey]; ok && s != "" {
			if err := json.Unmarshal([]byte(s), &c.Zipkin); err != nil {
				return nil, err
			}
		}
	}
	return c, nil
}

// ConfigMapName gets the name of the tracing ConfigMap
func ConfigMapName() string {
	if name := os.Getenv(configMapNameEnv); name != "" {
		return name
	}
	return defaultConfigMapName
}
