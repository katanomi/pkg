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
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/emicklei/go-restful/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var k8sClient client.Client
var k8sConfigSet *kubernetes.Clientset
var testEnv *envtest.Environment
var originalPropagator = otel.GetTextMapPropagator()
var originalTraceProvider = otel.GetTracerProvider()

func TestPkg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tracing Suite")
}

var _ = AfterEach(func() {
	if !reflect.DeepEqual(otel.GetTextMapPropagator(), originalPropagator) {
		otel.SetTextMapPropagator(originalPropagator)
	}
	if !reflect.DeepEqual(otel.GetTracerProvider(), originalTraceProvider) {
		otel.SetTracerProvider(originalTraceProvider)
	}
})

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{}

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	// +kubebuilder:scaffold:scheme
	k8sClient, err = client.New(cfg, client.Options{})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	k8sConfigSet, err = kubernetes.NewForConfig(cfg)
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sConfigSet).NotTo(BeNil())
})

var _ = AfterSuite(func() {
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

func testWebService(path string, filter restful.FilterFunction, handler restful.RouteFunction) (
	*restful.Container, *http.Request) {
	ws := &restful.WebService{}
	ws.Route(ws.GET(path).To(handler))

	container := restful.NewContainer()
	container.Filter(filter)
	container.Add(ws)

	r := httptest.NewRequest("GET", path, nil)
	return container, r
}

func getValidTracingConfig() *Config {
	return &Config{
		Enable:        true,
		SamplingRatio: 1,
		Backend:       ExporterBackendZipkin,
		Zipkin: ZipkinConfig{
			Url: "127.0.0.1:11211",
		},
	}
}
