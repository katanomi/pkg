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

package sharedmain

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	cloudeventsv2client "github.com/cloudevents/sdk-go/v2/client"

	"github.com/katanomi/pkg/config"
	storageroute "github.com/katanomi/pkg/plugin/storage/route"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
	_ "sigs.k8s.io/controller-runtime/pkg/metrics"

	"k8s.io/client-go/dynamic"

	"github.com/emicklei/go-restful/v3"
	"github.com/go-logr/zapr"
	"github.com/go-resty/resty/v2"
	kclient "github.com/katanomi/pkg/client"
	"github.com/katanomi/pkg/controllers"
	klogging "github.com/katanomi/pkg/logging"
	kmanager "github.com/katanomi/pkg/manager"
	"github.com/katanomi/pkg/multicluster"
	"github.com/katanomi/pkg/plugin/client"
	"github.com/katanomi/pkg/plugin/route"
	"github.com/katanomi/pkg/restclient"
	kscheme "github.com/katanomi/pkg/scheme"
	"github.com/katanomi/pkg/tracing"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/injection/sharedmain"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/system"

	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	ctrlcluster "sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	healthzRoutePath = "healthz"
	readyzRoutePath  = "readyz"
)

var (
	DefaultTimeout = kclient.DefaultTimeout
	DefaultQPS     = kclient.DefaultQPS
	DefaultBurst   = kclient.DefaultBurst
	Timeout        time.Duration
	Burst          int
	QPS            float64
	ConfigFile     string

	WebServerPort      int
	InsecureSkipVerify bool
)

// AppBuilder builds an app using multiple configuration options
type AppBuilder struct {
	// Basic options
	Name    string
	Context context.Context
	Config  *rest.Config
	sync.Once

	startInformers func()

	scheme *runtime.Scheme

	// Log related options
	Logger       *zap.SugaredLogger
	ZapConfig    *zap.Config
	LevelManager *klogging.LevelManager

	// Controllers
	Manager ctrl.Manager

	// ConfigWatch
	ConfigMapWatcher DefaultingWatcherWithOnChange

	// Profiling
	ProfilingServer *http.Server

	// client manager
	ClientManager *kclient.Manager

	// plugins
	plugins []client.Interface

	// storage plugins
	storagePlugins []client.Interface

	// restful container
	container *restful.Container
	filters   []restful.FilterFunction

	startFunc []func(context.Context) error

	initClientOnce sync.Once

	// newResourceLock override the default resourcelock of controller-runtime.
	newResourceLock kmanager.ResourceLockFunc
}

// ParseFlag parse flag needed for App
func ParseFlag() {
	flag.DurationVar(&Timeout, "kube-api-timeout", DefaultTimeout,
		"The maximum length of time to wait before giving up on a server request."+
			"A value of zero means no timeout. DefaultTimeOut: 10s")
	flag.Float64Var(&QPS, "kube-api-qps", float64(DefaultQPS),
		"qps indicates the maximum QPS to the master from this client."+
			"If it's zero, the created RESTClient will use DefaultQPS: 50")
	flag.IntVar(&Burst, "kube-api-burst", DefaultBurst,
		"Maximum burst for throttle."+
			"If it's zero, the created RESTClient will use DefaultBurst: 60.")
	flag.StringVar(&ConfigFile, "config", "",
		"The controller will load its initial configuration from this file. "+
			"Omit this flag to use the default configuration values. "+
			"Command-line flags override configuration from this file.")
	flag.BoolVar(&InsecureSkipVerify, "insecure-skip-tls-verify", false,
		"skip TLS verification and disable cert checking (default: false)")
	flag.IntVar(&WebServerPort, "web-server-port", 8100, "http web server port")
	flag.Parse()
}

// App main constructor entrypoint for AppBuilder
func App(name string) *AppBuilder {
	return &AppBuilder{Name: name, startFunc: []func(context.Context) error{}}
}

