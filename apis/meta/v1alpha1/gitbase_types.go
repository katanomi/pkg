/*
Copyright 2021 The Katanomi Authors.

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

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// GitUserBaseInfo user base info
type GitUserBaseInfo struct {
	// Name is the login Name
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GitRepo Repo base info
type GitRepo struct {
	// Project gitlab is empty string, github,gogs,gitea is owner name
	Project string `json:"project,omitempty" yaml:"project,omitempty"`
	// Repository is different between platforms. gitlab is number;github,gogs,gitea is repo name
	Repository string `json:"repository,omitempty" yaml:"repository,omitempty"`
}

// GitOperateLogBaseInfo a simple log for operate
type GitOperateLogBaseInfo struct {
	User *GitUserBaseInfo `json:"user,omitempty"`
	Time *metav1.Time     `json:"time,omitempty"`
}
