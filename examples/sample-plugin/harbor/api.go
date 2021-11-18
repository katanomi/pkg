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
	"os"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
	"go.uber.org/zap"
)

type harbor struct {
	*zap.SugaredLogger
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

func (h *harbor) Setup(ctx context.Context, logger *zap.SugaredLogger) error {
	// fetch any necessary data from here
	// client := client.Client(ctx)
	h.SugaredLogger = logger
	return nil
}

func (h *harbor) ListProjects(ctx context.Context, option metav1alpha1.ListOptions) (*metav1alpha1.ProjectList, error) {
	auth := client.ExtractAuth(ctx)

	basic, err := auth.Basic()
	if err != nil {
		h.Infof("get basic auth error: %s", err.Error())
		return nil, err
	}

	// list project with harbor api
	if basic != nil {
		h.Infof("has basic auth")
	}

	return nil, nil
}

func (h *harbor) CreateProject(ctx context.Context, project *metav1alpha1.Project) (*metav1alpha1.Project, error) {
	auth := client.ExtractAuth(ctx)

	oauth2, err := auth.OAuth2()
	if err != nil {
		h.Infof("get oauth2 auth error: %s", err.Error())
		return nil, err
	}

	// create project with harbor api
	if oauth2 != nil {
		h.Infof("has oauth2 auth")
	}

	return &metav1alpha1.Project{}, nil
}

func (h *harbor) ListRepositories(ctx context.Context, option metav1alpha1.ListOptions) (*metav1alpha1.RepositoryList, error) {
	meta := client.ExtraMeta(ctx)

	// list resource with harbor api

	h.Infof("meta: %v", meta)

	return nil, nil
}
