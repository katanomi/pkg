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

package options

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
)

//go:generate mockgen  -destination=../../testing/mock/k8s.io/cli-runtime/pkg/genericclioptions/config_flags.go -package=genericclioptions k8s.io/cli-runtime/pkg/genericclioptions RESTClientGetter
//go:generate mockgen  -destination=../../testing/mock/k8s.io/client-go/tools/clientcmd/client_config.go -package=clientcmd k8s.io/client-go/tools/clientcmd ClientConfig

// RESTClientGetterOption is the generic client option for k8s client
type RESTClientGetterOption struct {
	// ConfigFlags for interface of genericclioptions.ConfigFlags
	ConfigFlag genericclioptions.RESTClientGetter
}

// Setup is the setup function for RESTClientGetterOption
func (m *RESTClientGetterOption) Setup(ctx context.Context, cmd *cobra.Command, args []string) (err error) {
	m.ConfigFlag = &genericclioptions.ConfigFlags{}
	return nil
}

// GetClusterToken get token for kubeconfig.
func (m *RESTClientGetterOption) GetClusterToken(ctx context.Context) (string, error) {
	config, err := m.ConfigFlag.ToRESTConfig()
	if err != nil {
		return "", fmt.Errorf("get kubeconfig token failed, error is: %v", err)
	}
	if config.BearerToken != "" {
		return config.BearerToken, nil
	}
	config, err = rest.InClusterConfig()
	if err != nil {
		return "", fmt.Errorf("get kubeconfig token from cluster failed, error is: %v", err)
	}
	return config.BearerToken, nil
}

// GetNamespace get namespace from environment or incluster kubeconfig
func (m *RESTClientGetterOption) GetNamespace() (string, error) {
	cfgLoader := m.ConfigFlag.ToRawKubeConfigLoader()
	ns, _, err := cfgLoader.Namespace()
	return ns, err
}
