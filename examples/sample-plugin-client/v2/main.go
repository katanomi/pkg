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

package v2

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"

	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
	v2 "github.com/katanomi/pkg/plugin/client/v2"
)

func main() {
	// get a real meta from harbor ClusterIntegration instance
	meta := client.Meta{
		Version: "v0.0.1",
		BaseURL: "http://harbor.test",
	}

	// get a real secret using k8s client
	secret := v1.Secret{}

	// get a real harbor plugin url from harbor IntegrationClass instance
	address := &duckv1.Addressable{}

	// get a project client with harbor meta
	// need to implement more clients, such as repository, artifact
	pluginClient := v2.NewPluginClientV2(address, meta, secret)

	list, err := pluginClient.ListProjects(context.TODO(), v1alpha1.ListOptions{})

	fmt.Println(list, err)
}
