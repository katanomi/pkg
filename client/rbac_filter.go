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
	"net/http"
	"strings"

	"knative.dev/pkg/injection"

	"github.com/emicklei/go-restful/v3"
	authnv1 "k8s.io/api/authentication/v1"
	authv1 "k8s.io/api/authorization/v1"
	"k8s.io/apiserver/pkg/authentication/user"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/client"

	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kerrors "github.com/katanomi/pkg/errors"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GetResourceAttributesFunc helper function to warp a function to ResourceAttributeGetter
type GetResourceAttributesFunc func(ctx context.Context, req *restful.Request) (authv1.ResourceAttributes, error)

func (p GetResourceAttributesFunc) GetResourceAttributes(ctx context.Context, req *restful.Request) (authv1.ResourceAttributes, error) {
	return p(ctx, req)
}

// ResourceAttributeGetter describe an interface to get resource attributes form request
type ResourceAttributeGetter interface {
	// GetResourceAttributes get resource attributes from request
	GetResourceAttributes(ctx context.Context, req *restful.Request) (authv1.ResourceAttributes, error)
}

// SubjectAccessReviewClientGetter describe an interface to get client for subject access review
// It is usually used for cross-cluster authentication.
type SubjectAccessReviewClientGetter interface {
	// GetClient get k8s client according to request
	GetClient(ctx context.Context, req *restful.Request) (client.Client, error)
}

// DynamicSubjectReviewFilter makes a subject review and the ResourceAttribute can be dynamically obtained
func DynamicSubjectReviewFilter(ctx context.Context, resourceAttGetter ResourceAttributeGetter) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		resourceAtt, err := resourceAttGetter.GetResourceAttributes(ctx, req)
		if err != nil {
			kerrors.HandleError(req, resp, err)
			return
		}
		reqCtx := req.Request.Context()
		log := logging.FromContext(reqCtx).With(
			"resource", resourceAtt.Resource,
			"group", resourceAtt.Group,
			"verb", resourceAtt.Verb,
		)
		reqCtx = logging.WithLogger(reqCtx, log)

		var review subjectAccessReviewObjInterface
		if !isImpersonateRequest(reqCtx) {
			review = makeSelfSubjectAccessReview(resourceAtt)
		} else {
			u := User(reqCtx)
			if u == nil {
				err := fmt.Errorf("not found impersonate info in config")
				log.Error(err.Error())
				kerrors.HandleError(req, resp, err)
				return
			}
			review = makeSubjectAccessReview(resourceAtt, u)
		}

		var clt client.Client
		if clientGetter, ok := resourceAttGetter.(SubjectAccessReviewClientGetter); ok {
			clt, err = clientGetter.GetClient(ctx, req)
			if err != nil {
				log.Debugw("get custom client for authentication failed", "err", err)
				kerrors.HandleError(req, resp, err)
				return
			}
		}
		if clt == nil {
			clt = Client(reqCtx)
		}
		err = postSubjectAccessReview(reqCtx, clt, review)
		if err != nil {
			log.Debugw("error verifying user permissions", "err", err, "review", review.GetObject())
			kerrors.HandleError(req, resp, err)
			return
		}
		chain.ProcessFilter(req, resp)
	}
}

// SubjectReviewFilterForResource makes a self subject review based a configuration already present inside the
// request context using the user's bearer token
// also, it makes a subject review based on Impersonate User info in request header
func SubjectReviewFilterForResource(ctx context.Context, resourceAtt authv1.ResourceAttributes, namespaceParameter, nameParameter string) restful.FilterFunction {
	getter := GetResourceAttributesFunc(func(ctx context.Context, req *restful.Request) (authv1.ResourceAttributes, error) {
		attr := resourceAtt.DeepCopy()
		if ns := req.PathParameter(namespaceParameter); ns != "" {
			attr.Namespace = ns
		}
		if name := req.PathParameter(nameParameter); name != "" {
			attr.Name = name
		}
		return *attr, nil
	})
	return DynamicSubjectReviewFilter(ctx, getter)
}

func isImpersonateRequest(reqCtx context.Context) bool {
	var config = injection.GetConfig(reqCtx)
	if config == nil {
		return false
	}
	return config.Impersonate.UserName != "" || len(config.Impersonate.Groups) != 0 || len(config.Impersonate.Extra) != 0
}

