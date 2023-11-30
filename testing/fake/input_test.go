/*
Copyright 2023 The Katanomi Authors.

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

package fake

import (
	"net/http"
	"strings"
	"testing"

	"github.com/emicklei/go-restful/v3"
	"github.com/katanomi/pkg/plugin/client"
	. "github.com/onsi/gomega"
)

func TestInputFromRequest(t *testing.T) {
	bodyReader := strings.NewReader(`{"qwer" : "1234"}`)
	httpRequest, _ := http.NewRequest("POST", "/test?abc=def", bodyReader)
	httpRequest.Header.Set("Content-Type", "application/json")

	httpRequest.Header.Set(client.PluginSecretHeader, "eyJ1c2VybmFtZSI6Ik1YRmhlZz09IiwicGFzc3dvcmQiOiJNbmR6ZUE9PSJ9")
	httpRequest.Header.Set(client.PluginAuthHeader, "kubernetes.io/basic-auth")
	httpRequest.Header.Set(client.PluginMetaHeader, "eyJiYXNlVVJMIjoiaHR0cDovL2FiYy5jb20ifQ==")
	request := &restful.Request{Request: httpRequest}

	g := NewGomegaWithT(t)

	input, err := InputFromRequest(request)
	g.Expect(err).To(BeNil())
	g.Expect(input["body"]).To(HaveKeyWithValue("qwer", "1234"))
	g.Expect(input["query"]).To(HaveKeyWithValue("abc", "def"))

	auth := input["auth"].(map[string]interface{})
	g.Expect(auth["data"]).To(HaveKeyWithValue("username", []byte("1qaz")))
	g.Expect(auth["data"]).To(HaveKeyWithValue("password", []byte("2wsx")))

	meta := input["meta"].(map[string]interface{})
	g.Expect(meta).To(HaveKeyWithValue("baseURL", "http://abc.com"))
}
