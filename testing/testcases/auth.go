//go:build e2e
// +build e2e

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

package testcases

import (
	"context"

	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
	corev1 "k8s.io/api/core/v1"
)

// MakeBasicAuthContext make a basic auth context
func MakeBasicAuthContext(ctx context.Context, meta client.Meta, auth client.Auth) context.Context {
	return meta.WithContext(auth.WithContext(ctx))
}

// MakeBasicAuthContextFromEnv make a basic auth context from env
func MakeBasicAuthContextFromEnv(ctx context.Context, api, username, password string) context.Context {
	meta := client.Meta{
		BaseURL: api,
	}
	auth := client.Auth{
		Type: v1alpha1.AuthTypeBasic,
		Secret: map[string][]byte{
			corev1.BasicAuthUsernameKey: []byte(username),
			corev1.BasicAuthPasswordKey: []byte(password),
		},
	}

	return MakeBasicAuthContext(ctx, meta, auth)
}
