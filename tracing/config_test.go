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

package tracing

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestParseConfig(t *testing.T) {
	config := map[string]string{
		"tracer-config-key": `{
			"jaegerUrl":"127.0.0.1:1234",
			"sampleType":"const",
			"sampleParam": 1
		}`,
	}

	g := NewGomegaWithT(t)

	configManager, err := ParseConfig(config)
	g.Expect(err).To(BeNil())

	g.Expect(configManager.global).ToNot(BeNil())

	g.Expect(configManager.global.JaegerUrl).To(Equal("127.0.0.1:1234"))
	g.Expect(configManager.global.SampleType).To(Equal("const"))
	g.Expect(configManager.global.SampleParam).To(Equal(float64(1)))
}

func TestParseConfig_Empty(t *testing.T) {
	config := map[string]string{}

	g := NewGomegaWithT(t)

	configManager, err := ParseConfig(config)
	g.Expect(err).To(BeNil())

	g.Expect(configManager.global).To(BeNil())
}

func TestGet(t *testing.T) {
	config := map[string]string{
		"tracer-config-key": `{
			"jaegerUrl":"127.0.0.1:1234",
			"sampleType":"const",
			"sampleParam": 1
		}`,
		"sampler.test": `{
			"sampleType": "probabilistic",
			"sampleParam": 0.5
		}`,
	}

	g := NewGomegaWithT(t)

	configManager, err := ParseConfig(config)
	g.Expect(err).To(BeNil())

	testConfig := configManager.Get("test")

	g.Expect(testConfig).ToNot(BeNil())
	g.Expect(testConfig.JaegerUrl).To(Equal("127.0.0.1:1234"))
	g.Expect(testConfig.SampleType).To(Equal("probabilistic"))
	g.Expect(testConfig.SampleParam).To(Equal(float64(0.5)))
}
