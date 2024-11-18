/*
Copyright 2021 The AlaudaDevops Authors.

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
	goerrors "errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAsAPIError(t *testing.T) {
	g := NewGomegaWithT(t)

	err := AsAPIError(fmt.Errorf("anyerror"))
	status := errors.APIStatus(nil)
	g.Expect(goerrors.As(err, &status)).To(BeTrue())
	g.Expect(errors.ReasonForError(err)).To(Equal(metav1.StatusReasonInternalError))

	err = AsAPIError(errors.NewBadRequest("bad request"))
	status = errors.APIStatus(nil)
	g.Expect(goerrors.As(err, &status)).To(BeTrue())
	g.Expect(errors.ReasonForError(err)).To(Equal(metav1.StatusReasonBadRequest))
}

func TestAsStatusCode(t *testing.T) {
	g := NewGomegaWithT(t)

	g.Expect(AsStatusCode(fmt.Errorf("anyerror"))).To(Equal(http.StatusInternalServerError))
	g.Expect(AsStatusCode(errors.NewBadRequest("bad"))).To(Equal(http.StatusBadRequest))
}

var _ = Describe("Test AsStatusError", func() {
	var (
		server   *httptest.Server
		response *resty.Response
	)

	newTestServer := func(f func(w http.ResponseWriter, r *http.Request)) {
		mux := http.NewServeMux()
		server = httptest.NewServer(mux)
		mux.HandleFunc("/test", f)
	}

	BeforeEach(func() {
		server = &httptest.Server{}
		response = &resty.Response{}
	})

	JustBeforeEach(func() {
		response, _ = resty.New().R().Get(server.URL + "/test")
	})

	When("server return 200", func() {
		BeforeEach(func() {
			newTestServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			DeferCleanup(func() {
				server.Close()
			})
		})

		It("should return nil", func() {
			Expect(response.StatusCode()).Should(Equal(http.StatusOK))
		})
	})

	When("server return 301", func() {
		BeforeEach(func() {
			newTestServer(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Location", "https://127.0.0.1")
				w.WriteHeader(http.StatusMovedPermanently)
			})
			DeferCleanup(func() {
				server.Close()
			})
		})

		It("rawResponse should be nil", func() {
			Expect(response.RawResponse).To(BeNil())
		})
	})

	When("server return 400", func() {
		BeforeEach(func() {
			newTestServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("test body"))
			})
			DeferCleanup(func() {
				server.Close()
			})
		})

		It("should return 400", func() {
			Expect(response.StatusCode()).To(Equal(http.StatusBadRequest))
			statusError := AsStatusError(response).(*errors.StatusError)
			Expect(int(statusError.ErrStatus.Code)).To(Equal(http.StatusBadRequest))
			Expect(statusError.ErrStatus.Message).To(ContainSubstring("the server rejected our request for an unknown reason"))
		})
	})
	When("server return 403", func() {
		BeforeEach(func() {
			newTestServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("test forbidden error"))
			})
			DeferCleanup(func() {
				server.Close()
			})
		})

		It("should return 403", func() {
			Expect(response.StatusCode()).To(Equal(http.StatusForbidden))
			statusError := AsStatusError(response).(*errors.StatusError)
			Expect(int(statusError.ErrStatus.Code)).To(Equal(http.StatusForbidden))
			Expect(statusError.ErrStatus.Message).To(ContainSubstring("test forbidden error"))
		})
	})
	When("server return 500", func() {
		BeforeEach(func() {
			newTestServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})
			DeferCleanup(func() {
				server.Close()
			})
		})

		It("should return 500", func() {
			Expect(response.StatusCode()).To(Equal(http.StatusInternalServerError))
			statusError := AsStatusError(response).(*errors.StatusError)
			Expect(int(statusError.ErrStatus.Code)).To(Equal(http.StatusInternalServerError))
			Expect(statusError.ErrStatus.Message).To(ContainSubstring("an error on the server (\"\") has prevented the request from succeeding"))
		})
	})
})
