/*
Copyright 2023 The AlaudaDevops Authors.

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

package config

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	restful "github.com/emicklei/go-restful/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"knative.dev/pkg/logging"
)

var _ = Describe("ConfigFilter", func() {

	var (
		manager              *Manager
		ctx                  context.Context
		request              *restful.Request
		response             *restful.Response
		chain                *restful.FilterChain
		recorder             *httptest.ResponseRecorder
		key                  string
		expectedKeyValueFunc ConfigKeyExpectedValueFunc
	)

	BeforeEach(func() {
		ctx = context.Background()
		ctx = logging.WithLogger(ctx, log)
		req := &http.Request{
			Header: map[string][]string{
				restful.HEADER_AcceptEncoding: []string{restful.MIME_JSON},
			},
		}
		testUrl, _ := url.Parse("http://test.example/some/path")
		req.URL = testUrl
		request = &restful.Request{Request: req}
		recorder = httptest.NewRecorder()
		response = &restful.Response{ResponseWriter: recorder}
		response.SetRequestAccepts(restful.MIME_JSON)
		chain = &restful.FilterChain{
			Filters: []restful.FilterFunction{},
			Target: func(req *restful.Request, resp *restful.Response) {
				resp.WriteHeader(http.StatusOK)
			},
		}

		manager = &Manager{Config: &Config{Data: map[string]string{"test": "test"}}}
		key = "test"
	})

	JustBeforeEach(func() {
		request.Request = request.Request.WithContext(ctx)

		ConfigFilter(ctx, manager, key, expectedKeyValueFunc)(request, response, chain)
	})

	Context("Uses a \"test\" key with \"test\" value using some basic ConfigKeyExpectedValueFunc implementation", func() {
		BeforeEach(func() {
			expectedKeyValueFunc = func(ctx context.Context, req *restful.Request, key string, value FeatureValue) (err error) {
				ok, err := value.AsBool()
				if err != nil {
					return err
				} else if !ok {
					return fmt.Errorf("value is not true: %v", value)
				}
				return nil
			}
		})
		It("should have a internal error as status code with api error in response body", func() {
			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			Expect(strings.TrimSpace(recorder.Body.String())).To(Equal(`{"metadata":{},"status":"Failure","message":"Internal error occurred: failed parsing feature flags config \"test\": strconv.ParseBool: parsing \"test\": invalid syntax","reason":"InternalError","details":{"causes":[{"message":"failed parsing feature flags config \"test\": strconv.ParseBool: parsing \"test\": invalid syntax"}]},"code":500}`))
		})

		Context("Uses a \"test\" key with \"true\" value using some basic ConfigKeyExpectedValueFunc implementation", func() {
			BeforeEach(func() {
				manager.Config.Data["test"] = "true"
			})
			It("should pass filter", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})
		})
	})

	Context("Uses ConfigFilterNotFoundWhenNotTrue with false value", func() {
		BeforeEach(func() {
			expectedKeyValueFunc = ConfigFilterNotFoundWhenNotTrue
			manager.Config.Data["test"] = "false"
		})

		It("should have a not found error as status code with api error in response body", func() {
			Expect(recorder.Code).To(Equal(http.StatusNotFound))
			Expect(strings.TrimSpace(recorder.Body.String())).To(Equal(`{"metadata":{},"status":"Failure","message":"the server could not find the requested resource ( API.katanomi.dev http://test.example/some/path)","reason":"NotFound","details":{"name":"http://test.example/some/path","group":"katanomi.dev","kind":"API"},"code":404}`))
		})

	})

	Context("Uses ConfigFilterNotFoundWhenNotTrue with true value", func() {
		BeforeEach(func() {
			expectedKeyValueFunc = ConfigFilterNotFoundWhenNotTrue
			manager.Config.Data["test"] = "true"
		})

		It("should pass filter", func() {
			Expect(recorder.Code).To(Equal(http.StatusOK))
		})

	})
})
