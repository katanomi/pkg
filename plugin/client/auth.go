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
	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/errors"

	corev1 "k8s.io/api/core/v1"
)

const (
	OAuth2KeyAccessToken = "accessToken"

	AuthHeaderAuthorization = "Authorization"

	AuthPrefixBearer = "Bearer"
)

// Auth plugin auth
type Auth struct {
	// Type secret type as in kubernetes secret.type
	Type v1alpha1.AuthType `json:"type"`
	// Secret 's data value extracted from kubernetes
	Secret map[string][]byte `json:"data"`
}

type AuthMethod func(request *resty.Request)

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

// FromSecret generate auth from secret
func FromSecret(secret corev1.Secret) *Auth {
	secretType := string(secret.Type)
	if v, exist := secret.Labels[v1alpha1.SecretTypeLabelKey]; exist && v != "" {
		secretType = v
	}

	return &Auth{
		Type:   v1alpha1.AuthType(secretType),
		Secret: secret.Data,
	}
}

// IsBasic check auth is basic
func (a *Auth) IsBasic() bool {
	return a.Type == v1alpha1.AuthTypeBasic
}

// IsOAuth2 check auth is oauth2
func (a *Auth) IsOAuth2() bool {
	return a.Type == v1alpha1.AuthTypeOAuth2
}

// GetBasicInfo get basic auth username and password
func (a *Auth) GetBasicInfo() (userName string, password string, err error) {
	u, err := a.Get(corev1.BasicAuthUsernameKey)
	if err != nil {
		return
	}

	p, err := a.Get(corev1.BasicAuthPasswordKey)
	if err != nil {
		return
	}

	return u, p, nil
}

// GetOAuth2Token get oauth2 access token
func (a *Auth) GetOAuth2Token() (string, error) {
	return a.Get(OAuth2KeyAccessToken)
}

// Get get specific attribute from secret
func (a *Auth) Get(attribute string) (string, error) {
	v, ok := a.Secret[attribute]
	if !ok {
		return "", fmt.Errorf("attribute not found: %s", attribute)
	}

	return string(v), nil
}

// Basic return a Basic auth function
func (a *Auth) Basic() (AuthMethod, error) {
	if a.Type != v1alpha1.AuthTypeBasic {
		return nil, fmt.Errorf("auth type not match, expected: %s, current: %s", v1alpha1.AuthTypeBasic, a.Type)
	}

	userName, password, err := a.GetBasicInfo()
	if err != nil {
		return nil, err
	}

	return func(request *resty.Request) {
		request.SetBasicAuth(userName, password)
	}, nil
}

// OAuth2 return an oauth2 auth method
func (a *Auth) OAuth2() (AuthMethod, error) {
	if a.Type != v1alpha1.AuthTypeOAuth2 {
		return nil, fmt.Errorf("auth type not match, expected: %s, current: %s", v1alpha1.AuthTypeOAuth2, a.Type)
	}

	return a.BearerToken(OAuth2KeyAccessToken)
}

// BearerToken return an bearer token auth method
func (a *Auth) BearerToken(attribute string) (AuthMethod, error) {
	return a.HeaderWithPrefix(attribute, AuthHeaderAuthorization, AuthPrefixBearer)
}

// Header return an auth method which could append to header with specific attribute
func (a *Auth) Header(attribute string, header string) (AuthMethod, error) {
	return a.HeaderWithPrefix(attribute, header, "")
}

// HeaderWithPrefix return an auth method which could append to header with specific attribute and prefix
func (a *Auth) HeaderWithPrefix(attribute string, header string, prefix string) (AuthMethod, error) {
	value, err := a.Get(attribute)
	if err != nil {
		return nil, err
	}

	return func(request *resty.Request) {
		if prefix != "" {
			prefix += " "
		}

		request.Header.Set(header, prefix+value)
	}, nil
}

// Query return an auth method which could append to query with specific attribute
func (a *Auth) Query(attribute string, query string) (AuthMethod, error) {
	value, ok := a.Secret[attribute]
	if !ok {
		return nil, fmt.Errorf("attribute not found: %s", attribute)
	}

	return func(request *resty.Request) {
		request.SetQueryParam(query, string(value))
	}, nil
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
		Type:   v1alpha1.AuthType(method),
		Secret: data,
	}

	ctx := req.Request.Context()
	req.Request = req.Request.WithContext(auth.WithContext(ctx))

	chain.ProcessFilter(req, resp)
}
