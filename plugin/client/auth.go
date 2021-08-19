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
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/emicklei/go-restful/v3"
	"github.com/go-resty/resty/v2"
	"github.com/katanomi/pkg/errors"

	corev1 "k8s.io/api/core/v1"
)

type AuthType string

const (
	AuthTypeBasic  = AuthType(corev1.SecretTypeBasicAuth)
	AuthTypeOauth2 = AuthType("katanomi.dev/oauth2")
)

// Auth plugin auth
// Method auth method, such as basic, oauth2
// Secret a base64 encoded json object, with auth field included
type Auth struct {
	// Type secret type as in kubernetes secret.type
	Type AuthType `json:"type"`
	// Secret 's data value extracted from kubernetes
	Secret map[string][]byte `json:"data"`
}

// ToRequest set request header for resty.Request
func (a *Auth) ToRequest(request *resty.Request) error {
	method, err := a.authMethod()
	if err != nil {
		return err
	}

	method.ToRequest(request)

	return nil
}

func (a *Auth) authMethod() (AuthMethod, error) {
	switch a.Type {
	case AuthTypeBasic:
		return a.Basic()
	case AuthTypeOauth2:
		return a.Oauth2()
	default:
		return nil, fmt.Errorf("no auth method matched for %s", a.Type)
	}
}

type authContextKey struct{}

// ExtractAuth extract auth from a specific context
func ExtractAuth(ctx context.Context) *Auth {
	value := ctx.Value(authContextKey{})
	if v, ok := value.(*Auth); ok {
		return v
	}
	return nil
}

// WithContext returns a copy of parent include with the auth
func (a *Auth) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, authContextKey{}, a)
}

// Basic return a Basic auth struct
func (a *Auth) Basic() (*AuthBasic, error) {
	basic := &AuthBasic{
		Username: string(a.Secret[corev1.BasicAuthUsernameKey]),
		Password: string(a.Secret[corev1.BasicAuthPasswordKey]),
	}
	return basic, nil
}

// Oauth2 return an Oauth2 struct
func (a *Auth) Oauth2() (*AuthOauth2, error) {
	oauth2 := &AuthOauth2{}
	// TODO: needs to add specific const keys to do conversion here

	return oauth2, nil
}

// AuthMethod set request header for resty.Request
type AuthMethod interface {
	ToRequest(request *resty.Request)
}

type AuthBasic struct {
	Username string
	Password string
}

func (a *AuthBasic) ToRequest(request *resty.Request) {
	request.SetBasicAuth(a.Username, a.Password)
}

type AuthOauth2 struct {
	Token        string
	ClientID     string
	ClientSecret string
	RefreshToken string
}

func (a *AuthOauth2) ToRequest(request *resty.Request) {
	//TODO: check token expired and refresh
	request.Header.Set("Authorization", "Bearer "+a.Token)
}

const (
	// PluginAuthHeader header for auth type (kubernetes secret type)
	PluginAuthHeader = "X-Plugin-Auth"
	// PluginSecretHeader header to store data part of the secret
	PluginSecretHeader = "X-Plugin-Secret"
)

// AuthFilter auth filter for go restful, parsing plugin auth
func AuthFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	method := req.HeaderParameter(PluginAuthHeader)
	encodedSecret := req.HeaderParameter(PluginSecretHeader)

	if method == "" || encodedSecret == "" {
		chain.ProcessFilter(req, resp)
		return
	}

	decodedSecret, err := base64.StdEncoding.DecodeString(encodedSecret)
	if err != nil {
		errors.HandleError(req, resp, fmt.Errorf("decode secret error: %s", err.Error()))
		return
	}

	data := map[string][]byte{}
	if err = json.Unmarshal(decodedSecret, &data); err != nil {
		errors.HandleError(req, resp, fmt.Errorf("decode secret error: %s", err.Error()))
		return
	}

	auth := Auth{
		Type:   AuthType(method),
		Secret: data,
	}

	ctx := req.Request.Context()
	req.Request = req.Request.WithContext(auth.WithContext(ctx))

	chain.ProcessFilter(req, resp)
}
