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

package path

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestPathEscape(t *testing.T) {
	g := NewGomegaWithT(t)
	tests := []struct {
		path string
		want string
	}{
		{
			path: "abc",
			want: "abc",
		},
		{
			path: "a/b",
			want: "a%252Fb",
		},
		{
			path: "a/b/c",
			want: "a%252Fb%252Fc",
		},
		{
			path: "a.b",
			want: "a.b",
		},
		{
			path: "a://b",
			want: "a%253A%252F%252Fb",
		},
	}
	for _, tt := range tests {
		if got := Escape(tt.path); got != tt.want {
			g.Expect(got).To(Equal(tt.want))
		}
	}
}

func TestFormatPath(t *testing.T) {
	g := NewGomegaWithT(t)
	tests := []struct {
		tmpl   string
		params []string
		want   string
	}{
		{
			tmpl:   "%s/%s",
			params: []string{"a", "b"},
			want:   "a/b",
		},
		{
			tmpl:   "%s/%s",
			params: []string{"a/b", "c/d"},
			want:   "a%252Fb/c%252Fd",
		},
		{
			tmpl:   "/base/path/%s",
			params: []string{"minio:///path/to/file"},
			want:   "/base/path/minio%253A%252F%252F%252Fpath%252Fto%252Ffile",
		},
	}
	for _, tt := range tests {
		g.Expect(Format(tt.tmpl, tt.params...)).To(Equal(tt.want))
	}
}
