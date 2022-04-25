package tracing

import (
	"os"
	"testing"

	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
)

func Test_newTracingConfigFromConfigMap(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]string
		want    *Config
		wantErr bool
	}{
		{
			input:   map[string]string{},
			want:    &Config{},
			wantErr: false,
		},
		{
			input: map[string]string{
				enableKey:        "true",
				backendKey:       "jaeger",
				samplingRatioKey: "0.5",
			},
			want: &Config{
				Enable:        true,
				Backend:       ExporterBackendJaeger,
				SamplingRatio: 0.5,
			},
			wantErr: false,
		},
		{
			input: map[string]string{
				enableKey:        "true",
				backendKey:       "zipkin",
				samplingRatioKey: "1.1",
			},
			want: &Config{
				Enable:        true,
				Backend:       ExporterBackendZipkin,
				SamplingRatio: 1.1,
			},
			wantErr: false,
		},
		{
			input: map[string]string{
				enableKey:        "true",
				backendKey:       "custom",
				samplingRatioKey: "1.1",
			},
			want: &Config{
				Enable:        true,
				Backend:       ExporterBackendCustom,
				SamplingRatio: 1.1,
			},
			wantErr: false,
		},
		{
			input: map[string]string{
				enableKey:        "true typo",
				backendKey:       "custom",
				samplingRatioKey: "0.5",
			},
			want:    nil,
			wantErr: true,
		},
	}
	g := NewGomegaWithT(t)
	for _, tt := range tests {
		got, err := newTracingConfigFromConfigMap(&v1.ConfigMap{
			Data: tt.input,
		})
		if tt.wantErr {
			g.Expect(err).ShouldNot(BeNil())
		} else {
			g.Expect(err).Should(BeNil())
		}
		g.Expect(got).Should(Equal(tt.want))
	}
}

func TestConfigMapName(t *testing.T) {
	tests := []struct {
		setting func()
		want    string
	}{
		{
			setting: func() {
				os.Setenv(configMapNameEnv, "test-config-tracing-name")
			},
			want: "test-config-tracing-name",
		},
		{
			setting: func() {
				os.Unsetenv(configMapNameEnv)
			},
			want: defaultConfigMapName,
		},
	}
	g := NewGomegaWithT(t)
	for _, tt := range tests {
		if tt.setting != nil {
			tt.setting()
		}
		got := ConfigMapName()
		g.Expect(got).Should(Equal(tt.want))
	}
}