func makeSubjectAccessReview(resourceAtt authv1.ResourceAttributes, user user.Info) subjectAccessReviewObjInterface {
	review := &authv1.SubjectAccessReview{
		Spec: authv1.SubjectAccessReviewSpec{
			ResourceAttributes: &resourceAtt,
			User:               user.GetName(),
			Groups:             user.GetGroups(),
			UID:                user.GetUID(),
			Extra:              map[string]authv1.ExtraValue{},
		},
	}
	for key, value := range user.GetExtra() {
		review.Spec.Extra[key] = value
	}

	return subjectAccessReviewObject{review}
}

func impersonateUser(req *http.Request) user.Info {

	u := user.DefaultInfo{}
	u.Name = req.Header.Get(authnv1.ImpersonateUserHeader)

	u.Groups = req.Header.Values(authnv1.ImpersonateGroupHeader)
	u.UID = req.Header.Get(authnv1.ImpersonateUIDHeader)

	if u.Name == "" && len(u.Groups) == 0 && u.UID == "" {
		return nil
	}

	for key, value := range req.Header {
		if u.Extra == nil {
			u.Extra = map[string][]string{}
		}
		if strings.HasPrefix(key, authnv1.ImpersonateUserExtraHeaderPrefix) {
			u.Extra[key] = value
		}
	}

	return &u
}

// SelfSubjectReviewFilterForResource makes a self subject review based a configuration already present inside the
// request context using the user's bearer token
//
// Deprecated: use SubjectReviewFilterForResource
func SelfSubjectReviewFilterForResource(ctx context.Context, resourceAtt authv1.ResourceAttributes, namespaceParameter, nameParameter string) restful.FilterFunction {
	return SubjectReviewFilterForResource(ctx, resourceAtt, namespaceParameter, nameParameter)
}

func makeSelfSubjectAccessReview(resourceAtt authv1.ResourceAttributes) subjectAccessReviewObjInterface {
	review := &authv1.SelfSubjectAccessReview{
		Spec: authv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &resourceAtt,
		},
	}

	return subjectAccessReviewObject{review}
}

type subjectAccessReviewObject struct {
	Object client.Object
}

func (object subjectAccessReviewObject) GetObject() client.Object {
	return object.Object
}

func (object subjectAccessReviewObject) GetResourceAttribute() (authv1.ResourceAttributes, error) {
	switch o := object.Object.(type) {
	case *authv1.SelfSubjectAccessReview:
		return *o.Spec.ResourceAttributes, nil
	case *authv1.SubjectAccessReview:
		return *o.Spec.ResourceAttributes, nil
	case *authv1.LocalSubjectAccessReview:
		return *o.Spec.ResourceAttributes, nil
	}

	return authv1.ResourceAttributes{}, fmt.Errorf("object type %T is not expected", object)
}

func (object subjectAccessReviewObject) GetSubjectAccessReviewStatus() (authv1.SubjectAccessReviewStatus, error) {
	switch o := object.Object.(type) {
	case *authv1.SelfSubjectAccessReview:
		return o.Status, nil
	case *authv1.SubjectAccessReview:
		return o.Status, nil
	case *authv1.LocalSubjectAccessReview:
		return o.Status, nil
	}

	return authv1.SubjectAccessReviewStatus{}, fmt.Errorf("object type %T is not expected", object)
}

type subjectAccessReviewObjInterface interface {
	GetObject() client.Object
	GetResourceAttribute() (authv1.ResourceAttributes, error)
	GetSubjectAccessReviewStatus() (authv1.SubjectAccessReviewStatus, error)
}

func postSubjectAccessReview(ctx context.Context, clt client.Client, reviewObj subjectAccessReviewObjInterface) (err error) {
	log := logging.FromContext(ctx)

	if clt == nil {
		// return error
		log.Debugw("error fetching the client from the context. Make sure the ManagerFilter is added before this filter")
		// return internal server error
		err = errors.NewUnauthorized("SubjectAccessReview needs client")
		return
	}

	review := reviewObj.GetObject()
	err = clt.Create(ctx, review)
	if err != nil {
		log.Errorw("error evaluating SubjectAccessReview", "err", err, "review", review)
		return
	}

	reviewStatus, err := reviewObj.GetSubjectAccessReviewStatus()
	if err != nil {
		return
	}
	log.Infow("checking user permission against the resource",
		"allowed", reviewStatus.Allowed,
		"reason", reviewStatus.Reason,
		"evError", reviewStatus.EvaluationError,
	)

	resourceAtt, err := reviewObj.GetResourceAttribute()
	if err != nil {
		return
	}
	if !reviewStatus.Allowed {
		err = errors.NewForbidden(schema.GroupResource{
			Group:    resourceAtt.Group,
			Resource: resourceAtt.Resource,
		}, resourceAtt.Name, fmt.Errorf("access not allowed"))
		return
	}
	return
}
