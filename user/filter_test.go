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

	restful "github.com/emicklei/go-restful/v3"
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
}
