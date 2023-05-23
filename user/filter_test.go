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
	"net/http/httptest"
	"testing"

	"go.etcd.io/etcd/client/v2"

	"k8s.io/apiserver/pkg/authentication/user"

	"github.com/onsi/gomega"

	"knative.dev/pkg/logging"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"k8s.io/apimachinery/pkg/api/meta"

	kclient "github.com/katanomi/pkg/client"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"

	mockfakeclient "github.com/katanomi/pkg/testing/mock/sigs.k8s.io/controller-runtime/pkg/client"
	authv1 "k8s.io/api/authorization/v1"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/golang/mock/gomock"

	"k8s.io/apimachinery/pkg/runtime/schema"

	restful "github.com/emicklei/go-restful/v3"
	batchv1 "k8s.io/api/batch/v1"
)

func TestUserInfoFilter(t *testing.T) {
	tokenString := "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImVtYWlsIjoiam9obkB0ZXN0LmNvbSIsImlhdCI6MTY1MzM3NzM0MCwiZXhwIjoxNjUzMzgwOTQwfQ.tT5YWrFu2F_0NWg3JFbpIaC-HbaBgJtRcff_O3Ti225RZztON1gpZOfFf06EIzWolkUrtAfrTFXonHx2rh5w8mK82kFX5sLPxiV3yDfvxU9KuE66IW_48ykFU6puBcnVMZNWUtPLHcH4OAq51zbca9eh3AFiFldnJ3YRJ85Bigky8nc1uf3CCm5c9d7P8q5SGDJZvGLHOXqQ4_aSfpX0HDTHbgw_82d7FpckeDYRdFN89LORaLRChbBM1IdfH4JvI9WtO3PDU7Ce49nMmmidhsESAalxfMrN3LIPmMz7vY0FJaBW24oA0FwJvq1Q6O_jzupMHRGS-clkUImRw185cA"
	_req := &http.Request{
		Header: map[string][]string{},
	}
	_req.Header.Set("Authorization", tokenString)

	req := &restful.Request{Request: _req}
	res := &restful.Response{}
	target := func(req *restful.Request, resp *restful.Response) {}
	chain := &restful.FilterChain{Target: target}
	UserInfoFilter(req, res, chain)
	userinfo := UserInfoValue(req.Request.Context())
	if userinfo.GetName() != "John Doe" {
		t.Errorf("should name is John Doe, but got %s", userinfo.GetName())
	}
	if userinfo.GetEmail() != "john@test.com" {
		t.Errorf("should email is john@test.com, but got %s", userinfo.GetEmail())
	}

	req.Request.Header["Impersonate-User"] = []string{"jackson"}
	UserInfoFilter(req, res, chain)
	userinfo = UserInfoValue(req.Request.Context())
	if userinfo.GetName() != "jackson" {
		t.Errorf("should name is jackson, but got %s", userinfo.GetName())
	}

}

func TestUserOwnedResourcePermissionFilter(t *testing.T) {
	ctx := context.TODO()
	ctx = logging.WithLogger(ctx, logger)
	mockCtl := gomock.NewController(t)

	mockClient := mockfakeclient.NewMockClient(mockCtl)
	mockClient.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {

			sub, ok := obj.(*authv1.SubjectAccessReview)
			if ok {
				if sub.Spec.User == "admin" {
					sub.Status.Denied = false
					sub.Status.Allowed = true
					return nil
				}
				sub.Status.Denied = true
				sub.Status.Allowed = false
				return nil
			}

			sub.Status.Denied = true
			sub.Status.Allowed = false
			return nil
		}).AnyTimes()

	mockClient.EXPECT().Get(gomock.Any(), gomock.Eq(ctrlclient.ObjectKey{Name: "job1", Namespace: "default"}), gomock.Any()).DoAndReturn(
		func(ctx context.Context, objectKey ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...client.GetOptions) error {
			job := obj.(*unstructured.Unstructured)
			job.SetAnnotations(map[string]string{
				metav1alpha1.UserOwnedAnnotationKey: "jackson",
			})
			job.SetName("job1")
			job.SetNamespace("default")

			return nil
		}).AnyTimes()
	mockClient.EXPECT().RESTMapper().DoAndReturn(func() meta.RESTMapper {
		fmt.Println("in restmapper mock")
		mapper := meta.NewDefaultRESTMapper([]schema.GroupVersion{batchv1.SchemeGroupVersion})
		mapper.AddSpecific(schema.GroupVersionKind{
			Group:   "batch",
			Version: "v1",
			Kind:    "Job",
		}, schema.GroupVersionResource{
			Group:    "batch",
			Version:  "v1",
			Resource: "jobs",
		}, schema.GroupVersionResource{
			Group:    "batch",
			Version:  "v1",
			Resource: "job",
		}, meta.RESTScopeNamespace)
		return mapper
	}).AnyTimes()

	ctx = kclient.WithClient(ctx, mockClient)

	filter := UserOwnedResourcePermissionFilter(ctx, &schema.GroupVersionResource{
		Group:    "batch",
		Version:  "v1",
		Resource: "jobs",
	})

	target := func(req *restful.Request, resp *restful.Response) {
		resp.Write([]byte("OK"))
	}
	chain := &restful.FilterChain{Target: target}

	var data = []struct {
		desc          string
		permissionOK  bool
		userNameInReq string

		respBody string
		respCode int32
	}{
		{
			desc:          "user should only request his resource",
			userNameInReq: "jackson",
			respBody:      "OK",
		},
		{
			desc:          "user cannot request resource not owned by him",
			userNameInReq: "user",
			respCode:      http.StatusForbidden,
		},
		{
			desc:          "user could request any resource when he has permission to get resource",
			userNameInReq: "admin",
			respBody:      "OK",
		},
	}

	for _, item := range data {

		t.Run(item.desc, func(t *testing.T) {
			g := gomega.NewWithT(t)

			ctx = kclient.WithUser(ctx, &user.DefaultInfo{Name: item.userNameInReq})

			_req := httptest.NewRequest("GET", "http://localhost", nil)
			request := restful.NewRequest(_req)

			request.PathParameters()["name"] = "job1"
			request.PathParameters()["namespace"] = "default"

			request.Request = request.Request.WithContext(ctx)
			recorder := httptest.NewRecorder()
			response := restful.NewResponse(recorder)

			filter(request, response, chain)

			if item.respBody != "" {
				g.Expect(recorder.Body.String()).ShouldNot(gomega.Equal(200))
			}
			if item.respCode != 0 {
				g.Expect(recorder.Code).ShouldNot(gomega.Equal(200))
			}
		})
	}

}
