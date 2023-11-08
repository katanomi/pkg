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

package cluster

import (
	"strings"

	"github.com/katanomi/pkg/testing/framework/base"

	"github.com/katanomi/pkg/multicluster"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TestContext context for cluster testing
type TestContext struct {
	base.TestContext

	Config    *rest.Config
	Scheme    *runtime.Scheme
	Client    client.Client
	Namespace string

	MultiClusterClient multicluster.Interface
}

// TestContextOption options for building TestContext
type TestContextOption func(*TestContext)

// NamespaceOption customize the namespace name
func NamespaceOption(ns string) TestContextOption {
	return func(testCtx *TestContext) {
		testCtx.Namespace = ns
	}
}

// NamespacePrefixOption customize the prefix of the namespace name
func NamespacePrefixOption(prefix string) TestContextOption {
	return func(testCtx *TestContext) {
		testCtx.Namespace = strings.TrimSuffix(prefix, "-") + "-" + rand.String(5)
	}
}

type TestSpecFunc func(*TestContext)
