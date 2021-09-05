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

package route

import (
	"fmt"
	"strings"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/katanomi/pkg/plugin/client"
	"github.com/katanomi/pkg/plugin/component/metrics"
	"github.com/katanomi/pkg/plugin/component/tracing"
)

var DefaultFilters = []restful.FilterFunction{
	tracing.Filter,
	metrics.Filter,
	client.AuthFilter,
	client.MetaFilter,
}

// GetPluginWebPath returns a plugin
func GetPluginWebPath(c client.Interface) string {
	return fmt.Sprintf("/plugins/v1alpha1/%s", strings.TrimPrefix(c.Path(), "/"))
}

// Route a service should implement register func to register go restful webservice
type Route interface {
	Register(ws *restful.WebService)
}

// match math route with plugin client
func match(c client.Interface) []Route {
	routes := make([]Route, 0)
	if v, ok := c.(client.ProjectLister); ok {
		routes = append(routes, NewProjectList(v))
	}

	if v, ok := c.(client.ProjectCreator); ok {
		routes = append(routes, NewProjectCreate(v))
	}

	if v, ok := c.(client.ProjectGetter); ok {
		routes = append(routes, NewProjectGet(v))
	}

	if v, ok := c.(client.ResourceLister); ok {
		routes = append(routes, NewResourceList(v))
	}

	if v, ok := c.(client.RepositoryLister); ok {
		routes = append(routes, NewRepositoryList(v))
	}

	if v, ok := c.(client.ArtifactLister); ok {
		routes = append(routes, NewArtifactList(v))
	}

	if v, ok := c.(client.ArtifactGetter); ok {
		routes = append(routes, NewArtifactGet(v))
	}

	if v, ok := c.(client.ArtifactDeleter); ok {
		routes = append(routes, NewArtifactDelete(v))
	}

	if v, ok := c.(client.ScanImage); ok {
		routes = append(routes, NewScanImage(v))
	}

	if v, ok := c.(client.GitRepoFileGetter); ok {
		routes = append(routes, NewGitRepoFileGetter(v))
	}

	if v, ok := c.(client.GitRepoFileCreator); ok {
		routes = append(routes, NewGitRepoFileCreator(v))
	}

	if v, ok := c.(client.GitBranchLister); ok {
		routes = append(routes, NewGitBranchLister(v))
	}

	if v, ok := c.(client.GitBranchCreator); ok {
		routes = append(routes, NewGitBranchCreator(v))
	}

	if v, ok := c.(client.GitCommitGetter); ok {
		routes = append(routes, NewGitCommitGetter(v))
	}

	if v, ok := c.(client.GitPullRequestHandler); ok {
		routes = append(routes, NewGitPullRequestLister(v))
	}

	if v, ok := c.(client.GitPullRequestCommentCreator); ok {
		routes = append(routes, NewGitPullRequestNoteCreator(v))
	}

	if v, ok := c.(client.GitRepositoryLister); ok {
		routes = append(routes, NewGitRepositoryLister(v))
	}

	return routes
}

func GetMethods(c client.Interface) []string {
	// TODO: maybe there is a better way to do this without having
	// to manually add entries
	methods := make([]string, 0, 10)
	if _, ok := c.(client.ProjectLister); ok {
		methods = append(methods, "ListProjects")
	}
	if _, ok := c.(client.ProjectCreator); ok {
		methods = append(methods, "CreateProject")
	}
	if _, ok := c.(client.ResourceLister); ok {
		methods = append(methods, "ListResources")
	}
	if _, ok := c.(client.RepositoryLister); ok {
		methods = append(methods, "ListRepositories")
	}
	if _, ok := c.(client.ArtifactLister); ok {
		methods = append(methods, "ListArtifacts")
	}
	if _, ok := c.(client.ArtifactGetter); ok {
		methods = append(methods, "GetArtifact")
	}
	if _, ok := c.(client.ArtifactDeleter); ok {
		methods = append(methods, "DeleteArtifact")
	}
	if _, ok := c.(client.ScanImage); ok {
		methods = append(methods, "ScanImage")
	}
	if _, ok := c.(client.WebhookRegister); ok {
		methods = append(methods, "CreateWebhook", "UpdateWebhook", "DeleteWebhook")
	}
	if _, ok := c.(client.WebhookResourceDiffer); ok {
		methods = append(methods, "IsSameResource")
	}
	if _, ok := c.(client.WebhookReceiver); ok {
		methods = append(methods, "ReceiveWebhook")
	}
	if _, ok := c.(client.GitRepoFileGetter); ok {
		methods = append(methods, "GetGitRepoFile")
	}
	if _, ok := c.(client.GitRepoFileCreator); ok {
		methods = append(methods, "CreateGitRepoFile")
	}
	if _, ok := c.(client.GitBranchLister); ok {
		methods = append(methods, "ListGitBranch")
	}
	if _, ok := c.(client.GitBranchCreator); ok {
		methods = append(methods, "CreateGitBranch")
	}
	if _, ok := c.(client.GitCommitGetter); ok {
		methods = append(methods, "GetGitCommit")
	}
	if _, ok := c.(client.GitPullRequestHandler); ok {
		methods = append(methods, "ListGitPullRequest", "GetGitPullRequest", "CreatePullRequest")
	}
	if _, ok := c.(client.GitPullRequestCommentCreator); ok {
		methods = append(methods, "CreatePullRequestComment")
	}
	if _, ok := c.(client.GitRepositoryLister); ok {
		methods = append(methods, "ListGitRepository")
	}
	return methods
}

// NewService new service from plugin client
func NewService(c client.Interface, filters ...restful.FilterFunction) (*restful.WebService, error) {
	routes := match(c)
	if len(routes) == 0 {
		return nil, fmt.Errorf("no route for provider %s", c.Path())
	}

	group := &restful.WebService{}
	// adds standard prefix for plugins
	group.Path(GetPluginWebPath(c)).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	for _, filter := range filters {
		group.Filter(filter)
	}

	for _, r := range routes {
		r.Register(group)
	}

	return group, nil
}

// NewDefaultService default service included with metrics,pprof
func NewDefaultService() *restful.WebService {
	routes := []Route{
		NewSystem(),
		NewHealthz(),
	}

	ws := &restful.WebService{}
	for _, each := range routes {
		each.Register(ws)
	}

	return ws
}

//NewDocService go restful api doc
func NewDocService(webservices ...*restful.WebService) *restful.WebService {
	config := restfulspec.Config{
		WebServices: webservices,
		APIPath:     "/openapi.json",
	}
	return restfulspec.NewOpenAPIService(config)
}
