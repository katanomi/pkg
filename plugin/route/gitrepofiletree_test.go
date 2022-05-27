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

package route_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emicklei/go-restful/v3"
	kerrors "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/plugin/client"
	"github.com/katanomi/pkg/plugin/route"
	"k8s.io/apimachinery/pkg/api/errors"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	authv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestAboutGitRepositoryFileTreeTable(t *testing.T) {
	g := NewGomegaWithT(t)

	var meta string
	var c []byte
	var container *restful.Container
	var err error
	var ws *restful.WebService
	var httpRequest *http.Request
	var httpWriter *httptest.ResponseRecorder
	var metaData client.Meta
	var data []byte
	var FilePath string
	var params metav1alpha1.GitRepoFileTreeOption

	metaData = client.Meta{BaseURL: "http://api.test", Version: "v1"}
	data, _ = json.Marshal(metaData)
	FilePath = ".build%/build.yaml"
	params = metav1alpha1.GitRepoFileTreeOption{
		Path:      FilePath,
		TreeSha:   "master",
		Recursive: false,
	}
	meta = base64.StdEncoding.EncodeToString(data)
	c, _ = json.Marshal(params)
	httpRequest, _ = http.NewRequest("GET", "/plugins/v1alpha1/test-by/projects/1/coderepositories/1/tree", bytes.NewBuffer(c))
	httpRequest.Header.Set("content-type", "application/json")
	httpRequest.Header.Set(client.PluginMetaHeader, meta)
	resourceAtt := authv1.ResourceAttributes{}

	getters := []TestGitRepoFileTreeMockGetter{
		{
			MockResult: metav1alpha1.GitRepositoryFileTree{
				TypeMeta: metav1.TypeMeta{
					Kind: metav1alpha1.GitRepositoryFileTreeGVK.Kind,
				}},
		}, {
			MockError: errors.NewNotFound(schema.GroupResource{}, "error"),
		}, {
			MockError: errors.NewInternalError(fmt.Errorf("error")),
		}, {
			MockError: errors.NewTimeoutError("timeout", 0),
		}, {
			MockError: errors.NewForbidden(schema.GroupResource{
				Group:    resourceAtt.Group,
				Resource: resourceAtt.Resource},
				resourceAtt.Name, fmt.Errorf("access not allowed")),
		}, {
			MockError: errors.NewUnauthorized(""),
		},
	}
	for _, getter := range getters {
		ws, err = route.NewService(&getter, client.MetaFilter)
		g.Expect(err).To(BeNil())
		container = restful.NewContainer()
		container.Add(ws)

		httpWriter = httptest.NewRecorder()
		g.Expect(container).NotTo(BeNil())
		container.Dispatch(httpWriter, httpRequest)
		if getter.MockError != nil {
			g.Expect(container).NotTo(BeNil())
			g.Expect(httpWriter.Code).To(Equal(kerrors.AsStatusCode(getter.MockError)))
		} else {
			g.Expect(httpWriter.Code).To(Equal(http.StatusOK))
			fileTree := metav1alpha1.GitRepositoryFileTree{}
			err := json.Unmarshal(httpWriter.Body.Bytes(), &fileTree)
			g.Expect(err).To(BeNil())
			g.Expect(fileTree.Kind).To(Equal(metav1alpha1.GitRepositoryFileTreeGVK.Kind))
		}
	}

}

// mock ClientGitRepositoryFileTree interface
type TestGitRepoFileTreeMockGetter struct {
	MockResult metav1alpha1.GitRepositoryFileTree
	MockError  error
}

func (t *TestGitRepoFileTreeMockGetter) Path() string {
	return "test-by"
}

func (t *TestGitRepoFileTreeMockGetter) Setup(_ context.Context, _ *zap.SugaredLogger) error {
	return nil
}

func (t *TestGitRepoFileTreeMockGetter) GetGitRepositoryFileTree(ctx context.Context, repoOption metav1alpha1.GitRepoFileTreeOption, listOption metav1alpha1.ListOptions) (metav1alpha1.GitRepositoryFileTree, error) {
	return t.MockResult, t.MockError
}
