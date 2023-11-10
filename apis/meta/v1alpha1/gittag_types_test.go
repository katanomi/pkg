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

	. "github.com/onsi/gomega"
)

func TestGitTag_Validate(t *testing.T) {
	g := NewWithT(t)
	testCases := map[string]struct {
		gitTag     GitTag
		wantErr    bool
		wantErrMsg string
	}{
		"project is empty": {
			gitTag: GitTag{
				GitRepo: GitRepo{},
				Tag:     "",
			},
			wantErr:    true,
			wantErrMsg: "project is empty",
		},
		"repository is empty": {
			gitTag: GitTag{
				GitRepo: GitRepo{
					Project: "abc",
				},
				Tag: "",
			},
			wantErr:    true,
			wantErrMsg: "repository is empty",
		},
		"tag is empty": {
			gitTag: GitTag{
				GitRepo: GitRepo{
					Project:    "abc",
					Repository: "abc",
				},
				Tag: "",
			},
			wantErr:    true,
			wantErrMsg: "tag is empty",
		},
		"valid case": {
			gitTag: GitTag{
				GitRepo: GitRepo{
					Project:    "abc",
					Repository: "abc",
				},
				Tag: "v0.0.1",
			},
			wantErr: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.wantErrMsg, func(t *testing.T) {
			err := testCase.gitTag.Validate()
			if !testCase.wantErr {
				g.Expect(err).Should(Succeed())
			} else {
				g.Expect(err).Should(HaveOccurred())
				g.Expect(err.Error()).Should(ContainSubstring(testCase.wantErrMsg))
			}
		})
	}
}

func TestCreateGitTagPayload_Validate(t *testing.T) {
	tests := map[string]struct {
		GitTag  GitTag
		Ref     string
		Message string
		wantErr bool
	}{
		"validate pass": {
			GitTag: GitTag{
				GitRepo: GitRepo{
					Project:    "test",
					Repository: "test",
				},
				Tag: "0.1",
			},
			Ref: "master",
		},
		"ref is empty": {
			GitTag: GitTag{
				GitRepo: GitRepo{
					Project:    "test",
					Repository: "test",
				},
				Tag: "0.1",
			},
			wantErr: true,
		},
		"gitTag validate failed": {
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := &CreateGitTagPayload{
				GitTag:  tt.GitTag,
				Ref:     tt.Ref,
				Message: tt.Message,
			}
			if err := r.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("CreateGitTagPayload.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
