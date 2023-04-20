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
	"context"
	"fmt"
	"strings"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/katanomi/pkg/plugin/client"
	"github.com/katanomi/pkg/plugin/component/metrics"
	"github.com/katanomi/pkg/plugin/component/tracing"
)

const (
	// Defines the query key value for the search.
	SearchQueryKey = "name" // NOSONAR // ignore: "Key" detected here, make sure this is not a hard-coded credential
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
// Deprecated: replaced by ContextRoute
type Route interface {
	Register(ws *restful.WebService)
}

// ContextRoute is a service should implement register func to register go restful webservice
type ContextRoute interface {
	Register(ctx context.Context, ws *restful.WebService) error
}

// match math route with plugin client
func match(c client.Interface) []Route {
	routes := make([]Route, 0)

	defer func() {
		// if routes length is 0, NewService will return an error.
		if len(routes) != 0 {
			routes = append(routes, NewPluginMethodUnsupport())
		}
	}()

	if v, ok := c.(client.ProjectLister); ok {
		routes = append(routes, NewProjectList(v))
	}

	if v, ok := c.(client.ProjectCreator); ok {
		routes = append(routes, NewProjectCreate(v))
	}

	if v, ok := c.(client.ProjectGetter); ok {
		routes = append(routes, NewProjectGet(v))
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

	if v, ok := c.(client.ArtifactTagDeleter); ok {
		routes = append(routes, NewArtifactTagDelete(v))
	}

	if v, ok := c.(client.ScanImage); ok {
		routes = append(routes, NewScanImage(v))
	}

	if v, ok := c.(client.ImageConfigGetter); ok {
		routes = append(routes, NewImageConifgGetter(v))
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

	if v, ok := c.(client.GitBranchGetter); ok {
		routes = append(routes, NewGitBranchGetter(v))
	}

	if v, ok := c.(client.GitCommitGetter); ok {
		routes = append(routes, NewGitCommitGetter(v))
	}

	if v, ok := c.(client.GitCommitCreator); ok {
		routes = append(routes, NewGitCommitCreator(v))
	}

	if v, ok := c.(client.GitCommitLister); ok {
		routes = append(routes, NewGitCommitLister(v))
	}

	if v, ok := c.(client.GitPullRequestHandler); ok {
		routes = append(routes, NewGitPullRequestLister(v))
	}

	if v, ok := c.(client.GitPullRequestCommentCreator); ok {
		routes = append(routes, NewGitPullRequestNoteCreator(v))
	}

	if v, ok := c.(client.GitPullRequestCommentUpdater); ok {
		routes = append(routes, NewGitPullRequestNoteUpdater(v))
	}

	if v, ok := c.(client.GitPullRequestCommentLister); ok {
		routes = append(routes, NewGitPullRequestCommentLister(v))
	}

	if v, ok := c.(client.GitRepositoryLister); ok {
		routes = append(routes, NewGitRepositoryLister(v))
	}

	if v, ok := c.(client.GitRepositoryGetter); ok {
		routes = append(routes, NewGitRepositoryGetter(v))
	}

	if v, ok := c.(client.GitCommitCommentLister); ok {
		routes = append(routes, NewGitCommitCommentLister(v))
	}

	if v, ok := c.(client.GitCommitCommentCreator); ok {
		routes = append(routes, NewGitCommitCommentCreator(v))
	}

	if v, ok := c.(client.GitCommitStatusLister); ok {
		routes = append(routes, NewGitCommitStatusLister(v))
	}

	if v, ok := c.(client.GitCommitStatusCreator); ok {
		routes = append(routes, NewGitCommitStatusCreator(v))
	}

	if v, ok := c.(client.GitRepositoryTagGetter); ok {
		routes = append(routes, NewGitRepositoryTagGetter(v))
	}
	if v, ok := c.(client.GitRepositoryTagLister); ok {
		routes = append(routes, NewGitRepositoryTagLister(v))
	}

	if v, ok := c.(client.CodeQualityGetter); ok {
		routes = append(routes, NewCodeQualityGetter(v))
	}

	if v, ok := c.(client.BlobStoreLister); ok {
		routes = append(routes, NewBlobStoreLister(v))
	}

	if v, ok := c.(client.GitRepositoryFileTreeGetter); ok {
		routes = append(routes, NewGitRepositoryFileTreeGetter(v))
	}

	authCheck, ok := c.(client.AuthChecker)
	// uses a default implementation that returns an Unknown allowed result
	// with an NotImplemented reason
	if !ok {
		authCheck = NewDefaultAuthCheckImplementation()
	}
	routes = append(routes, NewAuthCheck(authCheck))

	if v, ok := c.(client.AuthTokenGenerator); ok {
		routes = append(routes, NewAuthToken(v))
	}

	if v, ok := c.(client.IssueLister); ok {
		routes = append(routes, NewIssueList(v))
	}

	if v, ok := c.(client.IssueGetter); ok {
		routes = append(routes, NewIssueGet(v))
	}

	if v, ok := c.(client.IssueAttributeGetter); ok {
		routes = append(routes, NewIssueAttributeGet(v))
	}

	if v, ok := c.(client.IssueBranchLister); ok {
		routes = append(routes, NewIssueBranchList(v))
	}

	if v, ok := c.(client.IssueBranchCreator); ok {
		routes = append(routes, NewIssueBranchCreate(v))
	}

	if v, ok := c.(client.IssueBranchDeleter); ok {
		routes = append(routes, NewIssueBranchDelete(v))
	}

	if v, ok := c.(client.ProjectUserLister); ok {
		routes = append(routes, NewProjectUserList(v))
	}
	if v, ok := c.(client.TestPlanLister); ok {
		routes = append(routes, NewTestPlanList(v))
	}
	if v, ok := c.(client.TestPlanGetter); ok {
		routes = append(routes, NewTestPlanGetter(v))
	}
	if v, ok := c.(client.TestCaseLister); ok {
		routes = append(routes, NewTestCaseLister(v))
	}
	if v, ok := c.(client.TestCaseGetter); ok {
		routes = append(routes, NewTestCaseGetter(v))
	}
	if v, ok := c.(client.TestModuleLister); ok {
		routes = append(routes, NewTestModuleLister(v))
	}
	if v, ok := c.(client.TestCaseExecutionLister); ok {
		routes = append(routes, NewTestCaseExecutionLister(v))
	}
	if v, ok := c.(client.TestCaseExecutionCreator); ok {
		routes = append(routes, NewTestCaseExecutionCreator(v))
	}
	if v, ok := c.(client.LivenessChecker); ok {
		routes = append(routes, NewLivenessCheck(v))
	}
	if v, ok := c.(client.Initializer); ok {
		routes = append(routes, NewInitializer(v))
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
	if _, ok := c.(client.ArtifactTagDeleter); ok {
		methods = append(methods, "DeleteArtifactTag")
	}
	if _, ok := c.(client.ScanImage); ok {
		methods = append(methods, "ScanImage")
	}
	if _, ok := c.(client.ImageConfigGetter); ok {
		methods = append(methods, "GetImageConfig")
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
	if _, ok := c.(client.GitBranchGetter); ok {
		methods = append(methods, "GetGitBranch")
	}
	if _, ok := c.(client.GitBranchCreator); ok {
		methods = append(methods, "CreateGitBranch")
	}
	if _, ok := c.(client.GitCommitGetter); ok {
		methods = append(methods, "GetGitCommit")
	}
	if _, ok := c.(client.GitCommitCreator); ok {
		methods = append(methods, "CreateGitCommit")
	}
	if _, ok := c.(client.GitCommitLister); ok {
		methods = append(methods, "ListGitCommit")
	}
	if _, ok := c.(client.GitPullRequestHandler); ok {
		methods = append(methods, "ListGitPullRequest", "GetGitPullRequest", "CreatePullRequest")
	}
	if _, ok := c.(client.GitPullRequestCommentCreator); ok {
		methods = append(methods, "CreatePullRequestComment")
	}
	if _, ok := c.(client.GitPullRequestCommentLister); ok {
		methods = append(methods, "ListPullRequestComment")
	}
	if _, ok := c.(client.GitRepositoryLister); ok {
		methods = append(methods, "ListGitRepository")
	}
	if _, ok := c.(client.GitRepositoryGetter); ok {
		methods = append(methods, "GetGitRepository")
	}
	if _, ok := c.(client.GitRepositoryFileTreeGetter); ok {
		methods = append(methods, "GetGitRepositoryFileTree")
	}
	if _, ok := c.(client.GitCommitCommentLister); ok {
		methods = append(methods, "ListGitCommitComment")
	}
	if _, ok := c.(client.GitCommitCommentCreator); ok {
		methods = append(methods, "CreateGitCommitComment")
	}
	if _, ok := c.(client.GitCommitStatusLister); ok {
		methods = append(methods, "ListGitCommitStatus")
	}
	if _, ok := c.(client.GitCommitStatusCreator); ok {
		methods = append(methods, "CreateGitCommitStatus")
	}
	if _, ok := c.(client.GitRepositoryTagGetter); ok {
		methods = append(methods, "GetGitRepositoryTag")
	}
	if _, ok := c.(client.GitRepositoryTagLister); ok {
		methods = append(methods, "ListGitRepositoryTag")
	}
	if _, ok := c.(client.CodeQualityGetter); ok {
		methods = append(methods, "GetCodeQuality", "GetCodeQualityOverviewByBranch", "GetCodeQualityLineCharts", "GetOverview", "GetSummaryByTaskID")
	}
	if _, ok := c.(client.BlobStoreLister); ok {
		methods = append(methods, "ListBlobStores")
	}
	if _, ok := c.(client.AuthChecker); ok {
		methods = append(methods, "AuthCheck")
	}
	if _, ok := c.(client.IssueLister); ok {
		methods = append(methods, "ListIssues")
	}
	if _, ok := c.(client.IssueGetter); ok {
		methods = append(methods, "GetIssue")
	}
	if _, ok := c.(client.IssueAttributeGetter); ok {
		methods = append(methods, "GetIssueAttribute")
	}
	if _, ok := c.(client.IssueBranchLister); ok {
		methods = append(methods, "ListIssueBranches")
	}
	if _, ok := c.(client.IssueBranchCreator); ok {
		methods = append(methods, "CreateIssueBranch")
	}
	if _, ok := c.(client.IssueBranchDeleter); ok {
		methods = append(methods, "DeleteIssueBranch")
	}
	if _, ok := c.(client.ProjectUserLister); ok {
		methods = append(methods, "ListProjectUsers")
	}
	if _, ok := c.(client.TestPlanLister); ok {
		methods = append(methods, "ListTestPlans")
	}
	if _, ok := c.(client.TestPlanGetter); ok {
		methods = append(methods, "GetTestPlan")
	}
	if _, ok := c.(client.TestCaseLister); ok {
		methods = append(methods, "ListTestCases")
	}
	if _, ok := c.(client.TestCaseGetter); ok {
		methods = append(methods, "GetTestCase")
	}
	if _, ok := c.(client.TestModuleLister); ok {
		methods = append(methods, "ListTestModules")
	}
	if _, ok := c.(client.TestCaseExecutionLister); ok {
		methods = append(methods, "ListTestCaseExecutions")
	}
	if _, ok := c.(client.TestCaseExecutionCreator); ok {
		methods = append(methods, "CreateTestCaseExecution")
	}
	if _, ok := c.(client.LivenessChecker); ok {
		methods = append(methods, "CheckAlive")
	}
	if _, ok := c.(client.Initializer); ok {
		methods = append(methods, "Initialize")
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

// NewDocService go restful api doc
func NewDocService(webservices ...*restful.WebService) *restful.WebService {
	config := restfulspec.Config{
		WebServices: webservices,
		APIPath:     "/openapi.json",
	}
	return restfulspec.NewOpenAPIService(config)
}
