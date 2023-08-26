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

package v1alpha1

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test.GetAuthFromDockerConfigJson", func() {
	var (
		emptyAuths = []byte("{}")
		mockAuths  = []byte(`
			{
				"auths": {
					"docker.io": {
						"username": "u1",
						"password": "p1"
					},
					"https://docker.io": {
						"username": "u2",
						"password": "p2"
					},
					"http://docker.io": {
						"username": "u3",
						"password": "p3"
					},
					"https://docker.io/user": {
						"username": "u4",
						"password": "p4"
					},
					"https://suffix.docker.io/////": {
						"username": "u5",
						"password": "p5"
					}
				}
			}
		`)
	)
	DescribeTable("GetAuthFromDockerConfigJson",
		func(registry string, dockerConfigJsonBytes []byte, username, password string, err error) {
			actualUsername, actualPassword, actualErr := GetAuthFromDockerConfigJson(registry, dockerConfigJsonBytes)
			Expect(username).To(Equal(actualUsername))
			Expect(password).To(Equal(actualPassword))
			errStr := ""
			if err != nil {
				errStr = err.Error()
			}
			actualErrStr := ""
			if actualErr != nil {
				actualErrStr = actualErr.Error()
			}
			Expect(errStr).To(Equal(actualErrStr))
		},
		Entry("bytes is nil", "", nil,
			"", "", errors.New("unexpected end of JSON input"),
		),
		Entry("auths is nil", "", emptyAuths,
			"", "", errors.New("no auths found"),
		),
		Entry("just matched registry", "docker.io", mockAuths,
			"u1", "p1", nil,
		),
		Entry("just matched registry", "https://docker.io", mockAuths,
			"u2", "p2", nil,
		),
		Entry("just matched registry", "http://docker.io", mockAuths,
			"u3", "p3", nil,
		),
		Entry("just matched registry", "https://docker.io/user", mockAuths,
			"u4", "p4", nil,
		),
		Entry("matched registry suffixed with /", "https://suffix.docker.io", mockAuths,
			"u5", "p5", nil,
		),
		Entry("fallback to host", "https://docker.io/not-exist", mockAuths,
			"u1", "p1", nil,
		),
		Entry("not auth found", "not.exist.com", mockAuths,
			"", "", errors.New("no auth found for registry: not.exist.com"),
		),
	)
})
