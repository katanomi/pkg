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

package client

import (
	"context"
	"time"

	"github.com/katanomi/pkg/tracing"
	"k8s.io/client-go/dynamic"

	"github.com/emicklei/go-restful/v3"
	kerrors "github.com/katanomi/pkg/errors"
	kscheme "github.com/katanomi/pkg/scheme"
	"go.uber.org/zap"
	"k8s.io/client-go/rest"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type GetBaseConfigFunc func() (*rest.Config, error)

// GetConfigFunc retrieves a configuration based on a request
type GetConfigFunc func(req *restful.Request, baseConfig GetBaseConfigFunc) (*rest.Config, error)

// Manager dynamically generates client based on user requests
type Manager struct {
	GetConfig      GetConfigFunc
	GetBasicConfig GetBaseConfigFunc
	*zap.SugaredLogger
}

// NewManager initializes a new manager based on func
func NewManager(ctx context.Context, get GetConfigFunc, baseConfig GetBaseConfigFunc) *Manager {
	if get == nil {
		get = FromBearerToken
	}
	configGetter := func(req *restful.Request, baseConfig GetBaseConfigFunc) (*rest.Config, error) {
		cfg, err := get(req, baseConfig)
		if cfg != nil {
			cfg.Wrap(tracing.WrapTransport)
		}
		return cfg, err
	}
	if baseConfig == nil {
		baseConfig = config.GetConfig
	}
	return &Manager{
		SugaredLogger:  logging.FromContext(ctx),
		GetConfig:      configGetter,
		GetBasicConfig: baseConfig,
	}
}

// Filter returns a filter to be used in
func (m *Manager) Filter(ctx context.Context) restful.FilterFunction {
	return ManagerFilter(ctx, m)
}

// ManagerFilter generates filter based on a manager to create a config based in a request and injects into context
func ManagerFilter(ctx context.Context, mgr *Manager) restful.FilterFunction {
	log := logging.FromContext(ctx).Named("manager-filter")
	scheme := kscheme.Scheme(ctx)
	serviceAccountClient := Client(ctx)

	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		start := time.Now()

		config, err := mgr.GetConfig(req, mgr.GetBasicConfig)
		log.Debugw("ManagerFilter, got config", "totalElapsed", time.Since(start).String())
		step := time.Now()
		if err != nil {
			log.Debugw("manager filter config get error", "err", err)
			err = kerrors.AsAPIError(err)
			resp.WriteHeaderAndJson(kerrors.AsStatusCode(err), err, restful.MIME_JSON)
			return
		}

		ctx := req.Request.Context()

		// setting defaults to the config
		config.Burst = DefaultBurst
		config.QPS = DefaultQPS
		config.Timeout = DefaultTimeout
		ctx = injection.WithConfig(ctx, config)

		directClient, err := client.New(config, client.Options{Scheme: scheme, Mapper: serviceAccountClient.RESTMapper()})
		log.Debugw("ManagerFilter, got direct client", "totalElapsed", time.Since(start).String(), "elapsed", time.Since(step).String())
		step = time.Now()
		if err != nil {
			log.Debugw("manager filter direct client create error", "err", err)
			err = kerrors.AsAPIError(err)
			resp.WriteHeaderAndJson(kerrors.AsStatusCode(err), err, restful.MIME_JSON)
			return
		}
		ctx = WithClient(ctx, directClient)

		dynamicClient, err := dynamic.NewForConfig(config)
		log.Debugw("ManagerFilter, got dynamic client", "totalElapsed", time.Since(start).String(), "elapsed", time.Since(step).String())
		if err != nil {
			log.Debugw("manager filter dynamic client create error", "err", err)
			err = kerrors.AsAPIError(err)
			resp.WriteHeaderAndJson(kerrors.AsStatusCode(err), err, restful.MIME_JSON)
			return
		}
		ctx = WithDynamicClient(ctx, dynamicClient)

		req.Request = req.Request.WithContext(ctx)

		log.Debugw("config,client,dynamicclient context done", "totalElapsed", time.Since(start).String())
		chain.ProcessFilter(req, resp)
	}
}
