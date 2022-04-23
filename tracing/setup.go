package tracing

import (
	"go.opentelemetry.io/otel/sdk/trace"
	"knative.dev/pkg/configmap"
)

type T interface {
	exporter(*Config) (trace.SpanExporter, error)
}

// SetupDynamicPublishing sets up trace publishing for the process, by watching a
// ConfigMap for the configuration. Note that other pieces still need to generate the traces, this
// just ensures that if generated, they are collected appropriately. This is normally done by using
// tracing.HTTPSpanMiddleware as a middleware HTTP handler. The configuration will be dynamically
// updated when the ConfigMap is updated.
func SetupDynamicPublishing(tracing *Tracing, configMapWatcher configmap.Watcher) error {
	tracerUpdater := func(name string, value interface{}) {
		if name != tracing.ConfigMapName {
			return
		}
		cfg := value.(*Config)
		tracing.ApplyConfig(cfg)
	}

	// Set up our config store.
	configStore := configmap.NewUntypedStore(
		"config-tracing-store",
		tracing.logger,
		configmap.Constructors{
			tracing.ConfigMapName: newTracingConfigFromConfigMap,
		},
		tracerUpdater)
	configStore.WatchConfigs(configMapWatcher)

	return nil
}
