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

package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emicklei/go-restful/v3"
	. "github.com/onsi/gomega"
)

func TestMeta_WithContext(t *testing.T) {
	g := NewGomegaWithT(t)

	meta := &Meta{
		BaseURL: "http://katanomi.dev",
		Version: "123",
	}

	metaCtx := meta.WithContext(context.TODO())

	g.Expect(metaCtx).ToNot(BeNil())
}

func TestExtractMeta(t *testing.T) {
	g := NewGomegaWithT(t)

	meta := &Meta{
		BaseURL: "http://katanomi.dev",
		Version: "123",
	}

	metaCtx := meta.WithContext(context.TODO())
	newMeta := ExtraMeta(metaCtx)

	g.Expect(newMeta).ToNot(BeNil())
	g.Expect(newMeta).To(Equal(meta))
}

func TestMetaFilter(t *testing.T) {
	httpRequest, _ := http.NewRequest("GET", "http://example.com/test", nil)
	httpRequest.Header.Set("Accept", "*/*")
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set(PluginMetaHeader, "eyJiYXNlVVJMIjoiaHR0cDovL2FiYy5jb20ifQ==")
	httpWriter := httptest.NewRecorder()

	ws := new(restful.WebService).Filter(MetaFilter)
	ws.Route(ws.Produces(restful.MIME_JSON).GET("test").To(func(request *restful.Request, response *restful.Response) {
		meta := ExtraMeta(request.Request.Context())

		g := NewGomegaWithT(t)
		g.Expect(meta).NotTo(BeNil())
		g.Expect(meta.BaseURL).To(Equal("http://abc.com"))
	}))

	c := restful.NewContainer().Add(ws)
	c.Dispatch(httpWriter, httpRequest)
}