func (a *AppBuilder) init() {
	a.Once.Do(func() {
		ParseFlag()
		a.Context = ctrl.SetupSignalHandler()
		a.Context, a.Config = GetConfigOrDie(a.Context)
		a.Config.Timeout = Timeout

		if a.Config.QPS < float32(QPS) {
			a.Config.QPS = float32(QPS)
		}
		if a.Config.Burst < Burst {
			a.Config.Burst = Burst
		}
		a.Context, a.startInformers = injection.EnableInjectionOrDie(a.Context, a.Config)
		a.Context = kclient.WithAppConfig(a.Context, a.Config)

		restyClient := resty.NewWithClient(kclient.NewHTTPClient())
		restyClient.SetDisableWarn(true)
		restyClient.SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: InsecureSkipVerify, // nolint: gosec // G402: TLS InsecureSkipVerify set true.
		})
		tracing.WrapTransportForRestyClient(restyClient)
		a.Context = restclient.WithRESTClient(a.Context, restyClient)

		a.ConfigMapWatcher = sharedmain.SetupConfigMapWatchOrDie(a.Context, a.Logger)
		a.startFunc = append(a.startFunc, func(ctx context.Context) error {
			return a.ConfigMapWatcher.Start(ctx.Done())
		})

		a.Context = multicluster.WithMultiCluster(a.Context, multicluster.NewClusterRegistryClientOrDie(a.Config))

		a.container = restful.NewContainer()
		a.container.Router(restful.RouterJSR311{})
		a.Context, a.ClientManager = GetClientManager(a.Context)
		a.filters = []restful.FilterFunction{}
	})
}

// Scheme provides a scheme to the app
func (a *AppBuilder) Scheme(scheme *runtime.Scheme) *AppBuilder {
	a.init()
	a.scheme = scheme
	a.Context = kscheme.WithScheme(a.Context, scheme)
	return a
}

func (a *AppBuilder) initClient(clientVar ctrlclient.Client) {
	a.initClientOnce.Do(func() {
		if clientVar == nil {
			cluster, err := ctrlcluster.New(a.Config, WithScheme(a.scheme))
			if err != nil {
				a.Logger.Fatalw("cluster client setup error", "err", err)
			}
			clientVar = cluster.GetClient()
			a.startFunc = append(a.startFunc, func(ctx context.Context) error {
				err := cluster.Start(ctx)
				if err != nil {
					a.Logger.Errorw("cluster start error", "err", err)
				}
				return err
			})
			a.Context = kclient.WithCluster(a.Context, cluster)
		}
		a.Context = kclient.WithClient(a.Context, clientVar)

		directClient, err := ctrlclient.New(a.Config, ctrlclient.Options{Scheme: a.scheme})
		if err != nil {
			a.Logger.Fatalw("direct client setup error", "err", err)
		}
		a.Context = kclient.WithDirectClient(a.Context, directClient)

		dynamicClient, err := dynamic.NewForConfig(a.Config)
		if err != nil {
			a.Logger.Fatalw("dynamic client setup error", "err", err)
		}
		a.Context = kclient.WithDynamicClient(a.Context, dynamicClient)

		ceClient, err := newCloudEventsClient(kclient.GetDefaultTransport())
		if err != nil {
			a.Logger.Fatalw("cloudevents client setup error", "err", err)
		}
		a.Context = kclient.WithCEClient(a.Context, ceClient)
	})
}

func newCloudEventsClient(roundTripper http.RoundTripper) (cloudevents.Client, error) {
	p, err := cloudevents.NewHTTP(cloudevents.WithRoundTripper(roundTripper))
	if err != nil {
		return nil, err
	}

	cloudEventClient, err := cloudevents.NewClient(p, cloudevents.WithUUIDs(), cloudevents.WithTimeNow(), cloudeventsv2client.WithForceStructured())
	if err != nil {
		return nil, err
	}

	return cloudEventClient, nil
}

// Tracing adds tracing and tracer to the app
func (a *AppBuilder) Tracing(ops ...tracing.TraceOption) *AppBuilder {
	a.init()
	ops = append([]tracing.TraceOption{
		tracing.WithServiceName(a.Name),
	}, ops...)
	t := tracing.NewTracing(a.Logger, ops...)
	err := tracing.SetupDynamicPublishing(t, a.ConfigMapWatcher)
	if err != nil {
		log.Fatal("Error reading/parsing tracing configuration: ", err)
	}
	a.filters = append(a.filters, tracing.RestfulFilter(a.Name, healthzRoutePath, readyzRoutePath))
	return a
}

// ConfigManager add katanomi manager to app context
func (a *AppBuilder) ConfigManager() *AppBuilder {
	a.init()

	name := config.Name()
	configMGR := config.NewManager(a.ConfigMapWatcher, a.Logger, name)
	a.Context = config.WithKatanomiConfigManager(a.Context, configMGR)

	return a
}

