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
	restful "github.com/emicklei/go-restful/v3"
	"knative.dev/pkg/logging"
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
