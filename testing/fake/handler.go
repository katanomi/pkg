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

package fake

import (
	"context"
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	kerrors "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/sharedmain"
	"github.com/katanomi/pkg/testing/fake/opa"
	"go.uber.org/zap"
	"knative.dev/pkg/logging"
)

// PolicyHandler is a struct that handles OPA policies requests.
type PolicyHandler struct {
	*zap.SugaredLogger
	Store Store
}

func (h *PolicyHandler) Name() string {
	return "Policyhandler"
}

func (h *PolicyHandler) Setup(ctx context.Context, add sharedmain.AddToRestContainer, logger *zap.SugaredLogger) error {
	h.SugaredLogger = logger

	if h.Store == nil {
		return fmt.Errorf("store required to save policy")
	}

	if err := h.Store.Setup(ctx); err != nil {
		return err
	}

	ws := new(restful.WebService).
		Path("mock").
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON)

	ws.Route(
		// makes the path as: /plugins/v1alpha1/<plugin path>/mock/policy
		ws.POST("policy").
			Reads(opa.Policy{}, "OPA Policy").
			Returns(http.StatusCreated, "OPA policy Created", opa.Policy{}).
			To(h.Handle),
	)
	add(ws)
	return nil
}

func (h *PolicyHandler) Handle(req *restful.Request, resp *restful.Response) {
	ctx := req.Request.Context()
	log := logging.FromContext(ctx)

	log.Debugw("starting to handle opa policy request")

	policy := &opa.Policy{}
	if err := req.ReadEntity(policy); err != nil {
		kerrors.HandleError(req, resp, err)
		return
	}

	err := h.Store.Create(ctx, policy)
	if err != nil {
		kerrors.HandleError(req, resp, err)
		return
	}
	h.Infow("handle opa policy", "policy", policy)
	if err := resp.WriteEntity(policy); err != nil {
		kerrors.HandleError(req, resp, err)
		return
	}

	return
}
