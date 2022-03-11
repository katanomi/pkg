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

package errors

import (
	"context"
	"fmt"
	"net/http"

	"k8s.io/apimachinery/pkg/api/errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ = Describe("ConvertResponseError", func() {
	var (
		ctx      context.Context
		response *http.Response
		gvk      schema.GroupVersionKind
		names    = []string{"abc"}
		err      error
		result   error
	)
	BeforeEach(func() {
		ctx = context.TODO()
		response = &http.Response{
			StatusCode: http.StatusBadGateway,
			Request:    &http.Request{Method: http.MethodGet},
		}

	})
	JustBeforeEach(func() {
		result = ConvertResponseError(ctx, response, err, gvk, names...)
	})

	Context("nil error", func() {
		BeforeEach(func() {
			err = nil
		})
		It("should be nil", func() {
			Expect(result).To(BeNil())
		})
	})
	Context("bad request error", func() {
		BeforeEach(func() {
			err = fmt.Errorf("some bad request")
			response.StatusCode = http.StatusBadRequest
		})
		It("should return kubernetes bad request error", func() {
			Expect(result).ToNot(BeNil())
			Expect(errors.IsBadRequest(result)).To(BeTrue())
		})
	})
	Context("internal server error", func() {
		BeforeEach(func() {
			err = fmt.Errorf("some internal error")
			response = nil
		})
		It("should return kubernetes internal server error", func() {
			Expect(result).ToNot(BeNil())
			Expect(errors.IsInternalError(result)).To(BeTrue())
			Expect(errors.IsBadRequest(result)).ToNot(BeTrue())
		})
	})
	Context("generic error", func() {
		BeforeEach(func() {
			err = fmt.Errorf("not allowed")
			response.StatusCode = http.StatusForbidden
		})
		It("should return kubernetes forbidden error", func() {
			Expect(result).ToNot(BeNil())
			Expect(errors.IsForbidden(result)).To(BeTrue())
		})
	})
})
