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
	"fmt"

	"github.com/emicklei/go-restful/v3"
	kerrors "github.com/katanomi/pkg/errors"
	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/logging"
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
	if baseConfig == nil {
		baseConfig = config.GetConfig
	}
	return &Manager{
		SugaredLogger:  logging.FromContext(ctx),
		GetConfig:      get,
		GetBasicConfig: baseConfig,
	}
}

// Filter returns a filter to be used in
func (m *Manager) Filter() restful.FilterFunction {
	return ManagerFilter(m)
}

func (m *Manager) DynamicClient(req *restful.Request) (dynamic.Interface, error) {
	config := injection.GetConfig(req.Request.Context())
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return dynamicClient, nil
}

// ManagerFilter generates filter based on a manager to create a config based in a request and injects into context
func ManagerFilter(mgr *Manager) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		config, err := mgr.GetConfig(req, mgr.GetBasicConfig)
		if err != nil {
			fmt.Println("err", err)
			err = kerrors.AsAPIError(err)
			resp.WriteHeaderAndJson(kerrors.AsStatusCode(err), err, restful.MIME_JSON)
			return
		}
		ctx := req.Request.Context()
		ctx = injection.WithConfig(ctx, config)
		req.Request = req.Request.WithContext(ctx)
		chain.ProcessFilter(req, resp)
	}
}
