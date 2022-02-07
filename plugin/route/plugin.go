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

package route

import (
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/emicklei/go-restful/v3"
	kerrors "github.com/katanomi/pkg/errors"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type pluginMethodUnsupport struct {
	tags []string
}

// NewPluginMethodUnsupport return an error with plugin client
func NewPluginMethodUnsupport() Route {
	return &pluginMethodUnsupport{
		tags: []string{"plugin"},
	}
}

// Register route
func (a *pluginMethodUnsupport) Register(ws *restful.WebService) {
	ws.Route(
		ws.POST("/{anything:*}").To(a.PluginMethodSupportError),
	)
	ws.Route(
		ws.DELETE("/{anything:*}").To(a.PluginMethodSupportError),
	)
	ws.Route(
		ws.GET("/{anything:*}").To(a.PluginMethodSupportError),
	)
	ws.Route(
		ws.PATCH("/{anything:*}").To(a.PluginMethodSupportError),
	)
}

func (a *pluginMethodUnsupport) PluginMethodSupportError(request *restful.Request, response *restful.Response) {
	err := errors.NewGenericServerResponse(
		http.StatusNotImplemented,
		request.Request.Method,
		schema.GroupResource{},
		request.Request.URL.String(),
		"The plugin has not implemented this function yet",
		0,
		false,
	)
	err.ErrStatus.Reason = metav1.StatusReasonMethodNotAllowed
	kerrors.HandleError(request, response, err)
}
