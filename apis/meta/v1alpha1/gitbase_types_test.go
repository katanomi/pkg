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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGitRepo_Validate(t *testing.T) {
	type fields struct {
		Project    string
		Repository string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "empty project",
			fields:  fields{Project: "", Repository: "repo"},
			wantErr: true,
		},
		{
			name:    "empty repository",
			fields:  fields{Project: "project", Repository: ""},
			wantErr: true,
		},
		{
			name:    "valid",
			fields:  fields{Project: "project", Repository: "repo"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &GitRepo{
				Project:    tt.fields.Project,
				Repository: tt.fields.Repository,
			}
			if err := r.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("GitRepo.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGitRepo_ProjectID(t *testing.T) {
	type fields struct {
		Project    string
		Repository string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "empty",
			fields: fields{Project: "", Repository: ""},
			want:   "/",
		},
		{
			name:   "valid",
			fields: fields{Project: "project", Repository: "repo"},
			want:   "project/repo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := &GitRepo{
				Project:    tt.fields.Project,
				Repository: tt.fields.Repository,
			}
			if got := in.ProjectID(); got != tt.want {
				t.Errorf("GitRepo.ProjectID() = %v, want %v", got, tt.want)
			}
		})
	}
}

var _ = DescribeTable("GitRepo.String", func(project, repository, expected string) {
	repo := GitRepo{Project: project, Repository: repository}
	Expect(repo.String()).To(Equal(expected), "for project %s and repository %s the expected value was %s", project, repository, expected)
},
	Entry("abc org using xyz repo", "abc", "xyz", "abc/xyz"),
	Entry("root org using my-repo repo", "root", "my-repo", "root/my-repo"),
)
