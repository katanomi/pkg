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

// Package env provides functions to get the cluster name from environment
package env

import (
	"os"
)

const (
	// ClusterNameEnvKey is the environment variable that specifies the system cluster name.
	// ClusterName can be attached to the created resource. When there are multiple clusters, the resource can be found by the api by the cluster name.
	ClusterNameEnvKey = "SYSTEM_CLUSTER_NAME" // NOSONAR // ignore: "Key" detected here, make sure this is not a hard-coded credential

	// BASE_DOMAIN used to get the data of BASE_DAMAIN from the environment variable.
	// Usually used as the base value of label or annotation. For example: {BASE_DOMAIN}/label: value.
	BaseDomainEnvKey = "BASE_DOMAIN"
)

// ClusterName returns the cluster name that describes the cluster name the instance runs on.
// You need to specify the relevant environment variables through deployment when deploying.
func ClusterName() string {
	return os.Getenv(ClusterNameEnvKey)
}

// BaseDomain return BASE_DOMAIN value from environment variable. when BASE_DOMAIN not set, will return "katanomi.dev"
func BaseDomain() string {
	return fromEnvStr(BaseDomainEnvKey, "katanomi.dev")
}

func fromEnvStr(key string, defaultStr string) string {
	str := os.Getenv(key)
	if str == "" {
		return defaultStr
	}
	return str
}
