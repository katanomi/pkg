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

package client

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Test.GetOptions", func() {
	var (
		opts *GetOptions
	)

	BeforeEach(func() {
		opts = NewGetOptions()
	})

	Describe("init", func() {
		Context("when GetOptions is not nil and Raw is not nil", func() {
			It("should return the input GetOptions", func() {
				raw := &metav1.GetOptions{}
				opt := &GetOptions{
					GetOptions: client.GetOptions{
						Raw: raw,
					},
				}
				Expect(opt.init()).To(Equal(opt))
			})
		})

		Context("when GetOptions is nil or Raw is nil", func() {
			It("should return a new GetOptions with Raw set to a new instance of metav1.GetOptions", func() {
				Expect(opts.init()).To(Equal(&GetOptions{
					GetOptions: client.GetOptions{
						Raw: &metav1.GetOptions{},
					},
				}))
			})
		})
	})

	Describe("WithCached", func() {
		It("should set the ResourceVersion to '0'", func() {
			Expect(opts.WithCached().Build().Raw.ResourceVersion).To(Equal("0"))
		})
	})

	Describe("Build", func() {
		It("should return a pointer to the underlying client.GetOptions", func() {
			Expect(opts.Build()).To(BeIdenticalTo(&opts.GetOptions))
		})
	})

})

var _ = Describe("GetCachedOptions", func() {
	It("should return client.GetOptions with ResourceVersion 0", func() {
		options := GetCachedOptions()
		Expect(options).ToNot(BeNil())
		Expect(options.Raw).To(Equal(&metav1.GetOptions{ResourceVersion: "0"}))
	})
})

var _ = Describe("Test.ListOptions", func() {
	var (
		opts *ListOptions
	)

	BeforeEach(func() {
		opts = NewListOptions()
	})

	Describe("init", func() {
		Context("when ListOptions is not nil and Raw is not nil", func() {
			It("should return the input ListOptions", func() {
				raw := &metav1.ListOptions{}
				opt := &ListOptions{
					ListOptions: client.ListOptions{
						Raw: raw,
					},
				}
				Expect(opt.init()).To(Equal(opt))
			})
		})

		Context("when ListOptions is nil or Raw is nil", func() {
			It("should return a new ListOptions with Raw set to a new instance of metav1.ListOptions", func() {
				Expect(opts.init()).To(Equal(&ListOptions{
					ListOptions: client.ListOptions{
						Raw: &metav1.ListOptions{},
					},
				}))
			})
		})
	})

	Describe("WithCached", func() {
		It("should set the ResourceVersion to '0'", func() {
			Expect(opts.WithCached().Build().Raw.ResourceVersion).To(Equal("0"))
		})
	})

	Describe("WithLimit", func() {
		It("should set the Limit", func() {
			Expect(opts.WithLimit(10).Build().Limit).To(Equal(int64(10)))
		})
	})

	Describe("WithNamespace", func() {
		It("should set the Namespace", func() {
			Expect(opts.WithNamespace("NS").Build().Namespace).To(Equal("NS"))
		})
	})

	Describe("WithUnsafeDisableDeepCopy", func() {
		It("should set the UnsafeDisableDeepCopy to true", func() {
			Expect(opts.WithUnsafeDisableDeepCopy().Build().UnsafeDisableDeepCopy).NotTo(BeNil())
			Expect(*opts.WithUnsafeDisableDeepCopy().Build().UnsafeDisableDeepCopy).To(BeTrue())
		})
	})

	Describe("Build", func() {
		It("should return a pointer to the underlying client.ListOptions", func() {
			Expect(opts.Build()).To(BeIdenticalTo(&opts.ListOptions))
		})
	})

})

var _ = Describe("ListCachedOptions", func() {
	It("should return client.ListOptions with ResourceVersion 0", func() {
		options := ListCachedOptions()
		Expect(options).ToNot(BeNil())
		Expect(options.Raw).To(Equal(&metav1.ListOptions{ResourceVersion: "0"}))
	})
})
