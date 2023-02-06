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

package v1alpha1

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestPager_GetLimit(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	tests := map[string]struct {
		Limit int
		want  int
	}{
		"limit is not set": {
			Limit: 0,
			want:  20,
		},
		"limit is set": {
			Limit: 1,
			want:  1,
		},
	}
	for name, tt := range tests {
		p := Pager{
			ItemsPerPage: tt.Limit,
		}
		t.Run(name, func(t *testing.T) {
			g.Expect(p.GetPageLimit()).To(gomega.Equal(tt.want))
		})
	}
}

func TestPager_GetPage(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	tests := map[string]struct {
		Page int
		want int
	}{
		"page is not set": {
			Page: 0,
			want: 1,
		},
		"page is set": {
			Page: 2,
			want: 2,
		},
	}
	for name, tt := range tests {
		p := Pager{
			Page: tt.Page,
		}
		t.Run(name, func(t *testing.T) {
			g.Expect(p.GetPage()).To(gomega.Equal(tt.want))
		})
	}
}

func TestPager_GetOffset(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	tests := map[string]struct {
		Page  int
		Limit int
		want  int
	}{
		"neither page nor limit is set": {
			want: 20,
		},
		"page is set and limit is not set": {
			Page: 2,
			want: 40,
		},
		"page is not set and limit is set": {
			Limit: 2,
			want:  2,
		},
		"page and limit are set": {
			Limit: 2,
			Page:  2,
			want:  4,
		},
	}
	for name, tt := range tests {
		p := Pager{
			Page:         tt.Page,
			ItemsPerPage: tt.Limit,
		}
		t.Run(name, func(t *testing.T) {
			g.Expect(p.GetOffset()).To(gomega.Equal(tt.want))
		})
	}
}
