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

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	kerrors "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/plugin/client"
)

type testCaseLister struct {
	impl client.TestCaseLister
	tags []string
}

// NewTestCaseLister creates a list testCase route with plugin client
func NewTestCaseLister(impl client.TestCaseLister) Route {
	return &testCaseLister{
		tags: []string{"projects", "testCase"},
		impl: impl,
	}
}

func (r *testCaseLister) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "testCase belong to integraion")
	testPlanIDParam := ws.PathParameter("testplanid", "test plan id")
	buildIDParam := ws.QueryParameter("buildID", "test plan id")
	ws.Route(
		ListOptionsDocs(
			ws.GET("/projects/{project:*}/testplans/{testplanid:*}/testcases").To(r.ListTestCases).
				// docs
				Doc("ListTestCases").
				Param(projectParam).
				Param(testPlanIDParam).
				Param(buildIDParam).
				Metadata(restfulspec.KeyOpenAPITags, r.tags).
				Returns(http.StatusOK, "OK", metav1alpha1.TestCaseList{}),
		),
	)
}

// ListTestCases http handler for list testCase
func (r *testCaseLister) ListTestCases(request *restful.Request, response *restful.Response) {
	option := GetListOptionsFromRequest(request)

	pathParams := metav1alpha1.TestProjectOptions{
		Project:    request.PathParameter("project"),
		TestPlanID: request.PathParameter("testplanid"),
		BuildID:    request.QueryParameter("buildID"),
	}
	testCases, err := r.impl.ListTestCases(request.Request.Context(), pathParams, option)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, testCases)
}

type testCaseGetter struct {
	impl client.TestCaseGetter
	tags []string
}

// NewTestCaseGetter creates a list testCase route with plugin client
func NewTestCaseGetter(impl client.TestCaseGetter) Route {
	return &testCaseGetter{
		tags: []string{"projects", "testCase"},
		impl: impl,
	}
}

func (r *testCaseGetter) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "project belong to integraion")
	testPlanIDParam := ws.PathParameter("testplanid", "test plan id")
	testCaseIDParam := ws.PathParameter("testcaseid", "test case id")
	buildIDParam := ws.QueryParameter("buildID", "test case id")
	ws.Route(
		ListOptionsDocs(
			ws.GET("/projects/{project:*}/testplans/{testplanid:*}/testcases/{testcaseid:*}").To(r.GetTestCase).
				// docs
				Doc("GetTestCase").
				Param(projectParam).
				Param(testPlanIDParam).
				Param(testCaseIDParam).
				Param(buildIDParam).
				Metadata(restfulspec.KeyOpenAPITags, r.tags).
				Returns(http.StatusOK, "OK", metav1alpha1.TestCase{}),
		),
	)
}

// GetTestCase http handler for getting testCase
func (r *testCaseGetter) GetTestCase(request *restful.Request, response *restful.Response) {
	pathParams := metav1alpha1.TestProjectOptions{
		Project:    request.PathParameter("project"),
		TestPlanID: request.PathParameter("testplanid"),
		TestCaseID: request.PathParameter("testcaseid"),
		BuildID:    request.QueryParameter("buildID"),
	}
	testCase, err := r.impl.GetTestCase(request.Request.Context(), pathParams)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, testCase)
}
