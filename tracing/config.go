package tracing

import (
	"encoding/json"
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	cm "knative.dev/pkg/configmap"
)

const (
	configMapNameEnv     = "CONFIG_TRACING_NAME"
	defaultServiceName   = "katanomi"
	defaultConfigMapName = "config-tracing"

	enableKey       = "enable"
	backendKey      = "backend"
	jaegerConfigKey = "jaeger-config"
	zipkinConfigKey = "zipkin-config"
	customConfigKey = "custom-config"
)

type ExporterBackend string

var (
	ExporterBackendJaeger ExporterBackend = "jaeger"
	ExporterBackendZipkin ExporterBackend = "zipkin"
	ExporterBackendCustom ExporterBackend = "custom"
)

// Config tracing config
type Config struct {
	Enable bool `json:"enable" yaml:"enable" env:"TRACE_ENABLE"`

	Backend ExporterBackend `json:"backend" yaml:"backend"`
	Jaeger  JaegerConfig    `json:"jaeger" yaml:"jaeger"`
	Zipkin  ZipkinConfig    `json:"zipkin" yaml:"zipkin"`

	Custom string `json:"custom" yaml:"custom"`
	// todo support OTLP
}

type ZipkinConfig struct {
	Url string `json:"url" yaml:"url"`
}

type JaegerConfig struct {
	Host                       string        `json:"host" yaml:"host"`
	Port                       string        `json:"port" yaml:"port"`
	MaxPacketSize              int           `json:"max_packet_size" yaml:"maxPacketSize"`
	DisableAttemptReconnecting bool          `json:"disable_attempt_reconnecting" yaml:"disableAttemptReconnecting"`
	AttemptReconnectInterval   time.Duration `json:"attempt_reconnect_interval" yaml:"attemptReconnectInterval"`
}

// NewTracingConfigFromConfigMap returns a Config for the given configmap
func NewTracingConfigFromConfigMap(config *corev1.ConfigMap) (*Config, error) {
	if config == nil {
		return &Config{}, nil
	}
	c := &Config{}
	backend := ""
	err := cm.Parse(config.Data,
		cm.AsBool(enableKey, &c.Enable),
		cm.AsString(backendKey, &backend),
		cm.AsString(customConfigKey, &c.Custom),
	)
	if err != nil {
		return nil, err
	}
	c.Backend = ExporterBackend(backend)
	if !c.Enable {
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
