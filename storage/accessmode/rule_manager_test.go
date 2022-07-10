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

package accessmode

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/configmap/informer"
)

var _ = Describe("Test DynamicAccessModeManager", func() {
	var (
		manager AccessModeManager
		testNs  = "default"
		logger  = zap.NewNop().Sugar()
		ctx     = context.Background()
	)

	newTestSC := func(name, provisioner string) *v1.StorageClass {
		sc := &v1.StorageClass{}
		sc.Name = name
		sc.Provisioner = provisioner
		return sc
	}

	newConfigMap := func(data map[string][]corev1.PersistentVolumeAccessMode) *corev1.ConfigMap {
		ruleData, _ := json.Marshal(data)

		cm := dftCm()
		cm.Namespace = testNs
		cm.Data = map[string]string{
			rulesKey: string(ruleData),
		}
		return cm
	}

	BeforeEach(func() {
		stopCh := make(chan struct{})

		watcher := informer.NewInformedWatcher(k8sConfigSet, testNs)
		manager = NewDynamicAccessModeManager(logger, watcher)
		err := watcher.Start(stopCh)
		Expect(err).ToNot(HaveOccurred())
	})

	Context("there is no configmap defined in the namespace", func() {
		It("should get default rules", func() {
			acList := manager.SupportedAccessModes(newTestSC("xx", StorageosProvisioner))
			Expect(acList).To(BeEquivalentTo(dftAccessModeRules()[StorageosProvisioner]))

			acList = manager.SupportedAccessModes(nil)
			Expect(acList).To(BeEmpty())

			acList = manager.SupportedAccessModes(newTestSC("not-exist", "not-exist"))
			Expect(acList).To(BeEmpty())
		})
	})

	Context("when configmap is exist", func() {
		It("should get the rules defined in configmap", func() {
			cm := newConfigMap(map[string][]corev1.PersistentVolumeAccessMode{
				StorageosProvisioner: {corev1.ReadWriteOnce},
			})
			_, err := k8sConfigSet.CoreV1().ConfigMaps(testNs).Create(ctx, cm, metav1.CreateOptions{})
			Expect(err).ToNot(HaveOccurred())
			time.Sleep(time.Second * 2)

			acList := manager.SupportedAccessModes(newTestSC("xx", StorageosProvisioner))
			Expect(len(acList)).To(Equal(1))
			Expect(acList[0]).To(Equal(corev1.ReadWriteOnce))

			By("update the configmap")
			cm = newConfigMap(map[string][]corev1.PersistentVolumeAccessMode{
				StorageosProvisioner: {corev1.ReadWriteMany},
			})
			_, err = k8sConfigSet.CoreV1().ConfigMaps(testNs).Update(ctx, cm, metav1.UpdateOptions{})
			Expect(err).ToNot(HaveOccurred())
			time.Sleep(time.Second * 2)

			acList = manager.SupportedAccessModes(newTestSC("xx", StorageosProvisioner))
			Expect(len(acList)).To(Equal(1))
			Expect(acList[0]).To(Equal(corev1.ReadWriteMany))
		})
	})
})

func TestConfigMapName(t *testing.T) {
	tests := []struct {
		setting func()
		want    string
	}{
		{
			setting: func() {
				os.Setenv(configMapNameEnv, "test-config-name")
			},
			want: "test-config-name",
		},
		{
			setting: func() {
				os.Unsetenv(configMapNameEnv)
			},
			want: defaultConfigMapName,
		},
	}
	g := NewGomegaWithT(t)
	for _, tt := range tests {
		if tt.setting != nil {
			tt.setting()
		}
		got := ConfigMapName()
		g.Expect(got).Should(Equal(tt.want))
	}
}
