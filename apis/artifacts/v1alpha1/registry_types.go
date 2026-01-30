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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// ImageConfig define sever and multiple podman credentials related content.
type ImageConfig struct {
	ImageAuths ImageAuths `json:"auths"`
}

// ImageAuths server and auth item list
type ImageAuths map[string][]DockerAuthItem

// RegistryConfigJson podman credentials data
type RegistryConfigJson struct {
	Auths map[string]DockerAuthItem `json:"auths"`
}

// DockerAuthItem registry credential information for a single repository
type DockerAuthItem struct {
	// Username username
	Username string `json:"username"`
	// Password authentication password
	Password string `json:"password"`
	// Auth authentication token, generated with username:password base64
	Auth string `json:"auth"`
	// Eamil user eamil
	Eamil string `json:"email"`
	// Scope project scope for credential authentication
	Scope string `json:"scope"`
}

// GenerateDockerAuth generates auth data from username and password.
func GenerateDockerAuth(username, password []byte) string {
	authData := make([]byte, 0, len(username)+len(password)+1)
	return base64.StdEncoding.EncodeToString(fmt.Appendf(authData, "%s:%s", username, password))
}

// GetAuthFromRegistryConfigJson get podman credential information from podman config json
func GetAuthFromRegistryConfigJson(registry string, registryConfigJsonBytes []byte) (username, password string, err error) {
	var registryConfig RegistryConfigJson

	if err = json.Unmarshal(registryConfigJsonBytes, &registryConfig); err != nil {
		return "", "", err
	}
	if registryConfig.Auths == nil {
		return "", "", fmt.Errorf("no auths found")
	}

	candidate := []string{registry, "https://" + registry, "http://" + registry}
	tmpURL := "https://" + strings.TrimPrefix(strings.TrimPrefix(registry, "http://"), "https://")
	u, err := url.Parse(tmpURL)
	if err == nil {
		candidate = append(candidate, u.Host, "https://"+u.Host, "http://"+u.Host)
	}

	// generate all possible address
	for address, auth := range registryConfig.Auths {
		address = strings.TrimRight(address, "/")
		registryConfig.Auths[address] = auth
	}

	for _, address := range candidate {
		if auth, ok := registryConfig.Auths[address]; ok {
			return auth.Username, auth.Password, nil
		}
	}

	return "", "", fmt.Errorf("no auth found for registry: %s", registry)
}
