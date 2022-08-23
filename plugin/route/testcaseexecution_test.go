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
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emicklei/go-restful/v3"
	"github.com/katanomi/integrations-testlink/pkg/response"
	. "github.com/onsi/gomega"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestListTestCaseExecutions(t *testing.T) {
	g := NewGomegaWithT(t)

	ws, err := NewService(&TestTestCaseExecutionLister{}, client.MetaFilter)
	g.Expect(err).To(BeNil())

	container := restful.NewContainer()

	container.Add(ws)

	httpRequest, _ := http.NewRequest("GET",
		"/plugins/v1alpha1/test-testcaseexecution/projects/xxx/testplans/yyy/testcases/xxx/executions", nil)
	httpRequest.Header.Set("Accept", "application/json")

	metaData := client.Meta{BaseURL: "http://api.test", Version: "v1"}
	data, _ := json.Marshal(metaData)
	meta := base64.StdEncoding.EncodeToString(data)
	httpRequest.Header.Set(client.PluginMetaHeader, meta)

	httpWriter := httptest.NewRecorder()

	container.Dispatch(httpWriter, httpRequest)
	g.Expect(httpWriter.Code).To(Equal(http.StatusOK))

	testCaseExecutionList := metav1alpha1.TestCaseExecutionList{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &testCaseExecutionList)
	g.Expect(err).To(BeNil())
	g.Expect(testCaseExecutionList.Kind).To(Equal(metav1alpha1.TestCaseExecutionListGVK.Kind))
}

func TestCreateTestCaseExecution(t *testing.T) {
	g := NewGomegaWithT(t)

	ws, err := NewService(&TestTestCaseExecutionCreator{}, client.MetaFilter)
	g.Expect(err).To(BeNil())

	container := restful.NewContainer()

	container.Add(ws)

	body := metav1alpha1.TestCaseExecution{
		Spec: metav1alpha1.TestCaseExecutionSpec{
			Status: response.TestCaseExecutionStatusPassed,
		},
	}
	bodyMarshal, _ := json.Marshal(body)
	httpRequest, _ := http.NewRequest("POST",
		"/plugins/v1alpha1/test-testcaseexecution/projects/xxx/testplans/yyy/testcases/123/executions",
		bytes.NewBuffer(bodyMarshal))
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")

	metaData := client.Meta{BaseURL: "http://api.test", Version: "v1"}
	data, _ := json.Marshal(metaData)
	meta := base64.StdEncoding.EncodeToString(data)
	httpRequest.Header.Set(client.PluginMetaHeader, meta)

	httpWriter := httptest.NewRecorder()

	container.Dispatch(httpWriter, httpRequest)
	g.Expect(httpWriter.Code).To(Equal(http.StatusOK))

	testCaseExecution := metav1alpha1.TestCaseExecution{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &testCaseExecution)
	g.Expect(err).To(BeNil())
	g.Expect(testCaseExecution.Kind).To(Equal(metav1alpha1.TestCaseExecutionGVK.Kind))
}

type TestTestCaseExecutionLister struct {
}

func (t *TestTestCaseExecutionLister) ListTestCaseExecutions(ctx context.Context, params metav1alpha1.TestProjectOptions, options metav1alpha1.ListOptions) (*metav1alpha1.TestCaseExecutionList, error) {
	return &metav1alpha1.TestCaseExecutionList{
		TypeMeta: metav1.TypeMeta{
			Kind: metav1alpha1.TestCaseExecutionListGVK.Kind,
		},
	}, nil
}

func (t *TestTestCaseExecutionLister) Path() string {
	return "test-testcaseexecution"
}

func (t *TestTestCaseExecutionLister) Setup(_ context.Context, _ *zap.SugaredLogger) error {
	return nil
}

type TestTestCaseExecutionCreator struct {
}

func (t *TestTestCaseExecutionCreator) CreateTestCaseExecution(ctx context.Context, params metav1alpha1.TestProjectOptions, payload metav1alpha1.TestCaseExecution) (*metav1alpha1.TestCaseExecution, error) {
	return &metav1alpha1.TestCaseExecution{
		TypeMeta: metav1.TypeMeta{
			Kind: metav1alpha1.TestCaseExecutionGVK.Kind,
		},
	}, nil
}
func (t *TestTestCaseExecutionCreator) Path() string {
	return "test-testcaseexecution"
}

func (t *TestTestCaseExecutionCreator) Setup(_ context.Context, _ *zap.SugaredLogger) error {
	return nil
}
