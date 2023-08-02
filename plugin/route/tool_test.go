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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

func TestInitialize(t *testing.T) {
	g := NewGomegaWithT(t)

	ws, err := NewService(&TestToolServiceInitialize{}, client.MetaFilter)
	g.Expect(err).To(BeNil())

	container := restful.NewContainer()
	container.Add(ws)

	httpRequest, _ := http.NewRequest("GET",
		"/plugins/v1alpha1/test-initialize/tools/initialize", nil)

	httpWriter := httptest.NewRecorder()
	container.Dispatch(httpWriter, httpRequest)
	g.Expect(httpWriter.Code).To(Equal(http.StatusOK))
}

func TestCheckAlive(t *testing.T) {
	g := NewGomegaWithT(t)

	ws, err := NewService(&TestToolServiceCheckAlive{}, client.MetaFilter)
	g.Expect(err).To(BeNil())

	container := restful.NewContainer()
	container.Add(ws)

	httpRequest, _ := http.NewRequest("GET",
		"/plugins/v1alpha1/test-checkAlive/tools/liveness", nil)

	httpWriter := httptest.NewRecorder()
	container.Dispatch(httpWriter, httpRequest)
	g.Expect(httpWriter.Code).To(Equal(http.StatusOK))
}

type TestToolServiceInitialize struct {
}

func (t *TestToolServiceInitialize) Initialize(_ context.Context) error {
	return nil
}

func (t *TestToolServiceInitialize) Path() string {
	return "test-initialize"
}

func (t *TestToolServiceInitialize) Setup(_ context.Context, _ *zap.SugaredLogger) error {
	return nil
}

type TestToolServiceCheckAlive struct {
}

func (t *TestToolServiceCheckAlive) CheckAlive(_ context.Context) error {
	return nil
}
func (t *TestToolServiceCheckAlive) Path() string {
	return "test-checkAlive"
}

func (t *TestToolServiceCheckAlive) Setup(_ context.Context, _ *zap.SugaredLogger) error {
	return nil
}

func TestGetToolMetadata(t *testing.T) {
	g := NewGomegaWithT(t)

	ws, err := NewService(&TestToolMetadata{}, client.MetaFilter)
	g.Expect(err).To(BeNil())

	container := restful.NewContainer()
	container.Add(ws)

	httpRequest, _ := http.NewRequest("GET",
		"/plugins/v1alpha1/test-toolMetadata/tools/metadata", nil)

	httpWriter := httptest.NewRecorder()
	container.Dispatch(httpWriter, httpRequest)
	g.Expect(httpWriter.Code).To(Equal(http.StatusOK))
}

type TestToolMetadata struct {
}

func (t *TestToolMetadata) GetToolMetadata(_ context.Context) (*metav1alpha1.ToolMeta, error) {
	return nil, nil
}
func (t *TestToolMetadata) Path() string {
	return "test-toolMetadata"
}

func (t *TestToolMetadata) Setup(_ context.Context, _ *zap.SugaredLogger) error {
	return nil
}
