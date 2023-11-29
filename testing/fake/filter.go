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
	"net/http"

	"github.com/emicklei/go-restful/v3"
	kerrors "github.com/katanomi/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Filter is a method that intercepts and processes incoming HTTP requests using OPA policies.
// It evaluates the request against a specific policy identified by the request's method and path,
// and modifies the response based on the policy's evaluation results.
func (h *PolicyHandler) Filter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	ctx := req.Request.Context()

	policy, err := h.Store.Get(ctx, IDFromRequest(req))
	if client.IgnoreNotFound(err) != nil {
		h.Errorw("get opa policy failed", "error", err)

		chain.ProcessFilter(req, resp)
		return
	}

	if policy == nil {
		chain.ProcessFilter(req, resp)
		return
	}

	// Extract input for the OPA policy evaluation from the request.
	input, err := InputFromRequest(req)
	if err != nil {
		kerrors.HandleError(req, resp, err)
		return
	}
	if err = policy.Eval(ctx, input); err != nil {
		kerrors.HandleError(req, resp, err)
		return
	}

	matched := policy.BoolResult(matchedQuery)
	if !matched {
		chain.ProcessFilter(req, resp)
		return
	}

	// Get the status code from the policy evaluation result.
	// see template.rego
	status := policy.IntResult(statusQuery)
	if status == 0 {
		status = http.StatusOK
	}
	// Get the response body from the policy evaluation result.
	// see template.rego
	result := policy.MapResult(resultQuery)

	_ = resp.WriteHeaderAndEntity(status, result)
}
