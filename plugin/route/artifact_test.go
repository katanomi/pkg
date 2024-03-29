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
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

func TestListArtifacts(t *testing.T) {
	g := NewGomegaWithT(t)

	ws, err := NewService(&MockArtifactImpl{}, client.MetaFilter)
	g.Expect(err).To(BeNil())

	var projects = []string{
		"proj",
		"proj%2Fsub",
		"repositories",
	}

	container := restful.NewContainer()
	container.Router(restful.RouterJSR311{})
	container.Add(ws)
	for _, project := range projects {
		path := fmt.Sprintf("/plugins/v1alpha1/test-artifacts-1/projects/%s/repositories/katanomi/artifacts", project)
		httpRequest, _ := http.NewRequest("GET", path, nil)
		httpRequest.Header.Set("Accept", "application/json")

		metaData := client.Meta{BaseURL: "http://api.test", Version: "v1"}
		data, _ := json.Marshal(metaData)
		meta := base64.StdEncoding.EncodeToString(data)
		httpRequest.Header.Set(client.PluginMetaHeader, meta)
		httpWriter := httptest.NewRecorder()
		container.Dispatch(httpWriter, httpRequest)
		g.Expect(httpWriter.Code).To(Equal(http.StatusOK))

		branchList := metav1alpha1.GitBranchList{}
		err = json.Unmarshal(httpWriter.Body.Bytes(), &branchList)
		g.Expect(err).To(BeNil())
	}
}

func TestListProjectArtifacts(t *testing.T) {
	g := NewGomegaWithT(t)

	ws, err := NewService(&MockArtifactImpl{}, client.MetaFilter)
	g.Expect(err).To(BeNil())

	var projects = []string{
		"proj",
		"proj%2Fsub",
		"repositories",
	}

	container := restful.NewContainer()
	container.Router(restful.RouterJSR311{})
	container.Add(ws)
	for _, project := range projects {
		path := fmt.Sprintf("/plugins/v1alpha1/test-artifacts-1/projects/%s/artifacts", project)
		httpRequest, _ := http.NewRequest("GET", path, nil)
		httpRequest.Header.Set("Accept", "application/json")

		metaData := client.Meta{BaseURL: "http://api.test", Version: "v1"}
		data, _ := json.Marshal(metaData)
		meta := base64.StdEncoding.EncodeToString(data)
		httpRequest.Header.Set(client.PluginMetaHeader, meta)
		httpWriter := httptest.NewRecorder()
		container.Dispatch(httpWriter, httpRequest)
		g.Expect(httpWriter.Code).To(Equal(http.StatusOK))

		branchList := metav1alpha1.GitBranchList{}
		err = json.Unmarshal(httpWriter.Body.Bytes(), &branchList)
		g.Expect(err).To(BeNil())
	}
}

func TestGetProjectArtifact(t *testing.T) {
	g := NewGomegaWithT(t)

	ws, err := NewService(&MockArtifactImpl{}, client.MetaFilter)
	g.Expect(err).To(BeNil())

	var projects = []string{
		"proj",
		"proj%2Fsub",
		"repositories",
	}

	container := restful.NewContainer()
	container.Router(restful.RouterJSR311{})
	container.Add(ws)
	for _, project := range projects {
		path := fmt.Sprintf("/plugins/v1alpha1/test-artifacts-1/projects/%s/artifacts/test.jpg", project)
		httpRequest, _ := http.NewRequest("GET", path, nil)
		httpRequest.Header.Set("Accept", "application/json")

		metaData := client.Meta{BaseURL: "http://api.test", Version: "v1"}
		data, _ := json.Marshal(metaData)
		meta := base64.StdEncoding.EncodeToString(data)
		httpRequest.Header.Set(client.PluginMetaHeader, meta)
		httpWriter := httptest.NewRecorder()
		container.Dispatch(httpWriter, httpRequest)
		g.Expect(httpWriter.Code).To(Equal(http.StatusOK))

		var artifact metav1alpha1.Artifact
		err = json.Unmarshal(httpWriter.Body.Bytes(), &artifact)
		g.Expect(err).To(BeNil())
	}
}

func TestGetProjectArtifactFile(t *testing.T) {
	g := NewGomegaWithT(t)

	ws, err := NewService(&MockArtifactImpl{}, client.MetaFilter)
	g.Expect(err).To(BeNil())

	var projects = []string{
		"proj",
		"proj%2Fsub",
		"repositories",
	}

	container := restful.NewContainer()
	container.Router(restful.RouterJSR311{})
	container.Add(ws)
	for _, project := range projects {
		path := fmt.Sprintf("/plugins/v1alpha1/test-artifacts-1/projects/%s/artifacts/test.jpg/file", project)
		httpRequest, _ := http.NewRequest("GET", path, nil)

		metaData := client.Meta{BaseURL: "http://api.test", Version: "v1"}
		data, _ := json.Marshal(metaData)
		meta := base64.StdEncoding.EncodeToString(data)
		httpRequest.Header.Set(client.PluginMetaHeader, meta)
		httpWriter := httptest.NewRecorder()
		container.Dispatch(httpWriter, httpRequest)
		g.Expect(httpWriter.Code).To(Equal(http.StatusOK))

		g.Expect(httpWriter.Body.Bytes()).To(Equal([]byte("xxx")))
		g.Expect(err).To(BeNil())
	}
}

