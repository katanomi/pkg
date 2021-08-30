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

package client

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/go-resty/resty/v2"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSecretOpts(t *testing.T) {
	secretData := map[string][]byte{
		"username": []byte("123"),
		"password": []byte("456"),
	}

	secret := v1.Secret{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Type:       v1.SecretTypeBasicAuth,
		Data:       secretData,
	}

	opt := SecretOpts(secret)
	request := resty.New().R()
	opt(request)

	g := NewGomegaWithT(t)
	g.Expect(request.Header.Get(PluginAuthHeader)).To(Equal(string(v1.SecretTypeBasicAuth)))
	dataBytes, _ := json.Marshal(secret.Data)

	g.Expect(request.Header.Get(PluginSecretHeader)).To(Equal(base64.StdEncoding.EncodeToString(dataBytes)))
}

func TestSecretOptsWithLabel(t *testing.T) {
	secretData := map[string][]byte{
		"accessToken": []byte("123"),
	}

	secret := v1.Secret{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				metav1alpha1.SecretTypeLabelKey: string(metav1alpha1.AuthTypeOAuth2),
			},
		},
		Type: v1.SecretType("unknown"),
		Data: secretData,
	}

	opt := SecretOpts(secret)
	request := resty.New().R()
	opt(request)

	g := NewGomegaWithT(t)
	g.Expect(request.Header.Get(PluginAuthHeader)).To(Equal(string(metav1alpha1.AuthTypeOAuth2)))
	dataBytes, _ := json.Marshal(secret.Data)

	g.Expect(request.Header.Get(PluginSecretHeader)).To(Equal(base64.StdEncoding.EncodeToString(dataBytes)))
}
