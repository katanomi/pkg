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

package v2

import (
	"context"

	"github.com/jarcoal/httpmock"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var artifactRepoOption = metav1alpha1.RepositoryOptions{
	Project: "test-project",
}
var artifactOption = metav1alpha1.ArtifactOptions{
	RepositoryOptions: artifactRepoOption,
	Repository:        "test-repo",
	Artifact:          "test-artifact",
}

var _ = Describe("Test ListArtifacts", func() {
	It("should generate the correct url and got expected response", func() {
		expected := fakeStruct[metav1alpha1.ArtifactList]()
		httpmock.RegisterResponder(
			"GET",
			"https://example.com/projects/test-project/repositories/test-repo/artifacts?itemsPerPage=3&page=2",
			httpmock.NewJsonResponderOrPanic(200, expected),
		)

		artifacts, err := pluginClient.ListArtifacts(context.Background(), artifactOption, listOption)
		Expect(err).To(Succeed())
		Expect(diff(artifacts, expected)).To(BeEmpty())
	})
})

var _ = Describe("Test GetArtifact", func() {
	It("should generate the correct url and got expected response", func() {
		expected := fakeStruct[metav1alpha1.Artifact]()
		httpmock.RegisterResponder(
			"GET",
			"https://example.com/projects/test-project/repositories/test-repo/artifacts/test-artifact",
			httpmock.NewJsonResponderOrPanic(200, expected),
		)

		got, err := pluginClient.GetArtifact(context.Background(), artifactOption)
		Expect(err).To(Succeed())
		Expect(diff(got, expected)).To(BeEmpty())
	})
})

var _ = Describe("Test DeleteArtifact", func() {
	It("should generate the correct url and got expected response", func() {
		httpmock.RegisterResponder(
			"DELETE",
			"https://example.com/projects/test-project/repositories/test-repo/artifacts/test-artifact",
			httpmock.NewJsonResponderOrPanic(200, nil),
		)

		err := pluginClient.DeleteArtifact(context.Background(), artifactOption)
		Expect(err).To(Succeed())
	})
})
