/*
Copyright 2023 The Katanomi Authors.

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

package config

import (
	"context"
	"fmt"
	"net/http"

	"github.com/AlaudaDevops/pkg/errors"
	restful "github.com/emicklei/go-restful/v3"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"knative.dev/pkg/logging"
)

// ConfigKeyExpectedValueFunc is a helper function to check if configmap has expected value
// If the value is not as expected, an error is expected to be returned
type ConfigKeyExpectedValueFunc func(ctx context.Context, req *restful.Request, key string, value FeatureValue) (err error)

// ConfigFilter adds a restful filter to manager to watch configmap and and custom validation
// according a specific key value pair.
func ConfigFilter(ctx context.Context, manager *Manager, configKey string, expectedKeyValueFunc ConfigKeyExpectedValueFunc) func(*restful.Request, *restful.Response, *restful.FilterChain) {
	return func(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
		featureValue := manager.GetFeatureFlag(configKey)
		if err := expectedKeyValueFunc(ctx, req, configKey, featureValue); err != nil {
			log := logging.FromContext(ctx)
			log.Debugw("Error in ConfigFilter, will return", "err", err, "code", res.StatusCode())
			errors.HandleError(req, res, err)
			return
		}
		chain.ProcessFilter(req, res)
	}
}

// ConfigFilterNotFoundWhenNotTrue is a helper ConfigKeyExpectedValue implementation that checks if the value is a boolean true
// value, if not true will return a standard 404 not found error
func ConfigFilterNotFoundWhenNotTrue(ctx context.Context, req *restful.Request, key string, value FeatureValue) (err error) {
	if ok, _ := value.AsBool(); !ok {
		return apierrors.NewGenericServerResponse(
			http.StatusNotFound,
			req.Request.Method,
			errors.RESTAPIGroupResource,
			req.Request.URL.String(),
			fmt.Sprintf("%s Not Found", req.Request.URL.String()),
			0,
			false,
		)
	}
	return nil
}