// Log adds logging and logger to the app
func (a *AppBuilder) Log() *AppBuilder {
	a.init()

	lvlMGR := klogging.NewLevelManager(a.Context, a.Name)
	a.LevelManager = &lvlMGR
	loggingConfig, err := GetLoggingConfig(a.Context)
	if err != nil {
		log.Fatal("Error reading/parsing logging configuration: ", err)
	}
	a.ZapConfig, err = klogging.ZapConfigFromJSON(loggingConfig.LoggingConfig)
	if err != nil {
		log.Fatal("Error parsing logging zapConfig: ", err)
	}

	a.Logger, _ = SetupLoggerOrDie(a.Context, a.Name)

	// support dynamic adjustment the log levels
	a.Logger = a.Logger.Desugar().WithOptions(zap.UpdateCore(a.LevelManager.Get(a.Name), *a.ZapConfig)).Sugar()
	a.Context = logging.WithLogger(a.Context, a.Logger)
	lvlMGR.SetLogger(a.Logger)

	// this logger will not respect the automatic level update feature
	// and should not be used
	// its main purpose is to provide a logger to controller-runtime
	zaplogger := zapr.NewLogger(a.Logger.Desugar())
	ctrl.SetLogger(zaplogger)
	a.Context = ctrllog.IntoContext(a.Context, zaplogger)

	// watches logging config changes and reset the level manager
	WatchLoggingConfigOrDie(a.Context, a.ConfigMapWatcher, a.Logger, a.LevelManager)

	// call constructors
	systemConfigMap, err := kubeclient.Get(a.Context).CoreV1().ConfigMaps(system.Namespace()).Get(a.Context, logging.ConfigMapName(),
		metav1.GetOptions{})
	if err != nil {
		a.Logger.Fatalw("read logging configmap error", "err", err)
	}
	a.ConfigMapWatcher.OnChange(systemConfigMap)

	// adds filter for logger
	a.filters = append(a.filters, klogging.Filter(a.Logger))

	return a
}

// RESTClient injects a RESTClient
func (a *AppBuilder) RESTClient(client *resty.Client) *AppBuilder {
	a.Context = restclient.WithRESTClient(a.Context, client)
	return a
}

// MultiClusterClient injects a multi cluster client into the context
func (a *AppBuilder) MultiClusterClient(client multicluster.Interface) *AppBuilder {
	a.Context = multicluster.WithMultiCluster(a.Context, client)
	return a
}

// Controllers adds controllers to the app, will start a manager under the hood
func (a *AppBuilder) Controllers(ctors ...controllers.SetupChecker) *AppBuilder {
	a.init()

	var err error
	options := ctrl.Options{Scheme: a.scheme}
	if ConfigFile != "" {
		options, err = options.AndFrom(ctrl.ConfigFile().AtPath(ConfigFile))
		if err != nil {
			a.Logger.Fatalw("unable to load the config file", "err", err)
		}
	}
	options.Scheme = a.scheme
	// If the value is empty, the default behavior will still be used.
	// Ref: https://github.com/kubernetes-sigs/controller-runtime/blob/b9219528d95974cb4f5b06f86c9b1c9b7d3045a5/pkg/manager/manager.go#L551
	// make the struct as an interface to check if it implements the other interface
	optionsInterface := interface{}(&options)
	if setter, ok := optionsInterface.(kmanager.ResourceLockSetter); ok {
		setter.SetNewResourceLock(a.newResourceLock)
		a.Logger.Infow("inject resource lock")
	}

	// BaseContext provides Context values to Runnables
	options.BaseContext = func() context.Context {
		return a.Context
	}
	a.Manager, err = ctrl.NewManager(a.Config, options)
	if err != nil {
		a.Logger.Fatalw("unable to start manager", "err", err)
	}
	if err := a.Manager.AddHealthzCheck(healthzRoutePath, healthz.Ping); err != nil {
		a.Logger.Fatalw("unable to set up health check", "err", err)
	}
	if err := a.Manager.AddReadyzCheck(readyzRoutePath, healthz.Ping); err != nil {
		a.Logger.Fatalw("unable to set up ready check", "err", err)
	}

	a.Context = kmanager.WithManager(a.Context, a.Manager)
	// a manager implements all cluster.Cluster methods
	a.Context = kclient.WithCluster(a.Context, a.Manager)
	a.initClient(a.Manager.GetClient())

	// TODO: make this interval setup configurable
	lazyLoader := controllers.NewLazyLoader(a.Context, time.Minute)

	for i := range ctors {
		controller := ctors[i]
		name := controller.Name()
		controllerAtomicLevel := a.LevelManager.Get(name)
		controllerLogger := a.Logger.Desugar().WithOptions(zap.UpdateCore(controllerAtomicLevel, *a.ZapConfig)).Named(name).Sugar()

		if err := lazyLoader.LazyLoad(a.Context, a.Manager, controllerLogger, controller); err != nil {
			a.Logger.Fatalw("controller setup error", "ctrl", name, "err", err)
		}
	}

	a.startFunc = append(a.startFunc, func(ctx context.Context) error {
		return a.Manager.Start(ctx)
	})

	a.startFunc = append(a.startFunc, func(ctx context.Context) error {
		return lazyLoader.Start(ctx)
	})

	return a
}

