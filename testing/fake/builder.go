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
	"bytes"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"net/http"
	"strings"
	"text/template"

	"github.com/emicklei/go-restful/v3"
	"github.com/go-resty/resty/v2"
	"github.com/katanomi/pkg/testing/fake/opa"

	_ "embed"
)

// Query variables defines in template.rego
const (
	resultQuery  opa.Query = "data.fake.result"
	statusQuery  opa.Query = "data.fake.status"
	matchedQuery opa.Query = "data.fake.matched"
)

// Builder is a struct for building fake policies.
// It holds information about the method, path, input, value, status, and result of a policy.
type Builder struct {
	method string
	path   string

	when   string
	status int
	result string
}

func NewPolicyBuilder(method string, path string) *Builder {
	return &Builder{
		method: method,
		path:   path,
	}
}

// When sets the condition for the policy rule based on the input and value.
// Supported input: InputAuth | InputMeta | InputQuery | InputBody | InputPath
func (b *Builder) When(input Input, value interface{}) *Builder {
	var when interface{}
	switch v := value.(type) {
	case string:
		when = fmt.Sprintf(`"%s"`, v)
	default:
		when = value
	}

	b.when = fmt.Sprintf("input.%s == %v", input, when)

	return b
}

// Status sets the HTTP status code expected in the policy rule.
func (b *Builder) Status(status int) *Builder {
	b.status = status
	return b
}

// Result sets the HTTP response body expected in the policy rule.
func (b *Builder) Result(result interface{}) *Builder {
	if v, ok := result.([]byte); ok {
		result = string(v)
	}

	switch v := result.(type) {
	// result type is Input, uses the specific value at the key from input
	case Input:
		b.result = fmt.Sprintf("input.%s", string(v))
	case string:
		m := map[string]interface{}{}
		if err := json.Unmarshal([]byte(v), &m); err == nil {
			b.result = v
		} else {
			resultBytes, _ := json.Marshal(result)
			b.result = string(resultBytes)
		}
	// others, use the original result
	default:
		resultBytes, _ := json.Marshal(result)
		b.result = string(resultBytes)
	}

	return b
}

type templateData struct {
	When   string `json:"when"`
	Result string `json:"result"`
	Status int    `json:"status"`
}

//go:embed template.rego
var tmpl string

// Complete finalizes the policy construction and returns the policy.
// It uses a template to generate the policy based on the set conditions.
func (b *Builder) Complete() (*opa.Policy, error) {
	policy := &opa.Policy{
		Method: b.method,
		Path:   b.path,
		ID:     IDFromMethodPath(b.method, b.path),
	}

	data := templateData{
		When:   b.when,
		Result: b.result,
		Status: b.status,
	}

	t, err := template.New("rego").Parse(tmpl)
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	err = t.Execute(&buffer, data)
	if err != nil {
		return nil, err
	}

	policy.Query = opa.Queries{resultQuery, statusQuery, matchedQuery}
	policy.Policy = buffer.String()

	return policy, nil
}

// IDFromMethodPath generates a unique ID for the policy based on the method and path.
func IDFromMethodPath(method string, path string) string {
	path = strings.TrimPrefix(path, "/")
	f := fnv.New32()
	_, _ = f.Write([]byte(fmt.Sprintf("%s.%s", method, path)))

	return fmt.Sprintf("%x", f.Sum32())
}

// IDFromRequest generates a unique ID for the policy based on request
func IDFromRequest(req *restful.Request) string {
	return IDFromMethodPath(req.Request.Method, req.Request.URL.Path)
}

// Create sends the constructed policy to the specified client.
func Create(client *resty.Client, builder *Builder) error {
	policy, err := builder.Complete()
	if err != nil {
		return err
	}

	resp, err := client.R().SetBody(policy).Post("mock/policy")
	if err != nil {
		return err
	}

	if resp.StatusCode() >= http.StatusBadRequest {
		return fmt.Errorf("create opa policy failed: %s, status: %d", resp.String(), resp.StatusCode())
	}

	return nil
}
