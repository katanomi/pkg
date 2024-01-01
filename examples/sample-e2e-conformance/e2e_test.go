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

package sample_e2e_conformance

import (
	"testing"

	"github.com/katanomi/pkg/examples/sample-e2e-conformance/auth"
	"github.com/katanomi/pkg/testing/framework"
	. "github.com/katanomi/pkg/testing/framework/conformance"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var fmw = framework.New("business-e2e")

func TestMain(m *testing.M) {
	fmw.Config(Configure).MRun(m)
}

func TestConformance(t *testing.T) {
	fmw.Run(t)
}

var _ = func() bool {
	m := NewModuleCase("ClusterIntegration")
	authFeature := NewFeatureCase("Integration")
	m.AddFeatureCase(
		authFeature.AddTestCaseSet(
			auth.CaseSet.New().Focus(
				auth.TestPointBasicAuth.Bind(authFeature).AddAssertion(func(statusCode int) {
					GinkgoWriter.Println("run customize assertion")
					Expect(statusCode >= 200).To(BeTrue())
				}),
				auth.TestPointOauth2,
			),
		),
		NewFeatureCase("feature2",
			auth.CaseSet.New().Focus(
				auth.TestPointOauth2,
			),
		),
	)
	m.RegisterTestCase()
	return false
}()
