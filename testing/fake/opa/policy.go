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

package opa

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/open-policy-agent/opa/rego"
)

// Policy represents a policy object for OPA (Open Policy Agent) evaluation.
// It's designed to work with Rego, the language for writing policies and queries for OPA.
// For more information about OPA and Rego, visit: https://www.openpolicyagent.org/docs/latest/
type Policy struct {
	ID string `json:"id"`
	// Policy contains the actual Rego policy as a string.
	Policy string `json:"policy"`
	// Query holds a set of Queries to be evaluated against the policy.
	Query Queries `json:"Query"`
	// Result stores the results of the policy evaluation.
	Result rego.Vars `json:"-"`
}

// Eval evaluates the policy against a given input, using OPA.
func (p *Policy) Eval(ctx context.Context, input interface{}) error {
	r := rego.New(
		rego.Query(p.Query.Eval()),
		rego.Module("policy.rego", p.Policy))

	eval, err := r.PrepareForEval(ctx)
	if err != nil {
		return err
	}

	rs, err := eval.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return err
	}

	if len(rs) == 1 && rs[0].Bindings != nil {
		p.Result = rs[0].Bindings
	}

	return nil
}

// IntResult extracts a int result for a given Query.
func (p *Policy) IntResult(query Query) int {
	if len(p.Result) == 0 {
		return 0
	}

	value, exist := p.Result[query.Var()]
	if !exist {
		return 0
	}

	switch v := value.(type) {
	case json.Number:
		intVal, _ := v.Int64()
		return int(intVal)
	case string:
		intVal, _ := strconv.Atoi(v)
		return intVal
	default:
		return 0
	}
}

// BoolResult extracts a string result for a given Query.
func (p *Policy) BoolResult(query Query) bool {
	if len(p.Result) == 0 {
		return false
	}

	value, exist := p.Result[query.Var()]
	if !exist {
		return false
	}
	return value.(bool)
}

// StringResult extracts a string result for a given Query.
func (p *Policy) StringResult(query Query) string {
	if len(p.Result) == 0 {
		return ""
	}

	value, exist := p.Result[query.Var()]
	if !exist {
		return ""
	}
	return value.(string)
}

// MapResult extracts a map result for a given Query.
func (p *Policy) MapResult(query Query) map[string]interface{} {
	if len(p.Result) == 0 {
		return nil
	}

	value, exist := p.Result[query.Var()]
	if !exist {
		return nil
	}

	return value.(map[string]interface{})
}
