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
	"fmt"
	"net/http"

	restful "github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	kclient "github.com/katanomi/pkg/client"
	kerror "github.com/katanomi/pkg/errors"
	authv1 "k8s.io/api/authorization/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/injection"
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

func UserOwnedResourcePermissionFilter(appCtx context.Context, gvr *schema.GroupVersionResource) restful.FilterFunction {
	appClient := kclient.Client(appCtx)
	fmt.Printf("appclient: %#v", appCtx)
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

		log = log.With("gvr", *gvr)

		gvk, err := restMapper.KindFor(*gvr)
		if err != nil {
			log.Errorw("get kind from groupversionresource error", "err", err)
			kerror.HandleError(req, res, err)
			return
		}

		usernameOfResource, obj := getResourceUsername(req, res, gvk, appCtx)
		if usernameOfResource == "" {
			return
		}

		req.Request = req.Request.WithContext(WithEntity(ctxInReq, obj))

		verb := "get"
		if req.Request.Method == http.MethodPost {
			verb = "create"
		}
		if req.Request.Method == http.MethodPut {
			verb = "update"
		}

		clientInReq := kclient.Client(ctxInReq)
		status, err := resourcePermissionCheck(ctxInReq, clientInReq, verb, *gvr, obj.GetNamespace(), obj.GetName())
		if err != nil {
			log.Errorw("resource permission check error", "namespace", obj.GetNamespace(), "name", obj.GetName(), "err", err)
			kerror.HandleError(req, res, err)
			return
		}

		if status.Allowed {
			log.Debugf("allow to %s resource %s, check pass", verb, gvr.String())
			chain.ProcessFilter(req, res)
			return
		} else {
			log.Debugw("self permission check error", "message", status.String())
		}

		// if self has no permission
		// user should only could request resource of current user
		userInReq, ok := UserInfoFrom(ctxInReq)
		if !ok {
			log.Errorw("not found userinfo in context")
			kerror.HandleError(req, res, k8serrors.NewBadRequest("not found userinfo in request"))
			return
		}

		if userInReq.Username != usernameOfResource {
			log.Warnw(fmt.Sprintf("no permissions to %s %s %s, username is not equal", verb, gvr.Resource, obj.GetNamespace()+"/"+obj.GetName()),
				"usernameInReq", userInReq, "usernameOfResource", usernameOfResource)
			kerror.HandleError(req, res, k8serrors.NewForbidden(
				gvr.GroupResource(),
				obj.GetName(), fmt.Errorf("no permissions to %s %s %s", verb, gvr.Resource, obj.GetNamespace()+"/"+obj.GetName()),
			))
			return
		}

		log.Debugf("allow to %s resource %s, user only request user owned resource", verb, gvr.String())
		chain.ProcessFilter(req, res)
		return
	}
}

func resourcePermissionCheck(ctx context.Context, client ctrlclient.Client, verb string, gvr schema.GroupVersionResource, namespace string, name string) (*authv1.SubjectAccessReviewStatus, error) {
	log := logging.FromContext(ctx)
	user := kclient.User(ctx)

	if user == nil {
		review := &authv1.SelfSubjectAccessReview{
			Spec: authv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authv1.ResourceAttributes{
					Verb:      verb,
					Group:     gvr.Group,
					Version:   gvr.Version,
					Resource:  gvr.Resource,
					Namespace: namespace,
					Name:      name,
				},
			},
		}

		err := client.Create(ctx, review)
		if err != nil {
			log.Errorw("error evaluating SelfSubjectAccessReview", "err", err, "review", review)
			return nil, err
		}
		return &review.Status, nil
	}

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
	restConfig := injection.GetConfig(ctx)

	err := client.Create(ctx, review)
	if err != nil {
		log.Errorw("error evaluating SubjectAccessReview", "err", err, "review", review)
		return nil, err
	}
	log.Infow("review result", "review", review, "client", fmt.Sprintf("%#v", client), "config", restConfig)
	return &review.Status, nil

}

func getResourceUsername(req *restful.Request, res *restful.Response, gvk schema.GroupVersionKind, appCtx context.Context) (
	username string, obj *unstructured.Unstructured) {
	ctx := req.Request.Context()
	log := logging.FromContext(ctx)

	obj = &unstructured.Unstructured{}
	obj.SetGroupVersionKind(gvk)

	if req.Request.Method == http.MethodPost || req.Request.Method == http.MethodPut {
		err := req.ReadEntity(obj)
		if err != nil {
			log.Errorw("error read entity from body", "err", err)
			kerror.HandleError(req, res, errors.NewBadRequest(err.Error()))
			return "", nil
		}
	}
	if req.Request.Method == http.MethodGet {
		resourceName := req.PathParameter("name")
		resourceNamespace := req.PathParameter("namespace")

		clientInApp := kclient.Client(appCtx)
		err := clientInApp.Get(ctx, ctrlclient.ObjectKey{Name: resourceName, Namespace: resourceNamespace}, obj)
		if err != nil {
			log.Errorw("get resource error", "err", err, "gvk", gvk.String(), "name", resourceNamespace+"/"+resourceName)
			kerror.HandleError(req, res, err)
			return "", nil
		}
	}

	annotations := obj.GetAnnotations()
	if len(annotations) == 0 || annotations[metav1alpha1.UserOwnedAnnotationKey] == "" {
		kerror.HandleError(req, res, errors.NewBadRequest("request object does not contains anotation key: "+metav1alpha1.UserOwnedAnnotationKey))
		return "", nil
	}

	return annotations[metav1alpha1.UserOwnedAnnotationKey], obj
}
