/*
Copyright 2022 The Katanomi Authors.

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

package secret

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	uberzap "go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func TestSecret(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Secret Suite")
}

var (
	scheme = runtime.NewScheme()
	logger *uberzap.SugaredLogger
)

func init() {
	corev1.AddToScheme(scheme)
}

func init() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	logger = zap.NewRaw(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)).Sugar()
}

var testLogger = zap.NewRaw(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)).Sugar()

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
})

var _ = AfterSuite(func() {
})
