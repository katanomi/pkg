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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMatch(t *testing.T) {
	testCases := []struct {
		c   client.Interface
		len int
	}{
		{
			c:   &TestProjectList{},
			len: 1,
		},
		{
			c:   &TestProjectCreate{},
			len: 1,
		},
		{
			c:   &TestResourceList{},
			len: 1,
		},
		{
			c:   &TestProjectListCreate{},
			len: 2,
		},
	}

	g := NewGomegaWithT(t)

	for i, item := range testCases {
		t.Run(fmt.Sprintf("test-%d", i+1), func(t *testing.T) {
			routes := match(item.c)
			g.Expect(len(routes)).To(Equal(item.len))
		})
	}
}

func TestRegister(t *testing.T) {
	testCases := []struct {
		c    client.Interface
		path string
	}{
		{
			c:    &TestProjectList{},
			path: "/projects",
		},
		{
			c:    &TestProjectCreate{},
			path: "/projects",
		},
		{
			c:    &TestResourceList{},
			path: "/resources",
		},
		{
			c:    &TestProjectListCreate{},
			path: "/projects",
		},
	}

	g := NewGomegaWithT(t)

	for i, item := range testCases {
		t.Run(fmt.Sprintf("test-%d", i+1), func(t *testing.T) {
			routes := match(item.c)

			ws := &restful.WebService{}
			routes[0].Register(ws)

			g.Expect(ws.Routes()[0].Path).To(Equal(item.path))
		})
	}
}

func TestNewService(t *testing.T) {
	testCases := []client.Interface{
		&TestProjectList{},
		&TestProjectCreate{},
		&TestResourceList{},
		&TestProjectListCreate{},
	}

	g := NewGomegaWithT(t)

	for _, c := range testCases {
		t.Run(c.Path(), func(t *testing.T) {
			ws, err := NewService(c)

			g.Expect(err).To(BeNil())
			g.Expect(ws.RootPath()).To(Equal("/" + c.Path()))
		})
	}
}

func TestProjectListNoMeta(t *testing.T) {
	g := NewGomegaWithT(t)

	ws, err := NewService(&TestProjectList{}, client.MetaFilter)
	g.Expect(err).To(BeNil())

	container := restful.NewContainer()
	container.Add(ws)

	httpRequest, _ := http.NewRequest("GET", "/test-1/projects", nil)
	httpRequest.Header.Set("Accept", "*/*")
	httpWriter := httptest.NewRecorder()

	container.Dispatch(httpWriter, httpRequest)

	g.Expect(httpWriter.Code).To(Equal(http.StatusBadRequest))
}

func TestProjectListWithMeta(t *testing.T) {
	g := NewGomegaWithT(t)

	ws, err := NewService(&TestProjectList{}, client.MetaFilter)
	g.Expect(err).To(BeNil())

	container := restful.NewContainer()

	container.Add(ws)

	httpRequest, _ := http.NewRequest("GET", "/test-1/projects", nil)
	httpRequest.Header.Set("Accept", "application/json")

	metaData := client.Meta{BaseURL: "http://api.test", Version: "v1"}
	data, _ := json.Marshal(metaData)
	meta := base64.StdEncoding.EncodeToString(data)
	httpRequest.Header.Set(client.PluginMetaHeader, meta)

	httpWriter := httptest.NewRecorder()

	container.Dispatch(httpWriter, httpRequest)
	g.Expect(httpWriter.Code).To(Equal(http.StatusOK))

	list := metav1alpha1.ProjectList{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &list)
	g.Expect(err).To(BeNil())
	g.Expect(list.Items).ToNot(BeEmpty())
}

type TestProjectList struct {
}

func (t *TestProjectList) Path() string {
	return "test-1"
}

func (t *TestProjectList) ListProjects(ctx context.Context, option metav1alpha1.ListOptions) (*metav1alpha1.ProjectList, error) {
	return &metav1alpha1.ProjectList{
		Items: []metav1alpha1.Project{
			{
				TypeMeta: metav1.TypeMeta{
					Kind: "1",
				},
			},
			{
				TypeMeta: metav1.TypeMeta{
					Kind: "2",
				},
			},
		},
	}, nil
}

type TestProjectCreate struct {
}

func (t *TestProjectCreate) Path() string {
	return "test-2"
}

func (t *TestProjectCreate) CreateProject(ctx context.Context, project *metav1alpha1.Project) (*metav1alpha1.Project, error) {
	return &metav1alpha1.Project{}, nil
}

type TestResourceList struct {
}

func (t *TestResourceList) Path() string {
	return "test-3"
}

func (t *TestResourceList) ListResources(ctx context.Context, option metav1alpha1.ListOptions) (*metav1alpha1.ResourceList, error) {
	return &metav1alpha1.ResourceList{}, nil
}

type TestProjectListCreate struct {
}

func (t *TestProjectListCreate) Path() string {
	return "test-4"
}

func (t *TestProjectListCreate) ListProjects(ctx context.Context, option metav1alpha1.ListOptions) (*metav1alpha1.ProjectList, error) {
	return &metav1alpha1.ProjectList{}, nil
}

func (t *TestProjectListCreate) CreateProject(ctx context.Context, project *metav1alpha1.Project) (*metav1alpha1.Project, error) {
	return &metav1alpha1.Project{}, nil
}
