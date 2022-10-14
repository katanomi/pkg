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

type testPlanLister struct {
	impl client.TestPlanLister
	tags []string
}

// NewTestPlanList creates a list testPlan route with plugin client
func NewTestPlanList(impl client.TestPlanLister) Route {
	return &testPlanLister{
		tags: []string{"projects", "testPlan"},
		impl: impl,
	}
}

func (r *testPlanLister) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "testPlan belong to integraion")
	ws.Route(
		ListOptionsDocs(
			ws.GET("/projects/{project:*}/testplans").To(r.ListTestPlans).
				// docs
				Doc("ListTestPlans").Param(projectParam).
				Metadata(restfulspec.KeyOpenAPITags, r.tags).
				Returns(http.StatusOK, "OK", metav1alpha1.TestPlanList{}),
		),
	)
}

// ListTestPlans http handler for list testPlan
func (r *testPlanLister) ListTestPlans(request *restful.Request, response *restful.Response) {
	option := GetListOptionsFromRequest(request)
	pathParams := metav1alpha1.TestProjectOptions{
		Project: path.Parameter(request, "project"),
	}
	testPlans, err := r.impl.ListTestPlans(request.Request.Context(), pathParams, option)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, testPlans)
}

type testPlanGetter struct {
	impl client.TestPlanGetter
	tags []string
}

// NewTestPlanGetter creates a list testPlan route with plugin client
func NewTestPlanGetter(impl client.TestPlanGetter) Route {
	return &testPlanGetter{
		tags: []string{"projects", "testPlan"},
		impl: impl,
	}
}

func (r *testPlanGetter) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "testPlan project belong to integraion")
	testPlanIDParam := ws.PathParameter("testPlanID", "test plan id")
	buildIDParam := ws.QueryParameter("buildID", "test build id")
	ws.Route(
		ListOptionsDocs(
			ws.GET("/projects/{project:*}/testplans/{testplanid}").To(r.GetTestPlan).
				// docs
				Doc("GetTestCase").
				Param(projectParam).
				Param(testPlanIDParam).
				Param(buildIDParam).
				Metadata(restfulspec.KeyOpenAPITags, r.tags).
				Returns(http.StatusOK, "OK", metav1alpha1.TestCaseList{}),
		),
	)
}

// GetTestPlan http handler for getting testPlan
func (r *testPlanGetter) GetTestPlan(request *restful.Request, response *restful.Response) {
	pathParams := metav1alpha1.TestProjectOptions{
		Project:    path.Parameter(request, "project"),
		TestPlanID: path.Parameter(request, "testplanid"),
		BuildID:    request.QueryParameter("buildID"),
	}
	testPlan, err := r.impl.GetTestPlan(request.Request.Context(), pathParams)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, testPlan)
}
