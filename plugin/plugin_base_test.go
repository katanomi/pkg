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

package plugin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/katanomi/pkg/plugin/client"
	"github.com/katanomi/pkg/plugin/config"
	. "github.com/onsi/gomega"
	"knative.dev/pkg/system"
)

func TestPluginBase_CheckAlive(t *testing.T) {
	g := NewGomegaWithT(t)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer s.Close()

	metaData := client.Meta{BaseURL: s.URL, Version: "v1"}
	ctx := metaData.WithContext(context.Background())
	base := PluginBase{}
	err := base.CheckAlive(ctx)
	g.Expect(err).Should(Succeed())

	metaData = client.Meta{BaseURL: "invalid url", Version: "v1"}
	ctx = metaData.WithContext(context.Background())
	err = base.CheckAlive(ctx)
	g.Expect(err).ShouldNot(Succeed())
}

func TestPluginBase_GetAddressURL(t *testing.T) {
	g := NewGomegaWithT(t)
	oldWebhookAddress := os.Getenv(config.EnvWebhookAddress)
	oldServiceName := os.Getenv(config.EnvServiceName)
	oldServiceMethod := os.Getenv(config.EnvServiceMethod)
	oldNamespace := os.Getenv(system.NamespaceEnvKey)
	os.Setenv(config.EnvWebhookAddress, "https://abc.com")
	os.Setenv(config.EnvServiceName, "test-service")
	os.Setenv(config.EnvServiceMethod, "https")
	os.Setenv(system.NamespaceEnvKey, "test-ns")
	defer func() {
		os.Setenv(config.EnvWebhookAddress, oldWebhookAddress)
		os.Setenv(config.EnvServiceName, oldServiceName)
		os.Setenv(config.EnvServiceMethod, oldServiceMethod)
		os.Setenv(system.NamespaceEnvKey, oldNamespace)
	}()

	base := PluginBase{}
	url := base.GetAddressURL().String()
	g.Expect(url).To(Equal("https://test-service.test-ns.svc.cluster.local"))

	webhookURL, b := base.GetWebhookURL()
	g.Expect(webhookURL.String()).To(Equal("https://abc.com"))
	g.Expect(b).To(BeTrue())
}
