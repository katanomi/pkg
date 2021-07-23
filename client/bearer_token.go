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
	"strings"

	"github.com/emicklei/go-restful/v3"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

const (

	// UserConfigName configuration/context for user
	UserConfigName = "UserConfig"
	// AuthorizationHeader authorization header for http requests
	AuthorizationHeader = "Authorization"
	// BearerPrefix bearer token prefix for token
	BearerPrefix = "Bearer "

	// QueryParameterTokenName authorization token for http requests
	QueryParameterTokenName = "token"
)

// FromBearerToken  retrieves config based on the bearer token
func FromBearerToken(req *restful.Request, baseConfig GetBaseConfigFunc) (config *rest.Config, err error) {
	if config, err = baseConfig(); err != nil {
		return
	}
	token := GetToken(req)
	cmd := buildCmdConfig(&api.AuthInfo{Token: token}, config)
	config, err = cmd.ClientConfig()
	return
}

// GetToken get token from request headers or request query parameters.
// return emtry if no token find
func GetToken(req *restful.Request) (token string) {
	authHeader := req.HeaderParameter(AuthorizationHeader)

	if authHeader != "" && strings.HasPrefix(authHeader, BearerPrefix) && strings.TrimPrefix(authHeader, BearerPrefix) != "" {
		token = strings.TrimPrefix(authHeader, BearerPrefix)
		return
	}

	token = req.QueryParameter(QueryParameterTokenName)
	return
}

func buildCmdConfig(authInfo *api.AuthInfo, cfg *rest.Config) clientcmd.ClientConfig {
	cmdCfg := api.NewConfig()
	cmdCfg.Clusters[UserConfigName] = &api.Cluster{
		Server:                   cfg.Host,
		CertificateAuthority:     cfg.TLSClientConfig.CAFile,
		CertificateAuthorityData: cfg.TLSClientConfig.CAData,
		InsecureSkipTLSVerify:    cfg.TLSClientConfig.Insecure,
	}
	cmdCfg.AuthInfos[UserConfigName] = authInfo
	cmdCfg.Contexts[UserConfigName] = &api.Context{
		Cluster:  UserConfigName,
		AuthInfo: UserConfigName,
	}
	cmdCfg.CurrentContext = UserConfigName

	return clientcmd.NewDefaultClientConfig(
		*cmdCfg,
		&clientcmd.ConfigOverrides{},
	)
}
