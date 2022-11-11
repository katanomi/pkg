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

package sharedmain

import (
	"flag"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = BeforeSuite(func() {
	ParseFlag()
})

var _ = Describe("ParseFlag", func() {

	When("flag not provided", func() {
		It("return default values", func() {
			Expect(QPS).To(Equal(float64(DefaultQPS)))
			Expect(Burst).To(Equal(DefaultBurst))
			Expect(Timeout).To(Equal(DefaultTimeout))
			Expect(ConfigFile).To(Equal(""))
			Expect(InsecureSkipVerify).To(Equal(false))
		})
	})

	When("flag proviede", func() {
		BeforeEach(func() {
			flag.CommandLine.Parse([]string{
				"--kube-api-timeout", "20s",
				"--kube-api-qps", "80",
				"--kube-api-burst", "90",
				"--config", "config",
				"--insecure-skip-tls-verify", "true",
			})
		})
		It("return configured values", func() {
			Expect(QPS).To(Equal(float64(80)))
			Expect(Burst).To(Equal(90))
			Expect(Timeout).To(Equal(20 * time.Second))
			Expect(ConfigFile).To(Equal("config"))
			Expect(InsecureSkipVerify).To(Equal(true))
		})
	})

})
