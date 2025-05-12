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
	"errors"
	"fmt"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DefaultGitFilePath default git file path
const DefaultGitFilePath = ".git/katanomi.git.json"

// GitUserBaseInfo user base info
type GitUserBaseInfo struct {
	// Name is the login Name
	Name  string `json:"name" variable:"example=joedoe"`
	Email string `json:"email"  variable:"example=joedoe@example.com"`
}

// GitRepo Repo base info
type GitRepo struct {
	// Project gitlab is empty string, github,gogs,gitea is owner name
	Project string `json:"project,omitempty" yaml:"project,omitempty"`
	// Repository is different between platforms. gitlab is number;github,gogs,gitea is repo name
	Repository string `json:"repository,omitempty" yaml:"repository,omitempty"`
}

// Validate validate the git repo
func (r *GitRepo) Validate() error {
	if r.Project == "" {
		return errors.New("project is empty")
	}
	if r.Repository == "" {
		return errors.New("repository is empty")
	}
	return nil
}

// ProjectID return the project id of the repo
func (in *GitRepo) ProjectID() string {
	if in == nil {
		return ""
	}
	return in.String()
}

// ProjectIDOrProject returns the project string when the project is a pure number and repository is empty.
// Otherwise returns the result of ProjectID method.
func (in *GitRepo) ProjectIDOrProject() string {
	if in == nil {
		return ""
	}

	if in.Repository == "" {
		// Try to parse the project as an integer
		if _, err := strconv.Atoi(in.Project); err == nil {
			return in.Project
		}
	}

	return in.ProjectID()
}

// String prints the GitRepo data in a user friendly format
// Example:
//
//	For Project "abc" and repository "xyz" the output is "abc/xyz"
//
// Compatible with %s format and other string printing methods
func (in GitRepo) String() string {
	return fmt.Sprintf("%s/%s", in.Project, in.Repository)
}

// GitOperateLogBaseInfo a simple log for operate
type GitOperateLogBaseInfo struct {
	User *GitUserBaseInfo `json:"user,omitempty"`
	Time *metav1.Time     `json:"time,omitempty"`
}

type DownloadURL struct {
	Zip    string  `json:"zip"`
	TarGz  *string `json:"tar.gz"`
	TarBa2 *string `json:"tar.ba2"`
	Tar    *string `json:"tar"`
}
