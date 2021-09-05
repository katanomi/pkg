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
	"strconv"
	"testing"

	"github.com/emicklei/go-restful/v3"
	"github.com/katanomi/pkg/plugin/client"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAboutGitPullRequest(t *testing.T) {
	g := NewGomegaWithT(t)

	ws, err := NewService(&TestGitPullRequestHandler{}, client.MetaFilter)
	g.Expect(err).To(BeNil())

	container := restful.NewContainer()

	container.Add(ws)

	metaData := client.Meta{BaseURL: "http://api.test", Version: "v1"}
	data, _ := json.Marshal(metaData)
	meta := base64.StdEncoding.EncodeToString(data)

	httpWriter := httptest.NewRecorder()

	httpRequest1, _ := http.NewRequest("GET", "/plugins/v1alpha1/test-9/projects/1/coderepositories/1/pulls", nil)
	httpRequest1.Header.Set("Accept", "application/json")
	httpRequest1.Header.Set(client.PluginMetaHeader, meta)

	container.Dispatch(httpWriter, httpRequest1)
	g.Expect(httpWriter.Code).To(Equal(http.StatusOK))
	prList := metav1alpha1.GitPullRequestList{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &prList)
	g.Expect(err).To(BeNil())
	g.Expect(prList.Kind).To(Equal(metav1alpha1.GitPullrequestsListGVK.Kind))

	httpWriter2 := httptest.NewRecorder()
	httpRequest2, _ := http.NewRequest("GET", "/plugins/v1alpha1/test-9/projects/1/coderepositories/1/pulls/1", nil)
	httpRequest2.Header.Set("Accept", "application/json")
	httpRequest2.Header.Set(client.PluginMetaHeader, meta)
	container.Dispatch(httpWriter2, httpRequest2)
	g.Expect(httpWriter.Code).To(Equal(http.StatusOK))
	prObj := metav1alpha1.GitPullRequest{}
	err = json.Unmarshal(httpWriter2.Body.Bytes(), &prObj)
	g.Expect(err).To(BeNil())
	g.Expect(prObj.Name).To(Equal("1"))

	httpWriter3 := httptest.NewRecorder()
	prPayload := metav1alpha1.CreatePullRequestPayload{
		Source:      metav1alpha1.GitBranchBaseInfo{Name: "dev"},
		Target:      metav1alpha1.GitBranchBaseInfo{Name: "master"},
		Title:       "pr",
		Description: "describe",
	}
	content, _ := json.Marshal(prPayload)
	httpRequest3, _ := http.NewRequest("POST", "/plugins/v1alpha1/test-9/projects/1/coderepositories/1/pulls", bytes.NewBuffer(content))
	httpRequest3.Header.Set("Content-Type", "application/json")
	httpRequest3.Header.Set(client.PluginMetaHeader, meta)
	container.Dispatch(httpWriter3, httpRequest3)
	g.Expect(httpWriter3.Code).To(Equal(http.StatusOK))
	prObj2 := metav1alpha1.GitPullRequest{}
	err = json.Unmarshal(httpWriter3.Body.Bytes(), &prObj2)
	g.Expect(err).To(BeNil())
	g.Expect(prObj2.Spec.Title).To(Equal("pr"))

}

type TestGitPullRequestHandler struct {
}

func (t *TestGitPullRequestHandler) Path() string {
	return "test-9"
}

func (t *TestGitPullRequestHandler) Setup(_ context.Context, _ *zap.SugaredLogger) error {
	return nil
}

func (t *TestGitPullRequestHandler) ListGitPullRequest(ctx context.Context, option metav1alpha1.GitRepo, listOption metav1alpha1.ListOptions) (metav1alpha1.GitPullRequestList, error) {
	return metav1alpha1.GitPullRequestList{
		TypeMeta: metav1.TypeMeta{
			Kind: metav1alpha1.GitPullrequestsListGVK.Kind,
		},
	}, nil
}

func (t *TestGitPullRequestHandler) GetGitPullRequest(ctx context.Context, option metav1alpha1.GitPullRequestOption) (metav1alpha1.GitPullRequest, error) {
	index := strconv.Itoa(option.Index)
	return metav1alpha1.GitPullRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name: index,
		},
	}, nil
}

func (t *TestGitPullRequestHandler) CreatePullRequest(ctx context.Context, payload metav1alpha1.CreatePullRequestPayload) (metav1alpha1.GitPullRequest, error) {
	return metav1alpha1.GitPullRequest{
		Spec: metav1alpha1.GitPullRequestSpec{
			Title: payload.Title,
			Source: metav1alpha1.GitBranchBaseInfo{
				Name: payload.Source.Name,
			},
		},
	}, nil
}
