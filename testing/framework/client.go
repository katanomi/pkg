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

package framework

import (
	"encoding/base64"
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/katanomi/pkg/plugin/client"
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
