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
	"testing"

	corev1 "k8s.io/api/core/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"

	"github.com/go-resty/resty/v2"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	ktest "github.com/katanomi/pkg/testing"
	. "github.com/onsi/gomega"

	"github.com/jarcoal/httpmock"
	"knative.dev/pkg/apis"
)

func TestClientArtifactList(t *testing.T) {
	g := NewGomegaWithT(t)
	httpmock.Reset()

	list := &metav1alpha1.ArtifactList{}
	err := ktest.LoadYAML("testdata/artifactlist.yaml", &list)
	g.Expect(err).To(BeNil())

	responder, _ := httpmock.NewJsonResponder(200, &list)

	fakeUrl := "https://example.com/projects/devops/repositories/katanomi/artifacts"
	httpmock.RegisterResponder("GET", fakeUrl, responder)

	mate := Meta{Version: "v1.3.4", BaseURL: "http://plugin.com"}
	secret := corev1.Secret{
		Type: corev1.SecretTypeBasicAuth,
		Data: map[string][]byte{"username": []byte("username")},
	}

	RESTClient := resty.New()
	httpmock.ActivateNonDefault(RESTClient.GetClient())

	client := NewPluginClient(ClientOpts(RESTClient))
	artifactClient := client.Artifact(mate, secret)

	url, _ := apis.ParseURL("https://example.com/")
	artifactList, err := artifactClient.List(
		context.Background(),
		&duckv1.Addressable{URL: url},
		"devops",
		"katanomi",
		client.Secret(corev1.Secret{Type: corev1.SecretTypeBasicAuth, Data: map[string][]byte{"username": []byte("username")}}),
	)

	g.Expect(err).To(BeNil())
	g.Expect(artifactList.TotalItems).To(Equal(2))
	g.Expect(artifactList).To(Equal(list))
}

func TestClientArtifactDelete(t *testing.T) {
	g := NewGomegaWithT(t)
	httpmock.Reset()

	responder := httpmock.NewStringResponder(200, "OK")

	fakeUrl := "https://example.com/projects/devops/repositories/katanomi/artifacts/artifact"
	httpmock.RegisterResponder("DELETE", fakeUrl, responder)

	mate := Meta{Version: "v1.3.4", BaseURL: "http://plugin.com"}
	secret := corev1.Secret{
		Type: corev1.SecretTypeBasicAuth,
		Data: map[string][]byte{"username": []byte("username")},
	}

	RESTClient := resty.New()
	httpmock.ActivateNonDefault(RESTClient.GetClient())

	client := NewPluginClient(ClientOpts(RESTClient))
	artifactClient := client.Artifact(mate, secret)

	url, _ := apis.ParseURL("https://example.com/")
	err := artifactClient.Delete(
		context.Background(),
		&duckv1.Addressable{URL: url},
		"devops",
		"katanomi",
		"artifact",
		client.Secret(corev1.Secret{Type: corev1.SecretTypeBasicAuth, Data: map[string][]byte{"username": []byte("username")}}),
	)

	g.Expect(err).To(BeNil())
}

func TestClientArtifactDeleteTag(t *testing.T) {
	g := NewGomegaWithT(t)
	httpmock.Reset()

	responder := httpmock.NewStringResponder(200, "OK")

	fakeUrl := "https://example.com/projects/devops/repositories/katanomi/artifacts/artifact/tags/tag"
	httpmock.RegisterResponder("DELETE", fakeUrl, responder)

	mate := Meta{Version: "v1.3.4", BaseURL: "http://plugin.com"}
	secret := corev1.Secret{
		Type: corev1.SecretTypeBasicAuth,
		Data: map[string][]byte{"username": []byte("username")},
	}

	RESTClient := resty.New()
	httpmock.ActivateNonDefault(RESTClient.GetClient())

	client := NewPluginClient(ClientOpts(RESTClient))
	artifactClient := client.Artifact(mate, secret)

	url, _ := apis.ParseURL("https://example.com/")
	err := artifactClient.DeleteTag(
		context.Background(),
		&duckv1.Addressable{URL: url},
		"devops",
		"katanomi",
		"artifact",
		"tag",
		client.Secret(corev1.Secret{Type: corev1.SecretTypeBasicAuth, Data: map[string][]byte{"username": []byte("username")}}),
	)

	g.Expect(err).To(BeNil())
}
