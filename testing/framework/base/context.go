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

package base

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	"go.uber.org/zap"
)

// TestContext the context of a test case
type TestContext struct {
	Context context.Context
	*zap.SugaredLogger
}

// TestContextGetter interface to get test context
type TestContextGetter interface {
	// GetTestContext get test case context
	GetTestContext() *TestContext
}

// TestContextGetterFunc a function used to generate an implementation of TestContextGetter
type TestContextGetterFunc func() *TestContext

// GetTestContext get test case context
func (c TestContextGetterFunc) GetTestContext() *TestContext {
	return c()
}

// TestSpecFunc function used as describe
type TestSpecFunc func(testContext *TestContext)

var contextLabelKey = struct{}{}

// WithContextLabel inject labels into context
func WithContextLabel(ctx context.Context, labels Labels) context.Context {
	return context.WithValue(ctx, contextLabelKey, labels)
}

// ContextLabel get labels from context
func ContextLabel(ctx context.Context) Labels {
	val := ctx.Value(contextLabelKey)
	if val == nil {
		return nil
	}
	return val.(Labels)
}