func TestDeleteArtifactTag(t *testing.T) {
	g := NewGomegaWithT(t)

	ws, err := NewService(&MockArtifactImpl{}, client.MetaFilter)
	g.Expect(err).To(BeNil())

	var projects = []string{
		"proj",
		"proj%2Fsub",
		"repositories",
	}

	container := restful.NewContainer()
	container.Router(restful.RouterJSR311{})
	container.Add(ws)
	for _, project := range projects {
		path := fmt.Sprintf("/plugins/v1alpha1/test-artifacts-1/projects/%s/repositories/katanomi/artifacts/artifact/tags/tag", project)
		httpRequest, _ := http.NewRequest("DELETE", path, nil)
		httpRequest.Header.Set("Accept", "application/json")

		metaData := client.Meta{BaseURL: "http://api.test", Version: "v1"}
		data, _ := json.Marshal(metaData)
		meta := base64.StdEncoding.EncodeToString(data)
		httpRequest.Header.Set(client.PluginMetaHeader, meta)
		httpWriter := httptest.NewRecorder()
		container.Dispatch(httpWriter, httpRequest)
		g.Expect(httpWriter.Code).To(Equal(http.StatusOK))
		g.Expect(err).To(BeNil())
	}
}

func TestDeleteProjectArtifact(t *testing.T) {
	g := NewGomegaWithT(t)

	ws, err := NewService(&MockArtifactImpl{}, client.MetaFilter)
	g.Expect(err).To(BeNil())

	var projects = []string{
		"proj",
		"proj%2Fsub",
		"repositories",
	}

	container := restful.NewContainer()
	container.Router(restful.RouterJSR311{})
	container.Add(ws)
	for _, project := range projects {
		path := fmt.Sprintf("/plugins/v1alpha1/test-artifacts-1/projects/%s/artifacts/foo/bar.jpg", project)
		httpRequest, _ := http.NewRequest("DELETE", path, nil)
		httpRequest.Header.Set("Accept", "application/json")

		metaData := client.Meta{BaseURL: "http://api.test", Version: "v1"}
		data, _ := json.Marshal(metaData)
		meta := base64.StdEncoding.EncodeToString(data)
		httpRequest.Header.Set(client.PluginMetaHeader, meta)
		httpWriter := httptest.NewRecorder()
		container.Dispatch(httpWriter, httpRequest)
		g.Expect(httpWriter.Code).To(Equal(http.StatusOK))
		g.Expect(err).To(BeNil())
	}
}

type MockArtifactImpl struct {
}

func (t *MockArtifactImpl) DeleteArtifact(ctx context.Context, params metav1alpha1.ArtifactOptions) error {
	return nil
}

func (t *MockArtifactImpl) GetProjectArtifact(ctx context.Context, params metav1alpha1.ProjectArtifactOptions) (*metav1alpha1.Artifact, error) {
	return &metav1alpha1.Artifact{}, nil
}

func (t *MockArtifactImpl) GetProjectArtifactFile(ctx context.Context, params metav1alpha1.ProjectArtifactOptions) (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader([]byte("xxx"))), nil
}

func (t *MockArtifactImpl) Path() string {
	return "test-artifacts-1"
}

func (t *MockArtifactImpl) Setup(_ context.Context, _ *zap.SugaredLogger) error {
	return nil
}

func (t *MockArtifactImpl) ListArtifacts(ctx context.Context, params metav1alpha1.ArtifactOptions, option metav1alpha1.ListOptions) (*metav1alpha1.ArtifactList, error) {
	return &metav1alpha1.ArtifactList{}, nil
}

func (t *MockArtifactImpl) GetArtifact(ctx context.Context, params metav1alpha1.ArtifactOptions) (*metav1alpha1.Artifact, error) {
	return &metav1alpha1.Artifact{}, nil
}

func (t *MockArtifactImpl) ListProjectArtifacts(ctx context.Context, params metav1alpha1.ArtifactOptions, option metav1alpha1.ListOptions) (*metav1alpha1.ArtifactList, error) {
	return &metav1alpha1.ArtifactList{}, nil
}

func (t *MockArtifactImpl) DeleteArtifactTag(ctx context.Context, params metav1alpha1.ArtifactTagOptions) error {
	return nil
}

func (t *MockArtifactImpl) DeleteProjectArtifact(ctx context.Context, params metav1alpha1.ProjectArtifactOptions) error {
	return nil
}
