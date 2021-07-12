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
	"fmt"
	"net/http"
	"strings"

	"github.com/emicklei/go-restful/v3"
)

var (
	headerPluginMeta = "X-Plugin-Meta"
	pluginContextKey = struct{}{}
)

// Meta Plugin meta with base url and version info, for calling plugin api
type Meta struct {
	Version string
	BaseURL string
}

// WithContext returns a copy of parent include with the plugin meta
func (p *Meta) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, pluginContextKey, p)
}

// ExtraMeta extract meta from a specific context
func ExtraMeta(ctx context.Context) *Meta {
	value := ctx.Value(pluginContextKey)
	if v, ok := value.(*Meta); ok {
		return v
	}

	return nil
}

// MetaFilter meta filter for go restful, parsing plugin meta
func MetaFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	encodedMeta := req.HeaderParameter(headerPluginMeta)
	decodedMeta, err := base64.StdEncoding.DecodeString(encodedMeta)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, fmt.Errorf("decode meta error: %s", err.Error()))
		return
	}

	metaAttrs := strings.Split(string(decodedMeta), ":")
	if len(metaAttrs) < 2 {
		resp.WriteError(http.StatusBadRequest, fmt.Errorf("invalid plugin meta: %s", decodedMeta))
		return
	}

	meta := &Meta{
		Version: metaAttrs[1],
		BaseURL: metaAttrs[0],
	}

	ctx := req.Request.Context()
	req.Request = req.Request.WithContext(meta.WithContext(ctx))

	chain.ProcessFilter(req, resp)
}
