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
	"net/http"
	"time"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	kerrors "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/plugin/client"
	"k8s.io/apimachinery/pkg/api/errors"
)

type codeQualityGetter struct {
	impl client.CodeQualityGetter
	tags []string
}

//NewCodeQualityGetter create a get codeQuality route with plugin client
func NewCodeQualityGetter(impl client.CodeQualityGetter) Route {
	return &codeQualityGetter{
		tags: []string{"codeQuality"},
		impl: impl,
	}
}

func (c *codeQualityGetter) Register(ws *restful.WebService) {
	ws.Route(
		ListOptionsDocs(
			ws.GET("/codeQuality/{project-key}").To(c.GetCodeQuality).
				Param(ws.PathParameter("project-id", "identifier of the project").DataType("string")).
				// docs
				Doc("GetCodeQuality").
				Metadata(restfulspec.KeyOpenAPITags, c.tags).
				Returns(http.StatusOK, "OK", metav1alpha1.CodeQuality{}),
		),
	)
	ws.Route(
		ws.GET("/codeQuality/{project-key}/branches/{branch}").To(c.GetCodeQualityOverviewByBranch).
			Param(ws.PathParameter("project-id", "identifier of the project").DataType("string")).
			Param(ws.PathParameter("branch", "branch name").DataType("string")).
			Doc("GetCodeQualityOverviewByBranch").
			Metadata(restfulspec.KeyOpenAPITags, c.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.CodeQuality{}),
	)
	ws.Route(
		ws.GET("/codeQuality/{project-key}/branches/{branch}/lineCharts").To(c.GetCodeQualityLineCharts).
			Param(ws.PathParameter("project-id", "identifier of the project").DataType("string")).
			Param(ws.PathParameter("branch", "branch name").DataType("string")).
			Param(ws.QueryParameter("metricKeys", "metric keys").DataType("string")).
			Doc("GetCodeQualityLineCharts").
			Metadata(restfulspec.KeyOpenAPITags, c.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.CodeQualityLineChart{}),
	)
	ws.Route(
		ws.GET("/codeQuality").To(c.GetOverview).
			Doc("GetOverview").
			Metadata(restfulspec.KeyOpenAPITags, c.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.CodeQualityProjectOverview{}),
	)
}

// GetCodeQuality http handler for get code quality
func (c *codeQualityGetter) GetCodeQuality(request *restful.Request, response *restful.Response) {
	projectKey := request.PathParameter("project-key")
	codeQuality, err := c.impl.GetCodeQuality(request.Request.Context(), projectKey)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, codeQuality)
}

// GetCodeQuality http handler for get code quality
func (c *codeQualityGetter) GetCodeQualityOverviewByBranch(request *restful.Request, response *restful.Response) {
	projectKey := request.PathParameter("project-key")
	branchKey := request.PathParameter("branch")
	codeQuality, err := c.impl.GetCodeQualityOverviewByBranch(request.Request.Context(), metav1alpha1.CodeQualityBaseOption{ProjectKey: projectKey, BranchKey: branchKey})
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, codeQuality)
}

// GetCodeQuality http handler for get code quality
func (c *codeQualityGetter) GetCodeQualityLineCharts(request *restful.Request, response *restful.Response) {
	projectKey := request.PathParameter("project-key")
	branchKey := request.PathParameter("branch")
	metricKeys := request.QueryParameter("metricKeys")
	startTime := request.QueryParameter("startTime")
	endTime := request.QueryParameter("completionTime")
	param := metav1alpha1.CodeQualityLineChartOption{
		CodeQualityBaseOption: metav1alpha1.CodeQualityBaseOption{
			ProjectKey: projectKey,
			BranchKey:  branchKey,
		},
		Metrics:        metricKeys,
		StartTime:      nil,
		CompletionTime: nil,
	}
	if startTime != "" {
		start, err := time.Parse(time.RFC3339, startTime)
		if err != nil {
			kerrors.HandleError(request, response, errors.NewBadRequest(err.Error()))
			return
		}
		param.StartTime = &start
	}
	if endTime != "" {
		end, err := time.Parse(time.RFC3339, endTime)
		if err != nil {
			kerrors.HandleError(request, response, errors.NewBadRequest(err.Error()))
			return
		}
		param.CompletionTime = &end
	}
	codeQuality, err := c.impl.GetCodeQualityLineCharts(request.Request.Context(), param)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, codeQuality)
}

func (c *codeQualityGetter) GetOverview(request *restful.Request, response *restful.Response) {
	result, err := c.impl.GetOverview(request.Request.Context())
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}
