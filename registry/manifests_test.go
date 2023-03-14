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
	"github.com/google/go-cmp/cmp"
	"github.com/regclient/regclient"
	"github.com/regclient/regclient/config"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
			m := schema2.ManifestList{
				Versioned:   schema2.ManifestListSchemaVersion,
				Annotations: item.annotations,
			}
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

			host := strings.TrimPrefix(server.URL, "http://")

			client := NewManifestClient(regclient.WithConfigHost(config.Host{
				Name:     host,
				Hostname: host,
				TLS:      config.TLSDisabled,
			}))

			image := fmt.Sprintf("%s/%s:%s", host, item.repo, item.tag)
			annotations, err := client.GetAnnotations(context.TODO(), image)

			if err != nil {
				t.Errorf("get annotations error: %s", err.Error())
			}

			if diff := cmp.Diff(item.annotations, annotations); diff != "" {
				t.Errorf("diff: %s", diff)
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
			var body map[string]map[string]string
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

			host := strings.TrimPrefix(server.URL, "http://")

			client := NewManifestClient(regclient.WithConfigHost(config.Host{
				Name:     host,
				Hostname: host,
				TLS:      config.TLSDisabled,
			}))

			image := fmt.Sprintf("%s/%s:%s", host, item.repo, item.tag)
			err := client.PutEmptyIndex(context.TODO(), image, item.annotations)

			if err != nil {
				t.Errorf("put annotations error: %s", err.Error())
			}

			if diff := cmp.Diff(item.annotations, body["annotations"]); diff != "" {
				t.Errorf("diff: %s", diff)
			}
		})
	}
}
