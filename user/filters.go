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

// Package user package is to manage the context of user information
package user

import (
	"context"
	"errors"
	"net/http"

	"k8s.io/apimachinery/pkg/api/meta"

	restful "github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	kclient "github.com/katanomi/pkg/client"
	kerror "github.com/katanomi/pkg/errors"
	authv1 "k8s.io/api/authorization/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/logging"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// UserInfoFilter is to parse the user login information from the request header into userinfo, and store it in the context
func UserInfoFilter(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
	curContext := req.Request.Context()
	log := logging.FromContext(curContext)
	userinfo, err := getUserInfoFromReq(req)
	if err != nil {
		log.Errorw("get user info from request failed", "err", err)
		return
	}
	newCtx := WithUserInfo(curContext, userinfo)
	req.Request = req.Request.WithContext(newCtx)
	chain.ProcessFilter(req, res)
}

// UserOwnedResourcePermissionFilter provides a mechanism to check permissions for user owned resource
// the user owned resource could use annotation to annotated owner name
// and the filter will check if current user could execcute the `verb` for the resource
// more doc about it please see Spec: user owned resource permission check
func UserOwnedResourcePermissionFilter(appCtx context.Context, gvr *schema.GroupVersionResource) restful.FilterFunction {
	appClient := kclient.Client(appCtx)
	restMapper := appClient.RESTMapper()

	return func(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
		ctxInReq := req.Request.Context()
		log := logging.FromContext(ctxInReq).Named("user-owned-resource")
		logging.WithLogger(ctxInReq, log)

		if gvr == nil {
			gvr = &schema.GroupVersionResource{
				Group:    req.PathParameter("group"),
				Version:  req.PathParameter("version"),
				Resource: req.PathParameter("resource"),
			}
		}
		log = log.With("gvr", gvr)

		// check rbac for current user
		rbacReviewStatus, resourcecOwner, resource, err := resourecRBACAllowed(appCtx, req, gvr, restMapper)
		if err != nil {
			log.Errorw("resource rbac allowed error", "err", err)
			kerror.HandleError(req, res, err)
			return
		}

		if rbacReviewStatus.Allowed {
			log.Debugf("rbac permission allowed")
			chain.ProcessFilter(req, res)
			return
		}

		log.Debugw("rbac permission not allowed", "message", rbacReviewStatus.String())

		// if current user has no permission
		// user should only could request resource of current user
		equal, err := resourceOwnerEqualToUserInReq(ctxInReq, resourcecOwner)

		if err != nil {
			kerror.HandleError(req, res, err)
			return
		}

		if !equal {
			log.Warnw("no permissions, username is not equal to resource owner", "resource-owner", resourcecOwner)
			kerror.HandleError(req, res, k8serrors.NewForbidden(gvr.GroupResource(), resource.GetName(), errors.New("no permissions")))
			return
		}

		log.Debugf("allow to %s resource %s, user only request user owned resource", req.Request.Method, gvr.String())
		chain.ProcessFilter(req, res)
		return
	}
}

func verbFromReq(req *restful.Request) string {
	verb := "get"
	if req.Request.Method == http.MethodPost {
		verb = "create"
	}
	if req.Request.Method == http.MethodPut {
		verb = "update"
	}
	return verb
}

func resourecRBACAllowed(appCtx context.Context, req *restful.Request,
	gvr *schema.GroupVersionResource, restMapper meta.RESTMapper) (*authv1.SubjectAccessReviewStatus, string, *unstructured.Unstructured, error) {
	log := logging.FromContext(req.Request.Context())
	ctxInReq := req.Request.Context()

	verb := verbFromReq(req)
	log = log.With("gvr", *gvr).With("verb", verb)

	gvk, err := restMapper.KindFor(*gvr)
	if err != nil {
		log.Errorw("get kind from groupversionresource error", "err", err)
		return nil, "", nil, err
	}

	resourceOwner, obj, err := getResourceFromRequest(req, gvk, appCtx)
	if err != nil {
		return nil, "", nil, err
	}
	log = log.With("resource", obj.GetNamespace()+"/"+obj.GetName())
	req.Request = req.Request.WithContext(WithEntity(ctxInReq, obj))

	clientInApp := kclient.Client(appCtx)
	status, err := resourceRBACCheck(ctxInReq, clientInApp, verb, *gvr, obj.GetNamespace(), obj.GetName())
	if err != nil {
		log.Errorw("resource permission check error", "namespace", obj.GetNamespace(), "name", obj.GetName(), "err", err)
		return nil, "", nil, err
	}

	return status, resourceOwner, obj, nil
}

func resourceOwnerEqualToUserInReq(ctxInReq context.Context, resourceOwner string) (bool, error) {
	log := logging.FromContext(ctxInReq)

	userInReq := kclient.User(ctxInReq)
	if userInReq == nil {
		log.Errorw("not found userinfo in context")
		return false, k8serrors.NewBadRequest("not found userinfo in request")
	}

	log = log.With("usernameInReq", userInReq.GetName(), "resourceOwner", resourceOwner)

	if userInReq.GetName() != resourceOwner {
		log.Infow("usename is not equal")
		return false, nil
	}

	return true, nil
}

func resourceRBACCheck(ctx context.Context, clientInApp ctrlclient.Client, verb string, gvr schema.GroupVersionResource, namespace string, name string) (*authv1.SubjectAccessReviewStatus, error) {
	log := logging.FromContext(ctx)
	user := kclient.User(ctx)

	review := &authv1.SubjectAccessReview{
		Spec: authv1.SubjectAccessReviewSpec{
			ResourceAttributes: &authv1.ResourceAttributes{
				Verb:      verb,
				Group:     gvr.Group,
				Version:   gvr.Version,
				Resource:  gvr.Resource,
				Namespace: namespace,
				Name:      name,
			},
			User:   user.GetName(),
			Groups: user.GetGroups(),
			UID:    user.GetUID(),
			Extra:  map[string]authv1.ExtraValue{},
		},
	}

	err := clientInApp.Create(ctx, review)
	if err != nil {
		log.Errorw("error evaluating SubjectAccessReview", "err", err, "review", review)
		return nil, err
	}
	return &review.Status, nil

}

func getResourceFromRequest(req *restful.Request, gvk schema.GroupVersionKind, appCtx context.Context) (
	owner string, obj *unstructured.Unstructured, err error) {
	ctx := req.Request.Context()
	log := logging.FromContext(ctx)

	obj = &unstructured.Unstructured{}
	obj.SetGroupVersionKind(gvk)

	if req.Request.Method == http.MethodPost || req.Request.Method == http.MethodPut {
		err := req.ReadEntity(obj)
		if err != nil {
			log.Errorw("error read entity from body", "err", err)
			return "", nil, k8serrors.NewBadRequest(err.Error())
		}
	}

	if req.Request.Method == http.MethodGet {
		resourceName := req.PathParameter("name")
		resourceNamespace := req.PathParameter("namespace")

		clientInApp := kclient.Client(appCtx)
		err = clientInApp.Get(ctx, ctrlclient.ObjectKey{Name: resourceName, Namespace: resourceNamespace}, obj)
		if err != nil {
			log.Errorw("get resource error", "err", err, "gvk", gvk.String(), "name", resourceNamespace+"/"+resourceName)
			return "", nil, err
		}
	}

	annotations := obj.GetAnnotations()
	if len(annotations) == 0 || annotations[metav1alpha1.UserOwnedAnnotationKey] == "" {
		return "", nil, k8serrors.NewBadRequest("request object does not contains annotation key: " + metav1alpha1.UserOwnedAnnotationKey)
	}

	return annotations[metav1alpha1.UserOwnedAnnotationKey], obj, nil
}