// Webhooks adds webhook setup for objects in app
func (a *AppBuilder) Webhooks(objs ...runtime.Object) *AppBuilder {
	a.init()
	for _, obj := range objs {
		if obj == nil {
			continue
		}
		var err error
		if setup, ok := obj.(WebhookSetup); ok {
			err = setup.SetupWebhookWithManager(a.Manager)
		} else {
			err = ctrl.NewWebhookManagedBy(a.Manager).
				For(obj).
				Complete()
		}
		if err != nil {
			a.Logger.Fatalw("webhook setup error for obj", "obj", obj.GetObjectKind(), "err", err)
		}

		// if obj implements the register interface, call it
		if setup, ok := obj.(WebhookRegisterSetup); ok {
			name := setup.GetLoggerName()
			controllerAtomicLevel := a.LevelManager.Get(name)
			controllerLogger := a.Logger.Desugar().WithOptions(zap.UpdateCore(controllerAtomicLevel, *a.ZapConfig)).Named(name).Sugar()
			controllerLogger.Infow("setup register webhook")
			ctx := logging.WithLogger(a.Context, controllerLogger)
			setup.SetupRegisterWithManager(ctx, a.Manager)
		}
	}
	return a
}

// Filters customize filters to this app
func (a *AppBuilder) Filters(filters ...restful.FilterFunction) *AppBuilder {
	a.filters = append(a.filters, filters...)
	return a
}

// Container adds a containers
func (a *AppBuilder) Container(container *restful.Container) *AppBuilder {
	a.container = container
	return a
}

func (a *AppBuilder) Webservices(webServices ...WebService) *AppBuilder {
	a.init()

	// will init a client if not already initiated
	a.initClient(nil)

	for _, item := range webServices {
		name := item.Name()
		webserviceAtomicLevel := a.LevelManager.Get(name)
		webserviceLogger := a.Logger.Desugar().WithOptions(zap.UpdateCore(webserviceAtomicLevel, *a.ZapConfig)).Named(name).Sugar()

		err := item.Setup(a.Context, func(ws *restful.WebService) {
			a.container.Add(ws)
		}, webserviceLogger)
		if err != nil {
			a.Logger.Fatalw("webservice setup error", "weservice", name, "err", err)
		}
	}

	return a
}

// Plugins adds plugins to this app
func (a *AppBuilder) Plugins(plugins ...client.Interface) *AppBuilder {
	a.init()

	// will init a client if not already initiated
	a.initClient(nil)
	a.plugins = plugins
	for _, plugin := range a.plugins {
		if err := plugin.Setup(a.Context, a.Logger); err != nil {
			a.Logger.Fatalw("plugin could not be setup correctly", "err", err, "plugin", plugin.Path())
		}
		// MetaFilter and AuthFilter are dedicated to plugin api,
		// so register the filters when the service is initialized.
		ws, err := route.NewService(plugin, client.MetaFilter, client.AuthFilter)
		if err != nil {
			a.Logger.Fatalw("plugin could not start correctly", "err", err, "plugin", plugin.Path())
		}
		a.container.Add(ws)
	}
	return a
}

// StoragePlugins adds storage plugins to this app
func (a *AppBuilder) StoragePlugins(plugins ...client.Interface) *AppBuilder {
	a.init()

	// will init a client if not already initiated
	a.initClient(nil)
	a.storagePlugins = plugins
	for _, plugin := range a.storagePlugins {
		if err := plugin.Setup(a.Context, a.Logger); err != nil {
			a.Logger.Fatalw("plugin could not be setup correctly", "err", err, "plugin", plugin.Path())
		}
		filters, _ := a.ClientManager.Filters(a.Context)
		wss, err := storageroute.NewServicesWithContext(a.Context, plugin, filters...)
		if err != nil {
			a.Logger.Fatalw("plugin could not start correctly", "err", err, "plugin", plugin.Path())
		}

		for _, ws := range wss {
			a.container.Add(ws)
		}

	}
	return a
}

