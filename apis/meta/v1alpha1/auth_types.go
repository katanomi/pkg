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

package v1alpha1

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AuthType string

const (
	AuthTypeBasic   = AuthType(corev1.SecretTypeBasicAuth)
	AuthTypeOAuth2  = AuthType("katanomi.dev/oauth2")
	AuthTypeDynamic = AuthType("katanomi.dev/dynamic")
)

var AuthCheckGVK = GroupVersion.WithKind("AuthCheck")

// AuthCheckOptions options for AuthCheck
// most of the necessary data is already embed in headers
type AuthCheckOptions struct {
	RedirectURL string `json:"redirectURL"`
}

// AuthCheck consists of result for an auth check request
type AuthCheck struct {
	metav1.TypeMeta `json:",inline"`
	Spec            *AuthCheckSpec  `json:"spec,omitempty"`
	Status          AuthCheckStatus `json:"status"`
}

type AuthCheckSpec struct {
	RedirectURL          string                  `json:"redirectURL"`
	Secret               *corev1.ObjectReference `json:"secretRef,omitempty"`
	IntegrationClassName string                  `json:"integrationClassName,omitempty"`
	BaseURL              string                  `json:"baseURL,omitempty"`
	Version              string                  `json:"version,omitempty"`
}

type AuthCheckStatus struct {
	// Allowed describes if the headers used where accepted or not by the integrated system.
	// `True` when accepted,
	// `False` when not accepted or when secret data is incomplete (oAuth2 flow),
	// `Unknown` for a standard `not implemented` response.
	Allowed corev1.ConditionStatus `json:"allowed"`

	// Message a message explaining the response content
	Message string `json:"message,omitempty"`

	// Reason specific reason enum for the given status
	Reason string `json:"reason,omitempty"`

	// RedirectURL the provided redirect url for oauthFlow
	// +optional
	RedirectURL string `json:"redirectURL,omitempty"`
	// AuthorizeURL url for the oAuth2 flow
	// +optional
	AuthorizeURL string `json:"authorizeURL,omitempty"`

	// Provides the ability to convert username and password, and token refresh.
	// +optional
	RefreshData *AuthTokenStatus `json:"refreshData,omitempty"`
}

const (
	NotAllowedAuthCheckReason         = "NotAllowed"
	NotImplementedAuthCheckReason     = "NotImplemented"
	NeedsAuthorizationAuthCheckReason = "NeedsAuthorization"
)

// +k8s:deepcopy-gen=false
// AuthToken access token request response
type AuthToken struct {
	metav1.TypeMeta `json:",inline"`
	Status          AuthTokenStatus `json:"status"`
}

// +k8s:deepcopy-gen=false
// AuthTokenStatus access token request response status
type AuthTokenStatus struct {
	// AccessTokenKey store the key for accessToken it is mainly for git clone as userName
	AccessTokenKey string `json:"accessTokenKey"`
	// AccessToken access token used for oAuth2
	AccessToken string `json:"accessToken"`
	// RefreshToken to renew access token when expired
	RefreshToken string `json:"refreshToken,omitempty"`
	// TokenType provides a token type for consumers, usually is a bearer type
	TokenType string `json:"tokenType"`
	// ExpiresIn token expiration duration
	ExpiresIn time.Duration `json:"expiresIn"`
	// CreatedAt time which access token was generated.
	// Used with ExpiresIn to calculate expiration time
	CreatedAt time.Time `json:"createdAt"`
}
