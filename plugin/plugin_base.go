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
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/katanomi/pkg/plugin/client"
	"github.com/katanomi/pkg/plugin/config"
	"k8s.io/apimachinery/pkg/api/errors"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/network"
	"knative.dev/pkg/system"
)

// PluginBase is the base struct for all plugins
type PluginBase struct{}

// CheckAlive request the tool service with get method to check if it is alive
func (p PluginBase) CheckAlive(ctx context.Context) error {
	meta := client.ExtraMeta(ctx)
	if meta == nil {
		return errors.NewBadRequest("missing meta")
	}
	// Use http request to check whether the network is connected,
	// and ignore the http response.
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, meta.BaseURL, nil)
	_, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.NewInternalError(err)
	}

	return nil
}

// GetAddressURL get the access url of the plugin service(not tool service)
func (p PluginBase) GetAddressURL() *apis.URL {
	svcNamespace := system.Namespace()
	svcName := os.Getenv(config.EnvServiceName)
	svcMethod := os.Getenv(config.EnvServiceMethod)
	if svcMethod == "" {
		svcMethod = "http"
	}
	urlStr := fmt.Sprintf("%s://%s", svcMethod, network.GetServiceHostname(svcName, svcNamespace))
	url, _ := apis.ParseURL(urlStr)
	return url
}

// GetWebhookURL get the webhook url
// used to receive events from the tool service
func (p PluginBase) GetWebhookURL() (*apis.URL, bool) {
	webhookAddr := os.Getenv(config.EnvWebhookAddress)
	var url *apis.URL
	if webhookAddr != "" {
		url, _ = apis.ParseURL(webhookAddr)
	}
	return url, url != nil
}
