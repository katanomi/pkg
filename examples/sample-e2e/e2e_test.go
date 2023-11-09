//go:build e2e
// +build e2e

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

package e2e

import (
	"testing"

	_ "github.com/katanomi/pkg/examples/sample-e2e/another"
	"github.com/katanomi/pkg/testing/framework"
	"github.com/katanomi/pkg/testing/framework/cluster"
	"k8s.io/client-go/kubernetes/scheme"
)

var fmw = framework.New("sample-e2e")

func TestMain(m *testing.M) {
	fmw.SynchronizedBeforeSuite(nil).
		SynchronizedAfterSuite(nil).
		Extensions(cluster.ShareScheme(scheme.Scheme)).
		MRun(m)
}

func TestE2E(t *testing.T) {
	// start step to run e2e
	fmw.Run(t)
}
