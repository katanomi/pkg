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

package registry

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	artifactsv1alpha1 "github.com/katanomi/pkg/apis/artifacts/v1alpha1"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	kclient "github.com/katanomi/pkg/client"
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// RegistrySchemeDetectionBySecret detect registry scheme by secret
type RegistrySchemeDetectionBySecret struct {
	*DefaultRegistrySchemeDetection

	clt       client.Client
	secretRef *corev1.ObjectReference
}

// NewRegistrySchemeDetectionBySecret create a new RegistrySchemeDetectionBySecret
func NewRegistrySchemeDetectionBySecret(client *resty.Client, insecure, cache bool) *RegistrySchemeDetectionBySecret {
	return &RegistrySchemeDetectionBySecret{
		DefaultRegistrySchemeDetection: NewDefaultRegistrySchemeDetection(client, insecure, cache),
	}
}

// WithClient set client
func (d *RegistrySchemeDetectionBySecret) WithClient(clt client.Client) *RegistrySchemeDetectionBySecret {
	d.clt = clt
	return d
}

// WithSecretRef set secret reference
func (d *RegistrySchemeDetectionBySecret) WithSecretRef(secretRef *corev1.ObjectReference) *RegistrySchemeDetectionBySecret {
	d.secretRef = secretRef
	return d
}

// DetectScheme detect registry scheme
// If the secret exists, the parameter of username and password will be ignored.
// These parameters are reserved only to satisfy the interface `RegistrySchemeDetection`
func (d *RegistrySchemeDetectionBySecret) DetectScheme(ctx context.Context, registry string, auths ...AuthOption) (string, error) {
	if d.DefaultRegistrySchemeDetection == nil {
		return "", fmt.Errorf("default registry scheme detection is nil")
	}

	log := logging.FromContext(ctx).With("registry", registry)

	if d.secretRef != nil && d.secretRef.Name != "" {
		clt := kclient.Client(ctx)
		if clt == nil {
			return "", fmt.Errorf("failed to get kubernetes client")
		}

		secret := &corev1.Secret{}
		secretKey := metav1alpha1.GetNamespacedNameFromRef(d.secretRef)
		err := clt.Get(ctx, secretKey, secret)
		if err != nil {
			log.Errorw("failed to get secret", "secretKey", secretKey, "error", err)
			return "", fmt.Errorf("failed to get secret %s: %w", secretKey.String(), err)
		}
		username, password, err := getDockerAuthFromSecret(registry, secret)
		if err != nil {
			log.Debugw("failed to get username and password from secret", "error", err)
		}
		if username != "" && password != "" {
			auths = append(auths, WithBasicAuth(username, password))
		}
	}

	return d.DefaultRegistrySchemeDetection.DetectScheme(ctx, registry, auths...)
}

// DetectSchemeWithDefault detect registry scheme, if detect failed, return default scheme
func (d *RegistrySchemeDetectionBySecret) DetectSchemeWithDefault(ctx context.Context, registry, defaultScheme string, auths ...AuthOption) string {
	scheme, err := d.DetectScheme(ctx, registry, auths...)
	if err != nil {
		return defaultScheme
	}
	return scheme
}

func getDockerAuthFromSecret(registryHost string, secret *corev1.Secret) (string, string, error) {
	switch secret.Type {
	case corev1.SecretTypeBasicAuth:
		username := string(secret.Data[corev1.BasicAuthUsernameKey])
		password := string(secret.Data[corev1.BasicAuthPasswordKey])
		return username, password, nil
	case corev1.SecretTypeDockerConfigJson:
		return artifactsv1alpha1.GetAuthFromDockerConfigJson(registryHost, secret.Data[corev1.DockerConfigJsonKey])
	default:
		return "", "", fmt.Errorf("unsupported secret type %s", secret.Type)
	}
}
