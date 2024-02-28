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

package testing

import (
	"testing"

	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGetCaseName(t *testing.T) {
	g := NewWithT(t)
	tests := []struct {
		name string
		s    string
		want []string
	}{
		{
			name: "without case name",
			s:    "test",
			want: []string{},
		},
		{
			name: "with case name",
			s:    "{case:aaaa}",
			want: []string{"aaaa"},
		},
		{
			name: "case name contains underscore",
			s:    "{case:aaaa_bbbb}",
			want: []string{"aaaa_bbbb"},
		},
		{
			name: "case name contains special characters",
			s:    "{case:aaaa@bbbb}",
			want: []string{},
		},
		{
			name: "empty string",
			s:    "",
			want: []string{},
		},
		{
			name: "Multiple case name",
			s:    "{case:aaaa} {case:bbbb}",
			want: []string{"aaaa", "bbbb"},
		},
		{
			name: "go test case name",
			s:    "TestAbc_eee",
			want: []string{"Abc_eee"},
		},
		{
			name: "go test case name with sub name",
			s:    "TestAbc_eee/case1",
			want: []string{"Abc_eee"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g.Expect(GetCaseNames(tt.s)).To(Equal(tt.want))
		})
	}
}

func TestCase(t *testing.T) {
	g := NewWithT(t)
	tests := []struct {
		name     string
		caseName string
		want     ginkgo.Labels
	}{
		{
			name:     "empty case name",
			caseName: "",
			want:     ginkgo.Labels{},
		},
		{
			name:     "case name",
			caseName: "abc",
			want:     ginkgo.Labels{"{case:abc}"},
		},
		{
			name:     "case name contains underscore",
			caseName: "abc_eee",
			want:     ginkgo.Labels{"{case:abc_eee}"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g.Expect(Case(tt.caseName)).To(Equal(tt.want))
		})
	}
}
