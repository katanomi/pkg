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

package configmap

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/configmap/informer"
)

const (
	enableKey        = "enable"
	backendKey       = "backend"
	samplingRatioKey = "sampling-ratio"
)

func testConfigMap(name string, data map[string]string) *v1.ConfigMap {
	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Data: data,
	}
}

// legacy kubebuilder test environment is disabled, this test should be skipped
var _ = PDescribe("testing for configmap watcher", func() {
	// if a default configMap is provided, it will first trigger the update
	// handler when the informer starts
	var startInformer func(string, *v1.ConfigMap, func(cm *v1.ConfigMap))
	var createCm func(name string, cm *v1.ConfigMap)
	var ctx = context.Background()
	var logger = zap.NewNop().Sugar()
	var testCmName = "test-config-name"
	var testNs = "default"
	var testCm = testConfigMap(testCmName, map[string]string{
		enableKey:        "true",
		backendKey:       "jaeger",
		samplingRatioKey: "1",
	})
	var testCount = int64(0)
	var triggeredByDefaultCount = 1

	BeforeEach(func() {
		stopCh := make(chan struct{})
		testCount = int64(0)
		var createdConfigMaps []string

		createCm = func(name string, cm *v1.ConfigMap) {
			newCm := cm.DeepCopy()
			newCm.Name = name
			_, err := k8sConfigSet.CoreV1().ConfigMaps(testNs).Create(ctx, newCm, metav1.CreateOptions{})
			Expect(err).ToNot(HaveOccurred())
			createdConfigMaps = append(createdConfigMaps, name)
		}

		startInformer = func(cmName string, defaultCm *v1.ConfigMap, handler func(cm *v1.ConfigMap)) {
			watcher := informer.NewInformedWatcher(k8sConfigSet, testNs)
			dftCmWatcher := NewWatcher("config-test", watcher).WithLogger(logger)
			dftCmWatcher.AddWatch(cmName, NewConfigConstructor(defaultCm, func(cm *v1.ConfigMap) {
				atomic.AddInt64(&testCount, 1)
				if handler != nil {
					handler(cm)
				}
			}))
			dftCmWatcher.Run()
			err := watcher.Start(stopCh)
			Expect(err).ToNot(HaveOccurred())
		}

		DeferCleanup(func() {
			close(stopCh)
			for _, name := range createdConfigMaps {
				err := k8sConfigSet.CoreV1().ConfigMaps(testNs).Delete(ctx, name, metav1.DeleteOptions{})
				Expect(err).ToNot(HaveOccurred())
			}
		})
	})

	Context("when create new configmap", func() {
		It("configmap that is not listening will not trigger update handler", func() {
			startInformer(testCm.GetName(), testCm, nil)

			for i := 0; i < 10; i++ {
				createCm(fmt.Sprintf("test-cm-%d", i), testCm)
			}
			// wait for async trigger
			time.Sleep(time.Second * 2)
			Expect(testCount).To(Equal(int64(triggeredByDefaultCount)))
		})
	})

	Context("when configmap is existed", func() {
		It("The update handler should be triggered when the informer started", func() {
			createCm(testCmName, testCm)
			startInformer(testCmName, nil, func(cm *v1.ConfigMap) {
				Expect(cm).NotTo(BeNil())
				Expect(cm.Data).NotTo(BeEmpty())
				Expect(cm.Data[enableKey]).To(Equal("true"))
				Expect(cm.Data[backendKey]).To(Equal("jaeger"))
				Expect(cm.Data[samplingRatioKey]).To(Equal("1"))
			})
			// wait for async trigger
			time.Sleep(time.Second * 2)
			Expect(testCount).To(Equal(int64(1)))
		})
	})

	Context("when configmap is not exist", func() {
		It("The update handler should be triggered with empty config", func() {
			defaultCM := testConfigMap(testCmName, map[string]string{
				enableKey: "false",
			})
			var currentCm *v1.ConfigMap
			startInformer(testCmName, defaultCM, func(cm *v1.ConfigMap) {
				currentCm = cm
			})
			// wait for async trigger
			time.Sleep(time.Second * 2)
			Expect(testCount).To(Equal(int64(triggeredByDefaultCount)))
			Expect(currentCm).NotTo(BeNil())
			Expect(currentCm.Data).NotTo(BeEmpty())
			Expect(currentCm.Data[enableKey]).To(Equal("false"))
			Expect(currentCm.Data[backendKey]).To(BeEmpty())

			By("create the configmap")
			createCm(testCmName, testCm)
			// wait for async trigger
			time.Sleep(time.Second * 2)
			Expect(testCount).To(Equal(int64(1 + triggeredByDefaultCount)))
			Expect(currentCm).NotTo(BeNil())
			Expect(currentCm.Data).NotTo(BeEmpty())
			Expect(currentCm.Data[enableKey]).To(Equal("true"))
			Expect(currentCm.Data[backendKey]).To(Equal("jaeger"))
		})
	})

	Context("testing for multi configmap watcher", func() {
		var testCm1, testCm2 *v1.ConfigMap
		var testCount1, testCount2 int64
		BeforeEach(func() {
			testCount1, testCount2 = 0, 0
			stopCh := make(chan struct{})
			testCm1 = testConfigMap("test1", map[string]string{
				"key": "value1",
			})
			testCm2 = testConfigMap("test2", map[string]string{
				"key": "value2",
			})
			watcher := informer.NewInformedWatcher(k8sConfigSet, testNs)
			dftCmWatcher := NewWatcher("config-test", watcher).WithLogger(logger)
			dftCmWatcher.AddWatch(testCm1.Name, NewConfigConstructor(testCm1, func(cm *v1.ConfigMap) {
				atomic.AddInt64(&testCount1, 1)
				Expect(cm.GetName()).To(Equal(testCm1.Name))
			}))
			dftCmWatcher.AddWatch(testCm2.Name, NewConfigConstructor(testCm2, func(cm *v1.ConfigMap) {
				atomic.AddInt64(&testCount2, 1)
				Expect(cm.GetName()).To(Equal(testCm2.Name))
			}))
			dftCmWatcher.Run()
			err := watcher.Start(stopCh)
			Expect(err).ToNot(HaveOccurred())

			DeferCleanup(func() {
				close(stopCh)
			})
		})

		When("no configmap in the namespace", func() {
			It("The update handler should be triggered with empty config", func() {
				Expect(testCount1).To(Equal(int64(1)))
				Expect(testCount2).To(Equal(int64(1)))
			})
		})

		When("configmap exist in the namespace", func() {
			It("The update handler should be triggered", func() {
				createCm(testCm1.Name, testCm1)
				Eventually(func(g Gomega) error {
					g.Expect(testCount1).To(Equal(int64(2)))
					g.Expect(testCount2).To(Equal(int64(1)))
					return nil
				}).WithPolling(time.Second).WithTimeout(time.Second * 5).Should(Succeed())

				createCm(testCm2.Name, testCm2)
				Eventually(func(g Gomega) error {
					g.Expect(testCount1).To(Equal(int64(2)))
					g.Expect(testCount2).To(Equal(int64(2)))
					return nil
				}).WithPolling(time.Second).WithTimeout(time.Second * 5).Should(Succeed())
			})
		})
	})
})
