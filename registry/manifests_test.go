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

package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/regclient/regclient"
	"github.com/regclient/regclient/config"
	"github.com/regclient/regclient/types/docker/schema2"
)

func TestManifestClient_GetAnnotations(t *testing.T) {
	tests := []struct {
		repo        string
		tag         string
		annotations map[string]string
	}{
		{
			repo:        "abc",
			tag:         "def",
			annotations: nil,
		},
		{
			repo: "abc",
			tag:  "def",
			annotations: map[string]string{
				"abc": "def",
			},
		},
		{
			repo: "abc/def",
			tag:  "xyz",
			annotations: map[string]string{
				"abc/def": "xyz",
			},
		},
	}

	for i, item := range tests {
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			var (
				m = schema2.ManifestList{
					Versioned:   schema2.ManifestListSchemaVersion,
					Annotations: item.annotations,
				}

				g = NewGomegaWithT(t)
			)

			handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				path := fmt.Sprintf("v2/%s/manifests/%s", item.repo, item.tag)
				if strings.Contains(request.URL.String(), path) {
					resp, _ := json.Marshal(m)
					writer.Header().Set("Content-Type", m.MediaType)
					writer.Write(resp)
				}
			})

			server := httptest.NewUnstartedServer(handler)
			server.Start()
			defer server.Close()

			host := config.HostNewName(server.URL)
			client := NewManifestClient(regclient.WithConfigHost(*host))

			image := fmt.Sprintf("%s/%s:%s", host.Name, item.repo, item.tag)
			annotations, err := client.GetAnnotations(context.TODO(), image)

			g.Expect(err).To(BeNil())
			g.Expect(annotations).To(Equal(item.annotations))
		})
	}
}

func TestManifestClient_Insecure(t *testing.T) {
	tests := []struct {
		repo        string
		tag         string
		annotations map[string]string
		insecure    bool
		error       bool
	}{
		{
			repo: "abc",
			tag:  "def",
			annotations: map[string]string{
				"abc": "def",
			},
			insecure: true,
			error:    false,
		},
		{
			repo: "abc",
			tag:  "def",
			annotations: map[string]string{
				"abc": "def",
			},
			insecure: false,
			error:    true,
		},
	}

	for i, item := range tests {
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			var (
				m = schema2.ManifestList{
					Versioned:   schema2.ManifestListSchemaVersion,
					Annotations: item.annotations,
				}
				g = NewGomegaWithT(t)
			)
			handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				path := fmt.Sprintf("v2/%s/manifests/%s", item.repo, item.tag)
				if strings.Contains(request.URL.String(), path) {
					resp, _ := json.Marshal(m)
					writer.Header().Set("Content-Type", m.MediaType)
					writer.Write(resp)
				}
			})

			server := httptest.NewUnstartedServer(handler)
			server.StartTLS()
			defer server.Close()

			host := config.HostNewName(server.URL)
			client := NewManifestClient(regclient.WithRetryLimit(1))
			client.Insecure(item.insecure)

			image := fmt.Sprintf("%s/%s:%s", host.Name, item.repo, item.tag)
			annotations, err := client.GetAnnotations(context.TODO(), image)
			if item.error {
				g.Expect(err).NotTo(BeNil())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(annotations).To(Equal(item.annotations))
			}
		})
	}
}

func TestManifestClient_PutEmptyIndex(t *testing.T) {
	tests := []struct {
		repo        string
		tag         string
		annotations map[string]string
	}{
		{
			repo:        "abc",
			tag:         "def",
			annotations: nil,
		},

		{
			repo: "abc",
			tag:  "def",
			annotations: map[string]string{
				"abc": "def",
			},
		},
		{
			repo: "abc/def",
			tag:  "xyz",
			annotations: map[string]string{
				"abc/def": "xyz",
			},
		},
	}

	for i, item := range tests {
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			var (
				body map[string]map[string]string
				g    = NewGomegaWithT(t)
			)
			handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				path := fmt.Sprintf("v2/%s/manifests/%s", item.repo, item.tag)
				if strings.Contains(request.URL.String(), path) {
					resp, _ := io.ReadAll(request.Body)
					_ = json.Unmarshal(resp, &body)

					writer.WriteHeader(201)
				}
			})

			server := httptest.NewUnstartedServer(handler)
			server.Start()
			defer server.Close()

			host := config.HostNewName(server.URL)

			client := NewManifestClient(regclient.WithConfigHost(*host))

			image := fmt.Sprintf("%s/%s:%s", host.Name, item.repo, item.tag)
			err := client.PutEmptyIndex(context.TODO(), image, item.annotations)
			g.Expect(err).To(BeNil())
			g.Expect(body["annotations"]).To(Equal(item.annotations))
		})
	}
}
