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

package route

import (
	"net/http"

	"github.com/katanomi/pkg/plugin/path"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	kerrors "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/plugin/client"
)

type testCaseExecutionLister struct {
	impl client.TestCaseExecutionLister
	tags []string
}

// NewTestCaseExecutionLister creates a list testCaseExecution route with plugin client
func NewTestCaseExecutionLister(impl client.TestCaseExecutionLister) Route {
	return &testCaseExecutionLister{
		tags: []string{"projects", "testCaseExecution"},
		impl: impl,
	}
}

func (r *testCaseExecutionLister) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "project belong to integraion")
	testPlanParam := ws.PathParameter("testplanid", "testPlan belong to project")
	testCaseParam := ws.PathParameter("testcaseid", "testCase belong to testPlan")
	buildParam := ws.QueryParameter("buildID", "build id related to testPlan")
	ws.Route(
		ListOptionsDocs(
			ws.GET("/projects/{project:*}/testplans/{testplanid:*}/testcases/{testcaseid}/executions").To(r.
				ListTestCaseExecutions).
				// docs
				Doc("ListTestCaseExecutions").Param(projectParam).
				Param(testPlanParam).
				Param(testCaseParam).
				Param(buildParam).
				Metadata(restfulspec.KeyOpenAPITags, r.tags).
				Returns(http.StatusOK, "OK", metav1alpha1.TestCaseExecutionList{}),
		),
	)
}

// ListTestCaseExecutions http handler for list testCaseExecution
func (r *testCaseExecutionLister) ListTestCaseExecutions(request *restful.Request, response *restful.Response) {
	option := GetListOptionsFromRequest(request)

	pathParams := metav1alpha1.TestProjectOptions{
		Project:    path.Parameter(request, "project"),
		TestPlanID: path.Parameter(request, "testplanid"),
		TestCaseID: path.Parameter(request, "testcaseid"),
		BuildID:    request.QueryParameter("buildID"),
	}
	testCaseExecutions, err := r.impl.ListTestCaseExecutions(request.Request.Context(), pathParams, option)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, testCaseExecutions)
}

type testCaseExecutionCreator struct {
	impl client.TestCaseExecutionCreator
	tags []string
}

// NewTestCaseExecutionCreator creates a creating testCaseExecution route with plugin client
func NewTestCaseExecutionCreator(impl client.TestCaseExecutionCreator) Route {
	return &testCaseExecutionCreator{
		tags: []string{"projects", "testCaseExecution"},
		impl: impl,
	}
}

func (r *testCaseExecutionCreator) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "project belong to integraion")
	testPlanParam := ws.PathParameter("testplanid", "testPlan belong to project")
	testCaseParam := ws.PathParameter("testcaseid", "testCase belong to testPlan")
	buildParam := ws.QueryParameter("buildID", "build id related to testPlan")
	ws.Route(
		ws.POST("/projects/{project:*}/testplans/{testplanid:*}/testcases/{testcaseid}/executions").
			To(r.CreateTestCaseExecution).
			// docs
			Doc("CreateTestCaseExecution").
			Param(projectParam).
			Param(testPlanParam).
			Param(testCaseParam).
			Param(buildParam).
			Metadata(restfulspec.KeyOpenAPITags, r.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.TestCaseExecution{}),
	)
}

// CreateTestCaseExecution http handler for creating testCaseExecution
func (r *testCaseExecutionCreator) CreateTestCaseExecution(request *restful.Request, response *restful.Response) {
	params := metav1alpha1.TestProjectOptions{
		Project:    path.Parameter(request, "project"),
		TestPlanID: path.Parameter(request, "testplanid"),
		TestCaseID: path.Parameter(request, "testcaseid"),
		BuildID:    request.QueryParameter("buildID"),
	}

	var payload metav1alpha1.TestCaseExecution
	if err := request.ReadEntity(&payload); err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	// assign case status if step status is omitted
	for idx, step := range payload.Spec.Steps {
		if step.Status == "" {
			payload.Spec.Steps[idx].Status = payload.Spec.Status
		}
	}

	testCaseExecution, err := r.impl.CreateTestCaseExecution(request.Request.Context(), params, payload)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, testCaseExecution)
}
