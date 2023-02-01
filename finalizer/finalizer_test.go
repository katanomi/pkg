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

package finalizer

import (
	"context"
	"errors"

	testing2 "github.com/katanomi/pkg/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var testFinalizerKey = "test.dev/test"

var _ = Describe("SetFinalizer", func() {
	var (
		ctx context.Context
		clt client.Client
		cm  *corev1.ConfigMap
		err error
	)
	BeforeEach(func() {
		ctx = context.Background()
		clt = fake.NewClientBuilder().WithScheme(scheme).Build()
		cm = &corev1.ConfigMap{}
	})
	JustBeforeEach(func() {
		err = AddFinalizer(ctx, clt, cm, testFinalizerKey)
	})
	Context("resource has no finalizer", func() {
		BeforeEach(func() {
			testing2.MustLoadYaml("./testdata/configmap.yaml", cm)
			Expect(clt.Create(ctx, cm)).To(Succeed())
		})
		It("should add finalizer", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(cm.GetFinalizers()).To(ContainElement(testFinalizerKey))
		})
	})

	Context("resource has finalizer already", func() {
		When("contains the specified finalizer", func() {
			BeforeEach(func() {
				testing2.MustLoadYaml("./testdata/configmap.withFinalizer.yaml", cm)
				Expect(clt.Create(ctx, cm)).To(Succeed())
			})

			It("should not add finalizer", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(cm.GetFinalizers()).To(ContainElement(testFinalizerKey))
			})
		})
		When("not contains the specified finalizer", func() {
			BeforeEach(func() {
				testing2.MustLoadYaml("./testdata/configmap.withOtherFinalizer.yaml", cm)
				Expect(clt.Create(ctx, cm)).To(Succeed())
			})

			It("last finalizer is the specified finalizer", func() {
				finalizers := cm.GetFinalizers()
				Expect(err).NotTo(HaveOccurred())
				Expect(finalizers[len(finalizers)-1]).To(Equal(testFinalizerKey))
			})
		})
	})
})

var _ = Describe("PrependFinalizer", func() {
	var (
		ctx context.Context
		clt client.Client
		cm  *corev1.ConfigMap
		err error
	)
	BeforeEach(func() {
		ctx = context.Background()
		clt = fake.NewClientBuilder().WithScheme(scheme).Build()
		cm = &corev1.ConfigMap{}
	})
	JustBeforeEach(func() {
		err = PrependFinalizer(ctx, clt, cm, testFinalizerKey)
	})
	Context("resource has no finalizer", func() {
		BeforeEach(func() {
			testing2.MustLoadYaml("./testdata/configmap.yaml", cm)
			Expect(clt.Create(ctx, cm)).To(Succeed())
		})
		It("should add finalizer", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(cm.GetFinalizers()).To(ContainElement(testFinalizerKey))
		})
	})

	Context("resource has finalizer already", func() {
		When("contains the specified finalizer", func() {
			BeforeEach(func() {
				testing2.MustLoadYaml("./testdata/configmap.withFinalizer.yaml", cm)
				Expect(clt.Create(ctx, cm)).To(Succeed())
			})

			It("should not add finalizer", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(cm.GetFinalizers()).To(ContainElement(testFinalizerKey))
			})
		})
		When("not contains the specified finalizer", func() {
			BeforeEach(func() {
				testing2.MustLoadYaml("./testdata/configmap.withOtherFinalizer.yaml", cm)
				Expect(clt.Create(ctx, cm)).To(Succeed())
			})

			It("first finalizer is the specified finalizer", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(cm.GetFinalizers()[0]).To(Equal(testFinalizerKey))
			})
		})
	})
})

var _ = Describe("RemoveFinalizer", func() {
	var (
		ctx      context.Context
		clt      client.Client
		cm       *corev1.ConfigMap
		callback func() error
		err      error
	)

	BeforeEach(func() {
		ctx = context.Background()
		clt = fake.NewClientBuilder().WithScheme(scheme).Build()
		cm = &corev1.ConfigMap{}
		callback = nil
	})

	JustBeforeEach(func() {
		err = RemoveFinalizer(ctx, clt, cm, testFinalizerKey, callback)
	})

	Context("resource has no finalizer", func() {
		BeforeEach(func() {
			testing2.MustLoadYaml("./testdata/configmap.yaml", cm)
			Expect(clt.Create(ctx, cm)).To(Succeed())
			var executed bool
			callback = func() error {
				executed = true
				return nil
			}

			DeferCleanup(func() {
				Expect(executed).To(BeFalse())
			})
		})
		It("no need to remove finalizer", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(cm.GetFinalizers()).NotTo(ContainElement(testFinalizerKey))
		})
	})

	Context("resource has finalizer already", func() {
		BeforeEach(func() {
			testing2.MustLoadYaml("./testdata/configmap.withFinalizer.yaml", cm)
			Expect(clt.Create(ctx, cm)).To(Succeed())
		})
		When("callback succeed", func() {
			BeforeEach(func() {
				var executed bool
				callback = func() error {
					executed = true
					return nil
				}

				DeferCleanup(func() {
					Expect(executed).To(BeTrue())
				})
			})

			It("should remove finalizer", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(cm.GetFinalizers()).NotTo(ContainElement(testFinalizerKey))
			})
		})

		When("callback failed", func() {
			BeforeEach(func() {
				callback = func() error {
					return errors.New("test-error")
				}
			})
			It("no need to remove finalizer", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("test-error"))
			})
		})
	})
})
