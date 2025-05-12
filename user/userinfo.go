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
	"fmt"
	"strings"

	"github.com/emicklei/go-restful/v3"
	"github.com/golang-jwt/jwt/v4"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/client"
	authenticationv1 "k8s.io/api/authentication/v1"
	"k8s.io/apimachinery/pkg/util/errors"
)

var hmacSampleSecret []byte

// Parsing the jwt information in the request header
func parseReqjwtToClaims(req *restful.Request) (claims jwt.MapClaims, err error) {

	tokenString := req.HeaderParameter(client.AuthorizationHeader)
	ignoreSigningError := "unexpected signing method: RS256"

	if !strings.HasPrefix(tokenString, client.BearerPrefix) {
		err = fmt.Errorf("token string not has prefix %s", client.BearerPrefix)
		return
	}
	tokenString = tokenString[len(client.BearerPrefix):]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})
	if err != nil {
		if err.Error() != ignoreSigningError {
			// The unexpected signing method test fails and does not affect the parsing
			// This conclusion is reliable
			return claims, err
		}
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	} else {
		return jwt.MapClaims{}, fmt.Errorf("token.Claims is not jwt.MapClaims type")
	}
}

// Get the userinfo from the request object
func getUserInfoFromReq(req *restful.Request) (userinfo metav1alpha1.UserInfo, err error) {
	errorList := make([]error, 0)
	user := client.ImpersonateUser(req.Request)
	if user != nil {
		return metav1alpha1.UserInfo{
			UserInfo: authenticationv1.UserInfo{
				Username: user.GetName(),
				Groups:   user.GetGroups(),
				UID:      user.GetUID(),
			},
		}, nil
	}
	claims, err := parseReqjwtToClaims(req)
	if err != nil {
		return userinfo, err
	}
	userinfo = metav1alpha1.UserInfo{}
	userinfo.FromJWT(claims)
	return userinfo, errors.NewAggregate(errorList)
}

// GetBaseUserInfoFromReq get the base userinfo from the request object
func GetBaseUserInfoFromReq(req *restful.Request) *metav1alpha1.GitUserBaseInfo {
	userinfo, ok := UserInfoFrom(req.Request.Context())
	if !ok {
		return nil
	}
	authorInfo := &metav1alpha1.GitUserBaseInfo{}
	authorInfo.Name = userinfo.GetName()
	authorInfo.Email = userinfo.GetEmail()
	return authorInfo
}
