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

package main

import (
	"context"
	"net/http"
	"sync"

	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/logging"

	"github.com/emicklei/go-restful/v3"
	"github.com/katanomi/pkg/client"
	"github.com/katanomi/pkg/sharedmain"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes/scheme"
)

func main() {

	sharedmain.App("sample-api").
		Scheme(scheme.Scheme).
		Log().
		Webservices(&Sample{}).
		Filters(fooFilter()).
		APIDocs().
		Profiling().
		Run()
}

func fooFilter() restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		ctx := req.Request.Context()
		log := logging.FromContext(ctx)
		log.Debugw("Before request", "url", req.Request.RequestURI, "Method", req.Request.Method)
		chain.ProcessFilter(req, resp)
		log.Debugw("After request", "url", req.Request.RequestURI, "Method", req.Request.Method, "StatusCode", resp.StatusCode())
	}
}

var _ sharedmain.WebService = &Sample{}

type Sample struct {
	once sync.Once
	ctx  context.Context
}

func (s *Sample) Name() string {
	return "sample"
}

func (s *Sample) Setup(ctx context.Context, add sharedmain.AddToRestContainer, logger *zap.SugaredLogger) error {
	s.once.Do(func() {
		s.ctx = ctx

		ws := new(restful.WebService)
		ws.Path("/v1")

		ws.Route(
			ws.GET("/hello").To(s.Hello),
		)

		add(ws)
	})

	return nil
}

func (s *Sample) Hello(req *restful.Request, res *restful.Response) {
	reqCtx := req.Request.Context()

	dclient, err := client.DynamicClient(reqCtx)
	if err != nil {
		res.WriteError(http.StatusInternalServerError, err)
		return
	}
	list, err := dclient.Resource(schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "namespaces",
	}).List(reqCtx, metav1.ListOptions{ResourceVersion: "0"})

	if err != nil {
		res.WriteError(http.StatusInternalServerError, err)
		return
	}

	client := client.Client(s.ctx)
	ns := &corev1.NamespaceList{}
	err = client.List(reqCtx, ns)
	if err != nil {
		res.WriteError(http.StatusInternalServerError, err)
		return
	}

	res.WriteAsJson(map[string]int{"dynamicClient": len(list.Items), "client": len(ns.Items)})
}
