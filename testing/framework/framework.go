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

package framework

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	"knative.dev/pkg/logging"
	. "github.com/katanomi/pkg/testing/framework/base"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// fmw global variable to used by different test cases
var fmw = &Framework{}

// New sets a name to framework
func New(name string) *Framework {
	fmw.Name = name
	fmw.init()
	return fmw
}

// Framework base framework for running automated test cases
type Framework struct {
	Name string

	TestContext

	configures []Configure
}

func (f *Framework) init() {
	f.Context = context.Background()

	logger, err := zap.NewDevelopment(zap.ErrorOutput(zapcore.AddSync(GinkgoWriter)))
	if err != nil {
		panic(err)
	}
	f.SugaredLogger = logger.Sugar()
	f.Context = logging.WithLogger(f.Context, f.SugaredLogger)
}

// MRun main testing.M run
func (f *Framework) MRun(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	result := m.Run()
	os.Exit(result)
}

// Run start tests
func (f *Framework) Extensions(extensions ...SharedExtension) *Framework {
	for _, extension := range extensions {
		f.Context = extension.SetShardInfo(f.Context)
	}
	return f
}

// Config register configuration mutators which executed before case running
func (f *Framework) Config(configures ...Configure) *Framework {
	f.configures = append(f.configures, configures...)
	return f
}

// Run start tests
func (f *Framework) Run(t *testing.T) {
	sConfig, rConfig := GinkgoConfiguration()
	for _, configure := range f.configures {
		configure.Config(&sConfig, &rConfig)
	}

	RegisterFailHandler(Fail)
	RunSpecs(t, f.Name, sConfig, rConfig)
}

func (f *Framework) WithContext(ctx context.Context) *Framework {
	f.Context = ctx
	return f
}

// SynchronizedBeforeSuite basic before suite initialization
func (f *Framework) SynchronizedBeforeSuite(initFunc func()) *Framework {
	SynchronizedBeforeSuite(func() []byte {
		By("Setup")
		if initFunc != nil {
			By("Setup.Func")
			initFunc()
		}
		return nil
	}, func(_ []byte) {
		// no-op for now
	})
	return f
}

// SynchronizedAfterSuite destroys the whole environment
func (f *Framework) SynchronizedAfterSuite(destroyFunc func()) *Framework {
	SynchronizedAfterSuite(func() {}, func() {
		By("Teardown")
		if destroyFunc != nil {
			By("Teardown.Func")
			destroyFunc()
		}
	})
	return f
}
