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
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emicklei/go-restful/v3"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
	"go.uber.org/zap"
)

func TestListTestCases(t *testing.T) {
	g := NewGomegaWithT(t)

	ws, err := NewService(&TestTestCaseLister{}, client.MetaFilter)
	g.Expect(err).To(BeNil())

	container := restful.NewContainer()

	container.Add(ws)

	httpRequest, _ := http.NewRequest("GET",
		"/plugins/v1alpha1/test-testcase/projects/xxx/testplans/yyy/testcases", nil)
	httpRequest.Header.Set("Accept", "application/json")

	metaData := client.Meta{BaseURL: "http://api.test", Version: "v1"}
	data, _ := json.Marshal(metaData)
	meta := base64.StdEncoding.EncodeToString(data)
	httpRequest.Header.Set(client.PluginMetaHeader, meta)

	httpWriter := httptest.NewRecorder()

	container.Dispatch(httpWriter, httpRequest)
	g.Expect(httpWriter.Code).To(Equal(http.StatusOK))

	testCaseList := metav1alpha1.TestCaseList{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &testCaseList)
	g.Expect(err).To(BeNil())
	g.Expect(testCaseList.Kind).To(Equal(metav1alpha1.TestCaseListGVK.Kind))
}

func TestGetTestCase(t *testing.T) {
	g := NewGomegaWithT(t)

	ws, err := NewService(&TestTestCaseGetter{}, client.MetaFilter)
	g.Expect(err).To(BeNil())

	container := restful.NewContainer()

	container.Add(ws)

	httpRequest, _ := http.NewRequest("GET",
		"/plugins/v1alpha1/test-testcase/projects/xxx/testplans/yyy/testcases/xxx", nil)
	httpRequest.Header.Set("Accept", "application/json")

	metaData := client.Meta{BaseURL: "http://api.test", Version: "v1"}
	data, _ := json.Marshal(metaData)
	meta := base64.StdEncoding.EncodeToString(data)
	httpRequest.Header.Set(client.PluginMetaHeader, meta)

	httpWriter := httptest.NewRecorder()

	container.Dispatch(httpWriter, httpRequest)
	g.Expect(httpWriter.Code).To(Equal(http.StatusOK))

	testCase := metav1alpha1.TestCase{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &testCase)
	g.Expect(err).To(BeNil())
	g.Expect(testCase.Kind).To(Equal(metav1alpha1.TestCaseGVK.Kind))
}

type TestTestCaseLister struct {
}

func (t *TestTestCaseLister) ListTestCases(ctx context.Context, params metav1alpha1.TestProjectOptions, options metav1alpha1.ListOptions) (*metav1alpha1.TestCaseList, error) {
	return &metav1alpha1.TestCaseList{
		TypeMeta: metav1.TypeMeta{
			Kind: metav1alpha1.TestCaseListGVK.Kind,
		},
	}, nil
}

func (t *TestTestCaseLister) Path() string {
	return "test-testcase"
}

func (t *TestTestCaseLister) Setup(_ context.Context, _ *zap.SugaredLogger) error {
	return nil
}

type TestTestCaseGetter struct {
}

func (t *TestTestCaseGetter) GetTestCase(ctx context.Context, params metav1alpha1.TestProjectOptions) (*metav1alpha1.TestCase, error) {
	return &metav1alpha1.TestCase{
		TypeMeta: metav1.TypeMeta{
			Kind: metav1alpha1.TestCaseGVK.Kind,
		},
	}, nil
}
func (t *TestTestCaseGetter) Path() string {
	return "test-testcase"
}

func (t *TestTestCaseGetter) Setup(_ context.Context, _ *zap.SugaredLogger) error {
	return nil
}
