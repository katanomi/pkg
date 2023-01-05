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
	"net/http"
	"testing"

	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestGetUserInfo(t *testing.T) {

	_req := &http.Request{
		Header: map[string][]string{},
	}
	req := &restful.Request{Request: _req}

	t.Run("normal", func(t *testing.T) {
		// placeholder data generated at https://token.dev/
		tokenString := "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImVtYWlsIjoiam9obkB0ZXN0LmNvbSIsImlhdCI6MTY1MzM3NzM0MCwiZXhwIjoxNjUzMzgwOTQwfQ.tT5YWrFu2F_0NWg3JFbpIaC-HbaBgJtRcff_O3Ti225RZztON1gpZOfFf06EIzWolkUrtAfrTFXonHx2rh5w8mK82kFX5sLPxiV3yDfvxU9KuE66IW_48ykFU6puBcnVMZNWUtPLHcH4OAq51zbca9eh3AFiFldnJ3YRJ85Bigky8nc1uf3CCm5c9d7P8q5SGDJZvGLHOXqQ4_aSfpX0HDTHbgw_82d7FpckeDYRdFN89LORaLRChbBM1IdfH4JvI9WtO3PDU7Ce49nMmmidhsESAalxfMrN3LIPmMz7vY0FJaBW24oA0FwJvq1Q6O_jzupMHRGS-clkUImRw185cA"

		req.Request.Header.Set("Authorization", tokenString)
		userinfo, _ := getUserInfoFromReq(req)
		if userinfo.GetName() != "John Doe" {
			t.Errorf("should have John Doe")
		}
		if userinfo.GetEmail() != "john@test.com" {
			t.Errorf("should have test@test.com but got %s", userinfo.GetEmail())
		}
	})
	t.Run("not found jwt", func(t *testing.T) {
		tokenString := ""
		req.Request.Header.Set("Authorization", tokenString)
		userinfo, _ := getUserInfoFromReq(req)
		if userinfo.GetName() != "" {
			t.Errorf("should have nil string")
		}
		if userinfo.GetEmail() != "" {
			t.Errorf("should have nil string but got %s", userinfo.GetEmail())
		}
	})

}

var _ = Describe("Test GetBaseUserInfoFromReq", func() {
	var (
		tokenString string
		authorInfo  *metav1alpha1.GitUserBaseInfo
	)
	BeforeEach(func() {
		tokenString = ""
		authorInfo = nil
	})
	JustBeforeEach(func() {
		_req := &http.Request{
			Header: map[string][]string{},
		}
		_req.Header.Set("Authorization", tokenString)
		req := &restful.Request{Request: _req}

		res := &restful.Response{}
		target := func(req *restful.Request, resp *restful.Response) {}
		chain := &restful.FilterChain{Target: target}
		UserInfoFilter(req, res, chain)

		authorInfo = GetBaseUserInfoFromReq(req)
	})

	Context("when the request context contains the user info", func() {
		BeforeEach(func() {
			tokenString = "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImVtYWlsIjoiam9obkB0ZXN0LmNvbSIsImlhdCI6MTY1MzM3NzM0MCwiZXhwIjoxNjUzMzgwOTQwfQ.tT5YWrFu2F_0NWg3JFbpIaC-HbaBgJtRcff_O3Ti225RZztON1gpZOfFf06EIzWolkUrtAfrTFXonHx2rh5w8mK82kFX5sLPxiV3yDfvxU9KuE66IW_48ykFU6puBcnVMZNWUtPLHcH4OAq51zbca9eh3AFiFldnJ3YRJ85Bigky8nc1uf3CCm5c9d7P8q5SGDJZvGLHOXqQ4_aSfpX0HDTHbgw_82d7FpckeDYRdFN89LORaLRChbBM1IdfH4JvI9WtO3PDU7Ce49nMmmidhsESAalxfMrN3LIPmMz7vY0FJaBW24oA0FwJvq1Q6O_jzupMHRGS-clkUImRw185cA"
		})
		It("should get the username as excepted", func() {
			gomega.Expect(authorInfo.Name).To(gomega.Equal("John Doe"))
		})
	})
	Context("when the request context dose not contains the user info", func() {
		It("should get nil", func() {
			gomega.Expect(authorInfo).To(gomega.BeNil())
		})
	})
})
