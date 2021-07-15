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

package harbor

import (
	"context"
	"fmt"
	"os"

	"github.com/katanomi/pkg/plugin/client"
)

type harbor struct {
}

func New() *harbor {
	return &harbor{}
}

func (h *harbor) Path() string {
	if path := os.Getenv("HARBOR_PATH"); path != "" {
		return path
	}

	return "harbor"
}

func (h *harbor) ListProjects(ctx context.Context, option client.ListOption) (client.ProjectList, error) {
	auth := client.ExtractAuth(ctx)

	basic, err := auth.Basic()
	if err != nil {
		fmt.Println("get basic auth error: ", err.Error())
		return nil, err
	}

	// list project with harbor api

	fmt.Printf("basic auth: %v", basic)

	return nil, nil
}

func (h *harbor) CreateProject(ctx context.Context, project *client.Project) (*client.Project, error) {
	auth := client.ExtractAuth(ctx)

	oauth2, err := auth.Oauth2()
	if err != nil {
		fmt.Println("get basic auth error: ", err.Error())
		return nil, err
	}

	// create project with harbor api

	fmt.Printf("basic auth: %v", oauth2)

	return &client.Project{}, nil
}

func (h *harbor) ListResources(ctx context.Context, option client.ListOption) (client.ResourceList, error) {
	meta := client.ExtraMeta(ctx)

	// list resource with harbor api

	fmt.Printf("basic auth: %v", meta)

	return nil, nil
}
