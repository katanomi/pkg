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
	"reflect"
	"testing"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
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

func TestGetGitRepoInfo(t *testing.T) {
	tests := []struct {
		name        string
		gitAddress  string
		wantHost    string
		wantGitRepo metav1alpha1.GitRepo
		wantErr     bool
	}{
		{
			name:        "invalid git repo url",
			gitAddress:  "http:// github.com/katanomi/pkg.git",
			wantHost:    "",
			wantGitRepo: metav1alpha1.GitRepo{},
			wantErr:     true,
		},
		{
			name:       "shuffix with .git",
			gitAddress: "http://github.com/katanomi/pkg.git",
			wantHost:   "http://github.com",
			wantGitRepo: metav1alpha1.GitRepo{
				Project:    "katanomi",
				Repository: "pkg",
			},
			wantErr: false,
		},
		{
			name:       "shuffix without .git",
			gitAddress: "http://github.com/katanomi/pkg",
			wantHost:   "http://github.com",
			wantGitRepo: metav1alpha1.GitRepo{
				Project:    "katanomi",
				Repository: "pkg",
			},
			wantErr: false,
		},
		{
			name:       "shuffix without /",
			gitAddress: "http://github.com/katanomi/pkg/",
			wantHost:   "http://github.com",
			wantGitRepo: metav1alpha1.GitRepo{
				Project:    "katanomi",
				Repository: "pkg",
			},
			wantErr: false,
		},
		{
			name:        "path too short",
			gitAddress:  "http://github.com/katanomi/",
			wantHost:    "",
			wantGitRepo: metav1alpha1.GitRepo{},
			wantErr:     true,
		},
		{
			name:       "path with subgroup",
			gitAddress: "http://github.com/katanomi/subgroup/pkg/",
			wantHost:   "http://github.com",
			wantGitRepo: metav1alpha1.GitRepo{
				Project:    "katanomi/subgroup",
				Repository: "pkg",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHost, gotGitRepo, err := GetGitRepoInfo(tt.gitAddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGitRepoInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHost != tt.wantHost {
				t.Errorf("GetGitRepoInfo() gotHost = %v, want %v", gotHost, tt.wantHost)
			}
			if !reflect.DeepEqual(gotGitRepo, tt.wantGitRepo) {
				t.Errorf("GetGitRepoInfo() gotGitRepo = %v, want %v", gotGitRepo, tt.wantGitRepo)
			}
		})
	}
}