// PluginAttributes set plugin attributes from yaml file
func (a *AppBuilder) PluginAttributes(plugin client.PluginAttributes, file string) *AppBuilder {
	attributes := make(map[string][]string)
	a.readAttributes(file, &attributes)
	for key, value := range attributes {
		plugin.SetAttribute(key, value...)
	}

	return a
}

// PluginVersionAttributes set plugin attributes from yaml file
func (a *AppBuilder) PluginVersionAttributes(plugin client.PluginVersionAttributes, file string) *AppBuilder {
	attributes := make(map[string]map[string][]string)
	a.readAttributes(file, &attributes)
	for key, value := range attributes {
		plugin.SetVersionAttributes(key, value)
	}

	return a
}

func (a *AppBuilder) readAttributes(file string, attributes interface{}) {
	data, err := os.ReadFile(file)
	if err != nil {
		a.Logger.Fatalw("read plugin version attributes file error", "err", err, "file", file)
	}
	reader := bytes.NewReader(data)
	err = utilyaml.NewYAMLOrJSONDecoder(reader, len(data)).Decode(attributes)
	if err != nil {
		a.Logger.Fatalw("parse plugin version attributes file error", "err", err, "file", file)
	}
}

// APIDocs adds api docs to the server
func (a *AppBuilder) APIDocs() *AppBuilder {
	// NO-OP for compatibility, this function is now a standard once there are webservices added
	return a
}

// Profiling enables profiling http server
func (a *AppBuilder) Profiling() *AppBuilder {
	a.init()

	// there is an issue when deploying together with
	// other controllers like knative or tekton
	// just ignore these for now and find a better way to do this
	// TODO: find a better way to do profiling
	// there is already pprof for webservices, maybe just need to solve
	// how do controllers do metrics?

	// profilingHandler := profiling.NewHandler(a.Logger, false)
	// a.ProfilingServer = profiling.NewServer(profilingHandler)
	// sharedmain.WatchObservabilityConfigOrDie(a.Context, a.ConfigMapWatcher, profilingHandler, a.Logger, a.Name)
	// a.startFunc = append(a.startFunc, func(ctx context.Context) error {
	// 	return a.ProfilingServer.ListenAndServe()
	// })

	return a
}

// NewResourceLock set a new resource lock
// Used to change the default behavior in controller-runtime
func (a *AppBuilder) NewResourceLock(newResourceLock kmanager.ResourceLockFunc) *AppBuilder {
	a.newResourceLock = newResourceLock
	return a
}

// Run starts all
func (a *AppBuilder) Run(startFuncs ...func(context.Context) error) error {
	defer func() {
		if a.Logger != nil {
			a.Logger.Sync()
		}
	}()

	// adds a http server if there are any endpoints registered
	if a.container != nil {
		// adds profiling and health checks
		a.container.Add(route.NewDefaultService(a.Context))

		if len(a.container.RegisteredWebServices()) > 0 {
			a.container.Add(route.NewDocService(a.container.RegisteredWebServices()...))
		}

		a.startFunc = append(a.startFunc, func(ctx context.Context) error {
			// TODO: find a better way to get this configuration
			for _, filter := range a.filters {
				a.container.Filter(filter)
			}

			srv := &http.Server{
				Addr:    fmt.Sprintf(":%d", WebServerPort),
				Handler: a.container,
			}
			return srv.ListenAndServe()
		})
	}

	a.startInformers()

	eg, egCtx := errgroup.WithContext(a.Context)
	startFuncs = append(startFuncs, a.startFunc...)
	for i, st := range startFuncs {
		var index = i
		startFunc := st
		eg.Go(func() error {
			err := startFunc(egCtx)
			// TODO: Fatal here, because not all startFunc have completed cancel mechanism by ctx, this will cause eg.Wait() handled and never stopped even some errors happened.
			// it is not a good solution, we should change it soon.
			if err != nil {
				a.Logger.Fatalw("error to start func", "index", index, "err", err)
			}
			return nil
		})
	}

	// wait until all are done
	err := eg.Wait()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.Logger.Errorw("Error while running server", zap.Error(err))
	}

	return nil
}

// WebhookSetup method to inject and setup webhooks using the object
type WebhookSetup interface {
	SetupWebhookWithManager(mgr ctrl.Manager) error
}

// WebhookRegisterSetup method to inject and setup webhook register using the object
type WebhookRegisterSetup interface {
	GetLoggerName() string
	SetupRegisterWithManager(context.Context, ctrl.Manager)
}
