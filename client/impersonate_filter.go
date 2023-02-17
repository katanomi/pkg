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
	"time"

	apiserverrequest "k8s.io/apiserver/pkg/endpoints/request"

	"knative.dev/pkg/injection"

	"github.com/emicklei/go-restful/v3"
)

var (
	// k8s.io/apiserver/pkg/server/options/authorization.go
	allowCacheTTL = 10 * time.Second
	denyCacheTTL  = 10 * time.Second
)

// ImpersonateFilter will inject current user into context and inject impersonate information into rest.Config in request
func ImpersonateFilter(_ context.Context) restful.FilterFunction {

	return func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {

		reqCtx := request.Request.Context()

		user := impersonateUser(request.Request)
		if user == nil {
			fmt.Println("user is nil")

			chain.ProcessFilter(request, response)
			return
		}

		reqConfig := injection.GetConfig(reqCtx)

		// change config to impersonate config
		reqConfig.Impersonate.UID = user.GetUID()
		reqConfig.Impersonate.Groups = user.GetGroups()
		reqConfig.Impersonate.UserName = user.GetName()
		reqConfig.Impersonate.Extra = user.GetExtra()

		reqCtx = injection.WithConfig(reqCtx, reqConfig)
		reqCtx = apiserverrequest.WithUser(reqCtx, user)
		request.Request = request.Request.WithContext(reqCtx)

		chain.ProcessFilter(request, response)
	}
}
