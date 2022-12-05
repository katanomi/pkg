/*
Copyright 2022 The Katanomi Authors.

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

package route

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	kerrors "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/plugin/client"
)

type livenessCheck struct {
	impl client.LivenessChecker
	tags []string
}

// NewLivenessCheck create a liveness check route with plugin client
func NewLivenessCheck(impl client.LivenessChecker) Route {
	return &livenessCheck{
		tags: []string{"tools", "liveness"},
		impl: impl,
	}
}

// Register register route
func (i *livenessCheck) Register(ws *restful.WebService) {
	ws.Route(
		ws.GET("/tools/liveness").To(i.CheckAlive).
			Doc("LivenessCheck").
			Metadata(restfulspec.KeyOpenAPITags, i.tags).
			Returns(http.StatusOK, "OK", nil),
	)
}

// CheckAlive register check alive route
func (i *livenessCheck) CheckAlive(request *restful.Request, response *restful.Response) {
	err := i.impl.CheckAlive(request.Request.Context())
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
}

type initialize struct {
	impl client.Initializer
	tags []string
}

// NewInitializer create a route for the tool service initializer
func NewInitializer(impl client.Initializer) Route {
	return &initialize{
		tags: []string{"tools", "initialize"},
		impl: impl,
	}
}

// Register register route
func (i *initialize) Register(ws *restful.WebService) {
	ws.Route(
		ws.GET("/tools/initialize").To(i.Initialize).
			Doc("Initialize").
			Metadata(restfulspec.KeyOpenAPITags, i.tags).
			Returns(http.StatusOK, "OK", nil),
	)
}

// Initialize register initialize route
func (i *initialize) Initialize(request *restful.Request, response *restful.Response) {
	err := i.impl.Initialize(request.Request.Context())
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
}
