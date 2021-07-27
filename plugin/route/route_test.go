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
	"github.com/katanomi/pkg/plugin/client"
	. "github.com/onsi/gomega"
)

func TestMatch(t *testing.T) {
	testCases := []struct {
		c   client.PluginClient
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
		c    client.PluginClient
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
	testCases := []client.PluginClient{
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
			g.Expect(ws.RootPath()).To(Equal(c.Path()))
		})
	}
}

func TestProjectListNoMeta(t *testing.T) {
	g := NewGomegaWithT(t)

	ws, err := NewService(&TestProjectList{})
	g.Expect(err).To(BeNil())

	restful.DefaultContainer.Add(ws)

	httpRequest, _ := http.NewRequest("GET", "/test-1/projects", nil)
	httpRequest.Header.Set("Accept", "*/*")
	httpWriter := httptest.NewRecorder()

	restful.DefaultContainer.Dispatch(httpWriter, httpRequest)

	g.Expect(httpWriter.Code).To(Equal(http.StatusBadRequest))
}

func TestProjectListWithMeta(t *testing.T) {
	g := NewGomegaWithT(t)

	ws, err := NewService(&TestProjectList{})
	g.Expect(err).To(BeNil())

	restful.DefaultContainer.Add(ws)

	httpRequest, _ := http.NewRequest("GET", "/test-1/projects", nil)
	httpRequest.Header.Set("Accept", "application/json")

	meta := base64.StdEncoding.EncodeToString([]byte("http://api.test:v1"))
	httpRequest.Header.Set("X-plugin-Meta", meta)

	httpWriter := httptest.NewRecorder()

	restful.DefaultContainer.Dispatch(httpWriter, httpRequest)
	g.Expect(httpWriter.Code).To(Equal(http.StatusOK))

	list := client.ProjectList{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &list)
	g.Expect(err).To(BeNil())
	g.Expect(list.Items).ToNot(BeEmpty())
}

type TestProjectList struct {
}

func (t *TestProjectList) Path() string {
	return "test-1"
}

func (t *TestProjectList) ListProjects(ctx context.Context, option client.ListOption) (client.ProjectList, error) {
	return client.ProjectList{
		Items: []*client.Project{
			{
				Kind: "1",
			},
			{
				Kind: "2",
			},
		},
	}, nil
}

type TestProjectCreate struct {
}

func (t *TestProjectCreate) Path() string {
	return "test-2"
}

func (t *TestProjectCreate) CreateProject(ctx context.Context, project *client.Project) (*client.Project, error) {
	return &client.Project{}, nil
}

type TestResourceList struct {
}

func (t *TestResourceList) Path() string {
	return "test-3"
}

func (t *TestResourceList) ListResources(ctx context.Context, option client.ListOption) (client.ResourceList, error) {
	return client.ResourceList{}, nil
}

type TestProjectListCreate struct {
}

func (t *TestProjectListCreate) Path() string {
	return "test-4"
}

func (t *TestProjectListCreate) ListProjects(ctx context.Context, option client.ListOption) (client.ProjectList, error) {
	return client.ProjectList{}, nil
}

func (t *TestProjectListCreate) CreateProject(ctx context.Context, project *client.Project) (*client.Project, error) {
	return &client.Project{}, nil
}
