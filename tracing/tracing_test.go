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
	"go.uber.org/zap"
)

func TestManager_Tracer(t *testing.T) {
	config := map[string]string{
		"tracer-config-key": `{
			"jaegerUrl":"127.0.0.1:1234",
			"sampleType":"const",
			"sampleParam": 1
		}`,
		"sampler.test": `{
			"sampleType": "const",
			"sampleParam": 0.5
		}`,
	}

	configManager, _ := ParseConfig(config)

	manager := &Manager{
		configManager: configManager,
		logger:        zap.S(),
	}

	g := NewGomegaWithT(t)

	tracer, err := manager.Tracer("test")

	g.Expect(err).To(BeNil())
	g.Expect(tracer).ToNot(BeNil())
}

func TestManager_Global(t *testing.T) {
	config := map[string]string{
		"tracer-config-key": `{
			"jaegerUrl":"127.0.0.1:1234",
			"sampleType":"const",
			"sampleParam": 1
		}`,
	}

	configManager, _ := ParseConfig(config)

	manager := &Manager{
		configManager: configManager,
		logger:        zap.S(),
	}

	g := NewGomegaWithT(t)

	tracer, err := manager.Tracer("test")

	g.Expect(err).To(BeNil())
	g.Expect(tracer).ToNot(BeNil())
}

func TestManager_Empty(t *testing.T) {
	config := map[string]string{}

	configManager, _ := ParseConfig(config)

	manager := &Manager{
		configManager: configManager,
		logger:        zap.S(),
	}

	g := NewGomegaWithT(t)

	tracer, err := manager.Tracer("test")

	g.Expect(err).To(BeNil())
	g.Expect(tracer).To(BeNil())
}

func TestManager_Sync(t *testing.T) {
	config := map[string]string{
		"tracer-config-key": `{
			"jaegerUrl":"127.0.0.1:1234",
			"sampleType":"const",
			"sampleParam": 1
		}`,
		"sampler.test": `{
			"sampleType": "const",
			"sampleParam": 0.5
		}`,
	}

	configManager, _ := ParseConfig(config)

	manager := &Manager{
		configManager: configManager,
		logger:        zap.S(),
	}

	g := NewGomegaWithT(t)

	tracer, err := manager.Tracer("test")

	g.Expect(err).To(BeNil())
	g.Expect(tracer).ToNot(BeNil())

	manager.Sync()
	_, loaded := manager.tracers.Load("test")

	g.Expect(loaded).To(BeFalse())
}
