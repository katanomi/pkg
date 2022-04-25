package tracing

import (
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"knative.dev/pkg/configmap"
	cminformer "knative.dev/pkg/configmap/informer"
)

type T interface {
	exporter(*Config) (trace.SpanExporter, error)
}

func defaultConfigMap(name string) *v1.ConfigMap {
	cm := &v1.ConfigMap{}
	cm.Name = name
	return cm
}

// SetupDynamicPublishing sets up trace publishing for the process, by watching a
// ConfigMap for the configuration. Note that other pieces still need to generate the traces, this
// just ensures that if generated, they are collected appropriately. This is normally done by using
// tracing.HTTPSpanMiddleware as a middleware HTTP handler. The configuration will be dynamically
// updated when the ConfigMap is updated.
func SetupDynamicPublishing(tracing *Tracing, configMapWatcher *cminformer.InformedWatcher) error {
	tracerUpdater := func(name string, value interface{}) {
		if name != tracing.ConfigMapName {
			return
		}
		cfg := value.(*Config)
		tracing.ApplyConfig(cfg)
	}

	// Set up our config store.
	w := NewDftConfigMapWatcher("config-tracing-store", tracing.logger, configMapWatcher)
	w.AddWatch(tracing.ConfigMapName, newTracingConfigFromConfigMap, defaultConfigMap(tracing.ConfigMapName))
	w.Run(tracerUpdater)

	return nil
}

func NewDftConfigMapWatcher(name string, logger *zap.SugaredLogger, informedWatcher *cminformer.InformedWatcher) *dftConfigMapWatcher {
	return &dftConfigMapWatcher{
		name:            name,
		InformedWatcher: informedWatcher,
		constructors:    map[string]interface{}{},
		defaults:        map[string]v1.ConfigMap{},
		logger:          logger,
	}
}

type dftConfigMapWatcher struct {
	*cminformer.InformedWatcher

	name         string
	constructors map[string]interface{}
	defaults     map[string]v1.ConfigMap

	logger *zap.SugaredLogger
}

func (d *dftConfigMapWatcher) AddWatch(name string, constructor interface{}, defaultCM *v1.ConfigMap) {
	d.constructors[name] = constructor
	if defaultCM != nil {
		d.defaults[name] = *defaultCM
	}
}

func (d *dftConfigMapWatcher) Run(onAfterStore ...func(name string, value interface{})) {
	configStore := configmap.NewUntypedStore(d.name, d.logger, d.constructors, onAfterStore...)
	configStore.WatchConfigs(d)
}

func (d *dftConfigMapWatcher) Watch(name string, obs ...configmap.Observer) {
	if cm, ok := d.defaults[name]; ok {
		d.InformedWatcher.WatchWithDefault(cm, obs...)
	} else {
		d.InformedWatcher.Watch(name, obs...)
	}
}
