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

package multicluster

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	// placeHolderClusterName defines the placeholder for the cluster name in proxy paths.
	placeHolderClusterName = "{name}"
)

// ClusterProxyHost constructs a complete proxy URL by replacing the cluster name placeholder in the proxy path
// It takes the proxy host and path, replaces the "{name}" placeholder with the actual cluster name,
// and returns the formatted proxy URL.
func ClusterProxyHost(proxyHost string, proxyPath string, clusterName string) (string, error) {
	proxyPath = strings.ReplaceAll(proxyPath, placeHolderClusterName, clusterName)
	proxyPath = strings.TrimPrefix(proxyPath, "/")

	hostURL, err := url.Parse(proxyHost)
	if err != nil {
		return "", err
	}
	//ensuring the host URL uses the HTTPS scheme if not specified.
	if hostURL.Scheme == "" {
		hostURL.Scheme = "https"
	}

	return fmt.Sprintf("%s/%s", hostURL.String(), proxyPath), nil
}
