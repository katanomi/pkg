/*
Copyright 2022 The Katanomi Authors.

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
	"testing"

	. "github.com/onsi/gomega"
)

func TestUrlToProjectID(t *testing.T) {
	g := NewGomegaWithT(t)
	tests := []struct {
		gitURL  string
		want    string
		wantErr bool
	}{
		{
			gitURL:  "http://github.com/katanomi/catalog?abc",
			want:    "github.com-katanomi-catalog",
			wantErr: false,
		},
		{
			gitURL:  "http://127.0.0.1:8080/katanomi/catalog?abc",
			want:    "127.0.0.1:8080-katanomi-catalog",
			wantErr: false,
		},
		{
			gitURL:  "https://127.0.0.1:8080/katanomi/catalog.git",
			want:    "127.0.0.1:8080-katanomi-catalog",
			wantErr: false,
		},
		{
			gitURL:  "https://127.0.0.1:8080/katanomi/catalog.git",
			want:    "127.0.0.1:8080-katanomi-catalog",
			wantErr: false,
		},
		{
			gitURL:  "ssh://git@192.168.130.62:31211/katanomi/catalog.git",
			want:    "192.168.130.62:31211-katanomi-catalog",
			wantErr: false,
		},
		{
			gitURL:  "git@github.com/katanomi/catalog.git",
			want:    "github.com-katanomi-catalog",
			wantErr: false,
		},
		{
			gitURL:  "ssh://git@github.com/katanomi/catalog.git",
			want:    "github.com-katanomi-catalog",
			wantErr: false,
		},
		{
			gitURL:  "http://[2004::192:168:139:4]:32078/root/image",
			want:    "2004::192:168:139:4-:32078-root-image",
			wantErr: false,
		},
		{
			gitURL:  "ssh://git@[2004::192:168:139:4]:32078/root/image",
			want:    "2004::192:168:139:4-:32078-root-image",
			wantErr: false,
		},
		{
			gitURL:  "git@[2004::192:168:139:4]:32078/root/image",
			want:    "2004::192:168:139:4-:32078-root-image",
			wantErr: false,
		},
		{
			gitURL:  "[2004::4]:80/root/image",
			want:    "2004::4-:80-root-image",
			wantErr: false,
		},
		{
			gitURL:  "[2004::4:80]/root/image",
			want:    "2004::4:80-root-image",
			wantErr: false,
		},
		{
			gitURL:  "git@[2004::192:168:139:4]/root/image",
			want:    "2004::192:168:139:4-root-image",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := UrlToProjectID(tt.gitURL)
		if tt.wantErr {
			g.Expect(err).ShouldNot(Succeed())
		} else {
			g.Expect(err).Should(Succeed())
		}
		g.Expect(got).To(Equal(tt.want))
	}
}
