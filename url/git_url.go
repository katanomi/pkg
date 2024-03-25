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
	"net/url"
	"strings"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"knative.dev/pkg/apis"
)

// MatchGitURLPrefix determine if the git URL is a subset of the target URL
// This matches:
// gitURL: https://github.com/katanomi/pkg.git
// target:
//   - https://github.com/katanomi/pkg.git
//   - https://github.com/katanomi/pkg
//   - https://github.com/katanomi/
//   - https://github.com
//
// This mismatch:
// gitURL: https://github.com/katanomi/pkg.git
// target:
//   - https://github.com/katanomi/pkg/
//
// This mismatch:
// gitURL: https://github.com/katanomi/pkg
// target:
//   - https://github.com/katanomi/pkg/
func MatchGitURLPrefix(gitURL, target *apis.URL) bool {
	if gitURL == nil || target == nil {
		return false
	}
	if gitURL.Host != target.Host {
		return false
	}
	sPath := gitURL.Path
	tPath := target.Path
	if strings.HasSuffix(sPath, ".git") {
		sPath = strings.TrimSuffix(sPath, ".git")
		tPath = strings.TrimSuffix(tPath, ".git")
	}
	if strings.HasSuffix(tPath, ".git") {
		// If the exact match, return true
		if sPath == strings.TrimSuffix(tPath, ".git") {
			return true
		}
	}
	// determine if the prefix matches
	return strings.HasPrefix(sPath, tPath)
}

// GetGitRepoInfo try to get the project and repository from the git address.
func GetGitRepoInfo(gitAddress string) (host string, gitRepo metav1alpha1.GitRepo, err error) {
	gitAddress = strings.TrimSuffix(gitAddress, ".git")
	URL, err := url.Parse(gitAddress)
	if err != nil {
		return
	}

	projectRepo := strings.Split(strings.Trim(URL.Path, "/"), "/")
	if len(projectRepo) < 2 {
		err = fmt.Errorf("invaild git address %s which should have project and repository", gitAddress)
		return
	}

	host = fmt.Sprintf("%s://%s", URL.Scheme, URL.Host)
	gitRepo.Project = strings.Join(projectRepo[:len(projectRepo)-1], "/")
	gitRepo.Repository = projectRepo[len(projectRepo)-1]
	return
}
