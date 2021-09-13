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
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/emicklei/go-restful/v3"
	"github.com/katanomi/pkg/plugin/client"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAboutGitCommitCreate(t *testing.T) {
	g := NewGomegaWithT(t)

	ws, err := NewService(&TestGitRepoFileCreator{}, client.MetaFilter)
	g.Expect(err).To(BeNil())

	container := restful.NewContainer()

	container.Add(ws)
	metaData := client.Meta{BaseURL: "http://api.test", Version: "v1"}
	data, _ := json.Marshal(metaData)
	meta := base64.StdEncoding.EncodeToString(data)

	httpWriter := httptest.NewRecorder()

	FilePath := url.PathEscape(strings.Replace(".build%2Fbuild.yaml", ".", "%2E", -1))
	fmt.Println(FilePath)
	params := metav1alpha1.CreateRepoFileParams{
		Message: "a",
		Branch:  "b",
	}
	c, _ := json.Marshal(params)
	httpRequest1, _ := http.NewRequest("POST", "/plugins/v1alpha1/test-by/projects/1/coderepositories/1/content/"+FilePath, bytes.NewBuffer(c))
	httpRequest1.Header.Set("content-type", "application/json")
	httpRequest1.Header.Set(client.PluginMetaHeader, meta)

	container.Dispatch(httpWriter, httpRequest1)
	g.Expect(httpWriter.Code).To(Equal(http.StatusOK))
	prList := metav1alpha1.GitCommit{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &prList)
	g.Expect(err).To(BeNil())
	g.Expect(prList.Kind).To(Equal(metav1alpha1.GitCommitGVK.Kind))

}

type TestGitRepoFileCreator struct {
}

func (t *TestGitRepoFileCreator) Path() string {
	return "test-by"
}

func (t *TestGitRepoFileCreator) Setup(_ context.Context, _ *zap.SugaredLogger) error {
	return nil
}

func (t *TestGitRepoFileCreator) CreateGitRepoFile(ctx context.Context, payload metav1alpha1.CreateRepoFilePayload) (metav1alpha1.GitCommit, error) {
	return metav1alpha1.GitCommit{
		TypeMeta: metav1.TypeMeta{
			Kind: metav1alpha1.GitCommitGVK.Kind,
		},
	}, nil
}
