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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestListGitBranch(t *testing.T) {
	g := NewGomegaWithT(t)

	ws, err := NewService(&TestGitBranchLister{}, client.MetaFilter)
	g.Expect(err).To(BeNil())

	container := restful.NewContainer()

	container.Add(ws)

	httpRequest, _ := http.NewRequest("GET", "/plugins/v1alpha1/test-6/projects/1/coderepositories/1/branches", nil)
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
	g.Expect(branchList.Kind).To(Equal(metav1alpha1.GitBranchListGVK.Kind))
}

type TestGitBranchLister struct {
}

func (t *TestGitBranchLister) Path() string {
	return "test-6"
}

func (t *TestGitBranchLister) Setup(_ context.Context, _ *zap.SugaredLogger) error {
	return nil
}

func (t *TestGitBranchLister) ListGitBranch(ctx context.Context, repoOption metav1alpha1.GitBranchOption, option metav1alpha1.ListOptions) (metav1alpha1.GitBranchList, error) {
	return metav1alpha1.GitBranchList{
		TypeMeta: metav1.TypeMeta{
			Kind: metav1alpha1.GitBranchListGVK.Kind,
		},
	}, nil
}

func TestCreateGitBranch(t *testing.T) {
	g := NewGomegaWithT(t)

	ws, err := NewService(&TestGitBranchCreator{}, client.MetaFilter)
	g.Expect(err).To(BeNil())

	container := restful.NewContainer()

	container.Add(ws)
	payload := metav1alpha1.CreateBranchParams{
		Branch: "dev",
		Ref:    "master",
	}
	content, _ := json.Marshal(payload)
	httpRequest, _ := http.NewRequest("POST", "/plugins/v1alpha1/test-7/projects/1/coderepositories/1/branches", bytes.NewBuffer(content))
	httpRequest.Header.Set("Content-Type", "application/json")

	metaData := client.Meta{BaseURL: "http://api.test", Version: "v1"}
	data, _ := json.Marshal(metaData)
	meta := base64.StdEncoding.EncodeToString(data)
	httpRequest.Header.Set(client.PluginMetaHeader, meta)

	httpWriter := httptest.NewRecorder()

	container.Dispatch(httpWriter, httpRequest)
	g.Expect(httpWriter.Code).To(Equal(http.StatusOK))

	branch := metav1alpha1.GitBranch{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &branch)
	g.Expect(err).To(BeNil())
	g.Expect(branch.Name).To(Equal("dev"))
}

type TestGitBranchCreator struct {
}

func (t *TestGitBranchCreator) Path() string {
	return "test-7"
}

func (t *TestGitBranchCreator) Setup(_ context.Context, _ *zap.SugaredLogger) error {
	return nil
}

func (t *TestGitBranchCreator) CreateGitBranch(ctx context.Context, payload metav1alpha1.CreateBranchPayload) (metav1alpha1.GitBranch, error) {
	return metav1alpha1.GitBranch{
		ObjectMeta: metav1.ObjectMeta{Name: payload.Branch},
	}, nil
}
