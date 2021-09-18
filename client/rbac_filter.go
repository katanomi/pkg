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
	authv1 "k8s.io/api/authorization/v1"
	"knative.dev/pkg/logging"

	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kerrors "github.com/katanomi/pkg/errors"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// SelfSubjectReviewFilterForResource makes a self subject review based a configuration already present inside the
// request context using the user's bearer token
func SelfSubjectReviewFilterForResource(ctx context.Context, resourceAtt authv1.ResourceAttributes, namespaceParameter, nameParameter string) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		ctx := req.Request.Context()
		log := logging.FromContext(ctx).With("resource", resourceAtt.Resource, "group", resourceAtt.Group, "verb", resourceAtt.Verb)
		ctx = logging.WithLogger(ctx, log)
		err := SelfSubjectAccessReviewForResource(ctx, req.PathParameter(nameParameter), req.PathParameter(namespaceParameter), resourceAtt, false)
		if err != nil {
			log.Debugw("error veryfing user permissions", "err", err)
			kerrors.HandleError(req, resp, err)
			return
		}
		chain.ProcessFilter(req, resp)
	}
}

func SelfSubjectAccessReviewForResource(ctx context.Context, name, namespace string, resourceAtt authv1.ResourceAttributes, addName bool) (err error) {
	log := logging.FromContext(ctx)
	clt := Client(ctx)
	if clt == nil {
		// return error
		log.Debugw("error fetching the client from the context. Make sure the ManagerFilter is added before this filter")
		// return internal server error
		err = errors.NewUnauthorized("SelfSubjectReview needs user's client")
		return
	}
	if namespace != "" {
		resourceAtt.Namespace = namespace
	}
	if name != "" {
		resourceAtt.Name = name
	}

	review := &authv1.SelfSubjectAccessReview{
		Spec: authv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &resourceAtt,
		},
	}
	// this is only to be used within tests
	// the fake client somehow requests the name
	if addName {
		review.Name = resourceAtt.Name
		review.Namespace = resourceAtt.Namespace
	}
	err = clt.Create(ctx, review)
	if err != nil {
		log.Errorw("error evaluating SelfSubjectReview", "err", err)
		return
	}

	log.Debug("checking user permission against the resource",
		"allowed", review.Status.Allowed,
		"reason", review.Status.Reason,
		"evError", review.Status.EvaluationError,
	)

	if !review.Status.Allowed {
		err = errors.NewForbidden(schema.GroupResource{
			Group:    resourceAtt.Group,
			Resource: resourceAtt.Resource,
		}, resourceAtt.Name, fmt.Errorf("access not allowed"))
		return
	}
	return
}
