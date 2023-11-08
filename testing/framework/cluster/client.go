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

package cluster

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/go-resty/resty/v2"
	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
	"github.com/katanomi/pkg/plugin/route"
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// RequestWrapAuthMeta wrap auth and meta information into http request
func RequestWrapAuthMeta(req *resty.Request, auth client.Auth, meta client.Meta) *resty.Request {
	if req == nil || req.Header == nil {
		return nil
	}

	req.SetHeader(client.PluginAuthHeader, string(auth.Type))
	authData, _ := json.Marshal(auth.Secret)
	req.SetHeader(client.PluginSecretHeader, base64.StdEncoding.EncodeToString(authData))

	client.MetaOpts(meta)(req)
	return req
}

// ConvertToAuthSecret convert auth
func ConvertToAuthSecret(auth client.Auth) (secret corev1.Secret) {
	secret = corev1.Secret{}
	secret.Annotations = map[string]string{
		v1alpha1.SecretTypeAnnotationKey: string(auth.Type),
	}
	secret.Data = auth.Secret
	return secret
}

// GetPluginAddress get the base url of the plugin
func GetPluginAddress(client *resty.Client, plugin client.Interface) (address *duckv1.Addressable, err error) {
	if client == nil || plugin == nil {
		return nil, errors.New("get address of plugin failed, param error")
	}
	var url *apis.URL
	url, err = apis.ParseURL(client.HostURL)
	if err != nil {
		return
	}
	url.Path = route.GetPluginWebPath(plugin)
	address = &duckv1.Addressable{URL: url}
	return address, nil
}

// MustGetPluginAddress get the base url or panics if the parse fails.
func MustGetPluginAddress(client *resty.Client, plugin client.Interface) (address *duckv1.Addressable) {
	address, err := GetPluginAddress(client, plugin)
	if err != nil {
		panic(err)
	}
	return address
}
