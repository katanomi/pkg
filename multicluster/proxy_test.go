/*
Copyright 2024 The Katanomi Authors.

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

package multicluster

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestClusterProxyHost(t *testing.T) {
	tests := map[string]struct {
		host        string
		endpoint    string
		clusterName string
		expected    string
	}{
		"host with scheme": {
			host:        "http://abc.test",
			endpoint:    "/kubernetes",
			clusterName: "global",
			expected:    "http://abc.test/kubernetes",
		},
		"endpoint without slash": {
			host:        "http://abc.test",
			endpoint:    "kubernetes",
			clusterName: "global",
			expected:    "http://abc.test/kubernetes",
		},
		"host with scheme and port": {
			host:        "https://abc.test:443",
			endpoint:    "/kubernetes",
			clusterName: "global",
			expected:    "https://abc.test:443/kubernetes",
		},
		"host without scheme": {
			host:        "abc.test",
			endpoint:    "/kubernetes",
			clusterName: "global",
			expected:    "https://abc.test/kubernetes",
		},
		"endpoint with name placeholder": {
			host:        "abc.test",
			endpoint:    "/kubernetes/{name}",
			clusterName: "global",
			expected:    "https://abc.test/kubernetes/global",
		},
	}

	for name, item := range tests {
		t.Run(name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			g.Expect(ClusterProxyHost(item.host, item.endpoint, item.clusterName)).To(Equal(item.expected))
		})
	}
}
