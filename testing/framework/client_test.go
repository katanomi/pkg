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

package framework

import (
	"context"
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	"github.com/go-resty/resty/v2"
)

func TestRequestWrapAuthMeta(t *testing.T) {
	g := NewGomegaWithT(t)
	req := &resty.Request{Header: http.Header{}}
	auth := client.Auth{
		Type: "xx",
		Secret: map[string][]byte{
			"username": []byte("tom"),
		},
	}
	meta := client.Meta{
		Version: "",
		BaseURL: "",
	}

	g.Expect(RequestWrapAuthMeta(nil, auth, meta)).To(BeNil())
	newReq := RequestWrapAuthMeta(req, auth, meta)
	g.Expect(newReq).NotTo(BeNil())
	g.Expect(newReq.Header.Get(client.PluginAuthHeader)).To(Equal(string(auth.Type)))

	secretBase64 := base64.StdEncoding.EncodeToString([]byte(`{"username":"dG9t"}`))
	g.Expect(newReq.Header.Get(client.PluginSecretHeader)).To(Equal(secretBase64))

	metaBase64 := base64.StdEncoding.EncodeToString([]byte(`{}`))
	g.Expect(newReq.Header.Get(client.PluginMetaHeader)).To(Equal(metaBase64))
}

func TestConvertToAuthSecret(t *testing.T) {
	g := NewGomegaWithT(t)
	auth := client.Auth{
		Type: "xx",
		Secret: map[string][]byte{
			"username": []byte("tom"),
		},
	}
	secret := ConvertToAuthSecret(auth)
	g.Expect(secret.Data).To(Equal(auth.Secret))
	g.Expect(secret.Annotations[v1alpha1.SecretTypeAnnotationKey]).To(Equal(string(auth.Type)))
}

type testPlugin struct{}

func (t *testPlugin) Path() string {
	return "test-plugin"
}

func (t *testPlugin) Setup(_ context.Context, _ *zap.SugaredLogger) error {
	return nil
}

func TestGetPluginAddress(t *testing.T) {
	g := NewGomegaWithT(t)

	plugin := &testPlugin{}
	c := &resty.Client{HostURL: "http://127.0.0.1:8888"}

	address, err := GetPluginAddress(c, plugin)
	g.Expect(err).Should(Succeed())
	g.Expect(address.URL.String()).Should(Equal("http://127.0.0.1:8888/plugins/v1alpha1/test-plugin"))

	address, err = GetPluginAddress(nil, plugin)
	g.Expect(err).ShouldNot(Succeed())
	g.Expect(address).To(BeNil())

	address, err = GetPluginAddress(c, nil)
	g.Expect(err).ShouldNot(Succeed())
	g.Expect(address).To(BeNil())
}

func TestMustGetPluginAddress(t *testing.T) {
	g := NewGomegaWithT(t)

	plugin := &testPlugin{}
	c := &resty.Client{HostURL: "http://127.0.0.1:8888"}

	address := MustGetPluginAddress(c, plugin)
	g.Expect(address.URL.String()).Should(Equal("http://127.0.0.1:8888/plugins/v1alpha1/test-plugin"))

	g.Expect(func() {
		MustGetPluginAddress(nil, plugin)
	}).Should(Panic())

	g.Expect(func() {
		MustGetPluginAddress(c, nil)
	}).Should(Panic())
}
