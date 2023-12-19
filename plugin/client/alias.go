/*
Copyright 2023 The Katanomi Authors.

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

// Package client contains functions to add and retrieve auth from context
package client

import "github.com/katanomi/pkg/plugin/client/base"

const (
	OAuth2KeyAccessToken    = base.OAuth2KeyAccessToken
	AuthHeaderAuthorization = base.AuthHeaderAuthorization
	AuthPrefixBearer        = base.AuthPrefixBearer
	// OAuth2ClientIDKey is the key of the clientID for AuthTypeOAuth2 secrets
	OAuth2ClientIDKey = base.OAuth2ClientIDKey
	// OAuth2ClientSecretKey is the key of the clientSecret for AuthTypeOAuth2 secrets
	OAuth2ClientSecretKey = base.OAuth2ClientSecretKey
	// OAuth2CodeKey is the key of the code for AuthTypeOAuth2 secrets
	OAuth2CodeKey = base.OAuth2CodeKey
	// OAuth2AccessTokenKeyKey is the key of the accessTokenKey for AuthTypeOAuth2 secrets
	OAuth2AccessTokenKeyKey = base.OAuth2AccessTokenKeyKey
	// OAuth2AccessTokenKey is the key of the accessToken for AuthTypeOAuth2 secrets
	OAuth2AccessTokenKey = base.OAuth2AccessTokenKey
	// OAuth2ScopeKey is the key of the scope for AuthTypeOAuth2 secrets
	OAuth2ScopeKey = base.OAuth2ScopeKey
	// OAuth2RefreshTokenKey is the key of the refreshToken for AuthTypeOAuth2 secrets
	OAuth2RefreshTokenKey = base.OAuth2RefreshTokenKey
	// OAuth2ExpiresInKey is the key of the expiresIn for AuthTypeOAuth2 secrets
	OAuth2CreatedAtKey = base.OAuth2CreatedAtKey
	// OAuth2ExpiresInKey is the key of the expiresIn for AuthTypeOAuth2 secrets
	OAuth2ExpiresInKey = base.OAuth2ExpiresInKey
	// OAuth2RedirectURLKey is the key of the redirectURL for AuthTypeOAuth2 secrets
	OAuth2RedirectURLKey = base.OAuth2RedirectURLKey
	// OAuth2BaseURLKey is the key of the baseURL for AuthTypeOAuth2 secrets
	OAuth2BaseURLKey = base.OAuth2BaseURLKey

	// DynamicUsernameKey is the key of the username for  dynamic secrets.
	DynamicUsernameKey = base.DynamicUsernameKey
	// DynamicPasswordKey is the key of the password for  dynamic secrets.
	DynamicPasswordKey = base.DynamicPasswordKey

	// DynamicClientKeyKey is the key of the clientKey for dynamic secret
	DynamicClientKeyKey = base.DynamicClientKeyKey
	// DynamicClientSecretKey is the key of the clientSecret for dynamic secret
	DynamicClientSecretKey = base.DynamicClientSecretKey
	// redefine key for dynamic token refresh.
	DynamicAccessTokenKey  = base.DynamicAccessTokenKey
	DynamicRefreshTokenKey = base.DynamicRefreshTokenKey
	DynamicCreatedAtKey    = base.DynamicCreatedAtKey
	DynamicExpiresInKey    = base.DynamicExpiresInKey
	DynamicBaseURLKey      = base.OAuth2ClientIDKey
)

type Auth = base.Auth

type AuthMethod = base.AuthMethod

var ExtractAuth = base.ExtractAuth

var FromSecret = base.FromSecret

const (
	// PluginAuthHeader header for auth type (kubernetes secret type)
	PluginAuthHeader = base.PluginAuthHeader
	// PluginSecretHeader header to store data part of the secret
	PluginSecretHeader = base.PluginSecretHeader
)

// AuthFilter auth filter for go restful, parsing plugin auth
var AuthFilter = base.AuthFilter
var AuthFromRequest = base.AuthFromRequest

// IsNotImplementedError returns true if the plugin not implement the specified interface
var IsNotImplementedError = base.IsNotImplementedError

// ResponseStatusErr is an error with `Status` type,
// used to handle plugin response
type ResponseStatusErr = base.ResponseStatusErr

const (
	// PluginMetaHeader header to store metadata for the plugin
	PluginMetaHeader = base.PluginMetaHeader

	// PluginSubresourcesHeader subresources header parameter
	// used as a header to avoid overloading the url query parameters
	// and any url length limits
	PluginSubresourcesHeader = base.PluginSubresourcesHeader
)

// Meta Plugin meta with base url and version info, for calling plugin api
type Meta = base.Meta

// ExtraMeta extract meta from a specific context
var ExtraMeta = base.ExtraMeta

// MetaFilter meta filter for go restful, parsing plugin meta
var MetaFilter = base.MetaFilter

var MetaFromRequest = base.MetaFromRequest

// OptionFunc options for requests
type OptionFunc = base.OptionFunc

// SecretOpts provides a secret to be assigned to the request in the header
var SecretOpts = base.SecretOpts

// MetaOpts provides metadata for the request
var MetaOpts = base.MetaOpts

// ListOpts options for lists
var ListOpts = base.ListOpts

// SubResourcesOpts set subresources header for the request
var SubResourcesOpts = base.SubResourcesOpts

// QueryOpts query parameters for the request
var QueryOpts = base.QueryOpts

// BodyOpts request body
var BodyOpts = base.BodyOpts

// ResultOpts request result automatically marshalled into object
var ResultOpts = base.ResultOpts

// DoNotParseResponseOpts do not parse response
var DoNotParseResponseOpts = base.DoNotParseResponseOpts

// ErrorOpts error response object
var ErrorOpts = base.ErrorOpts

// HeaderOpts sets a header
var HeaderOpts = base.HeaderOpts

// BuildOptions Options to build the plugin client
type BuildOptions = base.BuildOptions

// ClientOpts adds a custom client build options for plugin client
var ClientOpts = base.ClientOpts

// DefaultOptions for default plugin client options
var DefaultOptions = base.DefaultOptions

// GetSubResourcesOptionsFromRequest returns SubResourcesOptions based on a request
var GetSubResourcesOptionsFromRequest = base.GetSubResourcesOptionsFromRequest
