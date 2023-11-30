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

package fake

import (
	"net/http"
	"strings"

	"github.com/emicklei/go-restful/v3"
	"github.com/katanomi/pkg/plugin/client"
)

// Input is a custom type representing different types of request inputs.
type Input string

const (
	InputAuth  Input = "auth"
	InputMeta  Input = "meta"
	InputQuery Input = "query"
	InputBody  Input = "body"
	InputPath  Input = "path"
)

// Field extends an Input by appending additional fields.
func (i Input) Field(fields ...Input) Input {
	strFields := []string{string(i)}
	for _, field := range fields {
		strFields = append(strFields, string(field))
	}
	return Input(strings.Join(strFields, "."))
}

// InputFromRequest extracts different types of inputs (auth, meta, etc.) from a restful request.
func InputFromRequest(req *restful.Request) (map[Input]interface{}, error) {
	auth, err := client.AuthFromRequest(req)
	if err != nil {
		return nil, err
	}
	meta, err := client.MetaFromRequest(req)
	if err != nil {
		return nil, err
	}
	paths := req.PathParameters()

	query := make(map[string]interface{})
	for key, value := range req.Request.URL.Query() {
		if len(value) == 1 {
			query[key] = value[0]
		} else {
			query[key] = value
		}
	}

	body := make(map[string]interface{})
	switch req.Request.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		if err := req.ReadEntity(&body); err != nil {
			return nil, err
		}
	default:
	}

	input := map[Input]interface{}{
		InputPath:  paths,
		InputQuery: query,
		InputBody:  body,
	}

	if auth != nil {
		input[InputAuth] = map[string]interface{}{
			"type": auth.Type,
			"data": auth.Secret,
		}
	}

	if meta != nil {
		input[InputMeta] = map[string]interface{}{
			"version": meta.Version,
			"baseURL": meta.BaseURL,
		}
	}

	return input, nil
}
