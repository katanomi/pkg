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
	"fmt"
	"strings"
)

// Query represents a string that can be used as a Query in OPA
type Query string

// Var converts a Query to a variable format usable in rego.
func (q Query) Var() string {
	// a.b.c to a_b_c
	return strings.ReplaceAll(string(q), ".", "_")
}

// Eval generates a rego Query string from the Query.
func (q Query) Eval() string {
	queryVar := q.Var()

	return fmt.Sprintf("%s = %s", queryVar, q)
}

type Queries []Query

func (q Queries) Eval() string {
	queries := make([]string, 0, len(q))
	for _, query := range q {
		queries = append(queries, query.Eval())
	}
	return strings.Join(queries, ";")
}
