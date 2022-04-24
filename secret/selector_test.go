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

package secret

import (
	"strings"
	"testing"

	"go.uber.org/zap"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"

	"k8s.io/apimachinery/pkg/types"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSelect(t *testing.T) {
	zaplog, _ := zap.NewDevelopment()
	log := zaplog.Sugar()

	var buildSecret = func(name, namespace, address string, secretType corev1.SecretType, scopes []string, applyNamespaces []string, isGlobal bool) corev1.Secret {
		return corev1.Secret{
			Type:     secretType,
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{

				Name:      name,
				Namespace: namespace,
				Annotations: map[string]string{
					metav1alpha1.IntegrationAddressAnnotation:     address,
					metav1alpha1.IntegrationSecretApplyNamespaces: strings.Join(applyNamespaces, ","),
					metav1alpha1.IntegrationResourceScope:         strings.Join(scopes, ","),
				},
			},
		}
	}

	var expect = func(err error, actual *corev1.Secret, secretNamespace, secretName string) {
		if err != nil {
			t.Errorf("err should be nil, but: %s", err.Error())
		}
		if actual == nil {
			t.Errorf("actual should not be nil")
		}
		if actual.Namespace != secretNamespace || actual.Name != secretName {
			t.Errorf("acutal should be %s/%s but: %s/%s", secretNamespace, secretName, actual.Namespace, actual.Name)
		}
	}

	t.Run("has namespaced credentials and correct scopes", func(t *testing.T) {
		secrets := []corev1.Secret{
			buildSecret("secret-1", "project-1", "https://gitlab.com/",
				"", []string{"/devops-1/"}, []string{}, false),
			buildSecret("secret-2", "project-2", "https://gitlab.com/",
				"", []string{"/devops-2/"}, []string{}, false),
		}
		actual, err := selectToolSecret(log, secrets, []corev1.Secret{}, "https://gitlab.com/devops-2/demo", SelectSecretOption{
			PerferredSecret:            types.NamespacedName{},
			ExcludedSecretTypes:        nil,
			Namespace:                  "project-2",
			GlobalCredentialsNamespace: "global-credentials",
		})
		expect(err, actual, "project-2", "secret-2")
	})

	t.Run("has namespaced credentials and correct scopes and perferred secret", func(t *testing.T) {
		secrets := []corev1.Secret{
			buildSecret("secret-1", "project-1", "https://gitlab.com/",
				"", []string{"/devops-1/"}, []string{}, false),
			buildSecret("secret-2", "project-2", "https://gitlab.com/",
				"", []string{"/devops-2/"}, []string{}, false),
			buildSecret("secret-2.1", "project-2", "https://gitlab.com/",
				"", []string{"/devops-2/"}, []string{}, false),
			buildSecret("secret-2.0", "project-2", "https://gitlab.com/",
				"", []string{"/devops-2/"}, []string{}, false),
		}
		actual, err := selectToolSecret(log, secrets, []corev1.Secret{}, "https://gitlab.com/devops-2/demo", SelectSecretOption{
			PerferredSecret:            types.NamespacedName{Namespace: "project-2", Name: "secret-2.1"},
			ExcludedSecretTypes:        nil,
			Namespace:                  "project-2",
			GlobalCredentialsNamespace: "global-credentials",
		})

		expect(err, actual, "project-2", "secret-2.1")
	})

	t.Run("has global credentials but with not apply namespaces", func(t *testing.T) {
		secrets := []corev1.Secret{
			buildSecret("secret-1", "project-1", "https://gitlab.com/",
				"", []string{"/devops-1/"}, []string{}, false),
		}
		globalSecrets := []corev1.Secret{
			buildSecret("secret-2", "global-credentials", "https://gitlab.com/",
				"", []string{"/devops-2/"}, []string{}, true),
		}
		actual, err := selectToolSecret(log, secrets, globalSecrets, "https://gitlab.com/devops-2/demo", SelectSecretOption{
			PerferredSecret:            types.NamespacedName{Namespace: "project-2", Name: "secret-2.1"},
			ExcludedSecretTypes:        nil,
			Namespace:                  "project-1",
			GlobalCredentialsNamespace: "global-credentials",
		})
		if err != nil {
			t.Errorf("err should be nil but: %s", err.Error())
		}

		if actual != nil {
			t.Errorf("should not select any secret, but: %s/%s", actual.Namespace, actual.Name)
		}
	})

	t.Run("has global credentials and with correct apply namespaces", func(t *testing.T) {
		secrets := []corev1.Secret{
			buildSecret("secret-1", "project-1", "https://gitlab.com/",
				"", []string{"/devops-1/"}, []string{}, false),
		}
		globalSecrets := []corev1.Secret{
			buildSecret("secret-2", "global-credentials", "https://gitlab.com/",
				"", []string{"/devops-2/"}, []string{"project-1"}, true),
		}
		actual, err := selectToolSecret(log, secrets, globalSecrets, "https://gitlab.com/devops-2/demo", SelectSecretOption{
			PerferredSecret:            types.NamespacedName{Namespace: "project-2", Name: "secret-what-ever"},
			ExcludedSecretTypes:        nil,
			Namespace:                  "project-1",
			GlobalCredentialsNamespace: "global-credentials",
		})
		expect(err, actual, "global-credentials", "secret-2")
	})

	t.Run("has global credentials but  with not correct apply namespaces", func(t *testing.T) {
		secrets := []corev1.Secret{
			buildSecret("secret-1", "project-1", "https://gitlab.com/",
				"", []string{"/devops-1/"}, []string{}, false),
		}
		globalSecrets := []corev1.Secret{
			buildSecret("secret-2", "global-credentials", "https://gitlab.com/",
				"", []string{"/devops-2/"}, []string{"project-2"}, true),
		}
		actual, err := selectToolSecret(log, secrets, globalSecrets, "https://gitlab.com/devops-2/demo", SelectSecretOption{
			PerferredSecret:            types.NamespacedName{Namespace: "project-2", Name: "secret-what-ever"},
			ExcludedSecretTypes:        nil,
			Namespace:                  "project-1",
			GlobalCredentialsNamespace: "global-credentials",
		})
		if err != nil {
			t.Errorf("err should be nil but: %s", err.Error())
		}

		if actual != nil {
			t.Errorf("should not select any secret, but: %s/%s", actual.Namespace, actual.Name)
		}
	})

	t.Run("exclude oauth2 type", func(t *testing.T) {
		s1 := corev1.Secret{
			Type: corev1.SecretTypeBasicAuth,
		}
		sList := []corev1.Secret{s1}
		option := SelectSecretOption{ExcludedSecretTypes: []corev1.SecretType{corev1.SecretTypeBasicAuth}}
		secret, err := selectToolSecret(log, sList, []corev1.Secret{}, "", option)
		if err != nil {
			t.Errorf("should not is nil")
		}
		if secret != nil {
			t.Errorf("should nil")
		}
		if !option.ExcludedSecretTypes.Contains(s1.Type) {
			t.Errorf("should contain")
		}
	})

	t.Run("only basicAuth type secret can be selected", func(t *testing.T) {
		sList := []corev1.Secret{
			buildSecret("secret-basic", "default", "https://1.2.3.4/",
				corev1.SecretTypeBasicAuth, []string{"/devops/"}, []string{""}, false),
			buildSecret("secret-basic1", "default", "https://1.2.3.4/",
				corev1.SecretTypeSSHAuth, []string{"/devops/"}, []string{""}, false),
			buildSecret("secret-basic2", "default", "https://1.2.3.4/",
				corev1.SecretTypeBootstrapToken, []string{"/devops/"}, []string{""}, false),
		}

		option := SelectSecretOption{SecretTypes: []corev1.SecretType{corev1.SecretTypeBasicAuth}}
		secret, err := selectToolSecret(log, sList, []corev1.Secret{}, "https://1.2.3.4/devops/test.git", option)
		if err != nil {
			t.Errorf("should be nil")
		}

		if !option.SecretTypes.Contains(sList[0].Type) {
			t.Errorf("should contain")
		}

		if secret == nil {
			t.Errorf("should not be nil")
		}

		if secret.Name != "secret-basic" {
			t.Errorf("find wrong secret")
		}
	})

	t.Run("only the correct basicAuth secret will be return", func(t *testing.T) {
		sList := []corev1.Secret{
			buildSecret("secret-basic", "default", "https://1.2.3.4/",
				corev1.SecretTypeSSHAuth, []string{"/devops/"}, []string{""}, false),
			buildSecret("secret-basic1", "default", "https://1.2.3.4/",
				corev1.SecretTypeSSHAuth, []string{"/devops1/"}, []string{""}, false),
			buildSecret("secret-basic2", "default", "https://1.2.3.4/",
				corev1.SecretTypeBasicAuth, []string{"/devops/"}, []string{""}, false),
		}

		option := SelectSecretOption{SecretTypes: []corev1.SecretType{corev1.SecretTypeSSHAuth, corev1.SecretTypeBasicAuth}}
		secret, err := selectToolSecret(log, sList, []corev1.Secret{}, "https://1.2.3.4/devops/test.git", option)
		if err != nil {
			t.Errorf("should be nil")
		}

		if !option.SecretTypes.Contains(sList[0].Type) {
			t.Errorf("should contain")
		}

		if secret == nil {
			t.Errorf("should not be nil")
		}

		if secret.Name != "secret-basic" {
			t.Errorf("find wrong secret")
		}
	})

	t.Run("no secret will be return", func(t *testing.T) {
		sList := []corev1.Secret{
			buildSecret("secret-basic", "default", "https://1.2.3.4/",
				corev1.SecretTypeSSHAuth, []string{"/devops0/"}, []string{""}, false),
			buildSecret("secret-basic1", "default", "https://1.2.3.4/",
				corev1.SecretTypeSSHAuth, []string{"/devops1/"}, []string{""}, false),
			buildSecret("secret-basic2", "default", "https://1.2.3.4/",
				corev1.SecretTypeBasicAuth, []string{"/devops/"}, []string{""}, false),
		}

		option := SelectSecretOption{SecretTypes: []corev1.SecretType{corev1.SecretTypeSSHAuth}}
		secret, err := selectToolSecret(log, sList, []corev1.Secret{}, "https://1.2.3.4/devops/test.git", option)
		if err != nil {
			t.Errorf("should be nil")
		}

		if !option.SecretTypes.Contains(sList[0].Type) {
			t.Errorf("should contain")
		}

		if secret != nil {
			t.Errorf("should not return secret because no sutiable secret proviede")
		}
	})

	t.Run("secretType is not provided so secret secret-basic2 will be return", func(t *testing.T) {
		sList := []corev1.Secret{
			buildSecret("secret-basic", "default", "https://1.2.3.4/",
				corev1.SecretTypeSSHAuth, []string{"/devops0/"}, []string{""}, false),
			buildSecret("secret-basic1", "default", "https://1.2.3.4/",
				corev1.SecretTypeSSHAuth, []string{"/devops1/"}, []string{""}, false),
			buildSecret("secret-basic2", "default", "https://1.2.3.4/",
				corev1.SecretTypeBasicAuth, []string{"/devops/"}, []string{""}, false),
		}

		option := SelectSecretOption{SecretTypes: []corev1.SecretType{}}
		secret, err := selectToolSecret(log, sList, []corev1.Secret{}, "https://1.2.3.4/devops/test.git", option)
		if err != nil {
			t.Errorf("should be nil")
		}

		if secret == nil {
			t.Errorf("should return secret")
		}

		if secret.Name != "secret-basic2" {
			t.Errorf("should return secret secret-basic2")
		}
	})

}

func TestSortSecretList(t *testing.T) {

	s1 := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "a",
		},
	}
	s2 := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "b",
		},
	}

	t.Run("test normal", func(t *testing.T) {
		secrets := make([]corev1.Secret, 0)
		secrets = append(secrets, s2)
		secrets = append(secrets, s1)
		if secrets[0].Name != "b" {
			t.Errorf("Expected b, get %s", secrets[0].Name)
		}
		secrets = sortSecretList(secrets)
		if secrets[0].Name != "a" {
			t.Errorf("Expected a, get %s", secrets[0].Name)
		}
		if secrets[1].Name != "b" {
			t.Errorf("Expected b, get %s", secrets[0].Name)
		}
	})

	t.Run("test empty", func(t *testing.T) {
		secrets := make([]corev1.Secret, 0)
		secrets = sortSecretList(secrets)
	})

	t.Run("test nil", func(t *testing.T) {
		var secrets []corev1.Secret
		secrets = sortSecretList(secrets)
	})
}
