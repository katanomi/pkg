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
	"io"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"knative.dev/pkg/injection"
	ctrl "sigs.k8s.io/controller-runtime"
)

// fmw global variable to used by different test cases
var fmw = &Framework{}

// Framework base framework for running automated test cases
type Framework struct {
	Name string

	Config *rest.Config

	Scheme *runtime.Scheme

	Context context.Context

	*zap.SugaredLogger

	Output io.Writer

	InitTimeout time.Duration

	sync.Once
}

// New sets a name to framework
func New(name string) *Framework {
	fmw.Name = name
	fmw.init()
	return fmw
}

// init is a do once initialization function to startup any necessary data
func (f *Framework) init() {
	f.Once.Do(func() {
		f.Context = context.TODO()
		cfg := ctrl.GetConfigOrDie()
		f.Context = injection.WithConfig(f.Context, cfg)
		f.Config = cfg

		logger, err := zap.NewDevelopment(zap.ErrorOutput(zapcore.AddSync(GinkgoWriter)))
		if err != nil {
			panic(err)
		}
		f.SugaredLogger = logger.Sugar()
	})
}

// MRun main testing.M run
func (f *Framework) MRun(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	result := m.Run()
	os.Exit(result)
}

// Run start tests
func (f *Framework) Run(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, f.Name)
}

// WithScheme adds a scheme object to the framework
func (f *Framework) WithScheme(scheme *runtime.Scheme) *Framework {
	f.Scheme = scheme
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

// DurationToFloat converts a duration into a float64 seconds, useful for Ginkgo methods
func DurationToFloat(dur time.Duration) float64 {
	return float64(dur) / float64(time.Second)
}
