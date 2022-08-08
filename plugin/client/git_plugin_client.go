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

package client

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	corev1 "k8s.io/api/core/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/logging"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	kclient "github.com/katanomi/pkg/client"
	ksecret "github.com/katanomi/pkg/secret"
)

// GitPluginClient client for plugins
type GitPluginClient struct {
	*PluginClient

	// GitRepo Repo base info, such as project, repository
	GitRepo metav1alpha1.GitRepo

	// ClassObject store the integration class object
	// +optional
	ClassObject ctrlclient.Object
}

func (p *GitPluginClient) WithGitRepo(gitRepo metav1alpha1.GitRepo) *GitPluginClient {
	p.GitRepo = gitRepo
	return p
}

func (p *GitPluginClient) WithClassObject(object ctrlclient.Object) *GitPluginClient {
	p.ClassObject = object
	return p
}

// GenerateGitPluginClient generate git plugin client params
func GenerateGitPluginClient(ctx context.Context, secretRef *corev1.ObjectReference,
	gitRepoURL, integrationClassName string, classAddress *duckv1.Addressable) (
	gpclient *GitPluginClient, err error) {

	log := logging.FromContext(ctx)
	pclient := PluginClientValue(ctx)
	if pclient == nil {
		pclient = NewPluginClient()
	} else {
		pclient = pclient.Clone()
	}

	if secretRef != nil {
		var secret *corev1.Secret
		clt := kclient.Client(ctx)
		if clt == nil {
			err = fmt.Errorf("cannot get client from ctx")
			return
		}
		secret, err = ksecret.GetSecretByRefOrLabel(ctx, clt, secretRef)
		if err != nil {
			err = fmt.Errorf("get secret for version steam failed: %w", err)
			return
		}
		pclient = pclient.WithSecret(*secret)
	}

	if integrationClassName != "" {
		pclient = pclient.WithIntegrationClassName(integrationClassName)
	}
	if classAddress != nil {
		pclient = pclient.WithClassAddress(classAddress)
	}

	gitAddress, gitRepo, err := GetGitRepoInfo(gitRepoURL)
	if err != nil {
		err = fmt.Errorf("get git repo info failed: %w", err)
		return
	}
	meta := Meta{
		BaseURL: gitAddress,
	}
	gpclient = pclient.
		WithMeta(meta).
		GitPluginClient().
		WithGitRepo(gitRepo)

	log.Debugw("generate git plugin client", "BaseURL", gitAddress, "GitRepo", gitRepo,
		"ClassAddress", classAddress, "IntegrationClassName", integrationClassName)
	return gpclient, nil
}

// GetGitRepoInfo try to get the project and repository from the git address.
func GetGitRepoInfo(gitAddress string) (host string, gitRepo metav1alpha1.GitRepo, err error) {
	gitAddress = strings.TrimSuffix(gitAddress, ".git")
	URL, err := url.Parse(gitAddress)
	if err != nil {
		return
	}

	projectRepo := strings.Split(strings.Trim(URL.Path, "/"), "/")
	if len(projectRepo) < 2 {
		err = fmt.Errorf("invaild git address %s which should have project and repository", gitAddress)
		return
	}

	host = fmt.Sprintf("%s://%s", URL.Scheme, URL.Host)
	gitRepo.Project = strings.Join(projectRepo[:len(projectRepo)-1], "/")
	gitRepo.Repository = projectRepo[len(projectRepo)-1]
	return
}
