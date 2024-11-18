/*
Copyright 2022 The AlaudaDevops Authors.

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

package validation

import (
	"testing"

	duckv1 "knative.dev/pkg/apis/duck/v1"

	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"knative.dev/pkg/apis"
)

func TestValidateAddressable(t *testing.T) {
	testCases := []struct {
		addr     duckv1.Addressable
		optional bool
		path     *field.Path
		wanErr   bool
	}{
		{
			addr: duckv1.Addressable{
				URL: &apis.URL{
					Scheme: "",
					Host:   "",
				},
			},
			optional: false,
			path:     field.NewPath("url"),
			wanErr:   true,
		},
		{
			addr: duckv1.Addressable{
				URL: &apis.URL{
					Scheme: "http",
					Host:   "127.0.0.1",
				},
			},
			optional: false,
			path:     field.NewPath("url"),
			wanErr:   false,
		},
	}
	g := NewGomegaWithT(t)
	for _, tt := range testCases {
		err := ValidateAddressable(tt.addr, tt.optional, tt.path)
		if tt.wanErr {
			g.Expect(err).ShouldNot(BeNil())
			g.Expect(err[0].Detail).To(Equal("value is required"))
		} else {
			g.Expect(len(err)).Should(Equal(0))
		}
	}
}

func TestValidateURL(t *testing.T) {
	testCases := []struct {
		url    *apis.URL
		path   *field.Path
		wanErr bool
	}{
		{
			url: &apis.URL{
				Scheme: "",
				Host:   "",
			},
			path:   field.NewPath("url"),
			wanErr: true,
		},
		{
			url: &apis.URL{
				Scheme: "http",
				Host:   "127.0.0.1",
				Path:   "/api",
			},
			path:   field.NewPath("url"),
			wanErr: false,
		},
	}

	g := NewGomegaWithT(t)

	for _, tt := range testCases {
		err := ValidateURL(tt.url, tt.path)
		if tt.wanErr {
			g.Expect(err).ShouldNot(BeNil())
			g.Expect(err[0].Detail).To(Equal("value is required"))
		} else {
			g.Expect(len(err)).Should(Equal(0))
		}
	}
}
