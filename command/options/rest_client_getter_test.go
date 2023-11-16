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

package options

import (
	"context"
	"fmt"

	"github.com/golang/mock/gomock"
	mockopt "github.com/katanomi/pkg/testing/mock/k8s.io/cli-runtime/pkg/genericclioptions"
	mockclientcmd "github.com/katanomi/pkg/testing/mock/k8s.io/client-go/tools/clientcmd"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
)

var _ = Describe("Test.RESTClientGetterOption.Setup", func() {

	var (
		ctx context.Context
		opt *RESTClientGetterOption
	)

	BeforeEach(func() {
		ctx = context.Background()
		opt = &RESTClientGetterOption{}
	})

	JustBeforeEach(func() {
		opt.Setup(ctx, &cobra.Command{}, []string{})
	})
	It("init configFlag", func() {
		Expect(opt.ConfigFlag).NotTo(BeNil())
	})

})

var _ = Describe("Test.RESTClientGetterOption.Funcs", func() {

	var (
		ctx         context.Context
		opt         *RESTClientGetterOption
		mockCfgFlag *mockopt.MockRESTClientGetter
		mockCliCfg  *mockclientcmd.MockClientConfig
		err         error
		ret         string
	)

	BeforeEach(func() {
		ctx = context.Background()
		opt = &RESTClientGetterOption{}
		mockCtl := gomock.NewController(GinkgoT())
		mockCfgFlag = mockopt.NewMockRESTClientGetter(mockCtl)
		mockCliCfg = mockclientcmd.NewMockClientConfig(mockCtl)
		opt.ConfigFlag = mockCfgFlag

		err = nil
		ret = ""
	})

	When("GetClusterToken", func() {

		JustBeforeEach(func() {
			ret, err = opt.GetClusterToken(ctx)
		})

		Context("get kubeconfig failed", func() {
			expErr := fmt.Errorf("foo err")
			BeforeEach(func() {
				mockCfgFlag.EXPECT().ToRESTConfig().Return(nil, expErr)
			})

			It("returns empty token and err", func() {
				Expect(ret).To(BeEmpty())
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(expErr.Error()))
			})
		})

		Context("kubeconfig contains bearer token", func() {
			BeforeEach(func() {
				mockCfgFlag.EXPECT().ToRESTConfig().Return(&rest.Config{
					BearerToken: "xxx",
				}, nil)
			})

			It("return token", func() {
				Expect(err).To(BeNil())
				Expect(ret).To(Equal("xxx"))
			})
		})
	})

	When("GetNamespace", func() {
		BeforeEach(func() {

			mockCfgFlag.EXPECT().ToRawKubeConfigLoader().Return(mockCliCfg)
			mockCliCfg.EXPECT().Namespace().Return("yyy", false, nil)

		})
		JustBeforeEach(func() {
			ret, err = opt.GetNamespace()
		})

		It("returns namespace result", func() {
			Expect(err).NotTo(HaveOccurred())

			Expect(ret).To(Equal("yyy"))
		})

	})

})
