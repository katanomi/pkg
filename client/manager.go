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
	"strings"
	"time"

	"github.com/katanomi/pkg/multicluster"
	apiserverrequest "k8s.io/apiserver/pkg/endpoints/request"

	"k8s.io/apiserver/pkg/authentication/user"

	"github.com/golang-jwt/jwt/v4"
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

// NewClusterRegistryClientFunc constructs a multi cluster client based on the specified config.
type NewClusterRegistryClientFunc func(*rest.Config) (multicluster.Interface, error)

// Manager dynamically generates client based on user requests
type Manager struct {
	// GetConfig retrieves a configuration based on a request
	GetConfig GetConfigFunc
	// GetBasicConfig retrieves a configuration based on a request
	GetBasicConfig GetBaseConfigFunc
	// NewClusterRegistryClient constructs a multi cluster client based on the specified config.
	NewClusterRegistryClient NewClusterRegistryClientFunc
	*zap.SugaredLogger
}

// WithCtxManagerFilters put the manager in the context into the webservice.
func WithCtxManagerFilters(ctx context.Context, ws *restful.WebService) error {
	if manager := ManagerCtx(ctx); manager != nil {
		filters, err := manager.Filters(ctx)
		if err != nil {
			return err
		}
		for _, filter := range filters {
			ws = ws.Filter(filter)
		}
	}
	return nil
}

// NewManager initializes a new manager based on func
func NewManager(ctx context.Context, get GetConfigFunc, baseConfig GetBaseConfigFunc, newClusterRegistryClient NewClusterRegistryClientFunc) *Manager {
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
		SugaredLogger:            logging.FromContext(ctx),
		GetConfig:                configGetter,
		GetBasicConfig:           baseConfig,
		NewClusterRegistryClient: newClusterRegistryClient,
	}
}

// Filter returns a filter to be used in
func (m *Manager) Filter(ctx context.Context) restful.FilterFunction {
	return ManagerFilter(ctx, m)
}

// Filters returns manager.Filter with ImpersonateFilter
func (m *Manager) Filters(ctx context.Context) (filters []restful.FilterFunction, err error) {
	return []restful.FilterFunction{
		ManagerFilter(ctx, m),
		ImpersonateFilter(ctx),
	}, nil
}

// ManagerFilter generates filter based on a manager to create a config based in a request and injects into context
func ManagerFilter(ctx context.Context, mgr *Manager) restful.FilterFunction {
	log := logging.FromContext(ctx).Named("manager-filter")
	scheme := kscheme.Scheme(ctx)
	serviceAccountClient := Client(ctx)
	configInApp := GetAppConfig(ctx)

	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		start := time.Now()
		config, err := mgr.GetConfig(req, mgr.GetBasicConfig)
		log.Debugw("ManagerFilter, got config", "totalElapsed", time.Since(start).String())
		step := time.Now()
		if err != nil {
			log.Debugw("manager filter config get error", "err", err)
			kerrors.HandleError(req, resp, err)
			return
		}

		reqCtx := req.Request.Context()

		// setting defaults to the config
		config.Burst = DefaultBurst
		config.QPS = DefaultQPS
		config.Timeout = DefaultTimeout

		reqCtx = injection.WithConfig(reqCtx, config)
		reqCtx = WithAppConfig(reqCtx, configInApp)

		user, err := UserFromBearerToken(strings.TrimPrefix(req.Request.Header.Get("Authorization"), "Bearer "))
		if err != nil {
			log.Errorw("cannot get user info from token", "err", err)
			kerrors.HandleError(req, resp, err)
			return
		}
		reqCtx = apiserverrequest.WithUser(reqCtx, user)

		directClient, err := client.New(config, client.Options{Scheme: scheme, Mapper: serviceAccountClient.RESTMapper()})
		log.Debugw("ManagerFilter, got direct client", "totalElapsed", time.Since(start).String(), "elapsed", time.Since(step).String())
		step = time.Now()
		if err != nil {
			log.Debugw("manager filter direct client create error", "err", err)
			kerrors.HandleError(req, resp, err)
			return
		}
		reqCtx = WithClient(reqCtx, directClient)

		dynamicClient, err := dynamic.NewForConfig(config)
		log.Debugw("ManagerFilter, got dynamic client", "totalElapsed", time.Since(start).String(), "elapsed", time.Since(step).String())
		if err != nil {
			log.Debugw("manager filter dynamic client create error", "err", err)
			kerrors.HandleError(req, resp, err)
			return
		}
		reqCtx = WithDynamicClient(reqCtx, dynamicClient)

		if mgr.NewClusterRegistryClient != nil {
			multiClusterClient, err := mgr.NewClusterRegistryClient(config)
			if err != nil {
				log.Errorw("cannot get multi cluster client", "err", err)
				kerrors.HandleError(req, resp, err)
				return
			}
			reqCtx = multicluster.WithMultiCluster(reqCtx, multiClusterClient)
			log.Debugw("get multi cluster", "totalElapsed", time.Since(start).String())
		}

		req.Request = req.Request.WithContext(reqCtx)

		log.Debugw("config,client,dynamicclient context done", "totalElapsed", time.Since(start).String())
		chain.ProcessFilter(req, resp)
	}
}

func UserFromBearerToken(rawToken string) (user.Info, error) {
	mapClaims := jwt.MapClaims{}

	_, _, err := new(jwt.Parser).ParseUnverified(rawToken, mapClaims)
	if err != nil {
		return nil, err
	}
	info := user.DefaultInfo{}
	// username is claim by sub when current token is generated serviceaccount
	if val, ok := mapClaims["sub"].(string); ok {
		info.Name = val
	}

	// username is claim by email
	if val, ok := mapClaims["email"].(string); ok {
		info.Name = val
	}

	if val, ok := mapClaims["groups"].([]string); ok {
		info.Groups = val
	}

	return &info, nil
}
