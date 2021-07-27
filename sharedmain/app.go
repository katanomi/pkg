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
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"sync"

	"github.com/katanomi/pkg/tracing"

	"github.com/go-logr/zapr"
	kclient "github.com/katanomi/pkg/client"
	klogging "github.com/katanomi/pkg/logging"
	kmanager "github.com/katanomi/pkg/manager"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	cminformer "knative.dev/pkg/configmap/informer"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/injection/sharedmain"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/profiling"
	"knative.dev/pkg/system"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
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

	//Controllers
	Manager ctrl.Manager

	// ConfigWatch
	ConfigMapWatcher *cminformer.InformedWatcher

	// Profiling
	ProfilingServer *http.Server

	//tracing
	tracingManager *tracing.Manager

	startFunc []func(context.Context) error
}

// App main constructor entrypoint for AppBuilder
func App(name string) *AppBuilder {
	return &AppBuilder{Name: name, startFunc: []func(context.Context) error{}}
}

func (a *AppBuilder) init() {
	a.Once.Do(func() {
		a.Context = ctrl.SetupSignalHandler()
		a.Context, a.Config = GetConfigOrDie(a.Context)
		a.Context, a.startInformers = injection.EnableInjectionOrDie(a.Context, a.Config)
		a.ConfigMapWatcher = sharedmain.SetupConfigMapWatchOrDie(a.Context, a.Logger)
		a.startFunc = append(a.startFunc, func(ctx context.Context) error {
			return a.ConfigMapWatcher.Start(ctx.Done())
		})
	})
}

// Scheme provides a scheme to the app
func (a *AppBuilder) Scheme(scheme *runtime.Scheme) *AppBuilder {
	a.init()
	a.scheme = scheme
	return a
}

// Log adds logging and logger to the app
func (a *AppBuilder) Log() *AppBuilder {
	a.init()

	lvlMGR := klogging.NewLevelManager()
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
	a.Context = logging.WithLogger(a.Context, a.Logger)

	// watches logging config changes and reset the level manager
	WatchLoggingConfigOrDie(a.Context, a.ConfigMapWatcher, a.Logger, a.LevelManager)

	// call constructors
	systemConfigMap, err := kubeclient.Get(a.Context).CoreV1().ConfigMaps(system.Namespace()).Get(a.Context, logging.ConfigMapName(),
		metav1.GetOptions{})
	if err != nil {
		a.Logger.Fatalw("read logging configmap error", "err", err)
	}
	a.ConfigMapWatcher.OnChange(systemConfigMap)

	return a
}

// Controllers adds controllers to the app, will start a manager under the hood
func (a *AppBuilder) Controllers(ctors ...Controller) *AppBuilder {
	a.init()

	var configFile string
	flag.StringVar(&configFile, "config", "",
		"The controller will load its initial configuration from this file. "+
			"Omit this flag to use the default configuration values. "+
			"Command-line flags override configuration from this file.")

	flag.Parse()

	var err error
	options := ctrl.Options{Scheme: a.scheme}
	if configFile != "" {
		options, err = options.AndFrom(ctrl.ConfigFile().AtPath(configFile))
		if err != nil {
			a.Logger.Fatalw("unable to load the config file", "err", err)
		}
	}
	options.Scheme = a.scheme

	// this logger will not respect the automatic level update feature
	// and should not be used
	// its main purpose is to provide a logger to controller-runtime
	zaplogger := zapr.NewLogger(a.Logger.Desugar())
	ctrl.SetLogger(zaplogger)
	a.Context = ctrllog.IntoContext(a.Context, zaplogger)

	a.Manager, err = ctrl.NewManager(a.Config, options)
	if err != nil {
		a.Logger.Fatalw("unable to start manager", "err", err)
	}
	if err := a.Manager.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		a.Logger.Fatalw("unable to set up health check", "err", err)
	}
	if err := a.Manager.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		a.Logger.Fatalw("unable to set up ready check", "err", err)
	}

	a.Context = kmanager.WithManager(a.Context, a.Manager)
	a.Context = kclient.WithClient(a.Context, a.Manager.GetClient())

	for _, controller := range ctors {
		name := controller.Name()
		controllerAtomicLevel := a.LevelManager.Get(name)
		controllerLogger := a.Logger.Desugar().WithOptions(zap.UpdateCore(controllerAtomicLevel, *a.ZapConfig)).Named(name).Sugar()
		if err := controller.Setup(a.Context, a.Manager, controllerLogger, a.tracingManager); err != nil {
			a.Logger.Fatalw("controller setup error", "ctrl", name, "err", err)
		}
	}

	a.startFunc = append(a.startFunc, func(ctx context.Context) error {
		return a.Manager.Start(ctx)
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
	}
	return a
}

func (a *AppBuilder) Tracing() *AppBuilder {
	manager, err := tracing.SetupTracingOrDie(a.Context, a.Logger)
	if err != nil {
		a.Logger.Fatalw("set up tracing error", "err", err)

		return a
	}

	manager.WatchTracingConfigOrDie(a.Context, a.ConfigMapWatcher, a.Logger)

	a.tracingManager = manager

	return a
}

// Profiling enables profiling http server
func (a *AppBuilder) Profiling() *AppBuilder {
	a.init()

	profilingHandler := profiling.NewHandler(a.Logger, false)
	a.ProfilingServer = profiling.NewServer(profilingHandler)

	sharedmain.WatchObservabilityConfigOrDie(a.Context, a.ConfigMapWatcher, profilingHandler, a.Logger, a.Name)

	a.startFunc = append(a.startFunc, func(ctx context.Context) error {
		return a.ProfilingServer.ListenAndServe()
	})

	return a
}

// Run starts all
func (a *AppBuilder) Run() error {
	a.startInformers()

	eg, egCtx := errgroup.WithContext(a.Context)
	for _, st := range a.startFunc {
		startFunc := st
		eg.Go(func() error {
			return startFunc(egCtx)
		})
	}

	// waituntil all are done
	err := eg.Wait()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.Logger.Errorw("Error while running server", zap.Error(err))
	}

	if a.Logger != nil {
		a.Logger.Sync()
	}

	if a.tracingManager != nil {
		a.tracingManager.Sync()
	}

	return nil
}

// WebhookSetup method to inject and setup webhooks using the object
type WebhookSetup interface {
	SetupWebhookWithManager(mgr ctrl.Manager) error
}
