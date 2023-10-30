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
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	artifactv1 "github.com/katanomi/pkg/apis/artifacts/v1alpha1"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

func mockRegistry(t *testing.T) *httptest.Server {
	digest := ""
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodHead:
			if accept := r.Header.Get("Accept"); !strings.Contains(accept, ocispec.MediaTypeImageIndex) {
				t.Logf("manifest not convertable: %s", accept)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if strings.HasPrefix(r.URL.Path, "/v2/failed/manifests/") {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			digest = filepath.Base(r.URL.Path)
			w.Header().Set("Content-Type", ocispec.MediaTypeImageIndex)
			w.Header().Set("Docker-Content-Digest", digest)
			w.Header().Set("Content-Length", strconv.Itoa(int(0)))
			return
		case r.Method == http.MethodPut && r.URL.Path == "/v2/repo/manifests/tag":
			if contentType := r.Header.Get("Content-Type"); contentType != ocispec.MediaTypeImageManifest {
				t.Logf("manifest contentType: %s", contentType)
				w.WriteHeader(http.StatusBadRequest)
				break
			}
			buf := bytes.NewBuffer(nil)
			if _, err := buf.ReadFrom(r.Body); err != nil {
				t.Logf("fail to read: %v", err)
			}

			w.Header().Set("Docker-Content-Digest", digest)
			w.WriteHeader(http.StatusCreated)
			return
		default:
			t.Logf("not mock url: %s %s", r.Method, r.URL)
			w.WriteHeader(http.StatusForbidden)
		}
		t.Errorf("unexpected access: %s %s", r.Method, r.URL)
	}))
	ts.StartTLS()
	return ts
}

func TestPushEmptyImage(t *testing.T) {
	tests := map[string]struct {
		uri     artifactv1.URI
		wantErr bool
	}{
		"input is not satisfied": {
			uri:     artifactv1.URI{Path: "repo"},
			wantErr: true,
		},
		"invalid repo": {
			uri:     artifactv1.URI{Path: "rep:^o", Tag: "tag"},
			wantErr: true,
		},
		"normal push": {
			uri: artifactv1.URI{Path: "repo", Tag: "tag"},
		},
		"push failed": {
			uri:     artifactv1.URI{Path: "failed", Tag: "tag"},
			wantErr: true,
		},
	}
	ts := mockRegistry(t)
	defer ts.Close()
	uri, _ := url.Parse(ts.URL)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.uri.Host = uri.Host
			if err := PushEmptyImage(context.Background(), tt.uri); (err != nil) != tt.wantErr {
				t.Errorf("PushEmptyImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
