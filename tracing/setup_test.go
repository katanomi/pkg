package tracing

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/configmap/informer"
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

var _ = Describe("testing for configmap watcher", func() {
	var clear func()

	// if a default configMap is provided, it will first trigger the update
	// handler when the informer starts
	var startInformer func(*v1.ConfigMap, func(name string, value interface{}))
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

	AfterEach(func() {
		if clear != nil {
			clear()
		}
	})

	BeforeEach(func() {
		stopCh := make(chan struct{})
		testCount = 0
		var createdConfigMaps []string

		clear = func() {
			close(stopCh)
			for _, name := range createdConfigMaps {
				err := k8sConfigSet.CoreV1().ConfigMaps(testNs).Delete(ctx, name, metav1.DeleteOptions{})
				Expect(err).ShouldNot(HaveOccurred())
			}
		}

		createCm = func(name string, cm *v1.ConfigMap) {
			newCm := cm.DeepCopy()
			newCm.Name = name
			_, err := k8sConfigSet.CoreV1().ConfigMaps(testNs).Create(ctx, newCm, metav1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			createdConfigMaps = append(createdConfigMaps, name)
		}

		startInformer = func(defaultCm *v1.ConfigMap, handler func(name string, value interface{})) {
			watcher := informer.NewInformedWatcher(k8sConfigSet, testNs)
			dftCmWatcher := NewDftConfigMapWatcher("config-tracing-store", logger, watcher)
			dftCmWatcher.AddWatch(testCmName, newTracingConfigFromConfigMap, defaultCm)
			dftCmWatcher.Run(func(name string, value interface{}) {
				if name == testCmName {
					atomic.AddInt64(&testCount, 1)
				}
				if handler != nil {
					handler(name, value)
				}
			})
			err := watcher.Start(stopCh)
			Expect(err).ShouldNot(HaveOccurred())
		}
	})

	Context("when create new configmap", func() {
		It("configmap that is not listening will not trigger update handler", func() {
			startInformer(testCm, nil)

			for i := 0; i < 10; i++ {
				createCm(fmt.Sprintf("test-cm-%d", i), testCm)
			}
			// wait for async trigger
			time.Sleep(time.Second * 2)
			Expect(testCount).Should(Equal(int64(triggeredByDefaultCount)))
		})

	})

	Context("when configmap is existed", func() {
		It("The update handler should be triggered when the informer started", func() {
			createCm(testCmName, testCm)
			startInformer(nil, func(name string, value interface{}) {
				cfg, ok := value.(*Config)
				Expect(ok).Should(BeTrue())
				Expect(cfg.Enable).Should(BeTrue())
				Expect(cfg.Backend).Should(Equal(ExporterBackendJaeger))
				Expect(cfg.SamplingRatio).Should(Equal(float64(1)))
			})
			// wait for async trigger
			time.Sleep(time.Second * 2)
			Expect(testCount).Should(Equal(int64(1)))
		})
	})

	Context("when configmap is not exist", func() {
		It("The update handler should be triggered with empty config", func() {
			defaultCM := testConfigMap(testCmName, map[string]string{
				enableKey: "false",
			})
			var currentConfig *Config
			startInformer(defaultCM, func(name string, value interface{}) {
				cfg, ok := value.(*Config)
				Expect(ok).Should(BeTrue())
				currentConfig = cfg
			})
			// wait for async trigger
			time.Sleep(time.Second * 2)
			Expect(testCount).Should(Equal(int64(triggeredByDefaultCount)))
			Expect(currentConfig.Enable).Should(BeFalse())
			Expect(currentConfig.Backend).Should(BeEmpty())

			By("create the configmap")
			createCm(testCmName, testCm)
			// wait for async trigger
			time.Sleep(time.Second * 2)
			Expect(testCount).Should(Equal(int64(1 + triggeredByDefaultCount)))
			Expect(currentConfig.Enable).Should(BeTrue())
			Expect(currentConfig.Backend).Should(Equal(ExporterBackendJaeger))
		})
	})
})
