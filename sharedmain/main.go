package sharedmain

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/injection/sharedmain"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/logging/logkey"
	"knative.dev/pkg/profiling"
	"knative.dev/pkg/system"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	"yunion.io/x/pkg/util/wait"

	"github.com/go-logr/zapr"
	kclient "github.com/katanomi/pkg/client"
	kmanager "github.com/katanomi/pkg/manager"
	kscheme "github.com/katanomi/pkg/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Main entrypoint for controllers basic main program
// it will:
// 1. load a config file given by the flag as the manager configuration
// 2. setup signal handlers, client configuration, and will call MainWithConfig to bootstrap a controller-manager
func Main(component string, scheme *runtime.Scheme, ctors ...Controller) {

	var configFile string
	flag.StringVar(&configFile, "config", "",
		"The controller will load its initial configuration from this file. "+
			"Omit this flag to use the default configuration values. "+
			"Command-line flags override configuration from this file.")

	flag.Parse()

	var err error
	options := ctrl.Options{Scheme: scheme}
	if configFile != "" {
		options, err = options.AndFrom(ctrl.ConfigFile().AtPath(configFile))
		if err != nil {
			fmt.Println(err, "unable to load the config file")
			os.Exit(1)
		}
	}

	ctx := ctrl.SetupSignalHandler()

	var config *rest.Config
	ctx, config = GetConfigOrDie(ctx)

	MainWithConfig(ctx, component, config, options, ctors...)
}

func GetConfigOrDie(ctx context.Context) (context.Context, *rest.Config) {
	cfg := injection.GetConfig(ctx)
	if cfg == nil {
		cfg = ctrl.GetConfigOrDie()
		ctx = injection.WithConfig(ctx, cfg)
	}
	return ctx, cfg
}

// MainWithConfig runs the generic main flow for controllers
// with the given config.
// TODO: needs to add support to webhooks and custom configuration
func MainWithConfig(ctx context.Context, component string, cfg *rest.Config, opts ctrl.Options, ctors ...Controller) {

	ctx, startInformers := injection.EnableInjectionOrDie(ctx, cfg)

	logger, atomicLevel := SetupLoggerOrDie(ctx, component)
	defer flush(logger)
	ctx = logging.WithLogger(ctx, logger)

	// this logger will not respect the automatic level update feature
	// and should not be used
	zaplogger := zapr.NewLogger(logger.Desugar())
	ctrl.SetLogger(zaplogger)
	ctx = ctrllog.IntoContext(ctx, zaplogger)

	mgr, err := ctrl.NewManager(cfg, opts)
	if err != nil {
		fmt.Println(err, "unable to start manager")
		os.Exit(1)
	}
	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		logger.Panic(err, "unable to set up health check")
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		logger.Panic(err, "unable to set up ready check")
	}

	ctx = kmanager.WithManager(ctx, mgr)
	ctx = kclient.WithClient(ctx, mgr.GetClient())

	// copy from main, lets test
	cmw := sharedmain.SetupConfigMapWatchOrDie(ctx, logger)

	profilingHandler := profiling.NewHandler(logger, false)
	profilingServer := profiling.NewServer(profilingHandler)

	sharedmain.WatchObservabilityConfigOrDie(ctx, cmw, profilingHandler, logger, component)

	// TODO: add a logging config observer that:
	// watches the logging configmap configuration while managing multiple atomicLevels
	// uses the controller constructor below to provide a specific logger for each
	// with the specified atomicLevel
	sharedmain.WatchLoggingConfigOrDie(ctx, cmw, logger, atomicLevel, component)
	// call constructors
	for _, controller := range ctors {
		name := controller.Name()
		// add here the logic to use atomicLevels
		controllerLogger := logger.Desugar().Named(name).Sugar()
		controller.Setup(ctx, mgr, controllerLogger)
	}

	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(profilingServer.ListenAndServe)

	logger.Info("Starting configuration manager...")
	if err := cmw.Start(ctx.Done()); err != nil {
		logger.Fatalw("Failed to start configuration manager", zap.Error(err))
	}

	// start informers for config loader
	startInformers()

	// This will block until either a signal arrives or one of the grouped functions
	// returns an error.
	if err := mgr.Start(egCtx); err != nil {
		logger.Errorw("problem running manager", "err", err)
	}

	profilingServer.Shutdown(context.Background())
	// Don't forward ErrServerClosed as that indicates we're already shutting down.
	if err := eg.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Errorw("Error while running server", zap.Error(err))
	}
}

// SetupLoggerOrDie sets up the logger using the config from the given context
// and returns a logger and atomic level, or dies by calling log.Fatalf.
func SetupLoggerOrDie(ctx context.Context, component string) (*zap.SugaredLogger, zap.AtomicLevel) {
	loggingConfig, err := GetLoggingConfig(ctx)
	if err != nil {
		log.Fatal("Error reading/parsing logging configuration: ", err)
	}
	l, level := logging.NewLoggerFromConfig(loggingConfig, component)

	// If PodName is injected into the env vars, set it on the logger.
	// This is needed for HA components to distinguish logs from different
	// pods.
	if pn := os.Getenv("POD_NAME"); pn != "" {
		l = l.With(zap.String(logkey.Pod, pn))
	}

	return l, level
}

func GetLoggingConfig(ctx context.Context) (*logging.Config, error) {
	loggingConfigMap := &corev1.ConfigMap{}

	key := types.NamespacedName{Name: logging.ConfigMapName(), Namespace: system.Namespace()}

	directClt, err := client.New(injection.GetConfig(ctx), client.Options{Scheme: kscheme.Scheme(ctx)})
	if err != nil {
		return nil, err
	}

	// These timeout and retry interval are set by heuristics.
	// e.g. istio sidecar needs a few seconds to configure the pod network.
	var lastErr error
	if err := wait.PollImmediate(1*time.Second, 5*time.Second, func() (bool, error) {
		lastErr := directClt.Get(ctx, key, loggingConfigMap)
		fmt.Println("err?", lastErr, "key", key)
		return lastErr == nil || apierrors.IsNotFound(lastErr), nil
	}); err != nil {
		return nil, fmt.Errorf("timed out waiting for the condition: %w", lastErr)

	}

	if loggingConfigMap == nil {
		return logging.NewConfigFromMap(nil)
	}
	return logging.NewConfigFromConfigMap(loggingConfigMap)
}

func flush(logger *zap.SugaredLogger) {
	logger.Sync()
}
