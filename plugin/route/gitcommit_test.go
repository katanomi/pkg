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

func TestGetGitCommit(t *testing.T) {
	g := NewGomegaWithT(t)

	ws, err := NewService(&TestCommitGetter{}, client.MetaFilter)
	g.Expect(err).To(BeNil())

	container := restful.NewContainer()

	container.Add(ws)

	httpRequest, _ := http.NewRequest("GET", "/plugins/v1alpha1/test-5/projects/1/coderepositories/1/commit/aaaaaaa", nil)
	httpRequest.Header.Set("Accept", "application/json")

	metaData := client.Meta{BaseURL: "http://api.test", Version: "v1"}
	data, _ := json.Marshal(metaData)
	meta := base64.StdEncoding.EncodeToString(data)
	httpRequest.Header.Set(client.PluginMetaHeader, meta)

	httpWriter := httptest.NewRecorder()

	container.Dispatch(httpWriter, httpRequest)
	g.Expect(httpWriter.Code).To(Equal(http.StatusOK))

	commit := metav1alpha1.GitCommit{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &commit)
	g.Expect(err).To(BeNil())
	g.Expect(commit.Name).To(Equal("aaaaaaa"))
}

type TestCommitGetter struct {
}

func (t *TestCommitGetter) Path() string {
	return "test-5"
}

func (t *TestCommitGetter) Setup(_ context.Context, _ *zap.SugaredLogger) error {
	return nil
}

func (t *TestCommitGetter) GetGitCommit(ctx context.Context, option metav1alpha1.GitCommitOption) (metav1alpha1.GitCommit, error) {
	return metav1alpha1.GitCommit{
		ObjectMeta: metav1.ObjectMeta{
			Name: *option.SHA,
		},
	}, nil
}
