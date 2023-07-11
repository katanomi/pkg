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

package url

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"knative.dev/pkg/apis"
)

var _ = Describe("Test.MatchGitURLPrefix", func() {

	DescribeTable("MatchGitURLPrefix",
		func(source, target string, expected bool) {
			s, err := apis.ParseURL(source)
			Expect(err).To(BeNil())
			t, err := apis.ParseURL(target)
			Expect(err).To(BeNil())
			actual := MatchGitURLPrefix(s, t)
			Expect(actual).Should(Equal(expected))
		},
		Entry("test suffix .git", "https://github.com/katanomi/pkg.git", "https://github.com/katanomi/pkg.git", true),
		Entry("test suffix .git", "https://github.com/katanomi/pkg.git", "https://github.com/katanomi/pkg", true),
		Entry("test suffix .git", "https://github.com/katanomi/pkg.git", "https://github.com/katanomi/", true),
		Entry("test suffix .git", "https://github.com/katanomi/pkg.git", "https://github.com/", true),
		Entry("test suffix .git", "https://github.com/katanomi/pkg.git", "http://github.com", true),
		Entry("test suffix .git", "https://github.com/katanomi/pkg.git", "github.com", false),
		//
		Entry("test suffix /", "https://github.com/katanomi/pkg.git", "https://github.com/katanomi/pkg/", false),
		Entry("test suffix /", "https://github.com/katanomi/pkg.git", "https://github.com/katanomi/pkg", true),
		Entry("test suffix /", "https://github.com/katanomi/pkg.git", "https://github.com/katanomi/", true),
		Entry("test suffix /", "https://github.com/katanomi/pkg.git", "https://github.com/", true),
		Entry("test suffix /", "https://github.com/katanomi/pkg", "https://github.com/katanomi/", true),
		Entry("test suffix /", "https://github.com/katanomi/pkg", "http://github.com", true),
		//
		Entry("test suffix /", "https://github.com/katanomi/pkg", "https://github.com/katanomi/pkg.git", true),
	)
})

func TestExpandURLIPv6(t *testing.T) {
	tests := map[string]struct {
		rawURL  string
		want    string
		wantErr error
	}{
		"ipv4 url": {
			rawURL: "http://172.26.168.90",
			want:   "http://172.26.168.90",
		},
		"ipv4 port url": {
			rawURL: "http://172.26.168.90:8080",
			want:   "http://172.26.168.90:8080",
		},
		"ipv6 url": {
			rawURL: "http://[20::172:26:168:90]",
			want:   "http://[0020:0000:0000:0000:0172:0026:0168:0090]",
		},
		"ipv6 port url": {
			rawURL: "http://[20::172:26:168:90]:8080",
			want:   "http://[0020:0000:0000:0000:0172:0026:0168:0090]:8080",
		},
		"ipv6 with ipv4 url": {
			rawURL: "http://[::FFFF:192.168.0.1]",
			want:   "http://[0000:0000:0000:0000:0000:ffff:c0a8:0001]",
		},
		"ipv6 with ipv4 port url": {
			rawURL: "http://[::FFFF:192.168.0.1]:8080",
			want:   "http://[0000:0000:0000:0000:0000:ffff:c0a8:0001]:8080",
		},
		"domain url": {
			rawURL: "http://test.com",
			want:   "http://test.com",
		},
		"domain port url": {
			rawURL: "http://test.com:8080",
			want:   "http://test.com:8080",
		},
		"paser failed": {
			rawURL:  "http\n://[20::172::26:168:90]:8080",
			want:    "",
			wantErr: fmt.Errorf("net/url: invalid control character in URL"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			got, err := ExpandURLIPv6(tt.rawURL)
			if tt.wantErr != nil {
				g.Expect(err).NotTo(BeNil())
				g.Expect(err.Error()).To(ContainSubstring(tt.wantErr.Error()))
			} else {
				g.Expect(err).To(BeNil())
			}
			g.Expect(got).To(Equal(tt.want), name)
		})
	}
}
